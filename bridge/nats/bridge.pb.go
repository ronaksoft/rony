// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: bridge.proto

package natsBridge

import (
	fmt "fmt"
	github_com_gogo_protobuf_proto "github.com/gogo/protobuf/proto"
	proto "github.com/gogo/protobuf/proto"
	io "io"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

// MessageContainer
type Container struct {
	MessageID uint64     `protobuf:"varint,1,req,name=MessageID" json:"MessageID"`
	Messages  []*Message `protobuf:"bytes,2,rep,name=Messages" json:"Messages,omitempty"`
}

func (m *Container) Reset()         { *m = Container{} }
func (m *Container) String() string { return proto.CompactTextString(m) }
func (*Container) ProtoMessage()    {}
func (*Container) Descriptor() ([]byte, []int) {
	return fileDescriptor_1d3ed31acb30cd14, []int{0}
}
func (m *Container) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Container) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Container.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Container) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Container.Merge(m, src)
}
func (m *Container) XXX_Size() int {
	return m.Size()
}
func (m *Container) XXX_DiscardUnknown() {
	xxx_messageInfo_Container.DiscardUnknown(m)
}

var xxx_messageInfo_Container proto.InternalMessageInfo

func (m *Container) GetMessageID() uint64 {
	if m != nil {
		return m.MessageID
	}
	return 0
}

func (m *Container) GetMessages() []*Message {
	if m != nil {
		return m.Messages
	}
	return nil
}

// Message
type Message struct {
	Payload []byte `protobuf:"bytes,2,req,name=Payload" json:"Payload"`
	AuthID  int64  `protobuf:"varint,3,req,name=AuthID" json:"AuthID"`
	UserID  int64  `protobuf:"varint,4,opt,name=UserID" json:"UserID"`
}

func (m *Message) Reset()         { *m = Message{} }
func (m *Message) String() string { return proto.CompactTextString(m) }
func (*Message) ProtoMessage()    {}
func (*Message) Descriptor() ([]byte, []int) {
	return fileDescriptor_1d3ed31acb30cd14, []int{1}
}
func (m *Message) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Message) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Message.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Message) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Message.Merge(m, src)
}
func (m *Message) XXX_Size() int {
	return m.Size()
}
func (m *Message) XXX_DiscardUnknown() {
	xxx_messageInfo_Message.DiscardUnknown(m)
}

var xxx_messageInfo_Message proto.InternalMessageInfo

func (m *Message) GetPayload() []byte {
	if m != nil {
		return m.Payload
	}
	return nil
}

func (m *Message) GetAuthID() int64 {
	if m != nil {
		return m.AuthID
	}
	return 0
}

func (m *Message) GetUserID() int64 {
	if m != nil {
		return m.UserID
	}
	return 0
}

// Notifier
type Notifier struct {
	MessageID     uint64   `protobuf:"varint,1,req,name=MessageID" json:"MessageID"`
	ConnectionIDs []uint64 `protobuf:"varint,2,rep,name=ConnectionIDs" json:"ConnectionIDs,omitempty"`
}

func (m *Notifier) Reset()         { *m = Notifier{} }
func (m *Notifier) String() string { return proto.CompactTextString(m) }
func (*Notifier) ProtoMessage()    {}
func (*Notifier) Descriptor() ([]byte, []int) {
	return fileDescriptor_1d3ed31acb30cd14, []int{2}
}
func (m *Notifier) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Notifier) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Notifier.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Notifier) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Notifier.Merge(m, src)
}
func (m *Notifier) XXX_Size() int {
	return m.Size()
}
func (m *Notifier) XXX_DiscardUnknown() {
	xxx_messageInfo_Notifier.DiscardUnknown(m)
}

var xxx_messageInfo_Notifier proto.InternalMessageInfo

func (m *Notifier) GetMessageID() uint64 {
	if m != nil {
		return m.MessageID
	}
	return 0
}

func (m *Notifier) GetConnectionIDs() []uint64 {
	if m != nil {
		return m.ConnectionIDs
	}
	return nil
}

// DeliveryReport
type DeliveryReport struct {
	MessageIDs []uint64 `protobuf:"varint,1,rep,name=MessageIDs" json:"MessageIDs,omitempty"`
}

