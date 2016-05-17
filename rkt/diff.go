package rkt

func DiffPods(cur map[string]Pod, prev map[string]Pod) (add []string, remove []string) {
	// Check for pods to add to service discovery
	for uuid, _ := range cur {
		if _, present := prev[uuid]; !present {
			add = append(add, uuid)
		}
	}

	// Check for pods to remove from service discovery
	for uuid, _ := range prev {
		if _, present := cur[uuid]; !present {
			remove = append(remove, uuid)
		}
	}
	return
}
