// Code generated by protoc-gen-go. DO NOT EDIT.
// source: digimon.proto

package pbprotocol

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type LoginReq_LoginType int32

const (
	LoginReq_Visitor  LoginReq_LoginType = 0
	LoginReq_PassWord LoginReq_LoginType = 1
)

var LoginReq_LoginType_name = map[int32]string{
	0: "Visitor",
	1: "PassWord",
}

var LoginReq_LoginType_value = map[string]int32{
	"Visitor":  0,
	"PassWord": 1,
}

func (x LoginReq_LoginType) String() string {
	return proto.EnumName(LoginReq_LoginType_name, int32(x))
}

func (LoginReq_LoginType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_9f35044cfecd686d, []int{2, 0}
}

type RoomInfo_RoomType int32

const (
	RoomInfo_TWO RoomInfo_RoomType = 0
)

var RoomInfo_RoomType_name = map[int32]string{
	0: "TWO",
}

var RoomInfo_RoomType_value = map[string]int32{
	"TWO": 0,
}

func (x RoomInfo_RoomType) String() string {
	return proto.EnumName(RoomInfo_RoomType_name, int32(x))
}

func (RoomInfo_RoomType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_9f35044cfecd686d, []int{5, 0}
}

type MsgPack struct {
	Router               string   `protobuf:"bytes,1,opt,name=router,proto3" json:"router,omitempty"`
	Data                 []byte   `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MsgPack) Reset()         { *m = MsgPack{} }
func (m *MsgPack) String() string { return proto.CompactTextString(m) }
func (*MsgPack) ProtoMessage()    {}
func (*MsgPack) Descriptor() ([]byte, []int) {
	return fileDescriptor_9f35044cfecd686d, []int{0}
}

func (m *MsgPack) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MsgPack.Unmarshal(m, b)
}
func (m *MsgPack) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MsgPack.Marshal(b, m, deterministic)
}
func (m *MsgPack) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgPack.Merge(m, src)
}
func (m *MsgPack) XXX_Size() int {
	return xxx_messageInfo_MsgPack.Size(m)
}
func (m *MsgPack) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgPack.DiscardUnknown(m)
}

var xxx_messageInfo_MsgPack proto.InternalMessageInfo

func (m *MsgPack) GetRouter() string {
	if m != nil {
		return m.Router
	}
	return ""
}

func (m *MsgPack) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

type BaseAck struct {
	Result               int64    `protobuf:"varint,1,opt,name=Result,proto3" json:"Result,omitempty"`
	Msg                  string   `protobuf:"bytes,2,opt,name=Msg,proto3" json:"Msg,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BaseAck) Reset()         { *m = BaseAck{} }
func (m *BaseAck) String() string { return proto.CompactTextString(m) }
func (*BaseAck) ProtoMessage()    {}
func (*BaseAck) Descriptor() ([]byte, []int) {
	return fileDescriptor_9f35044cfecd686d, []int{1}
}

func (m *BaseAck) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BaseAck.Unmarshal(m, b)
}
func (m *BaseAck) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BaseAck.Marshal(b, m, deterministic)
}
func (m *BaseAck) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BaseAck.Merge(m, src)
}
func (m *BaseAck) XXX_Size() int {
	return xxx_messageInfo_BaseAck.Size(m)
}
func (m *BaseAck) XXX_DiscardUnknown() {
	xxx_messageInfo_BaseAck.DiscardUnknown(m)
}

var xxx_messageInfo_BaseAck proto.InternalMessageInfo

func (m *BaseAck) GetResult() int64 {
	if m != nil {
		return m.Result
	}
	return 0
}

func (m *BaseAck) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

