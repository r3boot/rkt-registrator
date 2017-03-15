package consul

type ConsulService struct {
	ID     string            `json:"ID"`
	Name   string            `json:"Name"`
	Tags   []string          `json:"Tags"`
	Ipaddr string            `json:"Address"`
	Port   int               `json:"Port"`
	Check  map[string]string `json:"Check"`
}

type CatalogServices map[string][]string

type CatalogService struct {
	ID          string `json:"ID"`
	Node        string `json:"Node"`
	ModifyIndex int    `json:"ModifyIndex"`
}
