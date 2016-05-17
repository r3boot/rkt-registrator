package rkt

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

const D_CONSUL_CHECK_INTERVAL = "10s"
const D_CONSUL_CHECK_TIMEOUT = "1s"

func GetPodManifest(uuid string) (name string, manifest RktManifest, err error) {
	var pod_dir string
	var manifest_path string
	var fs os.FileInfo
	var fd *os.File
	var data []byte

	pod_dir = Rkt_dir + "/pods/run/" + uuid

	apps, err := ioutil.ReadDir(pod_dir + "/appsinfo")
	if err != nil {
		err = errors.New("GetPodManifest(): Failed to get apps: " + err.Error())
		return
	}

	if len(apps) > 1 {
		Log.Warning("Only a single application per container is supported atm")
	}

	name = apps[0].Name()

	manifest_path = pod_dir + "/appsinfo/" + name + "/manifest"

	fs, err = os.Stat(manifest_path)
	if err != nil {
		err = errors.New("GetPodManifest(): Manifest file not found: " + err.Error())
		return
	}

	data = make([]byte, fs.Size())
	if fd, err = os.Open(manifest_path); err != nil {
		err = errors.New("GetPodManifest(): Failed to open manifest file: " + err.Error())
		return
	}
	defer fd.Close()

	if _, err = fd.Read(data); err != nil {
		err = errors.New("GetPodManifest(): Failed to read manfest file: " + err.Error())
		return
	}

	if err = json.Unmarshal(data, &manifest); err != nil {
		err = errors.New("GetPodManifest(): Failed to unmarshal to json: " + err.Error())
		return
	}

	return
}

func GetIpUuid(ip_file string) (uuid string, err error) {
	var fs os.FileInfo
	var fd *os.File
	var data []byte

	if fs, err = os.Stat(ip_file); err != nil {
		err = errors.New("GetIpUuid(): Failed to stat ip file: " + err.Error())
		return
	}

	data = make([]byte, fs.Size())
	if fd, err = os.Open(ip_file); err != nil {
		err = errors.New("GetIpUuid(): Failed to open ip file: " + err.Error())
		return
	}
	defer fd.Close()

	if _, err = fd.Read(data); err != nil {
		err = errors.New("GetIpUuid(): Failed to read ip file: " + err.Error())
		return
	}

	uuid = string(data)

	return
}

func GetNetworkData(net_dir string) (ipuuids map[string]string, err error) {
	var ip string
	var uuid string

	ipuuids = make(map[string]string)

	// this function assumes that net_dir exists
	ips, err := ioutil.ReadDir(net_dir)
	if err != nil {
		err = errors.New("GetNetworkData(): Failed to read list of ips: " + err.Error())
		return
	}

	for _, ipaddr := range ips {
		ip = ipaddr.Name()
		ip_file := net_dir + "/" + ip
		if uuid, err = GetIpUuid(ip_file); err != nil {
			return
		}
		ipuuids[uuid] = ip
	}

	return
}

func GetPods(netname *string) (pods map[string]Pod, err error) {
	var manifest RktManifest
	var netdata map[string]string

	pods = make(map[string]Pod)

	// Read manifests of all running pods
	uuids, err := ioutil.ReadDir(Rkt_dir + "/pods/run")
	if err != nil {
		err = errors.New("GetPods(): Failed to get UUIDs: " + err.Error())
		return
	}

	for _, uuid := range uuids {
		var pod Pod
		var annotation_map map[string]string

		annotation_map = make(map[string]string)

		pod.Uuid = uuid.Name()

		if pod.Name, manifest, err = GetPodManifest(pod.Uuid); err != nil {
			return
		}

		pod.Image = manifest.Name

		// Parse consul-specific annotations
		has_port := false
		for _, an := range manifest.Annotations {
			if !strings.HasPrefix(an.Name, "consul-") {
				continue
			}
			if an.Name == "consul-port" {
				has_port = true
			}
			annotation_map[an.Name] = an.Value
		}

		if has_port {
			pod.Consul.Port, err = strconv.Atoi(annotation_map["consul-port"])
			if err != nil {
				err = errors.New("GetPods(): Failed to convert port to int: " + err.Error())
				return
			}

			for key, value := range annotation_map {
				switch key {
				case "consul-dns":
					{
						pod.Consul.Dns = value
						break
					}
				case "consul-check-type":
					{
						pod.Consul.Check.Type = value
						break
					}
				case "consul-check-target":
					{
						pod.Consul.Check.Target = value
						break
					}
				case "consul-check-name":
					{
						pod.Consul.Check.Name = value
						break
					}
				case "consul-check-interval":
					{
						pod.Consul.Check.Interval = value
						break
					}
				case "consul-check-timeout":
					{
						pod.Consul.Check.Timeout = value
						break
					}
				}
			}
		}

		// Add consul tcp check if none is defined
		if pod.Consul.Check.Type == "" {
			pod.Consul.Check.Name = "TCP check on port " + strconv.Itoa(pod.Consul.Port)
			pod.Consul.Check.Type = "tcp"
		}

		// Set a default interval and timeout if not set
		if pod.Consul.Check.Interval == "" {
			pod.Consul.Check.Interval = D_CONSUL_CHECK_INTERVAL
		}
		if pod.Consul.Check.Timeout == "" {
			pod.Consul.Check.Timeout = D_CONSUL_CHECK_TIMEOUT
		}

		pods[uuid.Name()] = pod
	}

	// Parse networks for all pods
	if _, err = os.Stat(Cni_dir + "/networks/" + *netname); err != nil {
		err = errors.New("GetPods(): Failed to get network data: " + err.Error())
		return
	}

	net_dir := Cni_dir + "/networks/" + *netname

	if netdata, err = GetNetworkData(net_dir); err != nil {
		return
	}

	for _, pod := range pods {
		for net_uuid := range netdata {
			if net_uuid != pod.Uuid {
				continue
			}
			pod.IpAddress = netdata[pod.Uuid]
			pods[pod.Uuid] = pod
			break
		}
	}

	return
}
