package consul

type ConsulService struct {
	ID     string            `json:"ID"`
	Name   string            `json:"Name"`
	Tags   []string          `json:"Tags"`
	Ipaddr string            `json:"Address"`
	Port   int               `json:"Port"`
	Check  map[string]string `json:"Check"`
}
