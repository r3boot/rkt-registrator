package consul

import (
	"github.com/r3boot/rkt-registrator/utils"
)

var (
	Log       utils.Log
	Endpoint  string
	Worker    string
	Available bool = false
)

func Setup(l utils.Log, endpoint string, worker string) (err error) {
	Log = l
	Endpoint = endpoint
	Worker = worker

	if Available = AgentPing(); Available {
		Log.Debug("Consul api reachable")
	} else {
		Log.Warning("Failed to ping consul api")
	}

	return
}
