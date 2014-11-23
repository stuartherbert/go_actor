// Copyright (c) 2014-present Stuart Herbert
// Released under a 3-clause BSD license
package actor

import (
	"math/rand"
	"time"
)

// RunnableActor is an actor that runs in a go routine, in order to work
// concurrently with other parts of your app
type RunnableActor struct {
	// we reuse everything from our BaseActor
	BaseActor

	// a channel we monitor to prove that the actor is alive
	pingChan chan int

	// a channel the EventLoop monitors for when the actor needs to be stopped
	stopChan chan struct{}

	// a channel that the EventLoop needs to close when it exits
	stoppedChan chan struct{}
}

// PingChan() returns the channel to send ping messages to
func (self *RunnableActor) PingChan() chan int {
	return self.pingChan
}

// initPingChan() creates the channel we use to make sure this actor's
// EventLoop is still functional
func (self *RunnableActor) initPingChan() {
	self.mu.Lock()
	defer self.mu.Unlock()

	self.pingChan = make(chan int)
}

// StopChan() returns the channel to close to stop this actor
func (self *RunnableActor) StopChan() chan struct{} {
	return self.stopChan
}

// initStopChan() creates the channel we monitor for when it is time to
// shutdown this actor
func (self *RunnableActor) initStopChan() {
	self.mu.Lock()
	defer self.mu.Unlock()

	self.stopChan = make(chan struct{})
}

// initStoppedChan() creates the channel that the EventLoop must close when it
// exits
func (self *RunnableActor) initStoppedChan() {
	self.mu.Lock()
	defer self.mu.Unlock()

	self.stoppedChan = make(chan struct{})
}

// CloseStoppedChan() is called by the EventLoop when it exits
func (self *RunnableActor) CloseStoppedChan() {
	close(self.stoppedChan)
}

// Start() kicks off the EventLoop for this actor
func (self *RunnableActor) startEventLoop(eventLoop func(chan struct{})) {
	// are we already running?
	if self.IsRunning() {
		return
	}

	// initialise our internal channels
	self.initPingChan()
	self.initStopChan()
	self.initStoppedChan()

	// we use this to make sure the actor has started
	started := make(chan struct{})
	go eventLoop(started)
	<-started

	// we're up and running
	self.setIsRunning(true)
}

// Stop() shuts down this actor
func (self *RunnableActor) Stop() {
	// have we already stopped?
	if !self.isRunning {
		return
	}

	// mark us as not running
	//
	// this is to prevent any further requests into the actor
	self.setIsRunning(false)

	// close the stopChan
	go func() {
		close(self.stopChan)
	}()

	// wait for the EventLoop to have finished
	<-self.stoppedChan

	// all done
}

// Ping() helps us prove that the EventLoop is working
func (self *RunnableActor) Ping() bool {
	// if we're not running, nothing else to do
	if !self.IsRunning() {
		return false
	}

	// write a value to the ping channel
	expected := rand.Int()
	pingChan := self.PingChan()
	pingChan <- expected

	select {
	// did we get a response?
	case actual := <-pingChan:
		// did we get the response we expected?
		if actual == expected {
			return true
		} else {
			return false
		}
	// use a timeout to avoid blocking forever
	case <-time.After(time.Second):
		return false
	}
}
