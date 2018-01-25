// Code generated by protoc-gen-go. DO NOT EDIT.
// source: messages.proto

/*
Package proto is a generated protocol buffer package.

It is generated from these files:
	messages.proto
	registry_service.proto

It has these top-level messages:
	Empty
	ServiceType
	ServiceTypesList
	RegistryResponse
	Service
	ServiceList
	ServiceInfo
	InstanceInfo
	Health
*/
package proto

import proto1 "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto1.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto1.ProtoPackageIsVersion2 // please upgrade the proto package

type RegistryResponse_Code int32

const (
	RegistryResponse___DEFAULT       RegistryResponse_Code = 0
	RegistryResponse_OK              RegistryResponse_Code = 1
	RegistryResponse_INVALID         RegistryResponse_Code = 2
	RegistryResponse_FAIL            RegistryResponse_Code = 3
	RegistryResponse_CANCELED        RegistryResponse_Code = 4
	RegistryResponse_NOT_IMPLEMENTED RegistryResponse_Code = 5
	RegistryResponse_NULL            RegistryResponse_Code = 6
	RegistryResponse_EXISTS          RegistryResponse_Code = 7
)

var RegistryResponse_Code_name = map[int32]string{
	0: "__DEFAULT",
	1: "OK",
	2: "INVALID",
	3: "FAIL",
	4: "CANCELED",
	5: "NOT_IMPLEMENTED",
	6: "NULL",
	7: "EXISTS",
}
var RegistryResponse_Code_value = map[string]int32{
	"__DEFAULT":       0,
	"OK":              1,
	"INVALID":         2,
	"FAIL":            3,
	"CANCELED":        4,
	"NOT_IMPLEMENTED": 5,
	"NULL":            6,
	"EXISTS":          7,
}

func (x RegistryResponse_Code) String() string {
	return proto1.EnumName(RegistryResponse_Code_name, int32(x))
}
func (RegistryResponse_Code) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{3, 0} }

type Empty struct {
}

func (m *Empty) Reset()                    { *m = Empty{} }
func (m *Empty) String() string            { return proto1.CompactTextString(m) }
func (*Empty) ProtoMessage()               {}
func (*Empty) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type ServiceType struct {
	Type string `protobuf:"bytes,1,opt,name=type" json:"type,omitempty"`
}

func (m *ServiceType) Reset()                    { *m = ServiceType{} }
func (m *ServiceType) String() string            { return proto1.CompactTextString(m) }
func (*ServiceType) ProtoMessage()               {}
func (*ServiceType) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *ServiceType) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

type ServiceTypesList struct {
	Types []*ServiceType `protobuf:"bytes,1,rep,name=types" json:"types,omitempty"`
}

func (m *ServiceTypesList) Reset()                    { *m = ServiceTypesList{} }
func (m *ServiceTypesList) String() string            { return proto1.CompactTextString(m) }
func (*ServiceTypesList) ProtoMessage()               {}
func (*ServiceTypesList) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *ServiceTypesList) GetTypes() []*ServiceType {
	if m != nil {
		return m.Types
	}
	return nil
}

// RegistryResponse represents response from Registry service
// It contains status code and optional message
type RegistryResponse struct {
	Status  RegistryResponse_Code `protobuf:"varint,1,opt,name=status,enum=proto.RegistryResponse_Code" json:"status,omitempty"`
	Message string                `protobuf:"bytes,2,opt,name=message" json:"message,omitempty"`
}

func (m *RegistryResponse) Reset()                    { *m = RegistryResponse{} }
func (m *RegistryResponse) String() string            { return proto1.CompactTextString(m) }
func (*RegistryResponse) ProtoMessage()               {}
func (*RegistryResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *RegistryResponse) GetStatus() RegistryResponse_Code {
	if m != nil {
		return m.Status
	}
	return RegistryResponse___DEFAULT
}

func (m *RegistryResponse) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

type Service struct {
	Type           string   `protobuf:"bytes,1,opt,name=type" json:"type,omitempty"`
	Host           string   `protobuf:"bytes,2,opt,name=host" json:"host,omitempty"`
	Port           int32    `protobuf:"varint,3,opt,name=port" json:"port,omitempty"`
	BusinessRoutes []string `protobuf:"bytes,4,rep,name=business_routes,json=businessRoutes" json:"business_routes,omitempty"`
	HealthRoute    string   `protobuf:"bytes,5,opt,name=health_route,json=healthRoute" json:"health_route,omitempty"`
	Signature      string   `protobuf:"bytes,6,opt,name=signature" json:"signature,omitempty"`
}

func (m *Service) Reset()                    { *m = Service{} }
func (m *Service) String() string            { return proto1.CompactTextString(m) }
func (*Service) ProtoMessage()               {}
func (*Service) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *Service) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *Service) GetHost() string {
	if m != nil {
		return m.Host
	}
	return ""
}

