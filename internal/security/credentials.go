package security

import (
	"fmt"
	"strings"
)

type ConnectionInfo struct {
	SourceType string
	Host       string
	Port       int
	Database   string
	Username   string
	Password   string
	Table      string
}

func ParseConnectionString(connStr string) (*ConnectionInfo, error) {
	info := &ConnectionInfo{}
	parts := strings.Split(connStr, "@")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid connection string format")
	}
	userPass := strings.Split(parts[0], ":")
	if len(userPass) != 2 {
		return nil, fmt.Errorf("invalid username:password format")
	}
	info.Username = userPass[0]
	info.Password = userPass[1]
	hostDb := strings.Split(parts[1], "/")
	if len(hostDb) != 2 {
		return nil, fmt.Errorf("invalid host/database format")
	}
	hostPort := strings.Split(hostDb[0], ":")
	info.Host = hostPort[0]
	if len(hostPort) > 1 {
		fmt.Sscanf(hostPort[1], "%d", &info.Port)
	} else {
		info.Port = 3306
	}
	info.Database = hostDb[1]
	return info, nil
}

func (c *ConnectionInfo) Clear() {
	c.Password = ""
	c.Username = ""
}

func (c *ConnectionInfo) ToDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", c.Username, c.Password, c.Host, c.Port, c.Database)
}
