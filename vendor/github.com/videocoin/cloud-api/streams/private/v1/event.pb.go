// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: streams/private/v1/event.proto

package v1

import (
	fmt "fmt"
	_ "github.com/gogo/googleapis/google/api"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	golang_proto "github.com/golang/protobuf/proto"
	v1 "github.com/videocoin/cloud-api/streams/v1"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = golang_proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

type EventType int32

const (
	EventTypeUnknown      EventType = 0
	EventTypeCreate       EventType = 1
	EventTypeUpdate       EventType = 2
	EventTypeDelete       EventType = 3
	EventTypeUpdateStatus EventType = 4
)

var EventType_name = map[int32]string{
	0: "EVENT_TYPE_UNKNOWN",
	1: "EVENT_TYPE_CREATE",
	2: "EVENT_TYPE_UPDATE",
	3: "EVENT_TYPE_DELETE",
	4: "EVENT_TYPE_UPDATE_STATUS",
}

var EventType_value = map[string]int32{
	"EVENT_TYPE_UNKNOWN":       0,
	"EVENT_TYPE_CREATE":        1,
	"EVENT_TYPE_UPDATE":        2,
	"EVENT_TYPE_DELETE":        3,
	"EVENT_TYPE_UPDATE_STATUS": 4,
}

func (x EventType) String() string {
	return proto.EnumName(EventType_name, int32(x))
}

func (EventType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_00ff1fe306afdb32, []int{0}
}

type Event struct {
	Type                 EventType       `protobuf:"varint,1,opt,name=type,proto3,enum=cloud.api.streams.private.v1.EventType" json:"type,omitempty"`
	StreamID             string          `protobuf:"bytes,2,opt,name=stream_id,json=streamId,proto3" json:"stream_id,omitempty"`
	Status               v1.StreamStatus `protobuf:"varint,3,opt,name=status,proto3,enum=cloud.api.streams.v1.StreamStatus" json:"status,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *Event) Reset()         { *m = Event{} }
func (m *Event) String() string { return proto.CompactTextString(m) }
func (*Event) ProtoMessage()    {}
func (*Event) Descriptor() ([]byte, []int) {
	return fileDescriptor_00ff1fe306afdb32, []int{0}
}
func (m *Event) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Event) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Event.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Event) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Event.Merge(m, src)
}
func (m *Event) XXX_Size() int {
	return m.Size()
}
func (m *Event) XXX_DiscardUnknown() {
	xxx_messageInfo_Event.DiscardUnknown(m)
}

var xxx_messageInfo_Event proto.InternalMessageInfo

func (m *Event) GetType() EventType {
	if m != nil {
		return m.Type
	}
	return EventTypeUnknown
}

func (m *Event) GetStreamID() string {
	if m != nil {
		return m.StreamID
	}
	return ""
}

func (m *Event) GetStatus() v1.StreamStatus {
	if m != nil {
		return m.Status
	}
	return v1.StreamStatusNone
}

func (*Event) XXX_MessageName() string {
	return "cloud.api.streams.private.v1.Event"
}
func init() {
	proto.RegisterEnum("cloud.api.streams.private.v1.EventType", EventType_name, EventType_value)
	golang_proto.RegisterEnum("cloud.api.streams.private.v1.EventType", EventType_name, EventType_value)
	proto.RegisterType((*Event)(nil), "cloud.api.streams.private.v1.Event")
	golang_proto.RegisterType((*Event)(nil), "cloud.api.streams.private.v1.Event")
}

func init() { proto.RegisterFile("streams/private/v1/event.proto", fileDescriptor_00ff1fe306afdb32) }
func init() {
	golang_proto.RegisterFile("streams/private/v1/event.proto", fileDescriptor_00ff1fe306afdb32)
}

var fileDescriptor_00ff1fe306afdb32 = []byte{
	// 430 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x91, 0xc1, 0x6e, 0xd3, 0x30,
	0x18, 0xc7, 0xe7, 0x6e, 0x4c, 0xab, 0x85, 0x20, 0x33, 0x20, 0x85, 0x68, 0x32, 0xd1, 0x2e, 0x8c,
	0x89, 0xd9, 0x0a, 0x48, 0x20, 0xc1, 0x69, 0x5b, 0x7d, 0x98, 0x40, 0x65, 0x6a, 0x13, 0x10, 0x5c,
	0x2a, 0xb7, 0x31, 0xc1, 0xa2, 0xb3, 0xa3, 0xc6, 0x35, 0xda, 0x1b, 0x40, 0xdf, 0xa1, 0x27, 0xb8,
	0xf1, 0x06, 0x9c, 0x38, 0xee, 0xc8, 0x13, 0x20, 0x94, 0xbd, 0x08, 0xaa, 0x13, 0xca, 0x94, 0x21,
	0x6e, 0xdf, 0xf7, 0xf9, 0xf7, 0xff, 0xfe, 0x7f, 0xdb, 0x10, 0x17, 0x66, 0x22, 0xf8, 0x49, 0x41,
	0xf3, 0x89, 0xb4, 0xdc, 0x08, 0x6a, 0x23, 0x2a, 0xac, 0x50, 0x86, 0xe4, 0x13, 0x6d, 0x34, 0xda,
	0x1a, 0x8d, 0xf5, 0x34, 0x25, 0x3c, 0x97, 0xa4, 0x26, 0x49, 0x4d, 0x12, 0x1b, 0x05, 0x5b, 0x99,
	0xd6, 0xd9, 0x58, 0x50, 0x9e, 0x4b, 0xca, 0x95, 0xd2, 0x86, 0x1b, 0xa9, 0x55, 0x51, 0x69, 0x83,
	0xbd, 0x4c, 0x9a, 0x77, 0xd3, 0x21, 0x19, 0xe9, 0x13, 0x9a, 0xe9, 0x4c, 0x53, 0x37, 0x1e, 0x4e,
	0xdf, 0xba, 0xce, 0x35, 0xae, 0xaa, 0xf1, 0x47, 0x17, 0x70, 0x2b, 0x53, 0xa1, 0x47, 0x5a, 0x2a,
	0xea, 0xfc, 0xf7, 0x16, 0x06, 0x7f, 0x92, 0xda, 0xa8, 0x2e, 0x2b, 0xdd, 0xf6, 0x57, 0x00, 0xaf,
	0xb0, 0x45, 0x64, 0xf4, 0x14, 0xae, 0x99, 0xd3, 0x5c, 0xf8, 0x20, 0x04, 0x3b, 0xd7, 0x1e, 0xdc,
	0x25, 0xff, 0xcb, 0x4e, 0x9c, 0x24, 0x3e, 0xcd, 0x45, 0xcf, 0x89, 0xd0, 0x3d, 0xd8, 0xae, 0xa8,
	0x81, 0x4c, 0xfd, 0x56, 0x08, 0x76, 0xda, 0x07, 0x57, 0xcb, 0x9f, 0x77, 0x36, 0xfa, 0x6e, 0x78,
	0xd4, 0xe9, 0x6d, 0x54, 0xc7, 0x47, 0x29, 0x7a, 0x02, 0xd7, 0x0b, 0xc3, 0xcd, 0xb4, 0xf0, 0x57,
	0x9d, 0xd3, 0xf6, 0x3f, 0x9c, 0x6c, 0x44, 0x2a, 0x65, 0xdf, 0x91, 0xbd, 0x5a, 0xb1, 0xfb, 0xa9,
	0x05, 0xdb, 0x4b, 0x6b, 0x74, 0x1f, 0x22, 0xf6, 0x92, 0x75, 0xe3, 0x41, 0xfc, 0xfa, 0x98, 0x0d,
	0x92, 0xee, 0xb3, 0xee, 0x8b, 0x57, 0x5d, 0x6f, 0x25, 0xb8, 0x39, 0x9b, 0x87, 0xde, 0x12, 0x4b,
	0xd4, 0x7b, 0xa5, 0x3f, 0x28, 0xb4, 0x0b, 0x37, 0x2f, 0xd0, 0x87, 0x3d, 0xb6, 0x1f, 0x33, 0x0f,
	0x04, 0x37, 0x66, 0xf3, 0xf0, 0xfa, 0x12, 0x3e, 0x9c, 0x08, 0x6e, 0x44, 0x83, 0x4d, 0x8e, 0x3b,
	0x0b, 0xb6, 0xd5, 0x60, 0x93, 0x3c, 0xbd, 0xcc, 0x76, 0xd8, 0x73, 0x16, 0x33, 0x6f, 0xb5, 0xc1,
	0x76, 0xc4, 0x58, 0x18, 0x81, 0x1e, 0x43, 0xff, 0xd2, 0xde, 0x41, 0x3f, 0xde, 0x8f, 0x93, 0xbe,
	0xb7, 0x16, 0xdc, 0x9e, 0xcd, 0xc3, 0x5b, 0x8d, 0xf5, 0xd5, 0x03, 0x04, 0x9b, 0x1f, 0x3f, 0xe3,
	0x95, 0x6f, 0x5f, 0xf0, 0xdf, 0xdb, 0x1f, 0xf8, 0x67, 0x25, 0x06, 0x3f, 0x4a, 0x0c, 0x7e, 0x95,
	0x18, 0x7c, 0x3f, 0xc7, 0xe0, 0xec, 0x1c, 0x83, 0x37, 0x2d, 0x1b, 0x0d, 0xd7, 0xdd, 0xd7, 0x3e,
	0xfc, 0x1d, 0x00, 0x00, 0xff, 0xff, 0x16, 0x56, 0x6f, 0x03, 0x9f, 0x02, 0x00, 0x00,
}

func (m *Event) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Event) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Type != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintEvent(dAtA, i, uint64(m.Type))
	}
	if len(m.StreamID) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintEvent(dAtA, i, uint64(len(m.StreamID)))
		i += copy(dAtA[i:], m.StreamID)
	}
	if m.Status != 0 {
		dAtA[i] = 0x18
		i++
		i = encodeVarintEvent(dAtA, i, uint64(m.Status))
	}
	if m.XXX_unrecognized != nil {
		i += copy(dAtA[i:], m.XXX_unrecognized)
	}
	return i, nil
}

func encodeVarintEvent(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *Event) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Type != 0 {
		n += 1 + sovEvent(uint64(m.Type))
	}
	l = len(m.StreamID)
	if l > 0 {
		n += 1 + l + sovEvent(uint64(l))
	}
	if m.Status != 0 {
		n += 1 + sovEvent(uint64(m.Status))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovEvent(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozEvent(x uint64) (n int) {
	return sovEvent(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Event) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowEvent
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
			return fmt.Errorf("proto: Event: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Event: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Type", wireType)
			}
			m.Type = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Type |= EventType(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field StreamID", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthEvent
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvent
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.StreamID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Status", wireType)
			}
			m.Status = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Status |= v1.StreamStatus(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipEvent(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthEvent
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthEvent
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipEvent(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowEvent
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
					return 0, ErrIntOverflowEvent
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
					return 0, ErrIntOverflowEvent
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
				return 0, ErrInvalidLengthEvent
			}
			iNdEx += length
			if iNdEx < 0 {
				return 0, ErrInvalidLengthEvent
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowEvent
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
				next, err := skipEvent(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
				if iNdEx < 0 {
					return 0, ErrInvalidLengthEvent
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
	ErrInvalidLengthEvent = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowEvent   = fmt.Errorf("proto: integer overflow")
)
