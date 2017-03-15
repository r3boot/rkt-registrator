package rkt

type KeyValue struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type MountPoint struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type Port struct {
	Count           int    `json:"count"`
	Name            string `json:"name"`
	Port            int    `json:"port"`
	Protocol        string `json:"protocol"`
	SocketActivated bool   `json:"socketActivated"`
}

type Application struct {
	Exec        []string     `json:"exec"`
	Group       string       `json:"group"`
	MountPoints []MountPoint `json:"mountPoints"`
	Ports       []Port       `json:"ports"`
	User        string       `json:"user"`
}

type RktManifest struct {
	AcKind      string      `json:"acKind"`
	AcVersion   string      `json:"acVersion"`
	Annotations []KeyValue  `json:"annotations"`
	App         Application `json:"app"`
	Labels      []KeyValue  `json:"labels"`
	Name        string      `json:"name"`
}

type ConsulCheckSettings struct {
	Name     string
	Type     string
	Target   string
	Interval string
	Timeout  string
}

type ConsulSettings struct {
	Port  int
	Dns   string
	Check ConsulCheckSettings
}

type Pod struct {
	Uuid      string
	Name      string
	Image     string
	IpAddress string
	Consul    ConsulSettings
}
