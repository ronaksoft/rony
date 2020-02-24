// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: dev.proto

package msg

import (
	fmt "fmt"
	github_com_gogo_protobuf_proto "github.com/gogo/protobuf/proto"
	proto "github.com/gogo/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// ProtoMessage
// If AuthID == 0 then Payload is a MessageEnvelop otherwise Payload is a ProtoEncryptedPayload
type ProtoMessage struct {
	AuthID     int64  `protobuf:"varint,1,opt,name=AuthID" json:"AuthID"`
	MessageKey []byte `protobuf:"bytes,2,opt,name=MessageKey" json:"MessageKey"`
	Payload    []byte `protobuf:"bytes,3,req,name=Payload" json:"Payload"`
}

func (m *ProtoMessage) Reset()         { *m = ProtoMessage{} }
func (m *ProtoMessage) String() string { return proto.CompactTextString(m) }
func (*ProtoMessage) ProtoMessage()    {}
func (*ProtoMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_f7546b35e7a0cc16, []int{0}
}
func (m *ProtoMessage) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ProtoMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ProtoMessage.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ProtoMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProtoMessage.Merge(m, src)
}
func (m *ProtoMessage) XXX_Size() int {
	return m.Size()
}
func (m *ProtoMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_ProtoMessage.DiscardUnknown(m)
}

var xxx_messageInfo_ProtoMessage proto.InternalMessageInfo

func (m *ProtoMessage) GetAuthID() int64 {
	if m != nil {
		return m.AuthID
	}
	return 0
}

func (m *ProtoMessage) GetMessageKey() []byte {
	if m != nil {
		return m.MessageKey
	}
	return nil
}

func (m *ProtoMessage) GetPayload() []byte {
	if m != nil {
		return m.Payload
	}
	return nil
}

type EchoRequest struct {
	Int       int64 `protobuf:"varint,1,req,name=Int" json:"Int"`
	Bool      bool  `protobuf:"varint,2,req,name=Bool" json:"Bool"`
	Timestamp int64 `protobuf:"varint,3,req,name=Timestamp" json:"Timestamp"`
}

func (m *EchoRequest) Reset()         { *m = EchoRequest{} }
func (m *EchoRequest) String() string { return proto.CompactTextString(m) }
func (*EchoRequest) ProtoMessage()    {}
func (*EchoRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_f7546b35e7a0cc16, []int{1}
}
func (m *EchoRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *EchoRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_EchoRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *EchoRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EchoRequest.Merge(m, src)
}
func (m *EchoRequest) XXX_Size() int {
	return m.Size()
}
func (m *EchoRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_EchoRequest.DiscardUnknown(m)
}

var xxx_messageInfo_EchoRequest proto.InternalMessageInfo

func (m *EchoRequest) GetInt() int64 {
	if m != nil {
		return m.Int
	}
	return 0
}

func (m *EchoRequest) GetBool() bool {
	if m != nil {
		return m.Bool
	}
	return false
}

func (m *EchoRequest) GetTimestamp() int64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

type EchoResponse struct {
	Int       int64 `protobuf:"varint,1,req,name=Int" json:"Int"`
	Bool      bool  `protobuf:"varint,2,req,name=Bool" json:"Bool"`
	Timestamp int64 `protobuf:"varint,4,req,name=Timestamp" json:"Timestamp"`
	Delay     int64 `protobuf:"varint,5,req,name=Delay" json:"Delay"`
}

func (m *EchoResponse) Reset()         { *m = EchoResponse{} }
func (m *EchoResponse) String() string { return proto.CompactTextString(m) }
func (*EchoResponse) ProtoMessage()    {}
func (*EchoResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_f7546b35e7a0cc16, []int{2}
}
func (m *EchoResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *EchoResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_EchoResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *EchoResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EchoResponse.Merge(m, src)
}
func (m *EchoResponse) XXX_Size() int {
	return m.Size()
}
func (m *EchoResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_EchoResponse.DiscardUnknown(m)
}

var xxx_messageInfo_EchoResponse proto.InternalMessageInfo

func (m *EchoResponse) GetInt() int64 {
	if m != nil {
		return m.Int
	}
	return 0
}

func (m *EchoResponse) GetBool() bool {
	if m != nil {
		return m.Bool
	}
	return false
}

func (m *EchoResponse) GetTimestamp() int64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

func (m *EchoResponse) GetDelay() int64 {
	if m != nil {
		return m.Delay
	}
	return 0
}

func init() {
	proto.RegisterType((*ProtoMessage)(nil), "msg.ProtoMessage")
	proto.RegisterType((*EchoRequest)(nil), "msg.EchoRequest")
	proto.RegisterType((*EchoResponse)(nil), "msg.EchoResponse")
}

func init() { proto.RegisterFile("dev.proto", fileDescriptor_f7546b35e7a0cc16) }

