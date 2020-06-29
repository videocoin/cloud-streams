// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: users/v1/user.proto

package v1

import (
	fmt "fmt"
	_ "github.com/gogo/googleapis/google/api"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	_ "github.com/gogo/protobuf/types"
	golang_proto "github.com/golang/protobuf/proto"
	_ "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger/options"
	_ "github.com/videocoin/cloud-api/accounts/v1"
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
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type TokenType int32

const (
	TokenTypeRegular TokenType = 0
	TokenTypeAPI     TokenType = 1
)

var TokenType_name = map[int32]string{
	0: "TOKEN_TYPE_REGULAR",
	1: "TOKEN_TYPE_API",
}

var TokenType_value = map[string]int32{
	"TOKEN_TYPE_REGULAR": 0,
	"TOKEN_TYPE_API":     1,
}

func (x TokenType) String() string {
	return proto.EnumName(TokenType_name, int32(x))
}

func (TokenType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_622714e2df60ae10, []int{0}
}

type UserRole int32

const (
	UserRoleRegular UserRole = 0
	UserRoleMiner   UserRole = 1
	UserRoleQa      UserRole = 3
	UserRoleManager UserRole = 6
	UserRoleSuper   UserRole = 9
)

var UserRole_name = map[int32]string{
	0: "USER_ROLE_REGULAR",
	1: "USER_ROLE_MINER",
	3: "USER_ROLE_QA",
	6: "USER_ROLE_MANAGER",
	9: "USER_ROLE_SUPER",
}

var UserRole_value = map[string]int32{
	"USER_ROLE_REGULAR": 0,
	"USER_ROLE_MINER":   1,
	"USER_ROLE_QA":      3,
	"USER_ROLE_MANAGER": 6,
	"USER_ROLE_SUPER":   9,
}

func (x UserRole) String() string {
	return proto.EnumName(UserRole_name, int32(x))
}

func (UserRole) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_622714e2df60ae10, []int{1}
}

type UserUIRole int32

const (
	UserUIRoleBoth      UserUIRole = 0
	UserUIRoleMiner     UserUIRole = 1
	UserUIRolePublisher UserUIRole = 2
)

var UserUIRole_name = map[int32]string{
	0: "USER_ROLE_UI_BOTH",
	1: "USER_ROLE_UI_MINER",
	2: "USER_ROLE_UI_PUBLISHER",
}

var UserUIRole_value = map[string]int32{
	"USER_ROLE_UI_BOTH":      0,
	"USER_ROLE_UI_MINER":     1,
	"USER_ROLE_UI_PUBLISHER": 2,
}

func (x UserUIRole) String() string {
	return proto.EnumName(UserUIRole_name, int32(x))
}

func (UserUIRole) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_622714e2df60ae10, []int{2}
}

type UserProfile struct {
	ID                   string     `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Email                string     `protobuf:"bytes,2,opt,name=email,proto3" json:"email,omitempty"`
	FirstName            string     `protobuf:"bytes,3,opt,name=first_name,json=firstName,proto3" json:"first_name,omitempty"`
	LastName             string     `protobuf:"bytes,4,opt,name=last_name,json=lastName,proto3" json:"last_name,omitempty"`
	IsActive             bool       `protobuf:"varint,5,opt,name=is_active,json=isActive,proto3" json:"is_active,omitempty"`
	Role                 UserRole   `protobuf:"varint,6,opt,name=role,proto3,enum=cloud.api.users.v1.UserRole" json:"role,omitempty"`
	UIRole               UserUIRole `protobuf:"varint,7,opt,name=ui_role,json=uiRole,proto3,enum=cloud.api.users.v1.UserUIRole" json:"ui_role,omitempty"`
	Country              string     `protobuf:"bytes,8,opt,name=country,proto3" json:"country,omitempty"`
	Region               string     `protobuf:"bytes,9,opt,name=region,proto3" json:"region,omitempty"`
	City                 string     `protobuf:"bytes,10,opt,name=city,proto3" json:"city,omitempty"`
	Zip                  string     `protobuf:"bytes,11,opt,name=zip,proto3" json:"zip,omitempty"`
	Address_1            string     `protobuf:"bytes,12,opt,name=address_1,json=address1,proto3" json:"address_1,omitempty"`
	Address_2            string     `protobuf:"bytes,13,opt,name=address_2,json=address2,proto3" json:"address_2,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *UserProfile) Reset()         { *m = UserProfile{} }
func (m *UserProfile) String() string { return proto.CompactTextString(m) }
func (*UserProfile) ProtoMessage()    {}
func (*UserProfile) Descriptor() ([]byte, []int) {
	return fileDescriptor_622714e2df60ae10, []int{0}
}
func (m *UserProfile) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *UserProfile) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_UserProfile.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *UserProfile) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UserProfile.Merge(m, src)
}
func (m *UserProfile) XXX_Size() int {
	return m.Size()
}
func (m *UserProfile) XXX_DiscardUnknown() {
	xxx_messageInfo_UserProfile.DiscardUnknown(m)
}

