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
			Proto:     "http",
			Type:      "test",
			HttpPort:  uint32(i),
			Signature: strconv.Itoa(i),
			Weight:    uint32(i + 1),
		})
	}
	dummyRegistry.lastRR["test"] = &LastWRRState{Index: -1, Weight: 0}
	dummyRegistry.weightInfo.MaxWeight["test"] = getMaxWeight(dummyRegistry.Instances["test"])
	// find GCD for each service type
	dummyRegistry.weightInfo.GCD["test"] = getGreatestCommonDivisorForWeights(dummyRegistry.Instances["test"])
	// calculate collective weight
	var weightSum uint32 = 0
	for _, instance := range dummyRegistry.Instances["test"] {
		weightSum += instance.Weight
	}
	dummyRegistry.weightInfo.CollectiveWeights["test"] = weightSum

	return dummyRegistry
}

func TestRegistry_Register(t *testing.T) {
	mock := MockHealth{}
	server := httptest.NewServer(mock)
	defer server.Close()

	u, _ := url.Parse(server.URL)
	port, _ := strconv.Atoi(u.Port())

	reg := NewRegistry()
	res, err := reg.Register(context.Background(), &proto.Service{
		Proto:     "http",
		Type:      "test",
		Host:      u.Hostname(),
		HttpPort:  uint32(port),
		Health:    "/",
		Signature: "12345",
	})

	if err != nil {
		t.Fatal("Failed to register mock server - ", err.Error())
	}
	if res.Message != "OK" {
		t.Fatal("Registration should have been successful")
	}
}

func TestRegistry_RoundRobinBalancerFunc(t *testing.T) {
	reg := createDummyRegistry(3)
	var balanceHistory [9]uint32
	for i := 0; i < 9; i++ {
		best := reg.RoundRobinBalancerFunc(&proto.ServiceType{"test"})
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
			best := reg.WeightedRoundRobinBalancerFunc(&proto.ServiceType{"test"})
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
				best := reg.WeightedRoundRobinBalancerFunc(&proto.ServiceType{"test"})
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
