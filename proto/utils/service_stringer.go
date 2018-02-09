package utils

import (
	"fmt"
	"github.com/whyamiroot/micro-todo/proto"
)

type ServiceStringer proto.Service

func (ss ServiceStringer) String() string {
	return fmt.Sprintf("Host: %s, Port: %d, HttpPort: %d, Health route: %s", ss.Host, ss.Port, ss.HttpPort, ss.Health)
}