var fileDescriptor_f7546b35e7a0cc16 = []byte{
	// 245 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x4c, 0x49, 0x2d, 0xd3,
	0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0xce, 0x2d, 0x4e, 0x57, 0x2a, 0xe2, 0xe2, 0x09, 0x00,
	0xf1, 0x7c, 0x53, 0x8b, 0x8b, 0x13, 0xd3, 0x53, 0x85, 0x64, 0xb8, 0xd8, 0x1c, 0x4b, 0x4b, 0x32,
	0x3c, 0x5d, 0x24, 0x18, 0x15, 0x18, 0x35, 0x98, 0x9d, 0x58, 0x4e, 0xdc, 0x93, 0x67, 0x08, 0x82,
	0x8a, 0x09, 0xa9, 0x70, 0x71, 0x41, 0x15, 0x7a, 0xa7, 0x56, 0x4a, 0x30, 0x29, 0x30, 0x6a, 0xf0,
	0x40, 0x55, 0x20, 0x89, 0x0b, 0xc9, 0x71, 0xb1, 0x07, 0x24, 0x56, 0xe6, 0xe4, 0x27, 0xa6, 0x48,
	0x30, 0x2b, 0x30, 0xc1, 0x95, 0xc0, 0x04, 0x95, 0x92, 0xb9, 0xb8, 0x5d, 0x93, 0x33, 0xf2, 0x83,
	0x52, 0x0b, 0x4b, 0x53, 0x8b, 0x4b, 0x84, 0xc4, 0xb8, 0x98, 0x3d, 0xf3, 0x4a, 0x24, 0x18, 0x15,
	0x98, 0xe0, 0xf6, 0x81, 0x04, 0x84, 0x24, 0xb8, 0x58, 0x9c, 0xf2, 0xf3, 0x73, 0x24, 0x98, 0x14,
	0x98, 0x34, 0x38, 0xa0, 0x12, 0x60, 0x11, 0x21, 0x25, 0x2e, 0xce, 0x90, 0xcc, 0xdc, 0xd4, 0xe2,
	0x92, 0xc4, 0xdc, 0x02, 0xb0, 0x15, 0x30, 0x7d, 0x08, 0x61, 0xa5, 0x06, 0x46, 0x2e, 0x1e, 0x88,
	0x2d, 0xc5, 0x05, 0xf9, 0x79, 0xc5, 0xa9, 0x94, 0x5a, 0xc3, 0x82, 0xd5, 0x1a, 0x21, 0x29, 0x2e,
	0x56, 0x97, 0xd4, 0x9c, 0xc4, 0x4a, 0x09, 0x56, 0x24, 0x79, 0x88, 0x90, 0x93, 0xc4, 0x89, 0x47,
	0x72, 0x8c, 0x17, 0x1e, 0xc9, 0x31, 0x3e, 0x78, 0x24, 0xc7, 0x38, 0xe1, 0xb1, 0x1c, 0xc3, 0x85,
	0xc7, 0x72, 0x0c, 0x37, 0x1e, 0xcb, 0x31, 0x00, 0x02, 0x00, 0x00, 0xff, 0xff, 0x1a, 0x1f, 0x68,
	0xec, 0x86, 0x01, 0x00, 0x00,
}

func (m *ProtoMessage) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ProtoMessage) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ProtoMessage) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Payload != nil {
		i -= len(m.Payload)
		copy(dAtA[i:], m.Payload)
		i = encodeVarintDev(dAtA, i, uint64(len(m.Payload)))
		i--
		dAtA[i] = 0x1a
	}
	if m.MessageKey != nil {
		i -= len(m.MessageKey)
		copy(dAtA[i:], m.MessageKey)
		i = encodeVarintDev(dAtA, i, uint64(len(m.MessageKey)))
		i--
		dAtA[i] = 0x12
	}
	i = encodeVarintDev(dAtA, i, uint64(m.AuthID))
	i--
	dAtA[i] = 0x8
	return len(dAtA) - i, nil
}

func (m *EchoRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *EchoRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *EchoRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	i = encodeVarintDev(dAtA, i, uint64(m.Timestamp))
	i--
	dAtA[i] = 0x18
	i--
	if m.Bool {
		dAtA[i] = 1
	} else {
		dAtA[i] = 0
	}
	i--
	dAtA[i] = 0x10
	i = encodeVarintDev(dAtA, i, uint64(m.Int))
	i--
	dAtA[i] = 0x8
	return len(dAtA) - i, nil
}

