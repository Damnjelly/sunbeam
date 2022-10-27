package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/pomdtr/sunbeam/api"
)

type CommandContainer struct {
	width, height int
	currentView   string

	form    *Form
	loading *Loading
	list    *List
	detail  *Detail
	err     *Detail

	command api.SunbeamCommand
	params  map[string]string
}

func NewCommandContainer(command api.SunbeamCommand, params map[string]string) *CommandContainer {
	return &CommandContainer{
		command: command,
		params:  params,
	}
}

func (c *CommandContainer) Init() tea.Cmd {
	missing := c.command.CheckMissingParams(c.params)

	if len(missing) > 1 {
		c.currentView = "form"
		c.form = NewFormContainer(c.command.Title, missing)
		c.form.SetSize(c.width, c.height)
		return c.form.Init()
	}

	return NewSubmitCmd(c.params)
}

type ListOutput []ListItem
type DetailOutput string

func (c *CommandContainer) Run(input api.CommandInput) tea.Cmd {
	return func() tea.Msg {
		switch c.command.Mode {
		case "list":
			scriptItems, err := c.command.List.Run(input)
			if err != nil {
				return err
			}
			listItems := make([]ListItem, len(scriptItems))
			for i, scriptItem := range scriptItems {
				listItems[i] = NewListItem(scriptItem)
			}

			return ListOutput(listItems)
		case "detail":
			detail, err := c.command.Detail.Run(input)
			if err != nil {
				return err
			}
			return DetailOutput(detail)
		default:
			return fmt.Errorf("unknown command mode: %s", c.command.Mode)
		}
	}
}

func (c *CommandContainer) SetSize(width, height int) {
	c.width, c.height = width, height
	if c.loading != nil {
		c.loading.SetSize(width, height)
	}
	if c.form != nil {
		c.form.SetSize(width, height)
	}
	if c.list != nil {
		c.list.SetSize(width, height)
	}
	if c.detail != nil {
		c.detail.SetSize(width, height)
	}
	if c.err != nil {
		c.err.SetSize(width, height)
	}
}

func (c *CommandContainer) Update(msg tea.Msg) (Container, tea.Cmd) {
	switch msg := msg.(type) {
	case SubmitMsg:
		c.currentView = "loading"
		c.loading = NewLoading()
		c.loading.SetSize(c.width, c.height)
		runCmd := c.Run(api.NewCommandInput(msg.values))
		return c, tea.Batch(c.loading.Init(), runCmd)
	case ListOutput:
		if c.list == nil {
			c.currentView = "list"
			if c.command.List.Dynamic {
				c.list = NewDynamicList(msg)
			} else {
				c.list = NewStaticList(msg)
			}
			if c.command.List.ShowDetail {
				c.list.showDetail = true
			}
			c.list.SetSize(c.width, c.height)
			return c, c.list.Init()
		}
		c.list.SetItems(msg)
	case DetailOutput:
		if c.detail == nil {
			c.currentView = "detail"
			var content string
			switch c.command.Detail.Format {
			case "markdown":
				renderer, err := glamour.NewTermRenderer(glamour.WithAutoStyle(), glamour.WithEmoji())
				if err != nil {
					return c, NewErrorCmd("failed to init markdown renderer: %s", err)
				}
				content, err = renderer.Render(string(msg))
				if err != nil {
					return c, NewErrorCmd("failed to render markdown: %s", err)
				}
			default:
				content = string(msg)
			}

			actions := make([]Action, len(c.command.Detail.Actions))
			for i, action := range c.command.Detail.Actions {
				actions[i] = NewAction(action)
			}
			c.detail = NewDetail(content, actions)
			c.detail.SetSize(c.width, c.height)
			return c, c.detail.Init()
		}
		c.detail.SetContent(string(msg))
	case ReloadMsg:
		return c, c.Run(msg.input)
	case QueryUpdateMsg:
		if c.command.List.Dynamic {
			input := api.CommandInput{
				Query:  msg.query,
				Params: c.params,
			}
			return c, c.Run(input)
		}
	case tea.WindowSizeMsg:
		c.SetSize(msg.Width, msg.Height)
		return c, nil
	case error:
		c.currentView = "error"
		c.err = NewDetail(msg.Error(), nil)
		c.err.SetSize(c.width, c.height)
		return c, c.err.Init()
	}

	var cmd tea.Cmd
	switch c.currentView {
	case "form":
		c.form, cmd = c.form.Update(msg)
	case "loading":
		c.loading, cmd = c.loading.Update(msg)
	case "list":
		c.list, cmd = c.list.Update(msg)
	case "detail":
		c.detail, cmd = c.detail.Update(msg)
	case "error":
		c.err, cmd = c.err.Update(msg)
	}
	return c, cmd
}

func (c *CommandContainer) View() string {
	switch c.currentView {
	case "form":
		return c.form.View()
	case "loading":
		return c.loading.View()
	case "list":
		return c.list.View()
	case "detail":
		return c.detail.View()
	case "error":
		return c.err.View()
	default:
		return "Unknown view"
	}
}