type LoginReq struct {
	Username             string             `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
	Password             string             `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
	Type                 LoginReq_LoginType `protobuf:"varint,3,opt,name=type,proto3,enum=pbprotocol.LoginReq_LoginType" json:"type,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *LoginReq) Reset()         { *m = LoginReq{} }
func (m *LoginReq) String() string { return proto.CompactTextString(m) }
func (*LoginReq) ProtoMessage()    {}
func (*LoginReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_9f35044cfecd686d, []int{2}
}

func (m *LoginReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LoginReq.Unmarshal(m, b)
}
func (m *LoginReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LoginReq.Marshal(b, m, deterministic)
}
func (m *LoginReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LoginReq.Merge(m, src)
}
func (m *LoginReq) XXX_Size() int {
	return xxx_messageInfo_LoginReq.Size(m)
}
func (m *LoginReq) XXX_DiscardUnknown() {
	xxx_messageInfo_LoginReq.DiscardUnknown(m)
}

var xxx_messageInfo_LoginReq proto.InternalMessageInfo

func (m *LoginReq) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *LoginReq) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

func (m *LoginReq) GetType() LoginReq_LoginType {
	if m != nil {
		return m.Type
	}
	return LoginReq_Visitor
}

type PlayerInfo struct {
	Id                   uint64   `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Nickname             string   `protobuf:"bytes,2,opt,name=nickname,proto3" json:"nickname,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PlayerInfo) Reset()         { *m = PlayerInfo{} }
func (m *PlayerInfo) String() string { return proto.CompactTextString(m) }
func (*PlayerInfo) ProtoMessage()    {}
func (*PlayerInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_9f35044cfecd686d, []int{3}
}

func (m *PlayerInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PlayerInfo.Unmarshal(m, b)
}
func (m *PlayerInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PlayerInfo.Marshal(b, m, deterministic)
}
func (m *PlayerInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PlayerInfo.Merge(m, src)
}
func (m *PlayerInfo) XXX_Size() int {
	return xxx_messageInfo_PlayerInfo.Size(m)
}
func (m *PlayerInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_PlayerInfo.DiscardUnknown(m)
}

var xxx_messageInfo_PlayerInfo proto.InternalMessageInfo

func (m *PlayerInfo) GetId() uint64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *PlayerInfo) GetNickname() string {
	if m != nil {
		return m.Nickname
	}
	return ""
}

type LoginAck struct {
	Base                 *BaseAck    `protobuf:"bytes,1,opt,name=base,proto3" json:"base,omitempty"`
	PlayerInfo           *PlayerInfo `protobuf:"bytes,2,opt,name=player_info,json=playerInfo,proto3" json:"player_info,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *LoginAck) Reset()         { *m = LoginAck{} }
func (m *LoginAck) String() string { return proto.CompactTextString(m) }
func (*LoginAck) ProtoMessage()    {}
func (*LoginAck) Descriptor() ([]byte, []int) {
	return fileDescriptor_9f35044cfecd686d, []int{4}
}

func (m *LoginAck) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LoginAck.Unmarshal(m, b)
}
func (m *LoginAck) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LoginAck.Marshal(b, m, deterministic)
}
func (m *LoginAck) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LoginAck.Merge(m, src)
}
func (m *LoginAck) XXX_Size() int {
	return xxx_messageInfo_LoginAck.Size(m)
}
func (m *LoginAck) XXX_DiscardUnknown() {
	xxx_messageInfo_LoginAck.DiscardUnknown(m)
}

var xxx_messageInfo_LoginAck proto.InternalMessageInfo

func (m *LoginAck) GetBase() *BaseAck {
	if m != nil {
		return m.Base
	}
	return nil
}

func (m *LoginAck) GetPlayerInfo() *PlayerInfo {
	if m != nil {
		return m.PlayerInfo
	}
	return nil
}

