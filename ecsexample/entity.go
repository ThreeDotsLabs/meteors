package ecsexample

import "fmt"

type Entity struct {
	components map[ComponentType]ComponentTyper
}

// HasComponent returns true if the entity has the component of type cType associated
func (e *Entity) HasComponent(cType ComponentType) bool {
	_, ok := e.components[cType]
	return ok
}

// GetComponent returns the component, panics if not found
func (e *Entity) GetComponent(cType ComponentType) ComponentTyper {
	if c, ok := e.components[cType]; !ok {
		panic(fmt.Sprintf("expected entity to have component of type %s attached", cType))
	} else {
		return c
	}
}

// AddComponent adds a component to the entity, panics if a component of the same type already exists
func (e *Entity) AddComponent(c ComponentTyper) {
	if e.HasComponent(c.Type()) {
		panic(fmt.Sprintf("entity already has component of type %s attached", c))
	}
	e.components[c.Type()] = c
}
