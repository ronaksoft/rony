package qlog

import (
	"github.com/francoispqt/gojay"

	"github.com/lucas-clemente/quic-go/internal/protocol"
	"github.com/lucas-clemente/quic-go/internal/wire"
)

func getPacketTypeFromHeader(hdr *wire.ExtendedHeader) packetType {
	if !hdr.IsLongHeader {
		return packetType1RTT
	}
	if hdr.Version == 0 {
		return packetTypeVersionNegotiation
	}
	switch hdr.Type {
	case protocol.PacketTypeInitial:
		return packetTypeInitial
	case protocol.PacketTypeHandshake:
		return packetTypeHandshake
	case protocol.PacketType0RTT:
		return packetType0RTT
	case protocol.PacketTypeRetry:
		return packetTypeRetry
	default:
		panic("unknown packet type")
	}
}

func getPacketTypeFromEncryptionLevel(encLevel protocol.EncryptionLevel) packetType {
	switch encLevel {
	case protocol.EncryptionInitial:
		return packetTypeInitial
	case protocol.EncryptionHandshake:
		return packetTypeHandshake
	case protocol.Encryption0RTT:
		return packetType0RTT
	case protocol.Encryption1RTT:
		return packetType1RTT
	default:
		panic("unknown encryption level")
	}
}

func transformHeader(hdr *wire.Header) *packetHeader {
	return &packetHeader{
		PayloadLength:    hdr.Length,
		SrcConnectionID:  hdr.SrcConnectionID,
		DestConnectionID: hdr.DestConnectionID,
		Version:          hdr.Version,
	}
}

func transformExtendedHeader(hdr *wire.ExtendedHeader) *packetHeader {
	h := transformHeader(&hdr.Header)
	h.PacketNumber = hdr.PacketNumber
	return h
}

type packetHeader struct {
	PacketNumber  protocol.PacketNumber
	PayloadLength protocol.ByteCount
	// Size of the QUIC packet (QUIC header + payload).
	// See https://github.com/quiclog/internet-drafts/issues/40.
	PacketSize protocol.ByteCount

	Version          protocol.VersionNumber
	SrcConnectionID  protocol.ConnectionID
	DestConnectionID protocol.ConnectionID
}

func (h packetHeader) MarshalJSONObject(enc *gojay.Encoder) {
	enc.StringKey("packet_number", toString(int64(h.PacketNumber)))
	enc.Int64KeyOmitEmpty("payload_length", int64(h.PayloadLength))
	enc.Int64KeyOmitEmpty("packet_size", int64(h.PacketSize))
	if h.Version != 0 {
		enc.StringKey("version", versionNumber(h.Version).String())
	}
	if h.SrcConnectionID.Len() > 0 {
		enc.StringKey("scil", toString(int64(h.SrcConnectionID.Len())))
		enc.StringKey("scid", connectionID(h.SrcConnectionID).String())
	}
	if h.DestConnectionID.Len() > 0 {
		enc.StringKey("dcil", toString(int64(h.DestConnectionID.Len())))
		enc.StringKey("dcid", connectionID(h.DestConnectionID).String())
	}
}

func (packetHeader) IsNil() bool { return false }
