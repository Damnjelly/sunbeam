package pages

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jinzhu/copier"
	"github.com/pomdtr/sunbeam/bubbles"
	"github.com/pomdtr/sunbeam/bubbles/list"
	commands "github.com/pomdtr/sunbeam/commands"
)

type ListContainer struct {
	list      *list.Model
	textInput *textinput.Model
	width     int
	height    int
	response  *commands.ListResponse
	runAction ActionRunner
}

var listContainer = list.New([]list.Item{}, NewItemDelegate(), 0, 0)

func NewListContainer(res *commands.ListResponse, runAction ActionRunner) Page {
	var l list.Model
	_ = copier.Copy(&l, &listContainer)

	textInput := textinput.NewModel()
	textInput.Prompt = ""
	textInput.Placeholder = "Search..."
	textInput.Focus()

	listItems := make([]list.Item, len(res.Items))
	for i, item := range res.Items {
		listItems[i] = item
	}
	l.SetItems(listItems)

	return &ListContainer{
		list:      &l,
		textInput: &textInput,
		response:  res,
		runAction: runAction,
	}
}

func (c ListContainer) Init() tea.Cmd {
	return nil
}

func (c *ListContainer) SetSize(width, height int) {
	c.width, c.height = width, height
	c.list.SetSize(width, height-lipgloss.Height(c.footerView())-lipgloss.Height(c.headerView()))
}

func (c *ListContainer) headerView() string {
	input := c.textInput.View()
	line := strings.Repeat("─", c.width)
	return lipgloss.JoinVertical(lipgloss.Left, input, line)
}

func (c *ListContainer) footerView() string {
	selectedItem := c.list.SelectedItem()
	if selectedItem == nil {
		return bubbles.SunbeamFooter(c.width, c.response.Title)
	}

	if item, ok := selectedItem.(commands.ScriptItem); ok && len(item.Actions) > 0 {
		return bubbles.SunbeamFooterWithActions(c.width, c.response.Title, item.Actions[0].Title())
	} else {
		return bubbles.SunbeamFooter(c.width, c.response.Title)
	}
}

func (c *ListContainer) Update(msg tea.Msg) (Page, tea.Cmd) {
	var cmds []tea.Cmd

	selectedItem := c.list.SelectedItem()
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if selectedItem == nil {
				break
			}
			selectedItem := selectedItem.(commands.ScriptItem)
			primaryAction := selectedItem.Actions[0]
			return c, c.runAction(primaryAction)
		default:
			if selectedItem == nil {
				break
			}
			selectedItem := selectedItem.(commands.ScriptItem)
			for _, action := range selectedItem.Actions {
				if action.Keybind == msg.String() {
					return c, c.runAction(action)
				}
			}
		}
	case commands.ScriptResponse:
		items := make([]list.Item, len(msg.List.Items))
		for i, item := range msg.List.Items {
			items[i] = item
		}
		cmd := c.list.SetItems(items)
		return c, cmd
	}

	t, cmd := c.textInput.Update(msg)
	if c.response.OnQueryChange != nil && t.Value() != c.textInput.Value() {
		cmds = append(cmds, c.runAction(*c.response.OnQueryChange))
	}
	cmds = append(cmds, cmd)
	c.textInput = &t

	l, cmd := c.list.Update(msg)
	cmds = append(cmds, cmd)
	c.list = &l

	return c, tea.Batch(cmds...)
}

func (c *ListContainer) View() string {
	return lipgloss.JoinVertical(lipgloss.Left, c.headerView(), c.list.View(), c.footerView())
}
