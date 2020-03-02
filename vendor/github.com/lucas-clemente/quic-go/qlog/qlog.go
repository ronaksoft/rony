package qlog

import (
	"io"
	"net"
	"time"

	"github.com/lucas-clemente/quic-go/internal/congestion"

	"github.com/francoispqt/gojay"
	"github.com/lucas-clemente/quic-go/internal/protocol"
	"github.com/lucas-clemente/quic-go/internal/wire"
)

// A Tracer records events to be exported to a qlog.
type Tracer interface {
	Export() error
	StartedConnection(t time.Time, local, remote net.Addr, version protocol.VersionNumber, srcConnID, destConnID protocol.ConnectionID)
	SentPacket(t time.Time, hdr *wire.ExtendedHeader, packetSize protocol.ByteCount, ack *wire.AckFrame, frames []wire.Frame)
	ReceivedRetry(time.Time, *wire.Header)
	ReceivedPacket(t time.Time, hdr *wire.ExtendedHeader, packetSize protocol.ByteCount, frames []wire.Frame)
	UpdatedMetrics(t time.Time, rttStats *congestion.RTTStats, cwnd protocol.ByteCount, bytesInFLight protocol.ByteCount, packetsInFlight int)
	LostPacket(time.Time, protocol.EncryptionLevel, protocol.PacketNumber, PacketLossReason)
}

type tracer struct {
	w           io.WriteCloser
	odcid       protocol.ConnectionID
	perspective protocol.Perspective

	events []event
}

var _ Tracer = &tracer{}

// NewTracer creates a new tracer to record a qlog.
func NewTracer(w io.WriteCloser, p protocol.Perspective, odcid protocol.ConnectionID) Tracer {
	return &tracer{
		w:           w,
		perspective: p,
		odcid:       odcid,
	}
}

func (t *tracer) Active() bool { return true }

// Export writes a qlog.
func (t *tracer) Export() error {
	enc := gojay.NewEncoder(t.w)
	tl := &topLevel{
		traces: traces{
			{
				VantagePoint: vantagePoint{Type: t.perspective},
				CommonFields: commonFields{ODCID: connectionID(t.odcid), GroupID: connectionID(t.odcid)},
				EventFields:  eventFields[:],
				Events:       t.events,
			},
		}}
	if err := enc.Encode(tl); err != nil {
		return err
	}
	return t.w.Close()
}

func (t *tracer) StartedConnection(time time.Time, local, remote net.Addr, version protocol.VersionNumber, srcConnID, destConnID protocol.ConnectionID) {
	// ignore this event if we're not dealing with UDP addresses here
	localAddr, ok := local.(*net.UDPAddr)
	if !ok {
		return
	}
	remoteAddr, ok := remote.(*net.UDPAddr)
	if !ok {
		return
	}
	t.events = append(t.events, event{
		Time: time,
		eventDetails: eventConnectionStarted{
			SrcAddr:          localAddr,
			DestAddr:         remoteAddr,
			Version:          version,
			SrcConnectionID:  srcConnID,
			DestConnectionID: destConnID,
		},
	})
}

func (t *tracer) SentPacket(time time.Time, hdr *wire.ExtendedHeader, packetSize protocol.ByteCount, ack *wire.AckFrame, frames []wire.Frame) {
	numFrames := len(frames)
	if ack != nil {
		numFrames++
	}
	fs := make([]frame, 0, numFrames)
	if ack != nil {
		fs = append(fs, *transformFrame(ack))
	}
	for _, f := range frames {
		fs = append(fs, *transformFrame(f))
	}
	header := *transformExtendedHeader(hdr)
	header.PacketSize = packetSize
	t.events = append(t.events, event{
		Time: time,
		eventDetails: eventPacketSent{
			PacketType: getPacketTypeFromHeader(hdr),
			Header:     header,
			Frames:     fs,
		},
	})
}

func (t *tracer) ReceivedPacket(time time.Time, hdr *wire.ExtendedHeader, packetSize protocol.ByteCount, frames []wire.Frame) {
	fs := make([]frame, len(frames))
	for i, f := range frames {
		fs[i] = *transformFrame(f)
	}
	header := *transformExtendedHeader(hdr)
	header.PacketSize = packetSize
	t.events = append(t.events, event{
		Time: time,
		eventDetails: eventPacketReceived{
			PacketType: getPacketTypeFromHeader(hdr),
			Header:     header,
			Frames:     fs,
		},
	})
}

func (t *tracer) ReceivedRetry(time time.Time, hdr *wire.Header) {
	t.events = append(t.events, event{
		Time: time,
		eventDetails: eventRetryReceived{
			Header: *transformHeader(hdr),
		},
	})
}

func (t *tracer) UpdatedMetrics(time time.Time, rttStats *congestion.RTTStats, cwnd, bytesInFlight protocol.ByteCount, packetsInFlight int) {
	t.events = append(t.events, event{
		Time: time,
		eventDetails: eventMetricsUpdated{
			MinRTT:           rttStats.MinRTT(),
			SmoothedRTT:      rttStats.SmoothedRTT(),
			LatestRTT:        rttStats.LatestRTT(),
			RTTVariance:      rttStats.MeanDeviation(),
			CongestionWindow: cwnd,
			BytesInFlight:    bytesInFlight,
			PacketsInFlight:  packetsInFlight,
		},
	})
}

func (t *tracer) LostPacket(time time.Time, encLevel protocol.EncryptionLevel, pn protocol.PacketNumber, lossReason PacketLossReason) {
	t.events = append(t.events, event{
		Time: time,
		eventDetails: eventPacketLost{
			PacketType:   getPacketTypeFromEncryptionLevel(encLevel),
			PacketNumber: pn,
			Trigger:      lossReason,
		},
	})
}
