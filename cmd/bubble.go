package main

import (
	//"context"
	"fmt"
	//"encoding/json"

	//"github.com/moby/moby/client"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	count int
}

func newModel() model {
	return model{}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			m.count++
		}
	}
	return m, nil
}

func (m model) View() string {
	return fmt.Sprintf("Count: %d\n", m.count)
}

func bubbleCLI() {

	m := newModel()
	tea.NewProgram(m).Run()
	/*
		ctx := context.Background()
		apiClient, err := client.New(client.FromEnv)
		if err != nil {
			panic(err)
		}
		defer apiClient.Close()

		containers, err := apiClient.ContainerList(ctx, client.ContainerListOptions{})
		if err != nil {
			panic(err)
		}

		//for _, c := range containers.Items {
		//}
		data, _ := json.MarshalIndent(containers.Items[0], "", "  ")
		fmt.Println(string(data))
	*/
}
