package app

import (
	"encoding/json"
	"os/exec"
	"strings"
	"text/template"

	"github.com/alessio/shellescape"
	"github.com/pomdtr/sunbeam/utils"
)

type Script struct {
	Command     string                 `json:"command" yaml:"command"`
	Description string                 `json:"description" yaml:"description"`
	OnSuccess   string                 `json:"onSuccess" yaml:"onSuccess"`
	Cwd         string                 `json:"cwd" yaml:"cwd"`
	Params      map[string]ScriptParam `json:"params" yaml:"params"`
	Page        Page                   `json:"page" yaml:"page"`
}

type Page struct {
	Type        string `json:"type" yaml:"type"`
	Title       string `json:"title" yaml:"title"`
	Mode        string `json:"mode" yaml:"mode"`
	ShowPreview bool   `json:"showPreview" yaml:"showPreview"`
}

type ScriptParam struct {
	Type        string `json:"type"`
	Enum        []any  `json:"enum"`
	Default     any    `json:"default"`
	Description string `json:"description"`
}

type FormInput struct {
	Type string `json:"type"`

	Title       string `json:"title,omitempty"`
	Value       any    `json:"value,omitempty"`
	Placeholder string `json:"placeholder,omitempty"`

	Label string `json:"label,omitempty"`

	Data []struct {
		Title string `json:"title,omitempty"`
		Value string `json:"value,omitempty"`
	} `json:"data"`
}

func (s Script) CheckMissingParams(with ScriptArguments) []string {
	missing := make([]string, 0)
	for key := range s.Params {
		for i, arg := range with {
			if arg.Param == key {
				break
			}
			if i == len(with)-1 {
				missing = append(missing, key)
			}
		}
	}

	return missing
}

type CommandInput struct {
	Params map[string]any
	Query  string
}

func (s Script) Cmd(input CommandInput) (*exec.Cmd, error) {
	var err error

	funcMap := template.FuncMap{
		"query": func() string {
			return shellescape.Quote(input.Query)
		},
	}

	for key, value := range input.Params {
		value := value
		funcMap[key] = func() any {
			switch value := value.(type) {
			case string:
				return shellescape.Quote(value)
			default:
				return value
			}
		}
	}

	rendered, err := utils.RenderString(s.Command, funcMap)
	if err != nil {
		return nil, err
	}
	return exec.Command("sh", "-c", rendered), nil
}

type ScriptItem struct {
	Id          string         `json:"id"`
	Title       string         `json:"title"`
	Subtitle    string         `json:"subtitle"`
	Preview     string         `json:"preview"`
	PreviewCmd  string         `json:"previewCmd"`
	Accessories []string       `json:"accessories"`
	Actions     []ScriptAction `json:"actions"`
}

func (li ScriptItem) PreviewCommand() *exec.Cmd {
	if li.PreviewCmd == "" {
		return nil
	}
	return exec.Command("sh", "-c", li.PreviewCmd)
}

type ScriptAction struct {
	Title    string `json:"title" yaml:"title"`
	Type     string `json:"type" yaml:"type"`
	Shortcut string `json:"shortcut,omitempty" yaml:"shortcut"`

	Text string `json:"text,omitempty" yaml:"textfield"`

	Url         string `json:"url,omitempty" yaml:"url"`
	Path        string `json:"path,omitempty" yaml:"path"`
	Application string `json:"application,omitempty" yaml:"application"`

	Extension string `json:"extension,omitempty" yaml:"extension"`
	Script    string `json:"script,omitempty" yaml:"script"`
	Command   string `json:"command,omitempty" yaml:"command"`

	OnSuccess string          `json:"onSuccess" yaml:"onSuccess"`
	With      ScriptArguments `json:"with,omitempty" yaml:"with"`
}

type ScriptArgument struct {
	Param string `json:"param"`
	Value any    `json:"value"`
}

type ScriptArguments []ScriptArgument

func ParseAction(output string) (action ScriptAction, err error) {
	err = json.Unmarshal([]byte(output), &action)
	return action, err
}

func ParseListItems(output string) (items []ScriptItem, err error) {
	rows := strings.Split(output, "\n")
	for _, row := range rows {
		if row == "" {
			continue
		}
		var item ScriptItem
		err = json.Unmarshal([]byte(row), &item)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}
