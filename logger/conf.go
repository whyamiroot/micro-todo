package main

import (
	"net"
	"os"
	"strconv"
)

const (
	//LoggerPrefix is a string, which should be prepended to strings representing environment variables for diminishing
	//possible clash of environment names.
	LoggerPrefix = "LG_"
)

//EnvironmentConfig is a struct with Logger configuration parameters, which are read from the environment
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

	//Retries is a number of retries for network operations. Values - 0-255. Default value is 1, so a network action is performed only once
	Retries uint8

	//Timeout is a number of seconds service should wait before considering network operation to be timed out. Values - 0-65535.
	//Default value is 2 seconds
	Timeout uint16

	//RegistryHost is a host of the Registry service RPC server
	RegistryHost string

	//RegistryRPCPort is a RPC port of the Registry service RPC server
	RegistryRPCPort uint32
}

var config *EnvironmentConfig

//GetConfig loads configuration from environment variables, saves it and returns pointer to the configuration
func GetConfig() *EnvironmentConfig {
	if config == nil {
		config = loadConfigFromEnv()
	}
	return config
}

func loadConfigFromEnv() *EnvironmentConfig {
	env := &EnvironmentConfig{}
	//Getting RPC port
	rpcPort := os.Getenv(LoggerPrefix + "RPC_PORT")
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
	httpPort := os.Getenv(LoggerPrefix + "HTTP_PORT")
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
	httpsPort := os.Getenv(LoggerPrefix + "HTTPS_PORT")
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
	env.CertFile = os.Getenv(LoggerPrefix + "CERT_FILE")
	shouldUseTLS = env.CertFile != ""

	//Getting Private Key file path
	env.KeyFile = os.Getenv(LoggerPrefix + "KEY_FILE")
	shouldUseTLS = env.KeyFile != ""

	env.TLSEnabled = shouldUseTLS
	//Check if all HTTPS parameters were valid and we can use TLS
	if shouldUseTLS {
		//set HTTP port to 0, so the service won't start HTTP server
		env.HTTPPort = 0
	}

	//Getting number of retries
	retries := os.Getenv(LoggerPrefix + "RETRIES")
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
	timeout := os.Getenv(LoggerPrefix + "TIMEOUT")
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

	//Getting Registry address
	registry := os.Getenv(LoggerPrefix + "REGISTRY")
	if registry == "" {
		panic("Registry address not found in the environment")
	}
	if registryIP := net.ParseIP(registry); registryIP == nil {
		panic("Given registry address is not an IP address")
	}
	env.RegistryHost = registry

	//Getting Registry RPC port
	registryPort := os.Getenv(LoggerPrefix + "REGISTRY_PORT")
	var p uint32
	if p, err := strconv.ParseInt(registryPort, 10, 32); err != nil || p <= 0 {
		panic("Registry port cannot be zero or lower")
	}
	env.RegistryRPCPort = uint32(p)

	return env
}
