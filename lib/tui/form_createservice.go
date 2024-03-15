package tui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/viper"
	"strings"
)

type ServiceForm struct {
	Model
	compose string

	ContainerName string
	Image         string
	Env           string
	Expose        string
}

func NewServiceForm() ServiceForm {
	m := ServiceForm{
		Model: NewModel(0),
	}
	m.lg = lipgloss.DefaultRenderer()
	m.styles = NewStyles(m.lg)

	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("container_name").
				Title("Service name").
				Placeholder("api").
				Value(&m.ContainerName).
				Validate(
					VRequiredAndSave("container_name", "Service name is required"),
				),

			huh.NewInput().
				Key("image").
				Title("Image name").
				Placeholder("backend:latest").
				Value(&m.Image).
				Validate(
					VRequiredAndSave("image", "Image name is required"),
				),

			huh.NewText().
				Key("env").
				Title("Environment").
				Placeholder(`PORT=8080`).
				Editor("nano").
				Value(&m.Env).
				Validate(
					VSave("env"),
				),

			huh.NewText().
				Key("expose").
				Title("Expose").
				Placeholder(`8080:api.envme.bid`).
				Editor("nano").
				Value(&m.Expose).
				Validate(
					VSave("expose"),
				).
				Lines(2),

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

func (m ServiceForm) View() string {
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
				header    = s.Help.Render("docker-compose.yaml")
				service   = s.Help.Render("(none)")
				image     = s.Help.Render("(none)")
				additions = ""
				container = ""
			)

			containerName := m.form.GetString("container_name")
			if containerName != "" {
				header = s.Help.Render(containerName + "/" + header)
				service = s.Highlight.Render(containerName) + ":"
			} else {
				image = ""
			}

			imageName := m.form.GetString("image")
			envValue := m.form.GetString("env")
			if imageName != "" {
				image = "image: " + s.Highlight.Render(imageName)
				container = "container_name: " + s.Highlight.Render(containerName)
				additions += t + t + "restart: unless-stopped" + n
				if envValue != "" {
					additions += t + t + "environment:" + n
					envs := strings.Split(envValue, n)
					for _, env := range envs {
						additions += t + t + t + "- " + s.Highlight.Render(env) + n
					}
				}
				additions += t + t + "networks:" + n
				additions += t + t + t + "- envme" + n
				additions += t + t + "extra_hosts:" + n
				additions += t + t + t + "- host.docker.internal:host-gateway" + n
				additions += n
				additions += "networks:" + n
				additions += t + "envme:" + n
				additions += t + t + "external: true" + n
			}

			m.compose += header + n
			m.compose += n
			m.compose += "version: '3.8'" + n
			m.compose += n
			m.compose += "services:" + n
			m.compose += t + service + n
			m.compose += t + t + image + n
			m.compose += t + t + container + n
			m.compose += additions

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

func (m ServiceForm) Init() tea.Cmd {
	return m.form.Init()
}

func (m ServiceForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
