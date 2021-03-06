// Code generated by protoc-gen-go. DO NOT EDIT.
// source: registry_messages.proto

package proto

import proto1 "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto1.Marshal
var _ = fmt.Errorf
var _ = math.Inf

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
func (RegistryResponse_Code) EnumDescriptor() ([]byte, []int) { return fileDescriptor3, []int{2, 0} }

type ServiceType struct {
	Type string `protobuf:"bytes,1,opt,name=type" json:"type,omitempty"`
}

func (m *ServiceType) Reset()                    { *m = ServiceType{} }
func (m *ServiceType) String() string            { return proto1.CompactTextString(m) }
func (*ServiceType) ProtoMessage()               {}
func (*ServiceType) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{0} }

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
func (*ServiceTypesList) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{1} }

func (m *ServiceTypesList) GetTypes() []*ServiceType {
	if m != nil {
		return m.Types
	}
	return nil
}

// RegistryResponse represents response from Registry service
// It contains status code and optional message
type RegistryResponse struct {
	Status       RegistryResponse_Code `protobuf:"varint,1,opt,name=status,enum=proto.RegistryResponse_Code" json:"status,omitempty"`
	Message      string                `protobuf:"bytes,2,opt,name=message" json:"message,omitempty"`
	ServiceIndex uint32                `protobuf:"varint,3,opt,name=service_index,json=serviceIndex" json:"service_index,omitempty"`
}

func (m *RegistryResponse) Reset()                    { *m = RegistryResponse{} }
func (m *RegistryResponse) String() string            { return proto1.CompactTextString(m) }
func (*RegistryResponse) ProtoMessage()               {}
func (*RegistryResponse) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{2} }

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

func (m *RegistryResponse) GetServiceIndex() uint32 {
	if m != nil {
		return m.ServiceIndex
	}
	return 0
}

type Service struct {
	Proto     string   `protobuf:"bytes,1,opt,name=proto" json:"proto,omitempty"`
	Type      string   `protobuf:"bytes,2,opt,name=type" json:"type,omitempty"`
	Host      string   `protobuf:"bytes,3,opt,name=host" json:"host,omitempty"`
	Port      uint32   `protobuf:"varint,4,opt,name=port" json:"port,omitempty"`
	HttpPort  uint32   `protobuf:"varint,5,opt,name=httpPort" json:"httpPort,omitempty"`
	HttpsPort uint32   `protobuf:"varint,6,opt,name=httpsPort" json:"httpsPort,omitempty"`
	Routes    []string `protobuf:"bytes,7,rep,name=routes" json:"routes,omitempty"`
	Health    string   `protobuf:"bytes,8,opt,name=health" json:"health,omitempty"`
	Weight    uint32   `protobuf:"varint,9,opt,name=weight" json:"weight,omitempty"`
	Signature string   `protobuf:"bytes,10,opt,name=signature" json:"signature,omitempty"`
}

func (m *Service) Reset()                    { *m = Service{} }
func (m *Service) String() string            { return proto1.CompactTextString(m) }
func (*Service) ProtoMessage()               {}
func (*Service) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{3} }

func (m *Service) GetProto() string {
	if m != nil {
		return m.Proto
	}
	return ""
}

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

func (m *Service) GetPort() uint32 {
	if m != nil {
		return m.Port
	}
	return 0
}

func (m *Service) GetHttpPort() uint32 {
	if m != nil {
		return m.HttpPort
	}
	return 0
}

func (m *Service) GetHttpsPort() uint32 {
	if m != nil {
		return m.HttpsPort
	}
	return 0
}

func (m *Service) GetRoutes() []string {
	if m != nil {
		return m.Routes
	}
	return nil
}

func (m *Service) GetHealth() string {
	if m != nil {
		return m.Health
	}
	return ""
}

func (m *Service) GetWeight() uint32 {
	if m != nil {
		return m.Weight
	}
	return 0
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
func (*ServiceList) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{4} }

func (m *ServiceList) GetServices() []*Service {
	if m != nil {
		return m.Services
	}
	return nil
}

