package main

import (
	"os"
	"strconv"
)

const (
	//BalanceRoundRobin is a Round-Robin Load Balancer type constant. Services receive connections each after another.
	BalanceRoundRobin = "RR"

	//BalanceRandom is a Pseudo-Random Load Balancer type constant. Services receive connections in random order.
	BalanceRandom = "RND"

	//BalanceWeightedRoundRobin is Weighted Round-Robin Load Balancer type constant. Services receive connections each
	//after another according to the service's weight. Bigger weight - more connections.
	BalanceWeightedRoundRobin = "WRR"

	//BalanceWeightedRandom is a Weighted Pseudo-Random Load Balancer type constant. Services receive connections in
	//random order according to the service's weight. Bigger weight - higher chance to receive connections.
	BalanceWeightedRandom = "WRND"

	//REGISTRY_PREFIX is a string, which should be prepended to strings representing environment variables for diminishing
	//possible clash of environment names.
	RegistryPrefix = "RG_"
)

//EnvironmentConfig is a struct with Registry configuration parameters, which are read from the environment
type EnvironmentConfig struct {

	//RPCPort is a port number to bind RPC server to. This port should be available. This field is required to be greater than 0. Values - 0-65535.
	RPCPort uint16

	//TLSEnabled is a flag for for enabling or disabling TLS on HTTP server. If TLSEnabled is false, HTTPSPort, CertFile and KeyFile are ignored
	TLSEnabled bool

	//HTTPPort is a port number to bind HTTP server to. This port should be available.
	//If HTTPPort is 0, then the Registry should not start HTTP server. Values - 0-65535.
	HTTPPort uint16

	//HTTPSPort is a port number to bind HTTPS server to. This port should be available.
	//If HTTPPort is 0, then the Registry should not start HTTPS server. If TLSEnabled is true and HTTPS port, CertFile and keyFile are set, then HTTP port is ignored.
	//Values - 0-65535
	HTTPSPort uint16

	//CertFile is a full or relative path to the certificate file for TLS
	CertFile string

	//KeyFile is a full or relative path to the private key file for TLS
	KeyFile string

	//BalancerType is a type of load balancer Registry service should utilize. Possible values - "RR", "RND", "WRR", "WRND".
	//Default value is "RR", which stands for Round-Robin Load Balancer
	BalancerType string

	//Retries is a number of retries for network operations. Values - 0-255. Default value is 1, so a network action is performed only once
	Retries uint8

	//Timeout is a number of seconds service should wait before considering network operation to be timed out. Values - 0-65535.
	//Default value is 2 seconds
	Timeout uint16
}

func LoadConfigFromEnv() *EnvironmentConfig {
	env := &EnvironmentConfig{}
	//Getting RPC port
	rpcPort := os.Getenv(RegistryPrefix + "RPC_PORT")
	if rpcPort != "" {
		r, err := strconv.Atoi(rpcPort)
		if err != nil {
			env.RPCPort = 0
		} else {
			env.RPCPort = uint16(r)
		}
	} else {
		env.RPCPort = 0
	}

	//Getting HTTP port
	httpPort := os.Getenv(RegistryPrefix + "HTTP_PORT")
	if httpPort != "" {
		r, err := strconv.Atoi(httpPort)
		if err != nil {
			env.HTTPPort = 0
		} else {
			env.HTTPPort = uint16(r)
		}
	} else {
		env.HTTPPort = 0
	}

	//Getting HTTPS port
	shouldUseTLS := false
	httpsPort := os.Getenv(RegistryPrefix + "HTTPS_PORT")
	if httpsPort != "" {
		r, err := strconv.Atoi(httpsPort)
		if err != nil {
			env.HTTPSPort = 0
		} else {
			env.HTTPSPort = uint16(r)
			shouldUseTLS = true
		}
	} else {
		env.HTTPSPort = 0
	}

	//Getting Certificate file path
	env.CertFile = os.Getenv(RegistryPrefix + "CERT_FILE")
	shouldUseTLS = env.CertFile != ""

	//Getting Private Key file path
	env.KeyFile = os.Getenv(RegistryPrefix + "KEY_FILE")
	shouldUseTLS = env.KeyFile != ""

	env.TLSEnabled = shouldUseTLS
	//Check if all HTTPS parameters were valid and we can use TLS
	if shouldUseTLS {
		//set HTTP port to 0, so the service won't start HTTP server
		env.HTTPPort = 0
	}

	//Getting balancer type
	balancerType := os.Getenv(RegistryPrefix + "BALANCER")
	if balancerType == BalanceRandom || balancerType == BalanceRoundRobin ||
		balancerType == BalanceWeightedRandom || balancerType == BalanceWeightedRoundRobin {
		env.BalancerType = balancerType
	} else {
		env.BalancerType = BalanceRoundRobin
	}

	//Getting number of retries
	retries := os.Getenv(RegistryPrefix + "RETRIES")
	if retries != "" {
		r, err := strconv.Atoi(retries)
		if err != nil {
			env.Retries = 1
		} else {
			if r > 0 {
				env.Retries = uint8(r)
			} else {
				env.Retries = 1
			}
		}
	} else {
		env.Retries = 1
	}

	//Getting timeout
	timeout := os.Getenv(RegistryPrefix + "TIMEOUT")
	if timeout != "" {
		r, err := strconv.Atoi(timeout)
		if err != nil {
			env.Timeout = 2
		} else {
			if r > 0 {
				env.Timeout = uint16(r)
			} else {
				env.Timeout = 2
			}
		}
	}

	return env
}