type RoomInfo struct {
	RoomId               uint64            `protobuf:"varint,1,opt,name=room_id,json=roomId,proto3" json:"room_id,omitempty"`
	Type                 RoomInfo_RoomType `protobuf:"varint,2,opt,name=type,proto3,enum=pbprotocol.RoomInfo_RoomType" json:"type,omitempty"`
	IsStart              bool              `protobuf:"varint,3,opt,name=IsStart,proto3" json:"IsStart,omitempty"`
	CurrentPlayerNum     uint32            `protobuf:"varint,4,opt,name=current_player_num,json=currentPlayerNum,proto3" json:"current_player_num,omitempty"`
	PlayerInfos          []*PlayerInfo     `protobuf:"bytes,5,rep,name=player_infos,json=playerInfos,proto3" json:"player_infos,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *RoomInfo) Reset()         { *m = RoomInfo{} }
func (m *RoomInfo) String() string { return proto.CompactTextString(m) }
func (*RoomInfo) ProtoMessage()    {}
func (*RoomInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_9f35044cfecd686d, []int{5}
}

func (m *RoomInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RoomInfo.Unmarshal(m, b)
}
func (m *RoomInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RoomInfo.Marshal(b, m, deterministic)
}
func (m *RoomInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RoomInfo.Merge(m, src)
}
func (m *RoomInfo) XXX_Size() int {
	return xxx_messageInfo_RoomInfo.Size(m)
}
func (m *RoomInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_RoomInfo.DiscardUnknown(m)
}

var xxx_messageInfo_RoomInfo proto.InternalMessageInfo

func (m *RoomInfo) GetRoomId() uint64 {
	if m != nil {
		return m.RoomId
	}
	return 0
}

func (m *RoomInfo) GetType() RoomInfo_RoomType {
	if m != nil {
		return m.Type
	}
	return RoomInfo_TWO
}

func (m *RoomInfo) GetIsStart() bool {
	if m != nil {
		return m.IsStart
	}
	return false
}

func (m *RoomInfo) GetCurrentPlayerNum() uint32 {
	if m != nil {
		return m.CurrentPlayerNum
	}
	return 0
}

func (m *RoomInfo) GetPlayerInfos() []*PlayerInfo {
	if m != nil {
		return m.PlayerInfos
	}
	return nil
}

// join room
type JoinRoomReq struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *JoinRoomReq) Reset()         { *m = JoinRoomReq{} }
func (m *JoinRoomReq) String() string { return proto.CompactTextString(m) }
func (*JoinRoomReq) ProtoMessage()    {}
func (*JoinRoomReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_9f35044cfecd686d, []int{6}
}

func (m *JoinRoomReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_JoinRoomReq.Unmarshal(m, b)
}
func (m *JoinRoomReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_JoinRoomReq.Marshal(b, m, deterministic)
}
func (m *JoinRoomReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_JoinRoomReq.Merge(m, src)
}
func (m *JoinRoomReq) XXX_Size() int {
	return xxx_messageInfo_JoinRoomReq.Size(m)
}
func (m *JoinRoomReq) XXX_DiscardUnknown() {
	xxx_messageInfo_JoinRoomReq.DiscardUnknown(m)
}

var xxx_messageInfo_JoinRoomReq proto.InternalMessageInfo

type JoinRoomAck struct {
	Base                 *BaseAck  `protobuf:"bytes,1,opt,name=base,proto3" json:"base,omitempty"`
	RoomInfo             *RoomInfo `protobuf:"bytes,2,opt,name=room_info,json=roomInfo,proto3" json:"room_info,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *JoinRoomAck) Reset()         { *m = JoinRoomAck{} }
func (m *JoinRoomAck) String() string { return proto.CompactTextString(m) }
func (*JoinRoomAck) ProtoMessage()    {}
func (*JoinRoomAck) Descriptor() ([]byte, []int) {
	return fileDescriptor_9f35044cfecd686d, []int{7}
}

func (m *JoinRoomAck) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_JoinRoomAck.Unmarshal(m, b)
}
func (m *JoinRoomAck) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_JoinRoomAck.Marshal(b, m, deterministic)
}
func (m *JoinRoomAck) XXX_Merge(src proto.Message) {
	xxx_messageInfo_JoinRoomAck.Merge(m, src)
}
func (m *JoinRoomAck) XXX_Size() int {
	return xxx_messageInfo_JoinRoomAck.Size(m)
}
func (m *JoinRoomAck) XXX_DiscardUnknown() {
	xxx_messageInfo_JoinRoomAck.DiscardUnknown(m)
}

