package database

import (
	"database/sql"
	"errors"
	"xssh/models"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	*sql.DB
}

func NewDB(dbPath string) (*DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	if err := createTables(db); err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

func createTables(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS ssh_connections (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		host TEXT NOT NULL,
		port INTEGER NOT NULL,
		username TEXT NOT NULL,
		auth_type INTEGER NOT NULL, -- 1: password, 2: key
		password TEXT,
		private_key TEXT
	);`

	_, err := db.Exec(query)
	return err
}

func (db *DB) GetAllConnections() ([]models.ConnInfo, error) {
	rows, err := db.Query("SELECT * FROM ssh_connections order by name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var connections []models.ConnInfo
	for rows.Next() {
		var conn models.ConnInfo
		err := rows.Scan(
			&conn.ID,
			&conn.Name,
			&conn.Host,
			&conn.Port,
			&conn.Username,
			&conn.AuthType,
			&conn.Password,
			&conn.PrivateKey,
		)
		if err != nil {
			return nil, err
		}
		connections = append(connections, conn)
	}
	return connections, nil
}
func (db *DB) GetConnection(id int64) (models.ConnInfo, error) {
	row := db.QueryRow("SELECT * FROM ssh_connections WHERE id = ?", id)

	var conn models.ConnInfo
	err := row.Scan(
		&conn.ID,
		&conn.Name,
		&conn.Host,
		&conn.Port,
		&conn.Username,
		&conn.AuthType,
		&conn.Password,
		&conn.PrivateKey,
	)
	if err != nil {
		return models.ConnInfo{}, err
	}
	return conn, nil
}

func (db *DB) AddConnection(conn models.ConnInfo) error {
	if conn.AuthType == models.UsePass {
		if conn.Password == "" {
			return errors.New("password is required for password authentication")
		}
		var err error
		conn.Password, err = models.EncryptString(conn.Password)
		if err != nil {
			return err
		}
	}

	query := `
	INSERT INTO ssh_connections (name, host, port, username, auth_type, password, private_key)
	VALUES (?, ?, ?, ?, ?, ?, ?)`

	_, err := db.Exec(query,
		conn.Name,
		conn.Host,
		conn.Port,
		conn.Username,
		conn.AuthType,
		conn.Password,
		conn.PrivateKey,
	)
	return err
}

func (db *DB) UpdateConnection(conn models.ConnInfo) error {
	if conn.AuthType == models.UsePass && conn.Password != "" {
		var err error
		conn.Password, err = models.EncryptString(conn.Password)
		if err != nil {
			return err
		}
	} else {
		oldConn, err := db.GetConnection(conn.ID)
		if err != nil {
			return err
		}
		if conn.AuthType == models.UsePass && conn.Password == "" && oldConn.Password == "" {
			return errors.New("password is required for password authentication")
		}
		conn.Password = oldConn.Password
	}

	query := `
	UPDATE ssh_connections
	SET name = ?, host = ?, port = ?, username = ?, auth_type = ?, password = ?, private_key = ?
	WHERE id = ?`

	_, err := db.Exec(query,
		conn.Name,
		conn.Host,
		conn.Port,
		conn.Username,
		conn.AuthType,
		conn.Password,
		conn.PrivateKey,
		conn.ID,
	)
	return err
}

func (db *DB) DeleteConnection(id int64) error {
	_, err := db.Exec("DELETE FROM ssh_connections WHERE id = ?", id)
	return err
}