func (m *DeliveryReport) Reset()         { *m = DeliveryReport{} }
func (m *DeliveryReport) String() string { return proto.CompactTextString(m) }
func (*DeliveryReport) ProtoMessage()    {}
func (*DeliveryReport) Descriptor() ([]byte, []int) {
	return fileDescriptor_1d3ed31acb30cd14, []int{3}
}
func (m *DeliveryReport) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *DeliveryReport) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_DeliveryReport.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *DeliveryReport) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeliveryReport.Merge(m, src)
}
func (m *DeliveryReport) XXX_Size() int {
	return m.Size()
}
func (m *DeliveryReport) XXX_DiscardUnknown() {
	xxx_messageInfo_DeliveryReport.DiscardUnknown(m)
}

var xxx_messageInfo_DeliveryReport proto.InternalMessageInfo

func (m *DeliveryReport) GetMessageIDs() []uint64 {
	if m != nil {
		return m.MessageIDs
	}
	return nil
}

func init() {
	proto.RegisterType((*Container)(nil), "natsBridge.Container")
	proto.RegisterType((*Message)(nil), "natsBridge.Message")
	proto.RegisterType((*Notifier)(nil), "natsBridge.Notifier")
	proto.RegisterType((*DeliveryReport)(nil), "natsBridge.DeliveryReport")
}

func init() { proto.RegisterFile("bridge.proto", fileDescriptor_1d3ed31acb30cd14) }

var fileDescriptor_1d3ed31acb30cd14 = []byte{
	// 263 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x8f, 0xb1, 0x4e, 0xc3, 0x30,
	0x10, 0x86, 0xe3, 0x24, 0xa2, 0xed, 0x51, 0x18, 0xcc, 0xe2, 0x01, 0x99, 0x28, 0x62, 0xc8, 0x14,
	0x10, 0x6f, 0x40, 0xea, 0x25, 0x03, 0x08, 0x45, 0xb0, 0x63, 0xe8, 0x51, 0x2c, 0x55, 0x76, 0x65,
	0x1b, 0xa4, 0xbe, 0x05, 0x8f, 0xd5, 0xb1, 0x23, 0x13, 0x42, 0xc9, 0x8b, 0x20, 0x5a, 0xd3, 0xc0,
	0xc6, 0x76, 0xf7, 0x7d, 0x77, 0xbf, 0xf4, 0xc3, 0xf8, 0xc1, 0xaa, 0xe9, 0x0c, 0xcb, 0x85, 0x35,
	0xde, 0x50, 0xd0, 0xd2, 0xbb, 0x6a, 0x43, 0xf2, 0x7b, 0x18, 0x4d, 0x8c, 0xf6, 0x52, 0x69, 0xb4,
	0x34, 0x87, 0xd1, 0x15, 0x3a, 0x27, 0x67, 0x58, 0x0b, 0x46, 0xb2, 0xb8, 0x48, 0xab, 0x74, 0xf5,
	0x71, 0x12, 0x35, 0x3d, 0xa6, 0x67, 0x30, 0x0c, 0x8b, 0x63, 0x71, 0x96, 0x14, 0xfb, 0x17, 0x47,
	0x65, 0x9f, 0x57, 0x06, 0xd7, 0xec, 0x8e, 0x72, 0x84, 0x41, 0x98, 0x29, 0x87, 0xc1, 0x8d, 0x5c,
	0xce, 0x8d, 0x9c, 0xb2, 0x38, 0x8b, 0x8b, 0x71, 0x48, 0xff, 0x81, 0xf4, 0x18, 0xf6, 0x2e, 0x5f,
	0xfc, 0x73, 0x2d, 0x58, 0x92, 0xc5, 0x45, 0x12, 0x74, 0x60, 0xdf, 0xf6, 0xce, 0xa1, 0xad, 0x05,
	0x4b, 0x33, 0xd2, 0xdb, 0x2d, 0xcb, 0x6f, 0x61, 0x78, 0x6d, 0xbc, 0x7a, 0x52, 0xff, 0xec, 0x71,
	0x0a, 0x07, 0x13, 0xa3, 0x35, 0x3e, 0x7a, 0x65, 0x74, 0x2d, 0xb6, 0x65, 0xd2, 0xe6, 0x2f, 0xcc,
	0xcf, 0xe1, 0x50, 0xe0, 0x5c, 0xbd, 0xa2, 0x5d, 0x36, 0xb8, 0x30, 0xd6, 0x53, 0x0e, 0xb0, 0x0b,
	0x71, 0x8c, 0x6c, 0x9e, 0x7e, 0x91, 0x8a, 0xad, 0x5a, 0x4e, 0xd6, 0x2d, 0x27, 0x9f, 0x2d, 0x27,
	0x6f, 0x1d, 0x8f, 0xd6, 0x1d, 0x8f, 0xde, 0x3b, 0x1e, 0x7d, 0x05, 0x00, 0x00, 0xff, 0xff, 0x3c,
	0xf6, 0x51, 0x57, 0x85, 0x01, 0x00, 0x00,
}

