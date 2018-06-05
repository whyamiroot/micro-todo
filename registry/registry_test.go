package main

import (
	"strconv"
	"testing"

	"encoding/json"
	"net/http"
	"net/http/httptest"

	"net/url"

	"github.com/whyamiroot/micro-todo/proto"
	"golang.org/x/net/context"
)

type MockHealth struct{}

func (m MockHealth) ServeHTTP(wr http.ResponseWriter, r *http.Request) {
	res, err := json.Marshal(proto.Health{Up: true})
	if err != nil {
		panic("failed to marshal response")
	}
	wr.Write(res)
}

func createDummyRegistry(serviceCount int) *Registry {
	dummyRegistry := NewRegistry()
	for i := 0; i < serviceCount; i++ {
		dummyRegistry.Instances["test"] = append(dummyRegistry.Instances["test"], &proto.Service{
			Proto:  "http",
			Type:   "test",
			Port:   uint32(i),
			Weight: uint32(i + 1),
		})
	}
	dummyRegistry.LastRR["test"] = &LastWRRState{Index: -1, Weight: 0}
	dummyRegistry.WeightInfo.MaxWeight["test"] = getMaxWeight(dummyRegistry.Instances["test"])
	// find GCD for each service type
	dummyRegistry.WeightInfo.GCD["test"] = getGreatestCommonDivisorForWeights(dummyRegistry.Instances["test"])
	// calculate collective weight
	var weightSum uint32 = 0
	for _, instance := range dummyRegistry.Instances["test"] {
		weightSum += instance.Weight
	}
	dummyRegistry.WeightInfo.CollectiveWeights["test"] = weightSum

	return dummyRegistry
}

func customStupidBalancerFunc(registry *Registry, serviceType *proto.ServiceType) *proto.Service {
	if registry == nil || serviceType == nil {
		return nil
	}

	if len(registry.Instances[serviceType.Type]) == 0 {
		return nil
	}

	return registry.Instances[serviceType.Type][0]
}

func TestRegistry_Register(t *testing.T) {
	t.Run("Registering empty or non valid service should fail", func(t *testing.T) {
		t.Run("Registering service with empty hostname should fail", func(t *testing.T) {
			mock := MockHealth{}
			server := httptest.NewServer(mock)
			defer server.Close()

			u, _ := url.Parse(server.URL)
			port, _ := strconv.Atoi(u.Port())

			reg := NewRegistry()
			res, err := reg.Register(context.Background(), &proto.Service{
				Proto:  "http",
				Type:   "test",
				Host:   "",
				Port:   uint32(port),
				Health: "/",
			})

			if err == nil && res.Message == "OK" {
				t.Fatal("Registration with empty hostname should fail")
			}
		})

		t.Run("Registering with port bigger than 65535 should fail", func(t *testing.T) {
			mock := MockHealth{}
			server := httptest.NewServer(mock)
			defer server.Close()

			u, _ := url.Parse(server.URL)

			reg := NewRegistry()
			res, err := reg.Register(context.Background(), &proto.Service{
				Proto:  "http",
				Type:   "test",
				Host:   u.Hostname(),
				Port:   65536,
				Health: "/",
			})

			if err == nil && res.Message == "OK" {
				t.Fatal("Registration with port bigger than 65535 should fail")
			}
		})

		t.Run("Registering service, which does not respond to health check should fail", func(t *testing.T) {
			reg := NewRegistry()
			res, err := reg.Register(context.Background(), &proto.Service{
				Proto:  "http",
				Type:   "test",
				Host:   "localhost",
				Port:   12345,
				Health: "/",
			})

			if err == nil && res.Message == "OK" {
				t.Fatal("Registering dead service should fail")
			}
		})

		t.Run("Registering HTTPS service without HTTP port provided should fail", func(t *testing.T) {
			reg := NewRegistry()
			res, err := reg.Register(context.Background(), &proto.Service{
				Proto:  "http",
				Type:   "test",
				Host:   "localhost",
				Port:   0,
				Health: "/",
			})

			if err == nil && res.Message == "OK" {
				t.Fatal("Registering HTTPS service without HTTPS port should fail")
			}
		})

		t.Run("Registering service with unsupported protocol should fail", func(t *testing.T) {
			reg := NewRegistry()
			res, err := reg.Register(context.Background(), &proto.Service{
				Proto:  "tcp",
				Type:   "test",
				Host:   "localhost",
				Port:   100,
				Health: "/",
			})

			if err == nil && res.Message == "OK" {
				t.Fatal("Registering service with unsupported protocol should fail")
			}
		})
	})

	t.Run("Registering valid service should succeed", func(t *testing.T) {
		mock := MockHealth{}
		server := httptest.NewServer(mock)
		defer server.Close()

		u, _ := url.Parse(server.URL)
		port, _ := strconv.Atoi(u.Port())

		reg := NewRegistry()
		res, err := reg.Register(context.Background(), &proto.Service{
			Proto:  "http",
			Type:   "test",
			Host:   u.Hostname(),
			Port:   uint32(port),
			Health: "/",
		})

		if err != nil {
			t.Fatal("Failed to register mock server - ", err.Error())
		}
		if res.Message != "OK" {
			t.Fatal("Registration should have been successful")
		}
	})
}

