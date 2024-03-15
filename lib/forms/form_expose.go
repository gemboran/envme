package forms

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/viper"
	"strings"
)

type ExposeForm struct {
	Model
	compose string

	ContainerName string
	Port          string
	Hostname      string
}

func NewExposeForm() ExposeForm {
	m := ExposeForm{
		Model: NewModel(0),
	}
	m.lg = lipgloss.DefaultRenderer()
	m.styles = NewStyles(m.lg)

	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Key("container_name").
				Title("Service name").
				Value(&m.ContainerName).
				Validate(
					VSave("container_name"),
				),

			huh.NewInput().
				Key("port").
				Title("Port to expose").
				Placeholder("8080").
				Value(&m.Port).
				Validate(
					VRequiredAndSave("port", "Port is required"),
				),

			huh.NewInput().
				Key("hostname").
				Title("Access from").
				Placeholder("api-local.envme.bid").
				Value(&m.Hostname).
				Validate(
					VRequiredAndSave("port", "Port is required"),
				),

			huh.NewConfirm().
				Key("done").
				Title("All done?").
				Validate(func(v bool) error {
					if !v {
						return fmt.Errorf("Welp, finish up then")
					}
					return nil
				}).
				Affirmative("Yep").
				Negative("Wait, no"),
		),
	).
		WithWidth(55).
		WithShowHelp(false).
		WithShowErrors(false)
	return m
}

func (m ExposeForm) View() string {
	s := m.styles

	switch m.form.State {
	case huh.StateCompleted:
		return s.Status.Copy().Margin(0, 1).Width(50).Render(viper.GetString("compose")) + "\n\n"
	default:
		// Form (left side)
		v := strings.TrimSuffix(m.form.View(), "\n\n")
		form := m.lg.NewStyle().Margin(1, 0).Render(v)

		// Status (right side)
		var status string
		{
			const (
				t = "  " // tab
				n = "\n" // end of line
			)
			var (
				header = s.Help.Render("config.yaml")
			)

			containerName := m.form.GetString("container_name")
			port := m.form.GetString("port")
			hostname := m.form.GetString("hostname")

			m.compose += header + n
			m.compose += n
			m.compose += "tunnel: localstack" + n
			m.compose += "credentials-file: credentials.json" + n
			m.compose += "ingress:" + n
			m.compose += t + "hostname: " + hostname + n
			m.compose += t + "service: http://" + containerName + ":" + port + n
			m.compose += t + "service: http_status:404"

			viper.Set("compose", m.compose)

			const statusWidth = 50
			statusMarginLeft := m.width - statusWidth - lipgloss.Width(form) - s.Status.GetMarginRight()
			status = s.Status.Copy().
				Height(lipgloss.Height(form)).
				Width(statusWidth).
				MarginLeft(statusMarginLeft).
				Render(m.compose)
		}

		errors := m.form.Errors()
		header := m.appBoundaryView("Create Service Form")
		if len(errors) > 0 {
			header = m.appErrorBoundaryView(m.errorView())
		}
		body := lipgloss.JoinHorizontal(lipgloss.Top, form, status)

		footer := m.appBoundaryView(m.form.Help().ShortHelpView(m.form.KeyBinds()))
		if len(errors) > 0 {
			footer = m.appErrorBoundaryView("")
		}

		return s.Base.Render(header + "\n" + body + "\n\n" + footer)
	}
}

func (m ExposeForm) Init() tea.Cmd {
	return m.form.Init()
}

func (m ExposeForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = min(msg.Width, maxWidth) - m.styles.Base.GetHorizontalFrameSize()
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+c":
			return m, tea.Quit
		}
	}

	var commands []tea.Cmd

	// Process the form
	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
		commands = append(commands, cmd)
	}

	if m.form.State == huh.StateCompleted {
		// Quit when the form is done.
		commands = append(commands, tea.Quit)
	}

	return m, tea.Batch(commands...)
}