var xxx_messageInfo_UserProfile proto.InternalMessageInfo

func (m *UserProfile) GetID() string {
	if m != nil {
		return m.ID
	}
	return ""
}

func (m *UserProfile) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *UserProfile) GetFirstName() string {
	if m != nil {
		return m.FirstName
	}
	return ""
}

func (m *UserProfile) GetLastName() string {
	if m != nil {
		return m.LastName
	}
	return ""
}

func (m *UserProfile) GetIsActive() bool {
	if m != nil {
		return m.IsActive
	}
	return false
}

func (m *UserProfile) GetRole() UserRole {
	if m != nil {
		return m.Role
	}
	return UserRoleRegular
}

func (m *UserProfile) GetUIRole() UserUIRole {
	if m != nil {
		return m.UIRole
	}
	return UserUIRoleBoth
}

func (m *UserProfile) GetCountry() string {
	if m != nil {
		return m.Country
	}
	return ""
}

func (m *UserProfile) GetRegion() string {
	if m != nil {
		return m.Region
	}
	return ""
}

func (m *UserProfile) GetCity() string {
	if m != nil {
		return m.City
	}
	return ""
}

func (m *UserProfile) GetZip() string {
	if m != nil {
		return m.Zip
	}
	return ""
}

func (m *UserProfile) GetAddress_1() string {
	if m != nil {
		return m.Address_1
	}
	return ""
}

func (m *UserProfile) GetAddress_2() string {
	if m != nil {
		return m.Address_2
	}
	return ""
}

func (*UserProfile) XXX_MessageName() string {
	return "cloud.api.users.v1.UserProfile"
}
func init() {
	proto.RegisterEnum("cloud.api.users.v1.TokenType", TokenType_name, TokenType_value)
	golang_proto.RegisterEnum("cloud.api.users.v1.TokenType", TokenType_name, TokenType_value)
	proto.RegisterEnum("cloud.api.users.v1.UserRole", UserRole_name, UserRole_value)
	golang_proto.RegisterEnum("cloud.api.users.v1.UserRole", UserRole_name, UserRole_value)
	proto.RegisterEnum("cloud.api.users.v1.UserUIRole", UserUIRole_name, UserUIRole_value)
	golang_proto.RegisterEnum("cloud.api.users.v1.UserUIRole", UserUIRole_name, UserUIRole_value)
	proto.RegisterType((*UserProfile)(nil), "cloud.api.users.v1.UserProfile")
	golang_proto.RegisterType((*UserProfile)(nil), "cloud.api.users.v1.UserProfile")
}

func init() { proto.RegisterFile("users/v1/user.proto", fileDescriptor_622714e2df60ae10) }
func init() { golang_proto.RegisterFile("users/v1/user.proto", fileDescriptor_622714e2df60ae10) }

