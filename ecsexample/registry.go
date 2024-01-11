package ecsexample

type Registry struct {
	entities []*Entity
}

// NewEntity creates a new entity and adds it to the internal list of entities
func (r *Registry) NewEntity() *Entity {
	e := Entity{components: make(map[ComponentType]ComponentTyper)}
	r.entities = append(r.entities, &e)
	return &e
}

// Query returns all entities that have all types associated
func (r *Registry) Query(types ...ComponentType) []*Entity {
	candidates := []*Entity{}

	for _, e := range r.entities {
		matchCount := 0
		for _, c := range types {
			if e.HasComponent(c) {
				matchCount++
			}
		}

		if matchCount == len(types) {
			candidates = append(candidates, e)
		}
	}
	return candidates
}
