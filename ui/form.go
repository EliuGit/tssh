package ui

import (
	"strconv"
	"strings"
	"tssh/database"
	"tssh/models"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-playground/validator/v10"
)

const (
	keyCharH tea.KeyType = 'h'
	keyCharL tea.KeyType = 'l'
)

const (
	idxName = iota
	idxHost
	idxPort
	idxUsername
	idxAuthType
	idxPass
	idxPrivateKey
	idxEnter
	idxCancel
)

type formModel struct {
	inputs           []textinput.Model
	title            string
	focusIndex       int
	authTypeSelected models.AuthType
	mainModel        tea.Model
	db               *database.DB
	conn             models.ConnInfo
	err              error
	isEdit           bool
}

func newFormModel(mainModel tea.Model, db *database.DB, conn models.ConnInfo) formModel {
	m := formModel{
		inputs:           make([]textinput.Model, 7),
		focusIndex:       0,
		authTypeSelected: conn.AuthType,
		mainModel:        mainModel,
		db:               db,
		conn:             conn,
		isEdit:           conn.ID != 0,
	}
	if m.isEdit {
		m.title = "Modify SSH Connection"
	} else {
		m.title = "Add SSH Connection"
	}

	// 输入框-连接名称
	t := textinput.New()
	t.Placeholder = "Connection name"
	t.Prompt = "  Name "
	if m.isEdit {
		t.SetValue(conn.Name)
	}
	t.Focus()
	t.Width = 30
	t.CharLimit = 50
	m.inputs[idxName] = t

	// 输入框-连接地址
	t = textinput.New()
	t.Placeholder = "Connection host"
	t.Prompt = "  Host "
	if m.isEdit {
		t.SetValue(conn.Host)
	}
	t.Width = 30
	t.CharLimit = 30
	m.inputs[idxHost] = t

	// 输入框-连接端口
	t = textinput.New()
	t.Placeholder = "Connection port"
	t.Prompt = "  Port "
	if m.isEdit {
		t.SetValue(strconv.Itoa(conn.Port))
	}
	t.Width = 5
	t.CharLimit = 5
	m.inputs[idxPort] = t

	// 输入框-登录用户名
	t = textinput.New()
	t.Placeholder = "Username"
	t.Prompt = "  User "
	if m.isEdit {
		t.SetValue(conn.Username)
	}
	t.Width = 20
	t.CharLimit = 20
	m.inputs[idxUsername] = t

	// 输入框-密码
	t = textinput.New()
	t.Placeholder = "Don't anything when empty"
	t.Prompt = "  Pass "
	t.Width = 30
	t.CharLimit = 30
	m.inputs[idxPass] = t

	// 输入框-密码
	t = textinput.New()
	t.Placeholder = "PrivateKey"
	t.Prompt = "  Key "
	if m.isEdit {
		t.SetValue(conn.PrivateKey)
	}
	t.Width = 40
	t.CharLimit = 5000
	m.inputs[idxPrivateKey] = t

	return m
}

