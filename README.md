### Wth is rkt-registrator?
Rkt-registrator is a bridge between [rkt containers](https://coreos.com/rkt/) and [Consul](https://www.consul.io/). It reads metadata from pods running on a host, and pushes these to consul. Pod metadata is read from the local filesystem, and converted into a json object which is posted to the consul api. Consul checks can be configured on a per-pod basis using annotations added to the ACI (see below).

### rkt-registrator commandline options
rkt-registrator is fully configured from the commandline. See the options below for an overview:
```
  -consul-endpoint string
    	Uri of consul master endpoint (default "http://localhost:8500")
  -consul-worker string
    	Consul node on which we register services
  -d	Enable debug output
  -net string
    	Publish ip addresses from this network (default "default-restricted")
  -rkt-cni-dir string
    	Path to rkt cni directory (default "/var/lib/cni")
  -rkt-data-dir string
    	Path to rkt data directory (default "/var/lib/rkt")

```

### ACI annotations for consul checks
You can add the annotations below to your ACI's to influence the way consul handles the service. Note that for now your ACI must atleast have consul-port to be usable with rkt-registrator.

 Name | Description
  --- | ---
 consul-port | Port on which the service is listening
 consul-dns | Override DNS hostname
 consul-check-type | Check type (See [Consul Service Definition docs](https://www.consul.io/docs/agent/services.html))
 consul-check-target | Target of the check
 consul-check-name | Name for the check
 consul-check-interval | Interval of check (defaults to 10s)
 consul-check-timeout | Timeout of check (defaults to 1s)

### Roadmap
The current code is not even alpha quality code. Further work include more documentation and building proper releases.
