package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/whyamiroot/micro-todo/proto"
	"github.com/whyamiroot/micro-todo/proto/utils"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

//balancerFunc is a pointer to the load balancing function, which should return service type and service instance index
type BalancerFunc func(serviceType *proto.ServiceType) *proto.Service

//ServiceWeightInfo is a helper data structure for keeping information about weights
type ServiceWeightInfo struct {

	//CollectiveWeights is map of weight sums, where key is the service type and value is a sum of weights of all
	//instances of this service type
	CollectiveWeights map[string]int32

	//MaxWeight is a map of maximum weights, where key is the service type and value is the biggest weight of all
	//instances of this service type
	MaxWeight map[string]int32

	//GCD is a map of greatest common divisor for all weights, where key is the service type and value is the greatest
	//common divisor for all weights of instances of this service type
	GCD map[string]int32
}

//LastWRRState saves state of the last round of Weighted Round Robin balancing algorithm
type LastWRRState struct {
	Index  int32
	Weight int32
}

//Registry is a service registry which holds list of service instances by its type
type Registry struct {
	Instances map[string][]*proto.Service

	balancerFunc BalancerFunc
	lock         *sync.Mutex
	lastRR       map[string]*LastWRRState
	weightInfo   *ServiceWeightInfo
	rnd          *rand.Rand
}

//DeadService is a data structure which defines dead service instance which should be removed from registry
type DeadService struct {
	ServiceType string
	Index       int
}

//NewRegistry returns new Registry instance
func NewRegistry() *Registry {
	reg := &Registry{
		Instances: make(map[string][]*proto.Service),
		lock:      &sync.Mutex{},
		lastRR:    make(map[string]*LastWRRState),
		rnd:       rand.New(rand.NewSource(time.Now().UnixNano())),
		weightInfo: &ServiceWeightInfo{
			CollectiveWeights: make(map[string]int32),
			GCD:               make(map[string]int32),
			MaxWeight:         make(map[string]int32),
		},
	}
	return reg
}

