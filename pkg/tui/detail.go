package tui

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	//"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/moby/moby/api/types/container"
)

func (m model) updateDetail(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Back):
			m.view = viewList
			m.inspect = nil
			m.stats = nil
			return m, nil
		case key.Matches(msg, keys.Stop):
			return m, tea.Sequence(m.stopContainer, m.fetchContainerDetail)
		case key.Matches(msg, keys.Start):
			return m, tea.Sequence(m.startContainer, m.fetchContainerDetail)
		case key.Matches(msg, keys.Restart):
			return m, tea.Sequence(m.restartContainer, m.fetchContainerDetail)
		case key.Matches(msg, keys.Delete):
			m.view = viewList
			return m, m.deleteContainer
		case key.Matches(msg, keys.Refresh):
			return m, m.fetchContainerDetail
		}
	case inspectMsg:
		m.inspect = msg.inspect
		m.stats = msg.stats
		// Update viewport content
		content := m.renderDetailContent()
		m.viewport.SetContent(content)
		return m, nil
	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - 4 // leave room for help
		return m, nil
	}

	// Handle viewport scrolling
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m model) fetchContainerDetail() tea.Msg {
	if len(m.containers) == 0 || m.cursor >= len(m.containers) {
		return nil
	}
	id := m.containers[m.cursor].ID

	inspect, err := m.client.Inspect(context.Background(), id)
	if err != nil {
		return errMsg(err)
	}

	stats, _ := m.client.Stats(context.Background(), id)

	return inspectMsg{inspect: &inspect, stats: stats}
}

func (m model) viewDetail() string {
	if m.inspect == nil {
		return "Loading..."
	}

	var b strings.Builder

	// Viewport content
	b.WriteString(m.viewport.View())
	b.WriteString("\n")

	// Scroll indicator
	scrollPercent := int(m.viewport.ScrollPercent() * 100)
	scrollInfo := statusStyle.Render(fmt.Sprintf("  [%d%%]", scrollPercent))

	// Help
	help := "[↑↓] scroll  [s]top  [r]esume  [R]estart  [d]elete  [f]refresh  [esc]back  [q]uit"
	b.WriteString(helpStyle.Render(help) + scrollInfo)

	return b.String()
}

func (m model) renderDetailContent() string {
	if m.inspect == nil {
		return ""
	}

	ins := m.inspect
	var b strings.Builder

	// Header
	dot := statusDot(string(ins.State.Status))
	state := strings.ToUpper(string(ins.State.Status))
	name := strings.TrimPrefix(ins.Name, "/")

	header := fmt.Sprintf("⬡  %s", titleStyle.Render(name))
	statusText := fmt.Sprintf("%s %s", dot, state)

	b.WriteString(header + "  " + statusText)
	b.WriteString("\n")
	b.WriteString(statusStyle.Render(strings.Repeat("─", 40)))
	b.WriteString("\n\n")

	// Calculer la largeur disponible
	availableWidth := m.width - 2 // marge
	if availableWidth < 40 {
		availableWidth = 40
	}

	// Render les boxes pour mesurer leur taille
	boxWidth := (availableWidth - 3) / 2 // -3 pour le gap
	if boxWidth < 30 {
		boxWidth = 30
	}

	resourcesBox := m.renderResourcesBox(boxWidth)
	networkBox := m.renderNetworkBox(boxWidth)

	// Mesurer la largeur réelle
	//totalWidth := lipgloss.Width(resourcesBox) + lipgloss.Width(networkBox) + 2

	// Si ça rentre pas, stack vertical
	if  m.width < 80 {
		// Vertical layout - full width
		fullWidth := availableWidth
		b.WriteString(m.renderResourcesBox(fullWidth))
		b.WriteString("\n\n")
		b.WriteString(m.renderNetworkBox(fullWidth))
		b.WriteString("\n\n")
		b.WriteString(m.renderInfoBox(fullWidth))
		b.WriteString("\n\n")

		if ins.Mounts != nil && len(ins.Mounts) > 0 {
			b.WriteString(m.renderMountsBox(fullWidth))
			b.WriteString("\n\n")
		}

		if ins.Config != nil && len(ins.Config.Env) > 0 {
			b.WriteString(m.renderEnvBox(fullWidth))
			b.WriteString("\n\n")
		}

		if ins.Config != nil && len(ins.Config.Labels) > 0 {
			b.WriteString(m.renderLabelsBox(fullWidth))
			b.WriteString("\n\n")
		}
	} else {
		// Horizontal layout
		row1 := lipgloss.JoinHorizontal(lipgloss.Top, resourcesBox, "  ", networkBox)
		b.WriteString(row1)
		b.WriteString("\n\n")

		fullWidth := availableWidth
		b.WriteString(m.renderInfoBox(fullWidth))
		b.WriteString("\n\n")

		if ins.Mounts != nil && len(ins.Mounts) > 0 {
			b.WriteString(m.renderMountsBox(fullWidth))
			b.WriteString("\n\n")
		}

		if ins.Config != nil && len(ins.Config.Env) > 0 {
			b.WriteString(m.renderEnvBox(fullWidth))
			b.WriteString("\n\n")
		}

		if ins.Config != nil && len(ins.Config.Labels) > 0 {
			b.WriteString(m.renderLabelsBox(fullWidth))
			b.WriteString("\n\n")
		}
	}

	return b.String()
}

