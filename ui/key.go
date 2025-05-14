package ui

import "github.com/charmbracelet/bubbles/key"

type MainKeyMap struct {
	Filter       key.Binding
	Add          key.Binding
	Edit         key.Binding
	Delete       key.Binding
	Connect      key.Binding
	SftpConnect  key.Binding
	Quit         key.Binding
	FilterEnter  key.Binding
	FilterCancel key.Binding
}

func NewMainKeyMap() *MainKeyMap {
	return &MainKeyMap{
		Filter:       key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "filter")),
		Add:          key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "add")),
		Edit:         key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit")),
		Delete:       key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "delete")),
		Connect:      key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "connect")),
		SftpConnect:  key.NewBinding(key.WithKeys("p"), key.WithHelp("p", "connect sftp")),
		FilterEnter:  key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "filter enter")),
		FilterCancel: key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "filter cancel")),
		Quit:         key.NewBinding(key.WithKeys("q"), key.WithHelp("q", "quit")),
	}
}
