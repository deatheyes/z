package main

import (
	"fmt"
	"time"
	"github.com/samuel/go-zookeeper/zk"
)

var Meta = NewContext([]string{}, 30)

type Context struct {
	Hosts   []string
	Path    string
	Conn    *zk.Conn
	Timeout time.Duration
}

func NewContext(hosts []string, timeoutS int) *Context {
	return &Context {
		Path: "/",
		Hosts: hosts,
		Timeout: time.Duration(timeoutS) * time.Second,
	}
}

func (c *Context) Init() error {
	conn, session, err := zk.Connect(c.Hosts, c.Timeout)
	if err != nil {
		return err
	}

	event := <-session
	if event.State != zk.StateConnected && event.State != zk.StateConnecting {
		err = fmt.Errorf("connect to zookeeper %v  failed: %v", c.Hosts, event.State)
	}
	if err != nil {
		c.Conn.Close()
		return err
	}
	c.Conn = conn
	return nil
}
