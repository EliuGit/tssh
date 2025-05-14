package ui

import (
	"fmt"
	"reflect"
	"strings"

	"xssh/database"
	"xssh/models"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var filterFocusStyle = footStyle.Foreground(lipgloss.Color("229"))
var filterBlurStyle = footStyle.Foreground(lipgloss.Color("240"))

type FocusView int8

const (
	Table       FocusView = 0
	FilterInput FocusView = 1
	Confirm     FocusView = 2
)

func genRow(conn *models.ConnInfo) table.Row {
	return table.Row{
		fmt.Sprintf("%d", conn.ID),
		conn.Name,
		conn.Host,
		fmt.Sprintf("%d", conn.Port),
		conn.Username,
	}
}

type MainModel struct {
	table        table.Model
	connections  []models.ConnInfo
	currentItems []*models.ConnInfo
	WillConn     *models.RunContext
	filter       string
	filterInput  textinput.Model
	db           *database.DB
	keyMap       *MainKeyMap
}

func InitialModel(connections []models.ConnInfo, db *database.DB) tea.Model {
	columns := []table.Column{
		{Title: "ID", Width: 4},
		{Title: "Name", Width: 15},
		{Title: "Host", Width: 15},
		{Title: "Port", Width: 6},
		{Title: "Username", Width: 10},
	}

	rows := make([]table.Row, 0)
	currentItems := make([]*models.ConnInfo, 0)
	for _, conn := range connections {
		rows = append(rows, genRow(&conn))
		currentItems = append(currentItems, &conn)
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)

	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	ti := textinput.New()
	ti.Placeholder = ""
	ti.Prompt = "Filter: "
	ti.Width = 20
	ti.CharLimit = 20

	return &MainModel{
		table:        t,
		connections:  connections,
		currentItems: currentItems,
		filter:       "",
		db:           db,
		keyMap:       NewMainKeyMap(),
		filterInput:  ti,
	}
}
func (m *MainModel) Cursor() *models.ConnInfo {
	idx := m.table.Cursor()
	if idx >= len(m.currentItems) {
		return nil
	}
	return m.currentItems[idx]
}

func (m *MainModel) Init() tea.Cmd {
	m.SwitchFocus(Table)
	return nil
}

func (m *MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.Connect):
			m.WillConn = &models.RunContext{Context: m.Cursor(), Command: models.RunCommandSsh}
			return m, tea.Quit
		case key.Matches(msg, m.keyMap.SftpConnect):
			m.WillConn = &models.RunContext{Context: m.Cursor(), Command: models.RunCommandSftp}
			return m, tea.Quit
		case key.Matches(msg, m.keyMap.Add):
			// 创建新的添加表单
			newForm := newFormModel(m, m.db, models.ConnInfo{Port: 22, AuthType: models.UsePass})
			return &newForm, nil
		case key.Matches(msg, m.keyMap.Edit):
			current := m.Cursor()
			if current != nil {
				newForm := newFormModel(m, m.db, *current)
				return &newForm, nil
			}
		case key.Matches(msg, m.keyMap.Delete):
			current := m.Cursor()
			if current != nil {
				dm := newConfirmModel(m, fmt.Sprintf("Are you sure you want to delete %s ?", current.Name))
				return &dm, nil
			}
		case key.Matches(msg, m.keyMap.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keyMap.Filter):
			m.SwitchFocus(FilterInput)
			return m, nil
		case key.Matches(msg, m.keyMap.FilterEnter):
			m.SwitchFocus(Table)
			return m, nil
		case key.Matches(msg, m.keyMap.FilterCancel):
			m.filterInput.SetValue("")
			m.filter = ""
			m.SwitchFocus(Table)
			m.updateTable()
			return m, nil
		}
	case DeleteConfirmMsg:
		if msg.context != nil {
			err := m.db.DeleteConnection(msg.context.ID)
			if err != nil {
				panic(err)
			}
			m.connections, err = m.db.GetAllConnections()
			if err != nil {
				panic(err)
			}
			m.updateTable()
		}
	}
	var cmd tea.Cmd
	if m.filterInput.Focused() {
		m.filterInput, cmd = m.filterInput.Update(msg)
		m.filter = m.filterInput.Value()
		m.updateTable()
		return m, cmd
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}
func (m *MainModel) SwitchFocus(fi FocusView) {
	switch fi {
	case Table:
		m.keyMap.Add.SetEnabled(true)
		m.keyMap.Filter.SetEnabled(true)
		m.keyMap.Edit.SetEnabled(true)
		m.keyMap.Delete.SetEnabled(true)
		m.keyMap.Quit.SetEnabled(true)
		m.keyMap.Connect.SetEnabled(true)
		m.keyMap.SftpConnect.SetEnabled(true)
		m.keyMap.FilterEnter.SetEnabled(false)
		// m.keyMap.FilterCancel.SetEnabled(false)

		m.filterInput.Blur()
		m.table.Focus()
	case FilterInput:
		m.keyMap.Add.SetEnabled(false)
		m.keyMap.Filter.SetEnabled(false)
		m.keyMap.Edit.SetEnabled(false)
		m.keyMap.Delete.SetEnabled(false)
		m.keyMap.Quit.SetEnabled(false)
		m.keyMap.Connect.SetEnabled(false)
		m.keyMap.SftpConnect.SetEnabled(false)
		m.keyMap.FilterEnter.SetEnabled(true)
		// m.keyMap.FilterCancel.SetEnabled(true)

		m.filterInput.Focus()
		m.table.Blur()
	}
}

func (m *MainModel) updateTable() {
	rows := make([]table.Row, 0)
	m.currentItems = make([]*models.ConnInfo, 0)
	for _, conn := range m.connections {
		if m.filter == "" ||
			strings.Contains(strings.ToLower(conn.Name), strings.ToLower(m.filter)) ||
			strings.Contains(strings.ToLower(conn.Host), strings.ToLower(m.filter)) ||
			strings.Contains(strings.ToLower(conn.Username), strings.ToLower(m.filter)) {
			rows = append(rows, genRow(&conn))
			m.currentItems = append(m.currentItems, &conn)
		}
	}
	m.table.SetRows(rows)
	m.table.SetCursor(0)
}

func (m *MainModel) View() string {
	var s strings.Builder
	if m.filterInput.Focused() {
		s.WriteString(filterFocusStyle.Render(m.filterInput.View()))
	} else {
		s.WriteString(filterBlurStyle.Render(m.filterInput.View()))
	}
	s.WriteString("\n")
	s.WriteString(tableStyle.Render(m.table.View()))
	help := helpStyle.Render(m.getHelpStr())
	s.WriteString("\n" + help + "\n")
	return s.String()
}

func (m *MainModel) getHelpStr() string {
	v := reflect.ValueOf(m.keyMap).Elem()
	b := strings.Builder{}
	count := 0
	for i := 0; i < v.NumField(); i++ {
		key := v.Field(i).Interface().(key.Binding)
		if !key.Enabled() {
			continue
		}
		count++
		s := key.Help()
		b.WriteString("'" + s.Key + "'-" + s.Desc + "  ")
		if count%4 == 0 {
			b.WriteString("\n")
		}
	}
	return b.String()
}
