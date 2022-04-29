package pkg

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/fdingiit/mpl/pkg/simple"
)

const (
	ErrCodeNo = iota
	ErrCodeIncorrectRequest
)

func handleTrans(ctx context.Context, req *simple.Request) (*simple.Response, error) {
	return &simple.Response{
		Header: simple.Header{
			Type:        "RS",
			PageMark:    req.PageMark,
			ServiceCode: req.ServiceCode,
		},
		UnixTimestamp: time.Now().Unix(),
		SerialNo:      req.SerialNo,
		ErrCode:       ErrCodeNo,
		Message:       "ok",
	}, nil
}

func dispatch(ctx context.Context, req *simple.Request) (*simple.Response, error) {
	switch req.ServiceCode {
	case 1000501:
		return handleTrans(ctx, req)
	default:
		return nil, errors.New("should not be here")
	}
}

func handleRequest(conn net.Conn) {
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var response *simple.Response
	var ret []byte
	var err error

	// Make a buffer to hold incoming data.
	buf := make([]byte, 1024*4)
	// Read the incoming connection into the buffer.
	if _, err := conn.Read(buf); err != nil {
		fmt.Println("[SDBS] Error reading:", err.Error())
		return
	}

	request := &simple.Request{}
	if err := request.Decode(ctx, buf); err != nil {
		msg := fmt.Sprintf("Error Decode: %s", err.Error())
		fmt.Println("[SDBS] Error Decode:", err.Error())

		response = &simple.Response{
			Header: simple.Header{
				Type:        "RS",
				PageMark:    request.PageMark,
				ServiceCode: request.ServiceCode,
			},
			UnixTimestamp: time.Now().Unix(),
			SerialNo:      request.SerialNo,
			ErrCode:       ErrCodeIncorrectRequest,
			Message:       msg,
		}

		ret, err = response.Encode(ctx)
		if err != nil {
			fmt.Println("[SDBS] Error Encode:", err.Error())
			return
		}

		goto WRITE
	}

	response, err = dispatch(ctx, request)
	if err != nil {
		fmt.Println("[SDBS] Handle req err:", err.Error())
		return
	}

	ret, err = response.Encode(ctx)
	if err != nil {
		fmt.Println("[SDBS] Error Encode:", err.Error())
		return
	}

WRITE:
	// Send a response back to person contacting us.
	conn.Write(ret)
	fmt.Println("[SDBS] Rsp: ", string(ret))
}

func Serve() {
	l, err := net.Listen("tcp", ":9999")
	if err != nil {
		fmt.Println("[SDBS] Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()

	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("[SDBS] Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn)
	}
}
