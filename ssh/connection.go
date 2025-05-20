package ssh

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"tssh/models"
)

func GetValidPath(inputPath string, defaultPath string) string {
	if inputPath != "" {
		if info, err := os.Stat(inputPath); err == nil && !info.IsDir() {
			return inputPath
		}
	}
	return defaultPath
}
func Connect(rctx *models.RunContext) {
	conn := rctx.Context
	protArg := "-p"
	if rctx.Command == models.RunCommandSftp {
		protArg = "-P"
	}
	userHost := fmt.Sprintf("%s@%s", conn.Username, conn.Host)
	var cmd *exec.Cmd
	switch conn.AuthType {
	case models.UsePass:
		pass, err := models.DecryptString(conn.Password)
		if err != nil {
			fmt.Printf("Failed to decrypt password: %v\n", err)
			return
		}
		cmd = exec.Command("sshpass", "-p", pass, string(rctx.Command), "-o", "StrictHostKeyChecking=no", protArg, fmt.Sprintf("%d", conn.Port), userHost)
	case models.UseKey:
		keyPath := strings.TrimSpace(conn.PrivateKey)
		keyPath = GetValidPath(keyPath, "~/.ssh/id_rsa")
		cmd = exec.Command(string(rctx.Command), "-i", keyPath, "-o", "StrictHostKeyChecking=no", protArg, fmt.Sprintf("%d", conn.Port), userHost)
	}
	// 绑定标准输入、输出和错误到当前终端
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("Connecting to %s...\n", conn.Name)
	// 执行命令
	if err := cmd.Run(); err != nil {
		fmt.Printf("SSH connection failed: %v\n", err)
	}
}