//StartRegistryServiceAndListen starts gRPC and HTTP servers for Registry service
func (r *Registry) StartRegistryServiceAndListen() {
	r.StartHealthChecks()
	envConf := GetConfig()

	switch envConf.BalancerType {
	case BalanceRandom:
		r.balancerFunc = r.RandomBalancerFunc
	case BalanceRoundRobin:
		r.balancerFunc = r.RoundRobinBalancerFunc
	case BalanceWeightedRandom:
		r.balancerFunc = r.WeightedRandomBalancerFunc
	case BalanceWeightedRoundRobin:
		r.balancerFunc = r.WeightedRoundRobinBalancerFunc
	default:
		r.balancerFunc = r.WeightedRoundRobinBalancerFunc
	}

	if envConf.RPCPort == 0 {
		//TODO add logging to the logging service
		panic("No RPC port is specified, unable to start")
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", envConf.RPCPort))
	if err != nil {
		fmt.Printf("failed to listen: %v\n", err)
		os.Exit(1)
	}

	var opts []grpc.DialOption
	var cred credentials.TransportCredentials

	if !envConf.TLSEnabled {
		opts = append(opts, grpc.WithInsecure())
	} else {
		cred, err := credentials.NewServerTLSFromFile(envConf.CertFile, envConf.KeyFile)
		if err != nil {
			fmt.Println("Failed to load TLS credentials")
		}
		opts = append(opts, grpc.WithTransportCredentials(cred))
	}

	fmt.Println("Starting gRPC server...")
	var server *grpc.Server
	if !envConf.TLSEnabled {
		server = grpc.NewServer()
	} else {
		server = grpc.NewServer(grpc.Creds(cred))
	}
	proto.RegisterRegistryServiceServer(server, r)
	go server.Serve(lis)

	mux := runtime.NewServeMux()

	proto.RegisterRegistryServiceHandlerFromEndpoint(context.Background(), mux, fmt.Sprintf(":%d", envConf.RPCPort), opts)
	var httpErr error
	if envConf.TLSEnabled {
		fmt.Println("Starting HTTPS gateway...") //TODO add logging to the logging service
		httpErr = http.ListenAndServeTLS(fmt.Sprintf(":%d", envConf.HTTPSPort), envConf.CertFile, envConf.KeyFile, mux)
	} else {
		fmt.Println("Starting HTTP gateway...")
		httpErr = http.ListenAndServe(fmt.Sprintf(":%d", envConf.HTTPPort), mux)
	}
	if httpErr != nil {
		panic(httpErr.Error())
	}
}

func (r *Registry) StartHealthChecks() {
	go func() {
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

		//start goroutine for every service type; each goroutine checks health of all services and then sleeps for some time
		for {
			//TODO add logging to the logging service
			fmt.Println("Performing service health check...")
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
			r.Sanitize(corpses)
		}
	}()
}

//GetHealth returns health status of the Registry service. Good health requires running RPC and HTTP or HTTPS server
//
//Method: GET
//
//Resource: /registry/health
func (r *Registry) GetHealth(context.Context, *proto.Empty) (*proto.Health, error) {
	//check if RPC port and HTTP or HTTPS port are listened, which means that server is running
	isRPCUp := false
	config = GetConfig()
	_, err := net.Listen("tcp", ":"+fmt.Sprintf(":%d", config.RPCPort))
	if err != nil {
		isRPCUp = true
	}

	isHTTPUp := false
	if config.HTTPPort != 0 {
		_, err := net.Listen("tcp", ":"+fmt.Sprintf(":%d", config.HTTPPort))
		if err != nil {
			isHTTPUp = true
		}
	} else if config.HTTPSPort != 0 {
		_, err := net.Listen("tcp", ":"+fmt.Sprintf(":%d", config.HTTPSPort))
		if err != nil {
			isHTTPUp = true
		}
	}

	return &proto.Health{Up: isRPCUp && isHTTPUp}, nil
}

//GetInfo returns service instance information
//
//Method: GET
//
//Resource: /registry/service/types/{type}/{index}
func (r *Registry) GetInfo(c context.Context, si *proto.ServiceInfo) (*proto.Service, error) {
	if si == nil {
		return nil, fmt.Errorf("NULL")
	}
	res := make(chan *proto.Service)

	go func() {
		r.lock.Lock()
		defer r.lock.Unlock()
		if length := len(r.Instances[si.Type]); length != 0 && int(si.Index) < length && r.Instances[si.Type][si.Index] != nil {
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
//
//Method: GET
//
//Resource: /registry/service/{instanceName}
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

//ListServicesTypes lists all registered service types. E.g. `apigateway`, `auth` etc.
//
//Method: GET
//
//Resource: /registry/service/types
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
//
//Method: GET
//
//Resource: /registry/service/types/{type}
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

//BestInstance returns best instance of required service type according to the load balancer
//
//Method: GET
//
//Resource: /registry/service/types/{type}/best
func (r *Registry) BestInstance(c context.Context, st *proto.ServiceType) (*proto.Service, error) {
	instance := r.balancerFunc(st)
	if instance == nil {
		return nil, fmt.Errorf("NULL")
	}

	return instance, nil
}

//Register registers new service instance in the registry, making this instance available for serving requests and participating in the
//load balancing
//
//Method: POST
//
//Resource: /registry/service
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

		if len(r.Instances[s.Type]) == 0 {
			r.lastRR[s.Type] = &LastWRRState{Index: -1, Weight: 0}
		}
		r.Instances[s.Type] = append(r.Instances[s.Type], s)

		// add service weight to service type collective weight
		r.weightInfo.CollectiveWeights[s.Type] += s.Weight
		// find new max weight
		if r.weightInfo.MaxWeight[s.Type] < s.Weight {
			r.weightInfo.MaxWeight[s.Type] = s.Weight
		}
		// recalculate GCD
		r.weightInfo.GCD[s.Type] = gcd(r.weightInfo.GCD[s.Type], s.Weight)

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

//String returns string representation of service (without business routes and signature). Not synchronized
func (r *Registry) String(serviceType string, index int) string {
	instances := r.Instances[serviceType]
	if len(instances) == 0 || index > len(instances) || instances[index] == nil {
		return ""
	}

	return fmt.Sprintf("Service: %s, %s", instances[index].Type+strconv.Itoa(index), (utils.ServiceStringer(*instances[index])).String())
}

func getHealth(service *proto.Service) *proto.Health {
	//TODO use TLS if service has defined its HTTPS server
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

func (r *Registry) Sanitize(corpses []*DeadService) {
	r.lock.Lock()
	defer r.lock.Unlock()

	sanitizedRegistry := make(map[string][]*proto.Service) // contains slices of service slices without dead services for each service type
	deletes := 0
	var resulting []*proto.Service
	var serviceList []*proto.Service
	for _, deadService := range corpses {
		serviceList = r.Instances[deadService.ServiceType]
		if len(serviceList) == 0 || deadService.Index >= len(serviceList) {
			continue
		}

		if deadService.Index == 0 {
			sanitizedRegistry[deadService.ServiceType] = serviceList[1:]
		} else {
			//cut dead service from instance list
			resulting = serviceList[:deadService.Index-deletes]
			resulting = append(resulting, serviceList[deadService.Index-deletes+1:]...)
			sanitizedRegistry[deadService.ServiceType] = resulting
		}
		deletes++
	}

	r.Instances = sanitizedRegistry
	regeneratedTypes := make(map[string]bool)
	for _, deadService := range corpses {
		// skip regenerating weight information for service type, which is already processed
		if regeneratedTypes[deadService.ServiceType] {
			continue
		}

		// find max weight for each service type
		r.weightInfo.MaxWeight[deadService.ServiceType] = getMaxWeight(r.Instances[deadService.ServiceType])
		// find GCD for each service type
		r.weightInfo.GCD[deadService.ServiceType] = getGreatestCommonDivisorForWeights(r.Instances[deadService.ServiceType])
		// calculate collective weight
		var weightSum int32 = 0
		for _, instance := range r.Instances[deadService.ServiceType] {
			weightSum += instance.Weight
		}
		r.weightInfo.CollectiveWeights[deadService.ServiceType] = weightSum

		// mark this service type as processed
		regeneratedTypes[deadService.ServiceType] = true
	}
}

//RandomBalancerFunc returns random service from a list of service type services
func (r *Registry) RandomBalancerFunc(serviceType *proto.ServiceType) *proto.Service {
	if serviceType == nil {
		return nil
	}

	r.lock.Lock()
	defer r.lock.Unlock()

	instances := r.Instances[serviceType.Type]
	if len(instances) == 0 {
		return nil
	} else if len(instances) == 1 {
		return instances[0]
	}

	instanceIndex := r.rnd.Intn(len(instances))

	return instances[instanceIndex]
}

func (r *Registry) RoundRobinBalancerFunc(serviceType *proto.ServiceType) *proto.Service {
	if serviceType == nil {
		return nil
	}

	r.lock.Lock()
	defer r.lock.Unlock()

	instances := r.Instances[serviceType.Type]
	lastIndex := r.lastRR[serviceType.Type].Index

	if len(instances) == 0 {
		return nil
	} else if len(instances) == 1 {
		r.lastRR[serviceType.Type].Index = 0
		return instances[0]
	}

	// previous balancing op returned last instance in the list
	if lastIndex+1 >= int32(len(instances)) {
		r.lastRR[serviceType.Type].Index = 0
		return instances[0]
	}

	r.lastRR[serviceType.Type].Index = lastIndex + 1
	return instances[lastIndex+1]
}

//WeightedRandomBalancerFunc returns random service from a list of service type services
func (r *Registry) WeightedRandomBalancerFunc(serviceType *proto.ServiceType) *proto.Service {
	if serviceType == nil {
		return nil
	}

	instances := r.Instances[serviceType.Type]
	if len(instances) == 0 {
		return nil
	} else if len(instances) == 1 {
		return instances[0]
	}

	randVal := r.rnd.Int31n(r.weightInfo.CollectiveWeights[serviceType.Type])

	for _, instance := range instances {
		randVal -= instance.Weight
		if randVal <= 0 {
			return instance
		}
	}

	return nil
}

func (r *Registry) WeightedRoundRobinBalancerFunc(serviceType *proto.ServiceType) *proto.Service {
	if serviceType == nil {
		return nil
	}

	instances := r.Instances[serviceType.Type]
	length := int32(len(instances))
	if length == 0 {
		return nil
	} else if length == 1 {
		return instances[0]
	}

	currentWeight := r.lastRR[serviceType.Type].Weight
	i := r.lastRR[serviceType.Type].Index
	for {
		i = (i + 1) % length
		if i == 0 {
			currentWeight -= r.weightInfo.GCD[serviceType.Type]
			if currentWeight <= 0 {
				currentWeight = r.weightInfo.MaxWeight[serviceType.Type]
				if currentWeight == 0 {
					return nil
				}
			}
		}

		if r.Instances[serviceType.Type][i].Weight >= currentWeight {
			r.lastRR[serviceType.Type].Index = i
			r.lastRR[serviceType.Type].Weight = currentWeight
			return r.Instances[serviceType.Type][i]
		}
	}
}

func (r *Registry) exists(s *proto.Service) (bool, error) {
	if s == nil {
		return false, fmt.Errorf("NULL")
	}

	for _, v := range r.Instances[s.Type] {
		// Comparing signatures should be enough
		if v.Signature == s.Signature {
			return true, nil
		}
	}
	return false, nil
}

func (r *Registry) isValidSignature(s *proto.Service) (bool, error) {
	//TODO implement
	return true, nil
}

func getMaxWeight(list []*proto.Service) int32 {
	length := len(list)
	if length == 0 {
		return 0
	} else if length == 1 {
		return list[0].Weight
	}

	max := list[0].Weight
	for _, inst := range list {
		if inst.Weight > max {
			max = inst.Weight
		}
	}

	return max
}

func getGreatestCommonDivisorForWeights(list []*proto.Service) int32 {
	length := len(list)
	if length == 0 {
		return 0
	} else if length == 1 {
		return list[0].Weight
	}

	result := gcd(list[0].Weight, list[1].Weight)
	for index := 2; index < length; index++ {
		result = gcd(result, list[index].Weight)
	}

	return result
}

func gcd(x, y int32) int32 {
	for y != 0 {
		x, y = y, x%y
	}
	return x
}