func (m formModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m formModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			return m.mainModel, nil
		case tea.KeyEnter:
			if m.focusIndex == idxEnter {
				port, _ := strconv.Atoi(m.inputs[idxPort].Value())
				conn := models.ConnInfo{
					ID:         m.conn.ID,
					Name:       m.inputs[idxName].Value(),
					Host:       m.inputs[idxHost].Value(),
					Port:       port,
					Username:   m.inputs[idxUsername].Value(),
					AuthType:   m.authTypeSelected,
					Password:   m.inputs[idxPass].Value(),
					PrivateKey: m.inputs[idxPrivateKey].Value(),
				}
				v := validator.New(validator.WithRequiredStructEnabled())
				err := v.Struct(conn)
				if err != nil {
					m.err = err
					return m, nil
				}
				if m.isEdit {
					err = m.db.UpdateConnection(conn)
				} else {
					err = m.db.AddConnection(conn)
				}
				if err != nil {
					m.err = err
					return m, nil
				}

				connections, err := m.db.GetAllConnections()
				if err != nil {
					m.err = err
					return m, nil
				}
				mainModel := m.mainModel.(*MainModel)
				mainModel.connections = connections
				mainModel.updateTable()

				return m.mainModel, nil
			} else if m.focusIndex == idxCancel {
				return m.mainModel, nil
			} else {
				m.changeFoucs(1)
			}

		case tea.KeyTab, tea.KeyDown:
			m.changeFoucs(1)
		case tea.KeyShiftTab, tea.KeyUp:
			m.changeFoucs(-1)
		case tea.KeySpace, tea.KeyLeft, tea.KeyRight, keyCharH, keyCharL:
			if m.focusIndex == idxAuthType {
				if m.authTypeSelected == models.UsePass {
					m.authTypeSelected = models.UseKey
				} else {
					m.authTypeSelected = models.UsePass
				}
			}
		}
	}
	cmd := m.updateInputs(msg)
	return m, cmd
}
func (m *formModel) changeFoucs(i int) {
	if m.focusIndex <= idxPrivateKey && m.focusIndex != idxAuthType {
		m.inputs[m.focusIndex].Blur()
	}
	if i > 0 {
		m.focusIndex = (m.focusIndex + 1) % (len(m.inputs) + 2)
		if m.focusIndex <= idxPrivateKey && m.focusIndex != idxAuthType {
			if m.focusIndex == idxPass && m.authTypeSelected == models.UseKey {
				m.focusIndex = idxPrivateKey
			}
			m.inputs[m.focusIndex].Focus()
			if m.focusIndex == idxPrivateKey && m.authTypeSelected == models.UsePass {
				m.focusIndex = m.focusIndex + 1
			}
		}
	} else {
		l := len(m.inputs) + 2
		m.focusIndex = (m.focusIndex - 1 + l) % l
		if m.focusIndex <= idxPrivateKey && m.focusIndex != idxAuthType {
			if m.focusIndex == idxPrivateKey && m.authTypeSelected == models.UsePass {
				m.focusIndex = idxPass
			}
			m.inputs[m.focusIndex].Focus()
			if m.focusIndex == idxPass && m.authTypeSelected == models.UseKey {
				m.focusIndex = m.focusIndex - 1
			}
		}
	}
}

func (m formModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m formModel) View() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render(m.title) + "\n\n")
	b.WriteString(m.inputView(idxName) + "\n\n")
	b.WriteString(m.inputView(idxHost) + "\n\n")
	b.WriteString(m.inputView(idxPort) + "\n\n")
	b.WriteString(m.inputView(idxUsername) + "\n\n")
	b.WriteString(m.authTypeView())
	b.WriteString("\n\n")
	if m.authTypeSelected == models.UsePass {
		b.WriteString(m.inputView(idxPass) + "\n")
	} else {
		b.WriteString(m.inputView(idxPrivateKey) + "\n")
	}
	if m.err != nil {
		b.WriteString(errorStyle.Render(m.err.Error()) + "\n")
	} else {
		b.WriteString("\n")
	}

	b.WriteString(m.buttonView())
	b.WriteString("\n\n")
	return b.String()
}

func (m formModel) authTypeView() string {
	var b strings.Builder
	if m.focusIndex == idxAuthType {
		b.WriteString(focusedStyle.Render("> AuthType "))
	} else {
		b.WriteString(noStyle.Render("  AuthType "))
	}
	if m.authTypeSelected == models.UsePass {
		b.WriteString(focusedStyle.Render("(x)Password"))
		b.WriteString("   ( )PrivateKey")
	} else {
		b.WriteString("( )Password   ")
		b.WriteString(focusedStyle.Render("(x)PrivateKey"))
	}
	return b.String()
}

func (m formModel) inputView(idx int) string {
	if idx > idxPrivateKey || idx == idxAuthType {
		return ""
	}
	t := m.inputs[idx]
	if idx == m.focusIndex {
		t.Prompt = strings.Replace(t.Prompt, " ", ">", 1)
		t.PromptStyle = focusedStyle
		return t.View()
	}
	t.Prompt = strings.Replace(t.Prompt, ">", " ", 1)
	t.PromptStyle = noStyle
	return t.View()
}

func (m formModel) buttonView() string {
	var b strings.Builder
	if m.focusIndex == idxEnter {
		b.WriteString(focusedStyle.Render("  [ Enter ]   "))
	} else {
		b.WriteString(noStyle.Render("  [ Enter ]   "))
	}
	if m.focusIndex == idxCancel {
		b.WriteString(focusedStyle.Render("  [ Cancel ]   "))
	} else {
		b.WriteString(noStyle.Render("  [ Cancel ]   "))
	}
	return b.String()
}
