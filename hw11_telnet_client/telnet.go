package main

import (
	"errors"
	"io"
	"net"
	"time"
)

var (
	ErrSending   = errors.New("sending error")
	ErrReceiving = errors.New("receiving error")
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &client{address, timeout, in, out, nil}
}

type client struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
}

func (c *client) Connect() error {
	var err error
	if c.conn, err = net.DialTimeout("tcp", c.address, c.timeout); err != nil {
		return err
	}
	return nil
}

func (c *client) Close() error {
	return c.conn.Close()
}

func (c *client) Send() error {
	if _, err := io.Copy(c.conn, c.in); err != nil {
		return ErrSending
	}
	return nil
}

func (c *client) Receive() error {
	if _, err := io.Copy(c.out, c.conn); err != nil {
		return ErrReceiving
	}
	return nil
}
