package main

import (
	"code.google.com/p/go.crypto/ssh"
	"github.com/msurdi/alf/db"
	"strconv"
	"sync"
)

// HostConnection implements the connection to a remote host
type HostConnection struct {
	// Host is a pointer to the Host instance this connection is associated with
	Host            *db.Host
	connection      *ssh.ClientConn
	connectionMutex sync.Mutex
}

var (
	// We keep a hostname keyed map of connections, we never should have more than
	// one connection per host
	connections = make(map[string]HostConnection)

	// connectionsLock is used to prevent two different goroutines attemp to connect
	// simultaneously to the same host
	connectionsLock = make(chan int, 1)
)

// GetHostConnection returns a HostConnection instance for the given Host.
func GetHostConnection(host *db.Host) *HostConnection {
	// Avoid the lock if it already exists
	if connection, found := connections[host.Hostname]; found {
		return &connection
	}

	// If it doesn't exist, wait for the lock, task again, create if needed
	connectionsLock <- 1
	connection, found := connections[host.Hostname]
	if !found {
		connection = HostConnection{Host: host}
		connections[host.Hostname] = connection
	}
	<-connectionsLock
	return &connection
}

// Close closes the connection with the host.
func (c *HostConnection) close() {
	if c.connection != nil {
		c.connection.Close()
		c.connection = nil
	}
}

// Connect implements the ssh connection phase
func (c *HostConnection) connect() error {
	c.connectionMutex.Lock()
	if c.connection == nil {
		config := &ssh.ClientConfig{
			User: c.Host.Username,
			Auth: []ssh.ClientAuth{
				ssh.ClientAuthPassword(password(c.Host.Password)),
			},
		}
		var err error
		url := c.Host.Hostname + ":" + strconv.Itoa(c.Host.Port)
		c.connection, err = ssh.Dial("tcp", url, config)
		if err != nil {
			return err
		}
	}
	c.connectionMutex.Unlock()
	return nil
}

// Get a new session from the host, tries to reconnect once
// on failure.
func (c *HostConnection) GetSession() (*ssh.Session, error) {
	if err := c.connect(); err != nil {
		return nil, err
	}
	// Setup a new session for the task
	session, err := c.connection.NewSession()
	if err != nil {
		c.connection.Close()
		// On Error, try to reconnect once
		if err = c.connect(); err != nil {
			return nil, err
		} else {
			session, err = c.connection.NewSession()
			if err != nil {
				return nil, err
			}
		}
	}
	return session, nil
}