func (m *Container) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Container) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	dAtA[i] = 0x8
	i++
	i = encodeVarintBridge(dAtA, i, uint64(m.MessageID))
	if len(m.Messages) > 0 {
		for _, msg := range m.Messages {
			dAtA[i] = 0x12
			i++
			i = encodeVarintBridge(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	return i, nil
}

func (m *Message) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Message) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Payload != nil {
		dAtA[i] = 0x12
		i++
		i = encodeVarintBridge(dAtA, i, uint64(len(m.Payload)))
		i += copy(dAtA[i:], m.Payload)
	}
	dAtA[i] = 0x18
	i++
	i = encodeVarintBridge(dAtA, i, uint64(m.AuthID))
	dAtA[i] = 0x20
	i++
	i = encodeVarintBridge(dAtA, i, uint64(m.UserID))
	return i, nil
}

func (m *Notifier) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Notifier) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	dAtA[i] = 0x8
	i++
	i = encodeVarintBridge(dAtA, i, uint64(m.MessageID))
	if len(m.ConnectionIDs) > 0 {
		for _, num := range m.ConnectionIDs {
			dAtA[i] = 0x10
			i++
			i = encodeVarintBridge(dAtA, i, uint64(num))
		}
	}
	return i, nil
}

func (m *DeliveryReport) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *DeliveryReport) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.MessageIDs) > 0 {
		for _, num := range m.MessageIDs {
			dAtA[i] = 0x8
			i++
			i = encodeVarintBridge(dAtA, i, uint64(num))
		}
	}
	return i, nil
}

func encodeVarintBridge(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *Container) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	n += 1 + sovBridge(uint64(m.MessageID))
	if len(m.Messages) > 0 {
		for _, e := range m.Messages {
			l = e.Size()
			n += 1 + l + sovBridge(uint64(l))
		}
	}
	return n
}

func (m *Message) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Payload != nil {
		l = len(m.Payload)
		n += 1 + l + sovBridge(uint64(l))
	}
	n += 1 + sovBridge(uint64(m.AuthID))
	n += 1 + sovBridge(uint64(m.UserID))
	return n
}

func (m *Notifier) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	n += 1 + sovBridge(uint64(m.MessageID))
	if len(m.ConnectionIDs) > 0 {
		for _, e := range m.ConnectionIDs {
			n += 1 + sovBridge(uint64(e))
		}
	}
	return n
}

func (m *DeliveryReport) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.MessageIDs) > 0 {
		for _, e := range m.MessageIDs {
			n += 1 + sovBridge(uint64(e))
		}
	}
	return n
}