func (m model) renderResourcesBox(width int) string {
	var content strings.Builder
	content.WriteString(boxTitleStyle.Render("RESOURCES"))
	content.WriteString("\n\n")

	cpuPercent := 0.0
	memUsage := uint64(0)
	memLimit := uint64(0)
	memPercent := 0.0

	if m.stats != nil {
		cpuPercent = calculateCPUPercent(m.stats)
		memUsage = m.stats.MemoryStats.Usage
		memLimit = m.stats.MemoryStats.Limit
		if memLimit > 0 {
			memPercent = float64(memUsage) / float64(memLimit) * 100
		}
	}

	barWidth := width - 20
	if barWidth < 10 {
		barWidth = 10
	}

	cpuBar := renderProgressBar(cpuPercent, barWidth)
	content.WriteString(fmt.Sprintf("%s  %s  %5.1f%%\n", labelStyle.Render("CPU"), cpuBar, cpuPercent))

	memBar := renderProgressBar(memPercent, barWidth)
	content.WriteString(fmt.Sprintf("%s  %s  %s\n", labelStyle.Render("RAM"), memBar, formatBytes(memUsage)))

	content.WriteString(fmt.Sprintf("%s  %s\n", labelStyle.Render("Limit"), valueStyle.Render(formatBytes(memLimit))))

	if m.stats != nil {
		content.WriteString(fmt.Sprintf("%s  %s", labelStyle.Render("PIDs"), valueStyle.Render(fmt.Sprintf("%d", m.stats.PidsStats.Current))))
	}

	return boxStyle.Width(width).Render(content.String())
}

func (m model) renderNetworkBox(width int) string {
	var content strings.Builder
	content.WriteString(boxTitleStyle.Render("NETWORK"))
	content.WriteString("\n\n")

	ins := m.inspect

	if ins.NetworkSettings != nil && len(ins.NetworkSettings.Networks) > 0 {
		for netName, net := range ins.NetworkSettings.Networks {
			content.WriteString(fmt.Sprintf("%s  %s\n", labelStyle.Render("Network"), valueStyle.Render(netName)))
			content.WriteString(fmt.Sprintf("%s  %s\n", labelStyle.Render("IP"), valueStyle.Render(net.IPAddress.String())))
			content.WriteString(fmt.Sprintf("%s  %s\n", labelStyle.Render("Gateway"), valueStyle.Render(net.Gateway.String())))
			break
		}
	}

	if ins.NetworkSettings != nil && len(ins.NetworkSettings.Ports) > 0 {
		var ports []string
		for port, bindings := range ins.NetworkSettings.Ports {
			for _, b := range bindings {
				ports = append(ports, fmt.Sprintf("%s:%s", b.HostPort, port))
			}
		}
		if len(ports) > 0 {
			content.WriteString(fmt.Sprintf("%s  %s", labelStyle.Render("Ports"), valueStyle.Render(strings.Join(ports, ", "))))
		}
	}

	return boxStyle.Width(width).Render(content.String())
}

