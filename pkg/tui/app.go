package tui

import (
	"context"

	"github.com/aogirikarma/mini-stackr-cli/pkg/docker"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

func Run(client *docker.Client) error {
	m := newModel(client)
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}

func newModel(client *docker.Client) model {
	vp := viewport.New(80, 20)
	vp.SetContent("")

	return model{
		client:   client,
		view:     viewList,
		viewport: vp,
	}
}

func (m model) Init() tea.Cmd {
	return m.fetchContainers
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, keys.Quit) {
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - 4
		if m.view == viewDetail && m.inspect != nil {
			m.viewport.SetContent(m.renderDetailContent())
		}
		return m, nil

	case containersMsg:
		m.containers = msg
		return m, nil

	case errMsg:
		m.err = msg
		return m, nil
	}

	switch m.view {
	case viewList:
		return m.updateList(msg)
	case viewDetail:
		return m.updateDetail(msg)
	}

	return m, nil
}

func (m model) View() string {
	if m.err != nil {
		return "Error: " + m.err.Error() + "\n\nPress q to quit."
	}

	switch m.view {
	case viewList:
		return m.viewList()
	case viewDetail:
		return m.viewDetail()
	}

	return ""
}

func (m model) fetchContainers() tea.Msg {
	containers, err := m.client.ListContainers(context.Background())
	if err != nil {
		return errMsg(err)
	}
	return containersMsg(containers)
}