var xxx_messageInfo_JoinRoomAck proto.InternalMessageInfo

func (m *JoinRoomAck) GetBase() *BaseAck {
	if m != nil {
		return m.Base
	}
	return nil
}

func (m *JoinRoomAck) GetRoomInfo() *RoomInfo {
	if m != nil {
		return m.RoomInfo
	}
	return nil
}

func init() {
	proto.RegisterEnum("pbprotocol.LoginReq_LoginType", LoginReq_LoginType_name, LoginReq_LoginType_value)
	proto.RegisterEnum("pbprotocol.RoomInfo_RoomType", RoomInfo_RoomType_name, RoomInfo_RoomType_value)
	proto.RegisterType((*MsgPack)(nil), "pbprotocol.MsgPack")
	proto.RegisterType((*BaseAck)(nil), "pbprotocol.BaseAck")
	proto.RegisterType((*LoginReq)(nil), "pbprotocol.LoginReq")
	proto.RegisterType((*PlayerInfo)(nil), "pbprotocol.PlayerInfo")
	proto.RegisterType((*LoginAck)(nil), "pbprotocol.LoginAck")
	proto.RegisterType((*RoomInfo)(nil), "pbprotocol.RoomInfo")
	proto.RegisterType((*JoinRoomReq)(nil), "pbprotocol.JoinRoomReq")
	proto.RegisterType((*JoinRoomAck)(nil), "pbprotocol.JoinRoomAck")
}

func init() { proto.RegisterFile("digimon.proto", fileDescriptor_9f35044cfecd686d) }

