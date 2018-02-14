package main

import (
	"encoding/json"
	"fmt"
	"github.com/whyamiroot/micro-todo/proto"
	"github.com/whyamiroot/micro-todo/proto/utils"
	"golang.org/x/net/context"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Registry struct {
	lock      *sync.Mutex
	Instances map[string][]*proto.Service
}

type DeadService struct {
	ServiceType string
	Index       int
}

func getHealth(service *proto.Service) *proto.Health {
	healthURL := service.Proto + "://" + service.Host + ":" + strconv.Itoa(int(service.HttpPort)) + service.Health
	resp, err := http.Get(healthURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	health := &proto.Health{}
	err = json.Unmarshal(bytes, health)
	if err != nil {
		return nil
	}
	return health
}

func sanitize() {

}

func (r *Registry) StartHealthChecks() {
	corpsesChan := make(chan *DeadService)
	var corpses []*DeadService
	//read channel for dead services and append them to the dead services list
	go func() {
		for {
			select {
			case corpse := <-corpsesChan:
				corpses = append(corpses, corpse)
			default:
			}
		}
	}()

	//
	for {
		for sType := range r.Instances {
			go func(serviceType string) {
				for index, service := range r.Instances[serviceType] {
					health := getHealth(service)
					if health == nil || !health.Up {
						corpsesChan <- &DeadService{ServiceType: serviceType, Index: index}
					}
				}
			}(sType)
		}
		time.Sleep(3 * time.Minute) //TODO Move to the conf
	}
}

func (r *Registry) GetHealth(context.Context, *proto.Empty) (*proto.Health, error) {
	return &proto.Health{Up: true}, nil
}

//GetInfo returns service instance information
func (r *Registry) GetInfo(c context.Context, si *proto.ServiceInfo) (*proto.Service, error) {
	if si == nil {
		return nil, fmt.Errorf("NULL")
	}
	res := make(chan *proto.Service)

	go func() {
		r.lock.Lock()
		defer r.lock.Unlock()
		if l := len(r.Instances[si.Type]); l != 0 && int(si.Index) < l && r.Instances[si.Type][si.Index] != nil {
			res <- r.Instances[si.Type][si.Index]
		} else {
			res <- nil
		}
	}()

	for {
		select {
		case <-c.Done():
			//TODO add logging to Logger
			return nil, fmt.Errorf("CANCELED")
		case service := <-res:
			if service == nil {
				return nil, fmt.Errorf("INVALID")
			} else {
				return service, nil
			}
		default:
		}
	}
}

//GetInstanceInfo returns service instance information. Instance is specified in a following format - `type-index`
func (r *Registry) GetInstanceInfo(c context.Context, i *proto.InstanceInfo) (*proto.Service, error) {
	if i == nil {
		return nil, fmt.Errorf("NULL")
	}
	req := strings.Split(i.InstanceName, "-")
	if len(req) != 2 {
		return nil, fmt.Errorf("INVALID")
	}

	serviceType := req[0]
	index, err := strconv.Atoi(req[1])
	if err != nil {
		return nil, fmt.Errorf("INVALID")
	}

	res := make(chan *proto.Service)
	go func() {
		r.lock.Lock()
		defer r.lock.Unlock()
		if l := len(r.Instances[serviceType]); l != 0 && int(index) < l && r.Instances[serviceType][index] != nil {
			res <- r.Instances[serviceType][index]
		} else {
			res <- nil
		}
	}()

	for {
		select {
		case <-c.Done():
			//TODO add logging to Logger
			return nil, fmt.Errorf("CANCELED")
		case service := <-res:
			if service == nil {
				return nil, fmt.Errorf("INVALID")
			} else {
				return service, nil
			}
		default:
		}
	}
}

//NewRegistry returns new Registry instance
func NewRegistry() *Registry {
	reg := &Registry{Instances: make(map[string][]*proto.Service), lock: &sync.Mutex{}}
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
		return
	}()

	for {
		select {
		case <-c.Done():
			//TODO add logging to Logger
			return nil, fmt.Errorf("CANCELED")
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
		return
	}()

	for {
		select {
		case <-c.Done():
			//TODO add logging to Logger
			return nil, fmt.Errorf("CANCELED")
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
			return
		}

		if b, err := r.isValidSignature(s); !b {
			if err != nil {
				res <- &proto.RegistryResponse{Status: proto.RegistryResponse_NULL, Message: "No service received"}
			} else {
				res <- &proto.RegistryResponse{Status: proto.RegistryResponse_INVALID, Message: "Service signature is invalid"}
			}
			return
		}

		health := getHealth(s)
		if health == nil || !health.Up {
			res <- &proto.RegistryResponse{Status: proto.RegistryResponse_NOT_IMPLEMENTED, Message: "Health check failed"}
			return
		}

		r.Instances[s.Type] = append(r.Instances[s.Type], s)
		res <- &proto.RegistryResponse{Status: proto.RegistryResponse_OK, Message: "OK"}
	}()

	for {
		select {
		case <-c.Done():
			//TODO add logging to Logger
			return &proto.RegistryResponse{Status: proto.RegistryResponse_CANCELED, Message: "Registration cancelled"}, nil
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
		// Comparing host, port and signature should be enough
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
