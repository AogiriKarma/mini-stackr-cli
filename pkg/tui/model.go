package tui

import (
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/moby/moby/api/types/container"
	"github.com/aogirikarma/mini-stackr-cli/pkg/docker"
)

type viewState int

const (
	viewList viewState = iota
	viewDetail
)

type model struct {
	client     *docker.Client
	view       viewState
	containers []container.Summary
	cursor     int
	width      int
	height     int
	err        error

	// Detail view data
	inspect  *container.InspectResponse
	stats    *container.StatsResponse
	viewport viewport.Model
}

// Messages
type containersMsg []container.Summary
type inspectMsg struct {
	inspect *container.InspectResponse
	stats   *container.StatsResponse
}
type errMsg error
type actionDoneMsg struct{}