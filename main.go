package main

import (
	"fmt"
	"os"
	"path/filepath"

	"xssh/database"
	"xssh/ssh"
	"xssh/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// 获取用户主目录
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting home directory: %v\n", err)
		os.Exit(1)
	}

	// 创建配置目录
	configDir := filepath.Join(homeDir, ".xssh")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		fmt.Printf("Error creating config directory: %v\n", err)
		os.Exit(1)
	}

	// 初始化数据库
	dbPath := filepath.Join(configDir, "connections.db")
	db, err := database.NewDB(dbPath)
	if err != nil {
		fmt.Printf("Error initializing database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	// 获取所有连接
	connections, err := db.GetAllConnections()
	if err != nil {
		fmt.Printf("Error getting connections: %v\n", err)
		os.Exit(1)
	}

	// 启动 TUI
	p := tea.NewProgram(ui.InitialModel(connections, db))
	m, err := p.Run()
	if err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
	if mm, ok := m.(*ui.MainModel); ok && mm.WillConn != nil {
		ssh.Connect(mm.WillConn)
	}
}