type ServiceInfo struct {
	Type  string `protobuf:"bytes,1,opt,name=type" json:"type,omitempty"`
	Index uint32 `protobuf:"varint,2,opt,name=index" json:"index,omitempty"`
}

func (m *ServiceInfo) Reset()                    { *m = ServiceInfo{} }
func (m *ServiceInfo) String() string            { return proto1.CompactTextString(m) }
func (*ServiceInfo) ProtoMessage()               {}
func (*ServiceInfo) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{5} }

func (m *ServiceInfo) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *ServiceInfo) GetIndex() uint32 {
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
func (*InstanceInfo) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{6} }

func (m *InstanceInfo) GetInstanceName() string {
	if m != nil {
		return m.InstanceName
	}
	return ""
}

func init() {
	proto1.RegisterType((*ServiceType)(nil), "proto.ServiceType")
	proto1.RegisterType((*ServiceTypesList)(nil), "proto.ServiceTypesList")
	proto1.RegisterType((*RegistryResponse)(nil), "proto.RegistryResponse")
	proto1.RegisterType((*Service)(nil), "proto.Service")
	proto1.RegisterType((*ServiceList)(nil), "proto.ServiceList")
	proto1.RegisterType((*ServiceInfo)(nil), "proto.ServiceInfo")
	proto1.RegisterType((*InstanceInfo)(nil), "proto.InstanceInfo")
	proto1.RegisterEnum("proto.RegistryResponse_Code", RegistryResponse_Code_name, RegistryResponse_Code_value)
}

func init() { proto1.RegisterFile("registry_messages.proto", fileDescriptor3) }

