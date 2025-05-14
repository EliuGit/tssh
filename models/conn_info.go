package models

type AuthType int
type RunCommand string

const (
	UsePass AuthType = 1
	UseKey  AuthType = 2
)
const (
	RunCommandSsh  RunCommand = "ssh"
	RunCommandSftp RunCommand = "sftp"
)

type RunContext struct {
	Context *ConnInfo
	Command RunCommand
}

type ConnInfo struct {
	ID         int64    `json:"id"`
	Name       string   `json:"name" validate:"required"`
	Host       string   `json:"host" validate:"required"`
	Port       int      `json:"port" validate:"required"`
	Username   string   `json:"username" validate:"required"`
	AuthType   AuthType `json:"auth_type" validate:"required"`
	Password   string   `json:"password,omitempty"`
	PrivateKey string   `json:"private_key,omitempty"`
}
