package main

import (
	"fmt"
	"github.com/whyamiroot/micro-todo/proto"
	"github.com/whyamiroot/micro-todo/proto/utils"
	"golang.org/x/net/context"
	"strconv"
	"sync"
)

type Registry struct {
	lock      sync.Mutex
	Instances map[string][]*proto.Service
}

func NewRegistry() *Registry {
	reg := &Registry{Instances: make(map[string][]*proto.Service)}
	return reg
}

//String returns string representation of service (without business routes and signature). Not synchronized
func (r *Registry) String(serviceType string, index int) string {
	instances := r.Instances[serviceType]
	if len(instances) == 0 || index > len(instances) || instances[index] == nil {
		return ""
	}

	return fmt.Sprintf("Service: %s, %s", instances[index].Type+strconv.Itoa(index), (utils.ServiceStringer(*instances[index])).String())
}

//ListServicesTypes lists all registered service types. E.g. `apigateway`, `auth` etc.
func (r *Registry) ListServicesTypes(c context.Context, _ *proto.Empty) (*proto.ServiceTypesList, error) {
	res := make(chan []*proto.ServiceType)

	go func() {
		var types []*proto.ServiceType
		r.lock.Lock()
		defer r.lock.Unlock()
		for key := range r.Instances {
			types = append(types, &proto.ServiceType{Type: key})
		}
		res <- types
	}()

	for {
		select {
		case <-c.Done():
			//TODO add logging to Logger
			return nil, fmt.Errorf("CANCEL")
		case types := <-res:
			return &proto.ServiceTypesList{Types: types}, nil
		default:
		}
	}
}

//ListByType lists all service instances of specified type
func (r *Registry) ListByType(c context.Context, st *proto.ServiceType) (*proto.ServiceList, error) {
	if st == nil {
		return nil, fmt.Errorf("NULL")
	}

	res := make(chan []*proto.Service)

	go func() {
		r.lock.Lock()
		defer r.lock.Unlock()
		res <- r.Instances[st.Type]
	}()

	for {
		select {
		case <-c.Done():
			//TODO add logging to Logger
			return nil, fmt.Errorf("CANCEL")
		case services := <-res:
			return &proto.ServiceList{Services: services}, nil
		default:
		}
	}
}

func (r *Registry) BestInstance(c context.Context, st *proto.ServiceType) (*proto.Service, error) {
	panic("implement me")
}

func (r *Registry) Register(c context.Context, s *proto.Service) (*proto.RegistryResponse, error) {
	if s == nil {
		return &proto.RegistryResponse{Status: proto.RegistryResponse_NULL, Message: "No service received"}, nil
	}

	res := make(chan *proto.RegistryResponse)

	go func() {
		r.lock.Lock()
		defer r.lock.Unlock()

		if b, err := r.exists(s); b {
			if err != nil {
				res <- &proto.RegistryResponse{Status: proto.RegistryResponse_NULL, Message: "No service received"}
			} else {
				res <- &proto.RegistryResponse{Status: proto.RegistryResponse_EXISTS, Message: "Service is already registered"}
			}
		}

		if b, err := r.isValidSignature(s); !b {
			if err != nil {
				res <- &proto.RegistryResponse{Status: proto.RegistryResponse_NULL, Message: "No service received"}
			} else {
				res <- &proto.RegistryResponse{Status: proto.RegistryResponse_INVALID, Message: "Service signature is invalid"}
			}
		}

		r.Instances[s.Type] = append(r.Instances[s.Type], s)
		res <- &proto.RegistryResponse{Status: proto.RegistryResponse_OK, Message: "OK"}
	}()

	for {
		select {
		case <-c.Done():
			//TODO add logging to Logger
			return &proto.RegistryResponse{Status: proto.RegistryResponse_CANCELED, Message: "Registration canceled"}, nil
		case response := <-res:
			return response, nil
		default:
		}
	}
}

func (r *Registry) exists(s *proto.Service) (bool, error) {
	if s == nil {
		return false, fmt.Errorf("NULL")
	}

	for _, v := range r.Instances[s.Type] {
		// Comparing host:port + signature should be enough
		if v.Host == s.Host && v.Port == s.Port && v.Signature == s.Signature {
			return true, nil
		}
	}
	return false, nil
}

func (r *Registry) isValidSignature(s *proto.Service) (bool, error) {
	//TODO implement
	return true, nil
}
