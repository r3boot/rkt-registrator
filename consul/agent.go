package consul

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/r3boot/rkt-registrator/rkt"
	"net/http"
)

func AgentPing() (result bool) {
	_, err := http.Get(Endpoint + "/v1/catalog/nodes")
	return err == nil
}

func Pod2Service(pod rkt.Pod) ConsulService {
	var svc ConsulService
	svc.Check = make(map[string]string)
	svc.Tags = make([]string, 0)

	svc.ID = Worker + "-" + pod.Name
	svc.Ipaddr = pod.IpAddress
	svc.Port = pod.Consul.Port

	if pod.Consul.Dns != "" {
		svc.Name = pod.Consul.Dns
	} else {
		svc.Name = pod.Name
	}

	svc.Check["Name"] = pod.Consul.Check.Name
	svc.Check[pod.Consul.Check.Type] = pod.Consul.Check.Target
	svc.Check["Interval"] = pod.Consul.Check.Interval
	svc.Check["Timeout"] = pod.Consul.Check.Timeout

	Log.Debug(svc)

	return svc
}

func Register(pod rkt.Pod) (err error) {
	data := new(bytes.Buffer)

	service := Pod2Service(pod)

	if err = json.NewEncoder(data).Encode(service); err != nil {
		err = errors.New("Register(): Failed to marshal json: " + err.Error())
		return
	}

	_, err = http.Post(Endpoint+"/v1/agent/service/register", "application/json", data)
	if err != nil {
		err = errors.New("Register(): Failed to register service: " + err.Error())
		return
	}

	return
}

func Deregister(pod rkt.Pod) (err error) {
	service := Pod2Service(pod)

	uri := Endpoint + "/v1/agent/service/deregister/" + service.ID

	if _, err = http.Get(uri); err != nil {
		err = errors.New("Deregister(): Failed to deregister service: " + err.Error())
		return
	}

	return
}