func (m *Service) GetPort() int32 {
	if m != nil {
		return m.Port
	}
	return 0
}

func (m *Service) GetBusinessRoutes() []string {
	if m != nil {
		return m.BusinessRoutes
	}
	return nil
}

func (m *Service) GetHealthRoute() string {
	if m != nil {
		return m.HealthRoute
	}
	return ""
}

func (m *Service) GetSignature() string {
	if m != nil {
		return m.Signature
	}
	return ""
}

type ServiceList struct {
	Services []*Service `protobuf:"bytes,1,rep,name=services" json:"services,omitempty"`
}

func (m *ServiceList) Reset()                    { *m = ServiceList{} }
func (m *ServiceList) String() string            { return proto1.CompactTextString(m) }
func (*ServiceList) ProtoMessage()               {}
func (*ServiceList) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *ServiceList) GetServices() []*Service {
	if m != nil {
		return m.Services
	}
	return nil
}

type ServiceInfo struct {
	Type  string `protobuf:"bytes,1,opt,name=type" json:"type,omitempty"`
	Index int32  `protobuf:"varint,2,opt,name=index" json:"index,omitempty"`
}

func (m *ServiceInfo) Reset()                    { *m = ServiceInfo{} }
func (m *ServiceInfo) String() string            { return proto1.CompactTextString(m) }
func (*ServiceInfo) ProtoMessage()               {}
func (*ServiceInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *ServiceInfo) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *ServiceInfo) GetIndex() int32 {
	if m != nil {
		return m.Index
	}
	return 0
}

type InstanceInfo struct {
	InstanceName string `protobuf:"bytes,1,opt,name=instanceName" json:"instanceName,omitempty"`
}

func (m *InstanceInfo) Reset()                    { *m = InstanceInfo{} }
func (m *InstanceInfo) String() string            { return proto1.CompactTextString(m) }
func (*InstanceInfo) ProtoMessage()               {}
func (*InstanceInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *InstanceInfo) GetInstanceName() string {
	if m != nil {
		return m.InstanceName
	}
	return ""
}

type Health struct {
	Up bool `protobuf:"varint,1,opt,name=up" json:"up,omitempty"`
}

func (m *Health) Reset()                    { *m = Health{} }
func (m *Health) String() string            { return proto1.CompactTextString(m) }
func (*Health) ProtoMessage()               {}
func (*Health) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func (m *Health) GetUp() bool {
	if m != nil {
		return m.Up
	}
	return false
}

func init() {
	proto1.RegisterType((*Empty)(nil), "proto.Empty")
	proto1.RegisterType((*ServiceType)(nil), "proto.ServiceType")
	proto1.RegisterType((*ServiceTypesList)(nil), "proto.ServiceTypesList")
	proto1.RegisterType((*RegistryResponse)(nil), "proto.RegistryResponse")
	proto1.RegisterType((*Service)(nil), "proto.Service")
	proto1.RegisterType((*ServiceList)(nil), "proto.ServiceList")
	proto1.RegisterType((*ServiceInfo)(nil), "proto.ServiceInfo")
	proto1.RegisterType((*InstanceInfo)(nil), "proto.InstanceInfo")
	proto1.RegisterType((*Health)(nil), "proto.Health")
	proto1.RegisterEnum("proto.RegistryResponse_Code", RegistryResponse_Code_name, RegistryResponse_Code_value)
}

