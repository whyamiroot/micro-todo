package registry

import (
	"fmt"
	"github.com/whyamiroot/micro-todo/proto"
	"golang.org/x/net/context"
	"sync"
)

type Registry struct {
	Lock      sync.Mutex
	Instances map[string][]*proto.Service
}

func (r *Registry) ListServicesTypes(c context.Context, _ *proto.Empty) (*proto.ServiceTypesList, error) {
	var types []*proto.ServiceType

	go func() {
		r.Lock.Lock()
		defer r.Lock.Unlock()
		for key := range r.Instances {
			types = append(types, &proto.ServiceType{Type: key})
		}
	}()

	select {
	case <-c.Done():
		//TODO add logging to Logger
		return nil, fmt.Errorf("CANCEL")
	default:
	}

	return &proto.ServiceTypesList{Types: types}, nil
}

func (r *Registry) ListByType(c context.Context, st *proto.ServiceType) (*proto.ServiceList, error) {
	if st == nil {
		return nil, fmt.Errorf("NULL")
	}
	var services []*proto.Service

	go func() {
		r.Lock.Lock()
		defer r.Lock.Unlock()
		services = r.Instances[st.Type]
	}()

	select {
	case <-c.Done():
		//TODO add logging to Logger
		return nil, fmt.Errorf("CANCEL")
	default:
	}

	return &proto.ServiceList{Services: services}, nil
}

func (r *Registry) BestInstance(c context.Context, st *proto.ServiceType) (*proto.Service, error) {
	panic("implement me")
}

func (r *Registry) Register(c context.Context, s *proto.Service) (*proto.RegistryResponse, error) {
	if s == nil {
		return nil, fmt.Errorf("NULL")
	}

	res := &proto.RegistryResponse{}

	go func() {
		r.Lock.Lock()
		defer r.Lock.Unlock()

	}()
}

func (r *Registry) exists(c context.Context, s *proto.Service) (bool, error) {
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

func (r *Registry) isValidSignature(s *proto.Service) bool {
	return true
}