func TestRegistry_SetCustomBalancerFunc(t *testing.T) {
	reg := createDummyRegistry(3)

	reg.SetCustomBalancerFunc(customStupidBalancerFunc)

	for i := 0; i < 5; i++ {
		best := reg.BalancerFunc(reg, &proto.ServiceType{Type: "test"})
		if best == nil {
			t.Fatal("Instance should not be nil")
		}

		if best.Weight != 1 {
			t.Fatal("Custom stupid balancing function should always returns first instance with weight 1")
		}
	}
}

func TestRegistry_RoundRobinBalancerFunc(t *testing.T) {
	reg := createDummyRegistry(3)
	var balanceHistory [9]uint32
	for i := 0; i < 9; i++ {
		best := RoundRobinBalancerFunc(reg, &proto.ServiceType{Type: "test"})
		if best == nil {
			t.Fatal("Instance should not be nil")
		}
		balanceHistory[i] = best.Weight
	}

	for i := 0; i < 3; i++ {
		if balanceHistory[i*3] != 1 {
			t.Fatal("Should be first instance with 0 weight, received weight - ", balanceHistory[i*3])
		}
		if balanceHistory[i*3+1] != 2 {
			t.Fatal("Should be first instance with 1 weight, received weight - ", balanceHistory[i*3+1])
		}
		if balanceHistory[i*3+2] != 3 {
			t.Fatal("Should be first instance with 2 weight, received weight - ", balanceHistory[i*3+2])
		}
	}
}

func TestRegistry_WeightedRoundRobinBalancerFunc(t *testing.T) {
	t.Run("Balancing is performed according to weights", func(t *testing.T) {
		reg := createDummyRegistry(3)
		balanceHistory := make(map[uint32]int)
		for i := 0; i < 6; i++ {
			best := WeightedRoundRobinBalancerFunc(reg, &proto.ServiceType{Type: "test"})
			if best == nil {
				t.Fatal("Instance should not be nil")
			}
			balanceHistory[best.Weight] += 1
		}

		for i, w := range balanceHistory {
			if w != balanceHistory[i] {
				t.Fatal("Weight should be equal to the number of times the service with this weight has been selected")
			}
		}
	})

	t.Run("Balancing result is the same every time", func(t *testing.T) {
		for i := 0; i < 1000; i++ {
			reg := createDummyRegistry(3)
			balanceHistory := make(map[uint32]int)
			for i := 0; i < 6; i++ {
				best := WeightedRoundRobinBalancerFunc(reg, &proto.ServiceType{Type: "test"})
				if best == nil {
					t.Fatal("Instance should not be nil")
				}
				balanceHistory[best.Weight] += 1
			}

			for i, w := range balanceHistory {
				if w != balanceHistory[i] {
					t.Fatal("Weight should be equal to the number of times the service with this weight has been selected")
				}
			}
		}
	})
}