func init() { proto1.RegisterFile("messages.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 446 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x91, 0xcf, 0x6e, 0xd3, 0x40,
	0x10, 0xc6, 0xb1, 0xe3, 0x3f, 0xc9, 0x24, 0xb8, 0xab, 0x81, 0x83, 0x0f, 0x3d, 0xa4, 0x7b, 0xc1,
	0xe2, 0x90, 0x43, 0x40, 0x42, 0x48, 0x5c, 0xa2, 0xc4, 0x15, 0x16, 0xae, 0x8b, 0x36, 0x2e, 0xe2,
	0x16, 0xb9, 0x74, 0x49, 0x2c, 0x11, 0xdb, 0xf2, 0xac, 0x11, 0x79, 0x2a, 0x9e, 0x85, 0x37, 0x42,
	0x5e, 0x3b, 0xa1, 0x45, 0x39, 0x79, 0xe6, 0x9b, 0xdf, 0x8c, 0xe7, 0xdb, 0x01, 0x6f, 0x2f, 0x89,
	0xb2, 0xad, 0xa4, 0x59, 0x55, 0x97, 0xaa, 0x44, 0x5b, 0x7f, 0xb8, 0x0b, 0x76, 0xb8, 0xaf, 0xd4,
	0x81, 0x5f, 0xc1, 0x78, 0x2d, 0xeb, 0x9f, 0xf9, 0x37, 0x99, 0x1e, 0x2a, 0x89, 0x08, 0x96, 0x3a,
	0x54, 0xd2, 0x37, 0xa6, 0x46, 0x30, 0x12, 0x3a, 0xe6, 0x1f, 0x80, 0x3d, 0x42, 0x28, 0xce, 0x49,
	0x61, 0x00, 0x76, 0x5b, 0x23, 0xdf, 0x98, 0x0e, 0x82, 0xf1, 0x1c, 0xbb, 0xe9, 0xb3, 0x47, 0x9c,
	0xe8, 0x00, 0xfe, 0xc7, 0x00, 0x26, 0xe4, 0x36, 0x27, 0x55, 0x1f, 0x84, 0xa4, 0xaa, 0x2c, 0x48,
	0xe2, 0x5b, 0x70, 0x48, 0x65, 0xaa, 0x21, 0xfd, 0x23, 0x6f, 0x7e, 0xd9, 0xf7, 0xff, 0x0f, 0xce,
	0x96, 0xe5, 0x83, 0x14, 0x3d, 0x8b, 0x3e, 0xb8, 0xbd, 0x1b, 0xdf, 0xd4, 0xfb, 0x1d, 0x53, 0xbe,
	0x07, 0xab, 0x25, 0xf1, 0x39, 0x8c, 0x36, 0x9b, 0x55, 0x78, 0xbd, 0xb8, 0x8b, 0x53, 0xf6, 0x0c,
	0x1d, 0x30, 0x6f, 0x3f, 0x31, 0x03, 0xc7, 0xe0, 0x46, 0xc9, 0x97, 0x45, 0x1c, 0xad, 0x98, 0x89,
	0x43, 0xb0, 0xae, 0x17, 0x51, 0xcc, 0x06, 0x38, 0x81, 0xe1, 0x72, 0x91, 0x2c, 0xc3, 0x38, 0x5c,
	0x31, 0x0b, 0x5f, 0xc0, 0x45, 0x72, 0x9b, 0x6e, 0xa2, 0x9b, 0xcf, 0x71, 0x78, 0x13, 0x26, 0x69,
	0xb8, 0x62, 0x76, 0x0b, 0x27, 0x77, 0x71, 0xcc, 0x1c, 0x04, 0x70, 0xc2, 0xaf, 0xd1, 0x3a, 0x5d,
	0x33, 0x97, 0xff, 0x36, 0xc0, 0xed, 0xad, 0x9e, 0x7b, 0xb1, 0x56, 0xdb, 0x95, 0xa4, 0xfa, 0x2d,
	0x75, 0xdc, 0x6a, 0x55, 0x59, 0x2b, 0x7f, 0x30, 0x35, 0x02, 0x5b, 0xe8, 0x18, 0x5f, 0xc1, 0xc5,
	0x7d, 0x43, 0x79, 0x21, 0x89, 0x36, 0x75, 0xd9, 0x28, 0x49, 0xbe, 0x35, 0x1d, 0x04, 0x23, 0xe1,
	0x1d, 0x65, 0xa1, 0x55, 0xbc, 0x82, 0xc9, 0x4e, 0x66, 0x3f, 0xd4, 0xae, 0xc3, 0x7c, 0x5b, 0x0f,
	0x1e, 0x77, 0x9a, 0x66, 0xf0, 0x12, 0x46, 0x94, 0x6f, 0x8b, 0x4c, 0x35, 0xb5, 0xf4, 0x1d, 0x5d,
	0xff, 0x27, 0xf0, 0xf7, 0xa7, 0x33, 0xeb, 0xf3, 0xbd, 0x86, 0x21, 0x75, 0xe9, 0xf1, 0x82, 0xde,
	0xd3, 0x0b, 0x8a, 0x53, 0x9d, 0xbf, 0x3b, 0xb5, 0x46, 0xc5, 0xf7, 0xf2, 0xac, 0xdf, 0x97, 0x60,
	0xe7, 0xc5, 0x83, 0xfc, 0xa5, 0x0d, 0xdb, 0xa2, 0x4b, 0xf8, 0x1c, 0x26, 0x51, 0x41, 0x2a, 0x2b,
	0xfa, 0x4e, 0x0e, 0x93, 0xbc, 0xcf, 0x93, 0x6c, 0x7f, 0x9c, 0xf0, 0x44, 0xe3, 0x3e, 0x38, 0x1f,
	0xb5, 0x29, 0xf4, 0xc0, 0x6c, 0x2a, 0xcd, 0x0c, 0x85, 0xd9, 0x54, 0xf7, 0x8e, 0xde, 0xef, 0xcd,
	0xdf, 0x00, 0x00, 0x00, 0xff, 0xff, 0xf0, 0x6c, 0x06, 0xed, 0xd1, 0x02, 0x00, 0x00,
}