var fileDescriptor_622714e2df60ae10 = []byte{
	// 716 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x93, 0xdd, 0x4e, 0xdb, 0x48,
	0x14, 0xc7, 0x71, 0x02, 0x21, 0x19, 0x20, 0x98, 0x01, 0xb1, 0x96, 0x97, 0x35, 0xd6, 0x6a, 0xb5,
	0x62, 0xb3, 0x24, 0xde, 0xc0, 0xcd, 0xde, 0x26, 0xbb, 0x16, 0x44, 0x25, 0x21, 0x4c, 0x92, 0x8b,
	0xf6, 0xc6, 0x9a, 0x38, 0x83, 0x19, 0xd5, 0xf1, 0x58, 0x63, 0x3b, 0x15, 0x7d, 0x82, 0x2a, 0xef,
	0xe0, 0xab, 0xf6, 0xb6, 0x2f, 0xd0, 0xab, 0x5e, 0x72, 0xd9, 0x27, 0x40, 0x55, 0x90, 0xfa, 0x1c,
	0x95, 0x27, 0xce, 0x57, 0x3f, 0xae, 0x7c, 0xce, 0xf9, 0xff, 0xce, 0xfc, 0xcf, 0xf1, 0xd8, 0x60,
	0x3f, 0x0a, 0x08, 0x0f, 0x8c, 0x51, 0xd5, 0x48, 0x82, 0x8a, 0xcf, 0x59, 0xc8, 0x20, 0xb4, 0x5d,
	0x16, 0x0d, 0x2a, 0xd8, 0xa7, 0x15, 0x21, 0x57, 0x46, 0x55, 0xf5, 0x5f, 0x87, 0x86, 0x77, 0x51,
	0xbf, 0x62, 0xb3, 0xa1, 0x31, 0xa2, 0x03, 0xc2, 0x6c, 0x46, 0x3d, 0x43, 0x80, 0x65, 0xec, 0x53,
	0x03, 0xdb, 0x36, 0x8b, 0xbc, 0x50, 0x1c, 0x95, 0xc6, 0xd3, 0xd3, 0xd4, 0x63, 0x87, 0x31, 0xc7,
	0x25, 0x86, 0xc8, 0xfa, 0xd1, 0xad, 0x11, 0xd2, 0x21, 0x09, 0x42, 0x3c, 0xf4, 0x53, 0xe0, 0x28,
	0x05, 0xc4, 0x31, 0x9e, 0xc7, 0x42, 0x1c, 0x52, 0xe6, 0x05, 0xa9, 0x7a, 0x2a, 0x1e, 0x76, 0xd9,
	0x21, 0x5e, 0x39, 0x78, 0x85, 0x1d, 0x87, 0x70, 0x83, 0xf9, 0x82, 0xf8, 0x01, 0x5d, 0x5e, 0x1a,
	0xd3, 0x61, 0x0e, 0x5b, 0xb8, 0x26, 0x99, 0x48, 0x44, 0x34, 0xc5, 0x7f, 0x8f, 0xb3, 0x60, 0xab,
	0x17, 0x10, 0xde, 0xe6, 0xec, 0x96, 0xba, 0x04, 0x1e, 0x82, 0x0c, 0x1d, 0x28, 0x92, 0x2e, 0x9d,
	0x14, 0xea, 0xb9, 0xc9, 0xe3, 0x71, 0xa6, 0xf1, 0x3f, 0xca, 0xd0, 0x01, 0x3c, 0x00, 0x1b, 0x64,
	0x88, 0xa9, 0xab, 0x64, 0x12, 0x09, 0x4d, 0x13, 0xf8, 0x1b, 0x00, 0xb7, 0x94, 0x07, 0xa1, 0xe5,
	0xe1, 0x21, 0x51, 0xb2, 0x42, 0x2a, 0x88, 0x4a, 0x0b, 0x0f, 0x09, 0xfc, 0x15, 0x14, 0x5c, 0x3c,
	0x53, 0xd7, 0x85, 0x9a, 0x4f, 0x0a, 0x33, 0x91, 0x06, 0x16, 0xb6, 0x43, 0x3a, 0x22, 0xca, 0x86,
	0x2e, 0x9d, 0xe4, 0x51, 0x9e, 0x06, 0x35, 0x91, 0xc3, 0x7f, 0xc0, 0x3a, 0x67, 0x2e, 0x51, 0x72,
	0xba, 0x74, 0x52, 0x3c, 0x3b, 0xaa, 0x7c, 0x7f, 0x1f, 0x95, 0x64, 0x6a, 0xc4, 0x5c, 0x82, 0x04,
	0x09, 0xff, 0x03, 0x9b, 0x11, 0xb5, 0x44, 0xd3, 0xa6, 0x68, 0xd2, 0x7e, 0xd6, 0xd4, 0x6b, 0x24,
	0x6d, 0x75, 0x30, 0x79, 0x3c, 0xce, 0x4d, 0x63, 0x94, 0x8b, 0x68, 0xf2, 0x84, 0x0a, 0xd8, 0x14,
	0x17, 0xc7, 0xef, 0x95, 0xbc, 0x18, 0x77, 0x96, 0xc2, 0x43, 0x90, 0xe3, 0xc4, 0xa1, 0xcc, 0x53,
	0x0a, 0x42, 0x48, 0x33, 0x08, 0xc1, 0xba, 0x4d, 0xc3, 0x7b, 0x05, 0x88, 0xaa, 0x88, 0xa1, 0x0c,
	0xb2, 0xaf, 0xa9, 0xaf, 0x6c, 0x89, 0x52, 0x12, 0x26, 0xbb, 0xe2, 0xc1, 0x80, 0x93, 0x20, 0xb0,
	0xaa, 0xca, 0xf6, 0xf4, 0x45, 0xa4, 0x85, 0xea, 0xb2, 0x78, 0xa6, 0xec, 0xac, 0x88, 0x67, 0x25,
	0x0e, 0x0a, 0x5d, 0xf6, 0x92, 0x78, 0xdd, 0x7b, 0x9f, 0xc0, 0x53, 0x00, 0xbb, 0xd7, 0xcf, 0xcc,
	0x96, 0xd5, 0x7d, 0xde, 0x36, 0x2d, 0x64, 0x5e, 0xf4, 0xae, 0x6a, 0x48, 0x5e, 0x53, 0x0f, 0xc6,
	0xb1, 0x2e, 0xcf, 0x31, 0x44, 0x9c, 0xc8, 0xc5, 0x1c, 0xfe, 0x01, 0x8a, 0x4b, 0x74, 0xad, 0xdd,
	0x90, 0x25, 0x55, 0x1e, 0xc7, 0xfa, 0xf6, 0x9c, 0xac, 0xb5, 0x1b, 0xea, 0xde, 0x9b, 0xb7, 0xda,
	0xda, 0x87, 0x77, 0xda, 0xc2, 0xa6, 0xf4, 0x45, 0x02, 0xf9, 0xd9, 0xdb, 0x85, 0x25, 0xb0, 0xd7,
	0xeb, 0x98, 0xc8, 0x42, 0xd7, 0x57, 0xcb, 0x96, 0xfb, 0xe3, 0x58, 0xdf, 0x9d, 0x5f, 0x41, 0xea,
	0xf8, 0x27, 0xd8, 0x5d, 0xb0, 0xcd, 0x46, 0xcb, 0x44, 0xb2, 0xa4, 0xee, 0x8d, 0x63, 0x7d, 0x67,
	0x46, 0x36, 0xa9, 0x47, 0x38, 0xd4, 0xc1, 0xf6, 0x82, 0xbb, 0xa9, 0xc9, 0x59, 0xb5, 0x38, 0x8e,
	0x75, 0x30, 0x83, 0x6e, 0xf0, 0xaa, 0x6b, 0xb3, 0xd6, 0xaa, 0x5d, 0x98, 0x48, 0xce, 0xad, 0xba,
	0x36, 0xb1, 0x87, 0x1d, 0xf2, 0x8d, 0x6b, 0xa7, 0xd7, 0x36, 0x91, 0x5c, 0x58, 0x75, 0xed, 0x44,
	0x3e, 0xe1, 0xaa, 0x9c, 0x6e, 0x3a, 0xdf, 0xad, 0xf4, 0x5e, 0x02, 0x60, 0xf1, 0x45, 0xc0, 0xbf,
	0x96, 0x4d, 0x7b, 0x0d, 0xab, 0x7e, 0xdd, 0xbd, 0x94, 0xd7, 0x54, 0x38, 0x8e, 0xf5, 0xe2, 0xd2,
	0x87, 0xc3, 0xc2, 0x3b, 0xf8, 0x37, 0x80, 0x2b, 0xe8, 0x6c, 0xd9, 0xf9, 0x80, 0x53, 0x76, 0xba,
	0xee, 0x39, 0x38, 0x5c, 0x81, 0xdb, 0xbd, 0xfa, 0x55, 0xa3, 0x73, 0x69, 0x22, 0x39, 0xa3, 0xfe,
	0x32, 0x8e, 0xf5, 0xfd, 0x45, 0x43, 0x3b, 0xea, 0xbb, 0x34, 0xb8, 0x23, 0x5c, 0x85, 0xe9, 0xb4,
	0x4b, 0x03, 0xd6, 0x95, 0x87, 0x89, 0x26, 0x7d, 0x9a, 0x68, 0xd2, 0xe7, 0x89, 0x26, 0x7d, 0x7c,
	0xd2, 0xa4, 0x87, 0x27, 0x4d, 0x7a, 0x91, 0x19, 0x55, 0xfb, 0x39, 0xf1, 0x37, 0x9f, 0x7f, 0x0d,
	0x00, 0x00, 0xff, 0xff, 0x8f, 0x9d, 0xe7, 0xd0, 0xce, 0x04, 0x00, 0x00,
}