func (m *EchoResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *EchoResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *EchoResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	i = encodeVarintDev(dAtA, i, uint64(m.Delay))
	i--
	dAtA[i] = 0x28
	i = encodeVarintDev(dAtA, i, uint64(m.Timestamp))
	i--
	dAtA[i] = 0x20
	i--
	if m.Bool {
		dAtA[i] = 1
	} else {
		dAtA[i] = 0
	}
	i--
	dAtA[i] = 0x10
	i = encodeVarintDev(dAtA, i, uint64(m.Int))
	i--
	dAtA[i] = 0x8
	return len(dAtA) - i, nil
}

func encodeVarintDev(dAtA []byte, offset int, v uint64) int {
	offset -= sovDev(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *ProtoMessage) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	n += 1 + sovDev(uint64(m.AuthID))
	if m.MessageKey != nil {
		l = len(m.MessageKey)
		n += 1 + l + sovDev(uint64(l))
	}
	if m.Payload != nil {
		l = len(m.Payload)
		n += 1 + l + sovDev(uint64(l))
	}
	return n
}

func (m *EchoRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	n += 1 + sovDev(uint64(m.Int))
	n += 2
	n += 1 + sovDev(uint64(m.Timestamp))
	return n
}

func (m *EchoResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	n += 1 + sovDev(uint64(m.Int))
	n += 2
	n += 1 + sovDev(uint64(m.Timestamp))
	n += 1 + sovDev(uint64(m.Delay))
	return n
}

func sovDev(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozDev(x uint64) (n int) {
	return sovDev(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *ProtoMessage) Unmarshal(dAtA []byte) error {
	var hasFields [1]uint64
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowDev
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
			return fmt.Errorf("proto: ProtoMessage: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ProtoMessage: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field AuthID", wireType)
			}
			m.AuthID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDev
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
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MessageKey", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDev
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
				return ErrInvalidLengthDev
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthDev
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.MessageKey = append(m.MessageKey[:0], dAtA[iNdEx:postIndex]...)
			if m.MessageKey == nil {
				m.MessageKey = []byte{}
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Payload", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDev
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
				return ErrInvalidLengthDev
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthDev
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
		default:
			iNdEx = preIndex
			skippy, err := skipDev(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthDev
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthDev
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

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *EchoRequest) Unmarshal(dAtA []byte) error {
	var hasFields [1]uint64
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowDev
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
			return fmt.Errorf("proto: EchoRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: EchoRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Int", wireType)
			}
			m.Int = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDev
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Int |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			hasFields[0] |= uint64(0x00000001)
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Bool", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDev
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.Bool = bool(v != 0)
			hasFields[0] |= uint64(0x00000002)
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Timestamp", wireType)
			}
			m.Timestamp = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDev
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Timestamp |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			hasFields[0] |= uint64(0x00000004)
		default:
			iNdEx = preIndex
			skippy, err := skipDev(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthDev
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthDev
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}
	if hasFields[0]&uint64(0x00000001) == 0 {
		return github_com_gogo_protobuf_proto.NewRequiredNotSetError("Int")
	}
	if hasFields[0]&uint64(0x00000002) == 0 {
		return github_com_gogo_protobuf_proto.NewRequiredNotSetError("Bool")
	}
	if hasFields[0]&uint64(0x00000004) == 0 {
		return github_com_gogo_protobuf_proto.NewRequiredNotSetError("Timestamp")
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *EchoResponse) Unmarshal(dAtA []byte) error {
	var hasFields [1]uint64
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowDev
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
			return fmt.Errorf("proto: EchoResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: EchoResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Int", wireType)
			}
			m.Int = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDev
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Int |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			hasFields[0] |= uint64(0x00000001)
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Bool", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDev
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.Bool = bool(v != 0)
			hasFields[0] |= uint64(0x00000002)
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Timestamp", wireType)
			}
			m.Timestamp = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDev
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Timestamp |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			hasFields[0] |= uint64(0x00000004)
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Delay", wireType)
			}
			m.Delay = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDev
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Delay |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			hasFields[0] |= uint64(0x00000008)
		default:
			iNdEx = preIndex
			skippy, err := skipDev(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthDev
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthDev
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}
	if hasFields[0]&uint64(0x00000001) == 0 {
		return github_com_gogo_protobuf_proto.NewRequiredNotSetError("Int")
	}
	if hasFields[0]&uint64(0x00000002) == 0 {
		return github_com_gogo_protobuf_proto.NewRequiredNotSetError("Bool")
	}
	if hasFields[0]&uint64(0x00000004) == 0 {
		return github_com_gogo_protobuf_proto.NewRequiredNotSetError("Timestamp")
	}
	if hasFields[0]&uint64(0x00000008) == 0 {
		return github_com_gogo_protobuf_proto.NewRequiredNotSetError("Delay")
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipDev(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowDev
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
					return 0, ErrIntOverflowDev
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowDev
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
				return 0, ErrInvalidLengthDev
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupDev
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthDev
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthDev        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowDev          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupDev = fmt.Errorf("proto: unexpected end of group")
)
