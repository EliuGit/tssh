package ui

import (
	"xssh/models"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	idxDialogYes = iota
	idxDialogNo
)

type DeleteConfirmMsg struct {
	context *models.ConnInfo
}

type confirmModel struct {
	focusIndex int
	question   string
	mainModel  tea.Model
}

func newConfirmModel(mainModel tea.Model, question string) confirmModel {

	return confirmModel{
		focusIndex: idxDialogNo,
		question:   question,
		mainModel:  mainModel,
	}
}

func (m confirmModel) Init() tea.Cmd {
	return nil
}
func (m confirmModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "n", "esc", "y", "q", "enter":
			if msg.String() == "y" || (msg.String() == "enter" && m.focusIndex == idxDialogYes) {
				return m.mainModel, func() tea.Msg {
					return DeleteConfirmMsg{m.mainModel.(*MainModel).Cursor()}
				}
			}
			return m.mainModel, nil
		case "tab", "shift+tab", "h", "l":
			m.focusIndex = (m.focusIndex + 1) % 2
		}
	}
	return m, nil
}
func (m confirmModel) View() string {
	question := questionStyle.Render(m.question)

	var yesButton, noButton string
	if m.focusIndex == idxDialogYes {
		yesButton = focusedButtonStyle.Render("Yes")
		noButton = buttonStyle.Render("No")
	} else {
		yesButton = buttonStyle.Render("Yes")
		noButton = focusedButtonStyle.Render(" No")
	}
	buttons := lipgloss.JoinHorizontal(lipgloss.Top, yesButton, " ", noButton)

	dialogContent := lipgloss.JoinVertical(lipgloss.Center, question, buttons)
	ui := dialogBoxStyle.Render(dialogContent)
	return ui
}