func (m *UserProfile) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *UserProfile) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *UserProfile) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Address_2) > 0 {
		i -= len(m.Address_2)
		copy(dAtA[i:], m.Address_2)
		i = encodeVarintUser(dAtA, i, uint64(len(m.Address_2)))
		i--
		dAtA[i] = 0x6a
	}
	if len(m.Address_1) > 0 {
		i -= len(m.Address_1)
		copy(dAtA[i:], m.Address_1)
		i = encodeVarintUser(dAtA, i, uint64(len(m.Address_1)))
		i--
		dAtA[i] = 0x62
	}
	if len(m.Zip) > 0 {
		i -= len(m.Zip)
		copy(dAtA[i:], m.Zip)
		i = encodeVarintUser(dAtA, i, uint64(len(m.Zip)))
		i--
		dAtA[i] = 0x5a
	}
	if len(m.City) > 0 {
		i -= len(m.City)
		copy(dAtA[i:], m.City)
		i = encodeVarintUser(dAtA, i, uint64(len(m.City)))
		i--
		dAtA[i] = 0x52
	}
	if len(m.Region) > 0 {
		i -= len(m.Region)
		copy(dAtA[i:], m.Region)
		i = encodeVarintUser(dAtA, i, uint64(len(m.Region)))
		i--
		dAtA[i] = 0x4a
	}
	if len(m.Country) > 0 {
		i -= len(m.Country)
		copy(dAtA[i:], m.Country)
		i = encodeVarintUser(dAtA, i, uint64(len(m.Country)))
		i--
		dAtA[i] = 0x42
	}
	if m.UIRole != 0 {
		i = encodeVarintUser(dAtA, i, uint64(m.UIRole))
		i--
		dAtA[i] = 0x38
	}
	if m.Role != 0 {
		i = encodeVarintUser(dAtA, i, uint64(m.Role))
		i--
		dAtA[i] = 0x30
	}
	if m.IsActive {
		i--
		if m.IsActive {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x28
	}
	if len(m.LastName) > 0 {
		i -= len(m.LastName)
		copy(dAtA[i:], m.LastName)
		i = encodeVarintUser(dAtA, i, uint64(len(m.LastName)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.FirstName) > 0 {
		i -= len(m.FirstName)
		copy(dAtA[i:], m.FirstName)
		i = encodeVarintUser(dAtA, i, uint64(len(m.FirstName)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Email) > 0 {
		i -= len(m.Email)
		copy(dAtA[i:], m.Email)
		i = encodeVarintUser(dAtA, i, uint64(len(m.Email)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.ID) > 0 {
		i -= len(m.ID)
		copy(dAtA[i:], m.ID)
		i = encodeVarintUser(dAtA, i, uint64(len(m.ID)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintUser(dAtA []byte, offset int, v uint64) int {
	offset -= sovUser(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *UserProfile) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.ID)
	if l > 0 {
		n += 1 + l + sovUser(uint64(l))
	}
	l = len(m.Email)
	if l > 0 {
		n += 1 + l + sovUser(uint64(l))
	}
	l = len(m.FirstName)
	if l > 0 {
		n += 1 + l + sovUser(uint64(l))
	}
	l = len(m.LastName)
	if l > 0 {
		n += 1 + l + sovUser(uint64(l))
	}
	if m.IsActive {
		n += 2
	}
	if m.Role != 0 {
		n += 1 + sovUser(uint64(m.Role))
	}
	if m.UIRole != 0 {
		n += 1 + sovUser(uint64(m.UIRole))
	}
	l = len(m.Country)
	if l > 0 {
		n += 1 + l + sovUser(uint64(l))
	}
	l = len(m.Region)
	if l > 0 {
		n += 1 + l + sovUser(uint64(l))
	}
	l = len(m.City)
	if l > 0 {
		n += 1 + l + sovUser(uint64(l))
	}
	l = len(m.Zip)
	if l > 0 {
		n += 1 + l + sovUser(uint64(l))
	}
	l = len(m.Address_1)
	if l > 0 {
		n += 1 + l + sovUser(uint64(l))
	}
	l = len(m.Address_2)
	if l > 0 {
		n += 1 + l + sovUser(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovUser(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozUser(x uint64) (n int) {
	return sovUser(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *UserProfile) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowUser
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
			return fmt.Errorf("proto: UserProfile: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: UserProfile: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ID", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowUser
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
				return ErrInvalidLengthUser
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthUser
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Email", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowUser
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
				return ErrInvalidLengthUser
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthUser
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Email = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field FirstName", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowUser
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
				return ErrInvalidLengthUser
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthUser
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.FirstName = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field LastName", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowUser
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
				return ErrInvalidLengthUser
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthUser
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.LastName = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field IsActive", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowUser
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
			m.IsActive = bool(v != 0)
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Role", wireType)
			}
			m.Role = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowUser
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Role |= UserRole(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 7:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field UIRole", wireType)
			}
			m.UIRole = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowUser
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.UIRole |= UserUIRole(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Country", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowUser
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
				return ErrInvalidLengthUser
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthUser
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Country = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 9:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Region", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowUser
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
				return ErrInvalidLengthUser
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthUser
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Region = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 10:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field City", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowUser
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
				return ErrInvalidLengthUser
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthUser
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.City = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 11:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Zip", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowUser
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
				return ErrInvalidLengthUser
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthUser
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Zip = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 12:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Address_1", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowUser
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
				return ErrInvalidLengthUser
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthUser
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Address_1 = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 13:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Address_2", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowUser
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
				return ErrInvalidLengthUser
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthUser
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Address_2 = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipUser(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthUser
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthUser
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
func skipUser(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowUser
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
					return 0, ErrIntOverflowUser
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
					return 0, ErrIntOverflowUser
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
				return 0, ErrInvalidLengthUser
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupUser
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthUser
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthUser        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowUser          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupUser = fmt.Errorf("proto: unexpected end of group")
)
