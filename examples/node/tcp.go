package main

import (
	"bytes"
	"fmt"
	"github.com/gopherjs/gopherjs/js"
)

var net = js.Global.Call("require", "net")

type ChromeTcp struct {
	buffer  *bytes.Buffer
	gotData chan int
	socket  *js.Object
}

func (c *ChromeTcp) listen() {
	c.socket.Call("on", "data", func(data []byte) {
		go func() {
			bytes := data
			c.buffer.Write(bytes)
			select {
			case c.gotData <- 1:
			default:
			}
		}()
	})

	c.socket.Call("on", "error", func(data *js.Object) {
		go func() {
			resultCode := -1
			c.gotData <- resultCode
		}()
	})
}

func uint8ArrayToArrayBuffer(p *js.Object) *js.Object {
	buffer := p.Get("buffer")
	byteOffset := p.Get("byteOffset").Int()
	byteLength := p.Get("byteLength").Int()
	if byteOffset != 0 || byteLength != buffer.Get("byteLength").Int() {
		return buffer.Call("slice", byteOffset, byteOffset+byteLength)
	}
	return buffer
}

func MakeConn() (*ChromeTcp, error) {
	done := make(chan int)
	buf := make([]byte, 0, 4096)
	conn := &ChromeTcp{
		buffer:  bytes.NewBuffer(buf),
		gotData: make(chan int),
	}

	socket := net.Get("Socket").New()
	conn.socket = socket
	socket.Call("connect", 4070, "sjc1-accesspoint-a40.ap.spotify.com", func(result *js.Object) {
		done <- 1
	})

	<-done
	conn.listen()
	return conn, nil
}

func (c *ChromeTcp) Write(buf []byte) (int, error) {
	done := make(chan int)
	arrayBuffer := js.Global.Get("Buffer").New(buf)
	c.socket.Call("write", arrayBuffer, func(bytesWritten *js.Object) {
		done <- bytesWritten.Int()
	})

	res := <-done
	if res >= 0 {
		return res, nil
	} else {
		return 0, fmt.Errorf("Failed chrome.sockets.tcp write, error code: %v", res)
	}

}

func (c *ChromeTcp) Read(buf []byte) (int, error) {
	if c.buffer.Len() == 0 {
		resultCode := <-c.gotData
		if resultCode < 0 {
			return 0, fmt.Errorf("Failed chrome.sockets.tcp read, error code: %v", resultCode)
		}
	}
	return c.buffer.Read(buf)
}
