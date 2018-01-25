package utils

import (
	"fmt"
	"github.com/whyamiroot/micro-todo/proto"
)

type ServiceStringer proto.Service

func (ss ServiceStringer) String() string {
	return fmt.Sprintf("Host: %s, Port: %d, Health route: %s", ss.Host, ss.Port, ss.HealthRoute)
}
