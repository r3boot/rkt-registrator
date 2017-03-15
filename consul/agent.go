package consul

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/r3boot/rkt-registrator/rkt"
	"net/http"
)

func AgentPing() (result bool) {
	var (
		response *http.Response
		err      error
	)

	response, err = http.Get(Endpoint + "/v1/catalog/nodes")
	defer response.Body.Close()

	return err == nil
}

func Pod2Service(pod rkt.Pod) ConsulService {
	var (
		svc ConsulService
	)

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

	return svc
}

func Register(pod rkt.Pod) (err error) {
	var (
		data     *bytes.Buffer
		response *http.Response
		service  ConsulService
	)

	data = new(bytes.Buffer)

	service = Pod2Service(pod)

	if err = json.NewEncoder(data).Encode(service); err != nil {
		err = errors.New("Register(): Failed to marshal json: " + err.Error())
		return
	}

	response, err = http.Post(Endpoint+"/v1/agent/service/register", "application/json", data)
	response.Body.Close()

	if err != nil {
		err = errors.New("Register(): Failed to register service: " + err.Error())
		return
	}

	return
}

func DeRegister(pod rkt.Pod) (err error) {
	var (
		service  ConsulService
		response *http.Response
		uri      string
	)

	service = Pod2Service(pod)

	uri = Endpoint + "/v1/agent/service/deregister/" + service.ID

	response, err = http.Get(uri)
	defer response.Body.Close()

	if err != nil {
		err = errors.New("DeRegister(): Failed to deregister service: " + err.Error())
		return
	}

	return
}

func DeRegisterByID(ID string) (err error) {
	var (
		response *http.Response
		uri      string
	)

	uri = Endpoint + "/v1/agent/service/deregister/" + ID

	response, err = http.Get(uri)
	defer response.Body.Close()

	if err != nil {
		err = errors.New("DeRegister(): Failed to deregister service: " + err.Error())
		return
	}

	return
}

func FlushDuplicates() (err error) {
	var (
		all_services_uri string
		allServices      CatalogServices
		serviceDetails   []CatalogService
		serviceDetail    CatalogService
		lastModifiedID   string
		lastModified     int
		idToRemove       []string
		ID               string
		service          string
		service_uri      string
		uri              string
		response         *http.Response
	)

	all_services_uri = Endpoint + "/v1/catalog/services"
	service_uri = Endpoint + "/v1/catalog/service/"

	response, err = http.Get(all_services_uri)
	if err != nil {
		err = errors.New("FlushDuplicates(): Failed to list services: " + err.Error())
		return
	}

	if err = json.NewDecoder(response.Body).Decode(&allServices); err != nil {
		err = errors.New("FlushDuplicates(): Failed to decode json: " + err.Error())
		return
	}
	response.Body.Close()

	for service, _ = range allServices {
		uri = service_uri + "/" + service

		if response, err = http.Get(uri); err != nil {
			err = errors.New("FlushDuplicates(): Failed to fetch details for service")
			return
		}

		if err = json.NewDecoder(response.Body).Decode(&serviceDetails); err != nil {
			err = errors.New("FlushDuplicates(): Failed to decode json for service: " + err.Error())
			return
		}
		response.Body.Close()

		lastModified = 0
		lastModifiedID = ""
		idToRemove = make([]string, 100)
		for _, serviceDetail = range serviceDetails {
			if serviceDetail.Node != Worker {
				continue
			}

			if serviceDetail.ModifyIndex > lastModified {
				if lastModifiedID != "" {
					idToRemove = append(idToRemove, lastModifiedID)
				}
				lastModified = serviceDetail.ModifyIndex
				lastModifiedID = serviceDetail.ID
			} else {
				idToRemove = append(idToRemove, serviceDetail.ID)
			}
		}

		for _, ID = range idToRemove {
			if ID == "" {
				continue
			}
			Log.Debug("Deregistering " + ID)
			if err = DeRegisterByID(ID); err != nil {
				err = errors.New("FlushDuplicates(): Failed to deregister service: " + err.Error())
				return
			}
		}

	}

	return
}