func (m model) renderInfoBox(width int) string {
	var content strings.Builder
	content.WriteString(boxTitleStyle.Render("CONTAINER INFO"))
	content.WriteString("\n\n")

	ins := m.inspect

	content.WriteString(fmt.Sprintf("%-12s  %s\n", labelStyle.Render("ID"), valueStyle.Render(ins.ID[:12])))
	content.WriteString(fmt.Sprintf("%-12s  %s\n", labelStyle.Render("Image"), valueStyle.Render(ins.Config.Image)))
	content.WriteString(fmt.Sprintf("%-12s  %s\n", labelStyle.Render("Created"), valueStyle.Render(formatTime(ins.Created))))
	content.WriteString(fmt.Sprintf("%-12s  %s\n", labelStyle.Render("Started"), valueStyle.Render(ins.State.StartedAt)))

	if ins.HostConfig != nil {
		restart := fmt.Sprintf("%s (count: %d)", ins.HostConfig.RestartPolicy.Name, ins.RestartCount)
		content.WriteString(fmt.Sprintf("%-12s  %s\n", labelStyle.Render("Restart"), valueStyle.Render(restart)))
	}

	content.WriteString(fmt.Sprintf("%-12s  %s", labelStyle.Render("Platform"), valueStyle.Render(ins.Platform)))

	return boxStyle.Width(width).Render(content.String())
}

func (m model) renderMountsBox(width int) string {
	var content strings.Builder
	content.WriteString(boxTitleStyle.Render("MOUNTS"))
	content.WriteString("\n\n")

	for _, mount := range m.inspect.Mounts {
		src := truncate(mount.Source, 30)
		dst := truncate(mount.Destination, 30)
		mode := "rw"
		if !mount.RW {
			mode = "ro"
		}
		line := fmt.Sprintf("%-6s  %s  →  %s  [%s]", string(mount.Type), src, dst, mode)
		content.WriteString(valueStyle.Render(line))
		content.WriteString("\n")
	}

	return boxStyle.Width(width).Render(strings.TrimRight(content.String(), "\n"))
}

func (m model) renderEnvBox(width int) string {
	var content strings.Builder
	content.WriteString(boxTitleStyle.Render("ENVIRONMENT"))
	content.WriteString("\n\n")

	for _, env := range m.inspect.Config.Env {
		content.WriteString(valueStyle.Render(truncate(env, width-4)))
		content.WriteString("\n")
	}

	return boxStyle.Width(width).Render(strings.TrimRight(content.String(), "\n"))
}

func (m model) renderLabelsBox(width int) string {
	var content strings.Builder
	content.WriteString(boxTitleStyle.Render("LABELS"))
	content.WriteString("\n\n")

	for k, v := range m.inspect.Config.Labels {
		line := fmt.Sprintf("%s=%s", k, v)
		content.WriteString(valueStyle.Render(truncate(line, width-4)))
		content.WriteString("\n")
	}

	return boxStyle.Width(width).Render(strings.TrimRight(content.String(), "\n"))
}

// Helpers
func calculateCPUPercent(stats *container.StatsResponse) float64 {
	cpuDelta := float64(stats.CPUStats.CPUUsage.TotalUsage - stats.PreCPUStats.CPUUsage.TotalUsage)
	systemDelta := float64(stats.CPUStats.SystemUsage - stats.PreCPUStats.SystemUsage)

	if systemDelta > 0 && cpuDelta > 0 {
		return (cpuDelta / systemDelta) * float64(stats.CPUStats.OnlineCPUs) * 100.0
	}
	return 0
}

func formatBytes(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := uint64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}

func formatTime(t string) string {
	parsed, err := time.Parse(time.RFC3339Nano, t)
	if err != nil {
		return t
	}
	return parsed.Format("2006-01-02 15:04:05")
}