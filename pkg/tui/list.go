package tui

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	//"github.com/charmbracelet/lipgloss"
	"github.com/moby/moby/api/types/container"
)

func (m model) updateList(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Up):
			if m.cursor > 0 {
				m.cursor--
			}
		case key.Matches(msg, keys.Down):
			if m.cursor < len(m.containers)-1 {
				m.cursor++
			}
		case key.Matches(msg, keys.Enter):
			if len(m.containers) > 0 {
				m.view = viewDetail
				return m, m.fetchContainerDetail
			}
		case key.Matches(msg, keys.Refresh):
			return m, m.fetchContainers
		case key.Matches(msg, keys.Stop):
			return m, m.stopContainer
		case key.Matches(msg, keys.Start):
			return m, m.startContainer
		case key.Matches(msg, keys.Restart):
			return m, m.restartContainer
		case key.Matches(msg, keys.Delete):
			return m, m.deleteContainer
		}
	case actionDoneMsg:
		return m, m.fetchContainers
	}
	return m, nil
}

func (m model) viewList() string {
	var b strings.Builder

	// Title
	title := titleStyle.Render("⬡ STACKR")
	count := statusStyle.Render(fmt.Sprintf("  %d containers", len(m.containers)))
	b.WriteString(title + count + "\n\n")

	if len(m.containers) == 0 {
		b.WriteString(statusStyle.Render("  No containers found.\n"))
	} else {
		// Calculate visible area
		visibleLines := m.height - 6 // title + help + margins
		if visibleLines < 5 {
			visibleLines = 5
		}

		// Calculate scroll offset
		offset := 0
		if m.cursor >= visibleLines {
			offset = m.cursor - visibleLines + 1
		}

		// Render visible containers
		end := offset + visibleLines
		if end > len(m.containers) {
			end = len(m.containers)
		}

		for i := offset; i < end; i++ {
			c := m.containers[i]
			line := m.renderLine(c, i == m.cursor)
			b.WriteString(line)
			b.WriteString("\n")
		}

		// Scroll indicator
		if len(m.containers) > visibleLines {
			indicator := statusStyle.Render(fmt.Sprintf("\n  [%d/%d]", m.cursor+1, len(m.containers)))
			b.WriteString(indicator)
		}
	}

	// Help
	b.WriteString("\n\n")
	help := "[↑↓] select  [enter] details  [s]top  [r]esume  [R]estart  [d]elete  [f]refresh  [q]uit"
	b.WriteString(helpStyle.Render(help))

	return b.String()
}

func (m model) renderLine(c container.Summary, selected bool) string {
	// Indicator
	indicator := "  "
	if selected {
		indicator = "▸ "
	}

	// Status dot
	dot := statusDot(string(c.State))

	// Name (max 16 chars)
	name := truncate(containerName(c), 16)

	// Image (max 24 chars)
	image := truncate(c.Image, 24)

	// Status text (max 16 chars)
	status := truncate(shortStatus(c.Status), 16)

	// Ports
	ports := truncate(formatPorts(c.Ports), 16)

	// Build line
	line := fmt.Sprintf("%s%s %-16s  %-24s  %-16s  %s",
		indicator,
		dot,
		name,
		image,
		status,
		ports,
	)

	if selected {
		return selectedStyle.Render(line)
	}

	// Color based on state
	if c.State == "running" {
		return runningStyle.Render(line)
	}
	return stoppedStyle.Render(line)
}

func shortStatus(status string) string {
	// Simplify "Up 3 days" -> "Up 3d", "Exited (0) 2 hours ago" -> "Exited"
	if strings.HasPrefix(status, "Up") {
		parts := strings.Fields(status)
		if len(parts) >= 3 {
			return fmt.Sprintf("Up %s%c", parts[1], parts[2][0])
		}
		return status
	}
	if strings.HasPrefix(status, "Exited") {
		return "Exited"
	}
	if strings.HasPrefix(status, "Created") {
		return "Created"
	}
	return truncate(status, 12)
}

func formatPorts(ports []container.PortSummary) string {
	if len(ports) == 0 {
		return ""
	}
	var parts []string
	seen := make(map[string]bool)
	for _, p := range ports {
		var key string
		if p.PublicPort != 0 {
			key = fmt.Sprintf("%d:%d", p.PublicPort, p.PrivatePort)
		} else {
			key = fmt.Sprintf("%d", p.PrivatePort)
		}
		if !seen[key] {
			parts = append(parts, key)
			seen[key] = true
		}
	}
	return strings.Join(parts, ",")
}

// Actions
func (m model) stopContainer() tea.Msg {
	if len(m.containers) == 0 {
		return nil
	}
	_ = m.client.Stop(context.Background(), m.containers[m.cursor].ID)
	return actionDoneMsg{}
}

func (m model) startContainer() tea.Msg {
	if len(m.containers) == 0 {
		return nil
	}
	_ = m.client.Start(context.Background(), m.containers[m.cursor].ID)
	return actionDoneMsg{}
}

func (m model) restartContainer() tea.Msg {
	if len(m.containers) == 0 {
		return nil
	}
	_ = m.client.Restart(context.Background(), m.containers[m.cursor].ID)
	return actionDoneMsg{}
}

func (m model) deleteContainer() tea.Msg {
	if len(m.containers) == 0 {
		return nil
	}
	_ = m.client.Remove(context.Background(), m.containers[m.cursor].ID)
	return actionDoneMsg{}
}

// Helpers
func containerName(c container.Summary) string {
	if len(c.Names) > 0 {
		return strings.TrimPrefix(c.Names[0], "/")
	}
	return c.ID[:12]
}

func truncate(s string, max int) string {
	if max < 4 {
		max = 4
	}
	if len(s) > max {
		return s[:max-3] + "..."
	}
	return s
}
