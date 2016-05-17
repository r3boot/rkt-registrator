package main

import (
	"flag"
	"github.com/r3boot/rkt-registrator/consul"
	"github.com/r3boot/rkt-registrator/rkt"
	"github.com/r3boot/rkt-registrator/utils"
	"reflect"
	"time"
)

const D_RKT_NET string = "default-restricted"
const D_RKT_DATA_DIR string = "/var/lib/rkt"
const D_RKT_CNI_DIR string = "/var/lib/cni"
const D_CONSUL_ENDPOINT string = "http://localhost:8500"
const D_DEBUG bool = false

var rkt_network = flag.String("net", D_RKT_NET, "Publish ip addresses from this network")
var rkt_data_dir = flag.String("rkt-data-dir", D_RKT_DATA_DIR, "Path to rkt data directory")
var rkt_cni_dir = flag.String("rkt-cni-dir", D_RKT_CNI_DIR, "Path to rkt cni directory")
var consul_endpoint = flag.String("consul-endpoint", D_CONSUL_ENDPOINT, "Uri of consul master endpoint")
var consul_worker = flag.String("consul-worker", "", "Consul node on which we register services")
var debug = flag.Bool("d", D_DEBUG, "Enable debug output")

var Log utils.Log

func init() {
	var err error
	flag.Parse()

	if *consul_worker == "" {
		Log.Fatal("A consul worker to register services on must be specified")
	}

	Log.UseDebug = *debug
	Log.UseVerbose = *debug
	Log.UseTimestamp = false
	Log.Debug("Logging initialized")

	if err = rkt.Setup(Log, *rkt_data_dir, *rkt_cni_dir); err != nil {
		Log.Fatal("Failed to setup rkt parser: " + err.Error())
	}
	Log.Debug("rkt parser initialized")

	if err = consul.Setup(Log, *consul_endpoint, *consul_worker); err != nil {
		Log.Fatal("Failed to setup consul agent: " + err.Error())
	}
	Log.Debug("consul agent initialized")

}

func main() {
	var err error
	var cur_pods map[string]rkt.Pod
	var prev_pods map[string]rkt.Pod
	var to_add []string
	var to_remove []string

	cur_pods = make(map[string]rkt.Pod)
	prev_pods = make(map[string]rkt.Pod)

	for {
		if cur_pods, err = rkt.GetPods(rkt_network); err != nil {
			Log.Fatal(err)
		}

		if !reflect.DeepEqual(cur_pods, prev_pods) {
			to_add, to_remove = rkt.DiffPods(cur_pods, prev_pods)

			for _, uuid := range to_add {
				Log.Debug("Registering " + cur_pods[uuid].Name + " on " + *consul_worker)
				consul.Register(cur_pods[uuid])
			}
			for _, uuid := range to_remove {
				Log.Debug("Deregistering " + cur_pods[uuid].Name + " on " + *consul_worker)
				consul.Deregister(prev_pods[uuid])
			}

			prev_pods = cur_pods
		}

		time.Sleep(1 * time.Second)
	}
}
