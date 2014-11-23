// Copyright (c) 2014-present Stuart Herbert
// Released under a 3-clause BSD license
package actor

import (
	"sync"
)

// BaseActor is the building block to create different types of actor from
type BaseActor struct {
	// the name of this actor
	name string

	// we use this to help keep track of whether we are running or not
	isRunning bool

	// we use this to avoid race conditions when modifying the actor
	mu sync.RWMutex
}

// Name() returns the name of this actor
func (self *BaseActor) Name() string {
	return self.name
}

// SetName() sets the name of this actor
func (self *BaseActor) SetName(newName string) {
	self.name = newName
}

// IsRunning() tells us if it is safe to write to the pingChan et al or not
func (self *BaseActor) IsRunning() bool {
	self.mu.RLock()
	defer self.mu.RUnlock()

	return self.isRunning
}

// setIsRunning() sets whether the channels are safe to use or not
func (self *BaseActor) setIsRunning(isRunning bool) {
	self.mu.Lock()
	defer self.mu.Unlock()

	self.isRunning = isRunning
}
