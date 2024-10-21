package main

import (
	"bytes"
	"io"
	"net"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var localhost = "127.0.0.1:"

func TestTelnetClient(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()

			in := &bytes.Buffer{}
			out := &bytes.Buffer{}

			timeout, err := time.ParseDuration("10s")
			require.NoError(t, err)

			client := NewTelnetClient(l.Addr().String(), timeout, io.NopCloser(in), out)
			require.NoError(t, client.Connect())
			defer func() { require.NoError(t, client.Close()) }()

			in.WriteString("hello\n")
			err = client.Send()
			require.NoError(t, err)

			err = client.Receive()
			require.NoError(t, err)
			require.Equal(t, "world\n", out.String())
		}()

		go func() {
			defer wg.Done()

			conn, err := l.Accept()
			require.NoError(t, err)
			require.NotNil(t, conn)
			defer func() { require.NoError(t, conn.Close()) }()

			request := make([]byte, 1024)
			n, err := conn.Read(request)
			require.NoError(t, err)
			require.Equal(t, "hello\n", string(request)[:n])

			n, err = conn.Write([]byte("world\n"))
			require.NoError(t, err)
			require.NotEqual(t, 0, n)
		}()

		wg.Wait()
	})

	t.Run("Error: connection refused", func(t *testing.T) {
		in := &bytes.Buffer{}
		out := &bytes.Buffer{}

		client := NewTelnetClient(localhost, 10*time.Second, io.NopCloser(in), out)

		require.ErrorIs(t, client.Connect(), syscall.ECONNREFUSED)
	})

	t.Run("Error: connection timeout", func(t *testing.T) {
		l, err := net.Listen("tcp", localhost)
		require.NoError(t, err)
		in := &bytes.Buffer{}
		out := &bytes.Buffer{}

		client := NewTelnetClient(l.Addr().String(), 10*time.Nanosecond, io.NopCloser(in), out)
		err = client.Connect()
		require.Contains(t, err.Error(), "i/o timeout")
	})

	t.Run("Error: sending for closed connection", func(t *testing.T) {
		l, err := net.Listen("tcp", localhost)
		require.NoError(t, err)
		in := &bytes.Buffer{}
		out := &bytes.Buffer{}

		client := NewTelnetClient(l.Addr().String(), 10*time.Second, io.NopCloser(in), out)
		require.NoError(t, client.Connect())
		require.NoError(t, l.Close())

		in.WriteString("closed connection")
		require.ErrorIs(t, client.Send(), ErrSending)
	})

	t.Run("Error: receiving for closed connection", func(t *testing.T) {
		l, err := net.Listen("tcp", localhost)
		require.NoError(t, err)
		in := &bytes.Buffer{}
		out := &bytes.Buffer{}

		client := NewTelnetClient(l.Addr().String(), 10*time.Second, io.NopCloser(in), out)
		require.NoError(t, client.Connect())
		require.NoError(t, l.Close())

		require.ErrorIs(t, client.Receive(), ErrReceiving)
	})
}