func sovBridge(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozBridge(x uint64) (n int) {
	return sovBridge(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Container) Unmarshal(dAtA []byte) error {
	var hasFields [1]uint64
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowBridge
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Container: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Container: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MessageID", wireType)
			}
			m.MessageID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBridge
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MessageID |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			hasFields[0] |= uint64(0x00000001)
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Messages", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBridge
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthBridge
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthBridge
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Messages = append(m.Messages, &Message{})
			if err := m.Messages[len(m.Messages)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipBridge(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthBridge
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthBridge
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}
	if hasFields[0]&uint64(0x00000001) == 0 {
		return github_com_gogo_protobuf_proto.NewRequiredNotSetError("MessageID")
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Message) Unmarshal(dAtA []byte) error {
	var hasFields [1]uint64
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowBridge
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Message: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Message: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Payload", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBridge
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthBridge
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthBridge
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Payload = append(m.Payload[:0], dAtA[iNdEx:postIndex]...)
			if m.Payload == nil {
				m.Payload = []byte{}
			}
			iNdEx = postIndex
			hasFields[0] |= uint64(0x00000001)
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field AuthID", wireType)
			}
			m.AuthID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBridge
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.AuthID |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			hasFields[0] |= uint64(0x00000002)
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field UserID", wireType)
			}
			m.UserID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBridge
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.UserID |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipBridge(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthBridge
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthBridge
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}
	if hasFields[0]&uint64(0x00000001) == 0 {
		return github_com_gogo_protobuf_proto.NewRequiredNotSetError("Payload")
	}
	if hasFields[0]&uint64(0x00000002) == 0 {
		return github_com_gogo_protobuf_proto.NewRequiredNotSetError("AuthID")
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Notifier) Unmarshal(dAtA []byte) error {
	var hasFields [1]uint64
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowBridge
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Notifier: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Notifier: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MessageID", wireType)
			}
			m.MessageID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBridge
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MessageID |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			hasFields[0] |= uint64(0x00000001)
		case 2:
			if wireType == 0 {
				var v uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowBridge
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					v |= uint64(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				m.ConnectionIDs = append(m.ConnectionIDs, v)
			} else if wireType == 2 {
				var packedLen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowBridge
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					packedLen |= int(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if packedLen < 0 {
					return ErrInvalidLengthBridge
				}
				postIndex := iNdEx + packedLen
				if postIndex < 0 {
					return ErrInvalidLengthBridge
				}
				if postIndex > l {
					return io.ErrUnexpectedEOF
				}
				var elementCount int
				var count int
				for _, integer := range dAtA[iNdEx:postIndex] {
					if integer < 128 {
						count++
					}
				}
				elementCount = count
				if elementCount != 0 && len(m.ConnectionIDs) == 0 {
					m.ConnectionIDs = make([]uint64, 0, elementCount)
				}
				for iNdEx < postIndex {
					var v uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowBridge
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						v |= uint64(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					m.ConnectionIDs = append(m.ConnectionIDs, v)
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field ConnectionIDs", wireType)
			}
		default:
			iNdEx = preIndex
			skippy, err := skipBridge(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthBridge
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthBridge
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}
	if hasFields[0]&uint64(0x00000001) == 0 {
		return github_com_gogo_protobuf_proto.NewRequiredNotSetError("MessageID")
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *DeliveryReport) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowBridge
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: DeliveryReport: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: DeliveryReport: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType == 0 {
				var v uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowBridge
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					v |= uint64(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				m.MessageIDs = append(m.MessageIDs, v)
			} else if wireType == 2 {
				var packedLen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowBridge
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					packedLen |= int(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if packedLen < 0 {
					return ErrInvalidLengthBridge
				}
				postIndex := iNdEx + packedLen
				if postIndex < 0 {
					return ErrInvalidLengthBridge
				}
				if postIndex > l {
					return io.ErrUnexpectedEOF
				}
				var elementCount int
				var count int
				for _, integer := range dAtA[iNdEx:postIndex] {
					if integer < 128 {
						count++
					}
				}
				elementCount = count
				if elementCount != 0 && len(m.MessageIDs) == 0 {
					m.MessageIDs = make([]uint64, 0, elementCount)
				}
				for iNdEx < postIndex {
					var v uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowBridge
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						v |= uint64(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					m.MessageIDs = append(m.MessageIDs, v)
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field MessageIDs", wireType)
			}
		default:
			iNdEx = preIndex
			skippy, err := skipBridge(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthBridge
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthBridge
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipBridge(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowBridge
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowBridge
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowBridge
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthBridge
			}
			iNdEx += length
			if iNdEx < 0 {
				return 0, ErrInvalidLengthBridge
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowBridge
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipBridge(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
				if iNdEx < 0 {
					return 0, ErrInvalidLengthBridge
				}
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthBridge = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowBridge   = fmt.Errorf("proto: integer overflow")
)
