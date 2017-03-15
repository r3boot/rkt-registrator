package main

import (
	"flag"
	"github.com/r3boot/rkt-registrator/consul"
	"github.com/r3boot/rkt-registrator/rkt"
	"github.com/r3boot/rkt-registrator/utils"
	"os"
	"reflect"
	"time"
)

const (
	D_RKT_NET         string = "default-restricted"
	D_RKT_DATA_DIR    string = "/var/lib/rkt"
	D_RKT_CNI_DIR     string = "/var/lib/cni"
	D_CONSUL_ENDPOINT string = "http://localhost:8500"
	D_DEBUG           bool   = false
)

var (
	f_rkt_network     = flag.String("net", D_RKT_NET, "Publish ip addresses from this network")
	f_rkt_data_dir    = flag.String("rkt-data-dir", D_RKT_DATA_DIR, "Path to rkt data directory")
	f_rkt_cni_dir     = flag.String("rkt-cni-dir", D_RKT_CNI_DIR, "Path to rkt cni directory")
	f_consul_endpoint = flag.String("consul-endpoint", D_CONSUL_ENDPOINT, "Uri of consul master endpoint")
	f_consul_worker   = flag.String("consul-worker", "", "Consul node on which we register services")
	f_debug           = flag.Bool("d", D_DEBUG, "Enable debug output")
	rkt_network       string
	rkt_data_dir      string
	rkt_cni_dir       string
	consul_endpoint   string
	consul_worker     string
	debug             bool
	Log               utils.Log
)

func parseOptions() {
	var (
		value string
	)

	rkt_network = *f_rkt_network
	if *f_rkt_network == D_RKT_NET {
		if value = os.Getenv("RKT_NETWORK"); value != "" {
			rkt_network = value
		}
	}

	rkt_data_dir = *f_rkt_data_dir
	if *f_rkt_data_dir == D_RKT_DATA_DIR {
		if value = os.Getenv("RKT_DATA_DIR"); value != "" {
			rkt_data_dir = value
		}
	}

	rkt_cni_dir = *f_rkt_cni_dir
	if *f_rkt_cni_dir == D_RKT_CNI_DIR {
		if value = os.Getenv("RKT_CNI_DIR"); value != "" {
			rkt_cni_dir = value
		}
	}

	consul_endpoint = *f_consul_endpoint
	if *f_consul_endpoint == D_CONSUL_ENDPOINT {
		if value = os.Getenv("CONSUL_ENDPOINT"); value != "" {
			consul_endpoint = value
		}
	}

	consul_worker = *f_consul_worker
	if *f_consul_worker == "" {
		if value = os.Getenv("CONSUL_WORKER"); value != "" {
			rkt_cni_dir = value
		}
	}
	debug = *f_debug
	if *f_debug == D_DEBUG {
		if value = os.Getenv("REGISTRATOR_DEBUG"); value != "" {
			debug = true
		}
	}

}

func init() {
	var (
		err error
	)

	flag.Parse()

	parseOptions()

	Log.UseDebug = debug
	Log.UseVerbose = debug
	Log.UseTimestamp = false
	Log.Debug("Logging initialized")

	if err = rkt.Setup(Log, rkt_data_dir, rkt_cni_dir); err != nil {
		Log.Fatal("Failed to setup rkt parser: " + err.Error())
	}
	Log.Debug("rkt parser initialized")

	if err = consul.Setup(Log, consul_endpoint, consul_worker); err != nil {
		Log.Fatal("Failed to setup consul agent: " + err.Error())
	}
	Log.Debug("consul agent initialized")

}

func main() {
	var (
		err       error
		cur_pods  map[string]rkt.Pod
		prev_pods map[string]rkt.Pod
		to_add    []string
		to_remove []string
	)

	cur_pods = make(map[string]rkt.Pod)
	prev_pods = make(map[string]rkt.Pod)

	for {
		if cur_pods, err = rkt.GetPods(rkt_network); err != nil {
			Log.Fatal(err)
		}

		if !reflect.DeepEqual(cur_pods, prev_pods) {
			to_add, to_remove = rkt.DiffPods(cur_pods, prev_pods)

			for _, uuid := range to_add {
				Log.Debug("Registering " + cur_pods[uuid].Name + " on " + consul_worker)
				consul.Register(cur_pods[uuid])
			}
			for _, uuid := range to_remove {
				Log.Debug("Deregistering " + cur_pods[uuid].Name + " on " + consul_worker)
				consul.DeRegister(prev_pods[uuid])
			}

			prev_pods = cur_pods
		}

		if err = consul.FlushDuplicates(); err != nil {
			Log.Fatal(err)
			os.Exit(1)
		}

		time.Sleep(1 * time.Second)
	}
}