var fileDescriptor3 = []byte{
	// 472 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x52, 0xcf, 0x8f, 0x93, 0x40,
	0x14, 0x16, 0x4a, 0xf9, 0xf1, 0xb6, 0x5d, 0x27, 0x4f, 0xa3, 0x13, 0xb3, 0x87, 0x8a, 0x17, 0xe2,
	0xa1, 0x87, 0x6a, 0x62, 0x4c, 0xbc, 0x34, 0x2d, 0x9b, 0x10, 0x59, 0x76, 0x43, 0x59, 0xe3, 0x8d,
	0xa0, 0x3b, 0x16, 0x12, 0x0b, 0x84, 0x99, 0xaa, 0xfd, 0xc7, 0xbd, 0x78, 0x31, 0x33, 0x4c, 0xe9,
	0xae, 0xf1, 0xc4, 0xfb, 0xbe, 0xf7, 0xbd, 0x1f, 0xf3, 0xf1, 0xe0, 0x79, 0xc7, 0xb6, 0x15, 0x17,
	0xdd, 0x21, 0xdf, 0x31, 0xce, 0x8b, 0x2d, 0xe3, 0xf3, 0xb6, 0x6b, 0x44, 0x83, 0x63, 0xf5, 0xf1,
	0x5f, 0xc2, 0xd9, 0x86, 0x75, 0x3f, 0xaa, 0xaf, 0x2c, 0x3b, 0xb4, 0x0c, 0x11, 0x2c, 0x71, 0x68,
	0x19, 0x35, 0x66, 0x46, 0xe0, 0xa5, 0x2a, 0xf6, 0x3f, 0x00, 0xb9, 0x27, 0xe1, 0x71, 0xc5, 0x05,
	0x06, 0x30, 0x96, 0x39, 0x4e, 0x8d, 0xd9, 0x28, 0x38, 0x5b, 0x60, 0xdf, 0x74, 0x7e, 0x4f, 0x97,
	0xf6, 0x02, 0xff, 0xb7, 0x01, 0x24, 0xd5, 0x3b, 0xa4, 0x8c, 0xb7, 0x4d, 0xcd, 0x19, 0xbe, 0x05,
	0x9b, 0x8b, 0x42, 0xec, 0xb9, 0x1a, 0x74, 0xbe, 0xb8, 0xd0, 0xf5, 0xff, 0x0a, 0xe7, 0xab, 0xe6,
	0x8e, 0xa5, 0x5a, 0x8b, 0x14, 0x1c, 0xfd, 0x08, 0x6a, 0xaa, 0xfd, 0x8e, 0x10, 0x5f, 0xc1, 0x94,
	0xf7, 0xa3, 0xf3, 0xaa, 0xbe, 0x63, 0xbf, 0xe8, 0x68, 0x66, 0x04, 0xd3, 0x74, 0xa2, 0xc9, 0x48,
	0x72, 0xfe, 0x0e, 0x2c, 0xd9, 0x0e, 0xa7, 0xe0, 0xe5, 0xf9, 0x3a, 0xbc, 0x5c, 0xde, 0xc6, 0x19,
	0x79, 0x84, 0x36, 0x98, 0xd7, 0x1f, 0x89, 0x81, 0x67, 0xe0, 0x44, 0xc9, 0xa7, 0x65, 0x1c, 0xad,
	0x89, 0x89, 0x2e, 0x58, 0x97, 0xcb, 0x28, 0x26, 0x23, 0x9c, 0x80, 0xbb, 0x5a, 0x26, 0xab, 0x30,
	0x0e, 0xd7, 0xc4, 0xc2, 0x27, 0xf0, 0x38, 0xb9, 0xce, 0xf2, 0xe8, 0xea, 0x26, 0x0e, 0xaf, 0xc2,
	0x24, 0x0b, 0xd7, 0x64, 0x2c, 0xc5, 0xc9, 0x6d, 0x1c, 0x13, 0x1b, 0x01, 0xec, 0xf0, 0x73, 0xb4,
	0xc9, 0x36, 0xc4, 0xf1, 0xff, 0x18, 0xe0, 0x68, 0x3f, 0xf0, 0x29, 0xf4, 0x76, 0x6b, 0x5f, 0x7b,
	0x30, 0x98, 0x6d, 0x9e, 0xcc, 0x96, 0x5c, 0xd9, 0x70, 0xa1, 0x1e, 0xe0, 0xa5, 0x2a, 0x96, 0x5c,
	0xdb, 0x74, 0x82, 0x5a, 0xea, 0x51, 0x2a, 0xc6, 0x17, 0xe0, 0x96, 0x42, 0xb4, 0x37, 0x92, 0x1f,
	0x2b, 0x7e, 0xc0, 0x78, 0x01, 0x9e, 0x8c, 0xb9, 0x4a, 0xda, 0x2a, 0x79, 0x22, 0xf0, 0x19, 0xd8,
	0x5d, 0xb3, 0x17, 0x8c, 0x53, 0x67, 0x36, 0x0a, 0xbc, 0x54, 0x23, 0xc9, 0x97, 0xac, 0xf8, 0x2e,
	0x4a, 0xea, 0xaa, 0xd9, 0x1a, 0x49, 0xfe, 0x27, 0xab, 0xb6, 0xa5, 0xa0, 0x9e, 0x6a, 0xa5, 0x91,
	0x9c, 0xc2, 0xab, 0x6d, 0x5d, 0x88, 0x7d, 0xc7, 0x28, 0xa8, 0x92, 0x13, 0xe1, 0xbf, 0x1f, 0xee,
	0x4a, 0xdd, 0xcb, 0x6b, 0x70, 0xf5, 0xbf, 0x38, 0x9e, 0xcc, 0xf9, 0xc3, 0x93, 0x49, 0x87, 0xbc,
	0xff, 0x6e, 0x28, 0x8d, 0xea, 0x6f, 0xcd, 0xff, 0x4e, 0x52, 0xfa, 0xd9, 0xff, 0x67, 0x53, 0xad,
	0xd4, 0x03, 0x7f, 0x01, 0x93, 0xa8, 0xe6, 0xa2, 0xa8, 0x75, 0xa5, 0x0f, 0x93, 0x4a, 0xe3, 0xa4,
	0xd8, 0x1d, 0x3b, 0x3c, 0xe0, 0xbe, 0xd8, 0x6a, 0x8b, 0x37, 0x7f, 0x03, 0x00, 0x00, 0xff, 0xff,
	0x80, 0xc7, 0x35, 0xb1, 0x28, 0x03, 0x00, 0x00,
}
