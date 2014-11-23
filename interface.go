// Copyright (c) 2014-present Stuart Herbert
// Released under a 3-clause BSD license
package actor

type Actor interface {
	Name() string
	SetName(string)
	IsRunning() bool
	Start()
	Stop()
}