var fileDescriptor_9f35044cfecd686d = []byte{
	// 449 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x51, 0xc1, 0x8e, 0xd3, 0x30,
	0x14, 0x5c, 0x37, 0xa1, 0x49, 0x5f, 0xda, 0x55, 0xe4, 0x45, 0x4b, 0x84, 0x04, 0x8a, 0x7c, 0x80,
	0x1e, 0x50, 0xa5, 0xed, 0x0a, 0x01, 0x47, 0xb8, 0x15, 0xb1, 0x50, 0x99, 0x15, 0x7b, 0xac, 0xdc,
	0xc6, 0x5b, 0x59, 0xdb, 0xd8, 0x59, 0xdb, 0x11, 0xea, 0xcf, 0xf0, 0x89, 0x7c, 0x03, 0x8a, 0xe3,
	0x26, 0x39, 0x20, 0x24, 0x6e, 0x33, 0x79, 0x6f, 0x9e, 0x67, 0x26, 0x30, 0x2b, 0xc4, 0x5e, 0x94,
	0x4a, 0x2e, 0x2a, 0xad, 0xac, 0xc2, 0x50, 0x6d, 0x1d, 0xd8, 0xa9, 0x03, 0x79, 0x0b, 0xd1, 0x8d,
	0xd9, 0xaf, 0xd9, 0xee, 0x01, 0x5f, 0xc2, 0x58, 0xab, 0xda, 0x72, 0x9d, 0xa1, 0x1c, 0xcd, 0x27,
	0xd4, 0x33, 0x8c, 0x21, 0x2c, 0x98, 0x65, 0xd9, 0x28, 0x47, 0xf3, 0x29, 0x75, 0x98, 0x5c, 0x43,
	0xf4, 0x89, 0x19, 0xfe, 0xb1, 0x95, 0x51, 0x6e, 0xea, 0x83, 0x75, 0xb2, 0x80, 0x7a, 0x86, 0x53,
	0x08, 0x6e, 0xcc, 0xde, 0xa9, 0x26, 0xb4, 0x81, 0xe4, 0x17, 0x82, 0xf8, 0x8b, 0xda, 0x0b, 0x49,
	0xf9, 0x23, 0x7e, 0x0e, 0x71, 0x6d, 0xb8, 0x96, 0xac, 0xe4, 0xfe, 0xbd, 0x8e, 0x37, 0xb3, 0x8a,
	0x19, 0xf3, 0x53, 0xe9, 0xc2, 0xeb, 0x3b, 0x8e, 0x97, 0x10, 0xda, 0x63, 0xc5, 0xb3, 0x20, 0x47,
	0xf3, 0xf3, 0xe5, 0xcb, 0x45, 0x9f, 0x65, 0x71, 0xba, 0xdd, 0x82, 0xdb, 0x63, 0xc5, 0xa9, 0xdb,
	0x25, 0xaf, 0x60, 0xd2, 0x7d, 0xc2, 0x09, 0x44, 0x3f, 0x84, 0x11, 0x56, 0xe9, 0xf4, 0x0c, 0x4f,
	0x21, 0x5e, 0x33, 0x63, 0xee, 0x94, 0x2e, 0x52, 0x44, 0xde, 0x03, 0xac, 0x0f, 0xec, 0xc8, 0xf5,
	0x4a, 0xde, 0x2b, 0x7c, 0x0e, 0x23, 0x51, 0x38, 0x6f, 0x21, 0x1d, 0x89, 0xa2, 0x71, 0x25, 0xc5,
	0xee, 0xc1, 0x39, 0xf6, 0xae, 0x4e, 0x9c, 0x1c, 0x7c, 0xb2, 0xa6, 0x90, 0xd7, 0x10, 0x6e, 0x99,
	0x69, 0x53, 0x25, 0xcb, 0x8b, 0xa1, 0x43, 0xdf, 0x19, 0x75, 0x0b, 0xf8, 0x1d, 0x24, 0x95, 0x7b,
	0x6e, 0x23, 0xe4, 0xbd, 0x72, 0x37, 0x93, 0xe5, 0xe5, 0x70, 0xbf, 0x77, 0x43, 0xa1, 0xea, 0x30,
	0xf9, 0x8d, 0x20, 0xa6, 0x4a, 0x95, 0xce, 0xe6, 0x33, 0x88, 0xb4, 0x52, 0xe5, 0xa6, 0xf3, 0x3a,
	0x6e, 0xe8, 0xaa, 0xc0, 0x57, 0xbe, 0xa9, 0x91, 0x6b, 0xea, 0xc5, 0xf0, 0xee, 0x49, 0xec, 0x40,
	0x5f, 0x14, 0xce, 0x20, 0x5a, 0x99, 0xef, 0x96, 0x69, 0xeb, 0xfa, 0x8d, 0xe9, 0x89, 0xe2, 0x37,
	0x80, 0x77, 0xb5, 0xd6, 0x5c, 0xda, 0x8d, 0xf7, 0x2c, 0xeb, 0x32, 0x0b, 0x73, 0x34, 0x9f, 0xd1,
	0xd4, 0x4f, 0x5a, 0xb7, 0x5f, 0xeb, 0x12, 0x7f, 0x80, 0xe9, 0x20, 0x99, 0xc9, 0x9e, 0xe4, 0xc1,
	0x3f, 0xa2, 0x25, 0x7d, 0x34, 0x43, 0x2e, 0xda, 0x68, 0xee, 0x57, 0x45, 0x10, 0xdc, 0xde, 0x7d,
	0x4b, 0xcf, 0xc8, 0x0c, 0x92, 0xcf, 0x4a, 0xc8, 0x66, 0x40, 0xf9, 0x23, 0x11, 0x3d, 0xfd, 0xaf,
	0xc2, 0xaf, 0x60, 0xd2, 0x56, 0xd5, 0xd7, 0xfd, 0xf4, 0x6f, 0xb5, 0xd0, 0x58, 0x7b, 0xb4, 0x1d,
	0xbb, 0xe1, 0xf5, 0x9f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x61, 0x87, 0xa2, 0xc6, 0x43, 0x03, 0x00,
	0x00,
}