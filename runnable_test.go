package actor

import (
	"github.com/bmizerany/assert"
	"testing"
)

type RunnableTestActor struct {
	RunnableActor
}

func (self *RunnableTestActor) EventLoop(started chan struct{}) {
	defer self.CloseStoppedChan()

	// the internal channels that we need to monitor
	pingChan := self.PingChan()
	stopChan := self.StopChan()

	// we're up and running
	close(started)

	// monitor all the things
	select {
	case pingReq := <-pingChan:
		pingChan <- pingReq
	case <-stopChan:
		return
	}
}

func (self *RunnableTestActor) Start() {
	self.startEventLoop(self.EventLoop)
}

func TestCanCreateRunnableActor(t *testing.T) {
	_ = RunnableActor{}
}

func TestNewRunnableActorIsNotRunning(t *testing.T) {
	a := RunnableActor{}

	assert.Equal(t, false, a.IsRunning())
}

func TestNewRunnableActorDoesNotPing(t *testing.T) {
	a := RunnableActor{}

	assert.Equal(t, false, a.Ping())
}

func TestCanStartRunnableActor(t *testing.T) {
	a := RunnableTestActor{}
	assert.Equal(t, false, a.IsRunning())

	a.Start()
	defer a.Stop()

	assert.Equal(t, true, a.IsRunning())
}

func TestCanPingStartedRunnableActor(t *testing.T) {
	a := RunnableTestActor{}
	assert.Equal(t, false, a.Ping())

	a.Start()
	defer a.Stop()

	assert.Equal(t, true, a.Ping())
}

func TestCanStopStartedRunnableActor(t *testing.T) {
	a := RunnableTestActor{}

	a.Start()
	assert.Equal(t, true, a.Ping())
	assert.Equal(t, true, a.IsRunning())

	a.Stop()
	assert.Equal(t, false, a.Ping())
	assert.Equal(t, false, a.IsRunning())
}
