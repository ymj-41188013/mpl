package test

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/fdingiit/mpl/pkg/plugin/demo/codec"
	"github.com/stretchr/testify/assert"
	v2 "mosn.io/mosn/pkg/config/v2"
)

const (
	mosnExecPath = "../build/mosn/mosnd"
	mosnVersion  = "v0.26.0"

	pluginCompiler = "../pkg/protocol/demo/make_codec.sh"

	mosnConfTaskBServer = "../build/mosn/mosn_config_lab0_taskb_server.json"
	mosnConfTaskBClient = "../build/mosn/mosn_config_lab0_taskb_client.json"

	mosnConfTaskCServer = "../build/mosn/mosn_config_lab0_taskc_server.json"
	mosnConfTaskCClient = "../build/mosn/mosn_config_lab0_taskc_client.json"
)

var random int

func Test_TaskA(t *testing.T) {
	cmd := exec.Command("/bin/bash", "-c", mosnExecPath)
	out, err := cmd.Output()
	if !assert.Nil(t, err) {
		t.Errorf("[failed] run mosn err: %+v", err)
		t.FailNow()
	}

	t.Logf("[pass] correct mosn executable")

	if !assert.True(t, strings.Contains(string(out), mosnVersion)) {
		t.Errorf("[failed] incorrect mosn version: %+v", string(out))
		t.FailNow()
	}

	t.Logf("[pass] correct mosn version")
}

func Test_TaskB(t *testing.T) {
	_, err := os.Stat(mosnConfTaskBServer)
	if !assert.Nil(t, err) {
		t.Errorf("[failed] mosn server conf err: %+v", err)
		t.FailNow()
	}

	t.Logf("[pass] mosn server config exists")

	_, err = os.Stat(mosnConfTaskBClient)
	if !assert.Nil(t, err) {
		t.Errorf("[failed] mosn client conf err: %+v", err)
		t.FailNow()
	}

	t.Logf("[pass] mosn client config exists")

	// start http server
	httpServe()

	// start server side mosn
	serverCmd := startMosn(mosnConfTaskBServer)
	defer serverCmd.Process.Kill()
	if !assert.NotNil(t, serverCmd) {
		t.Errorf("[failed] run server mosn failed")
		t.FailNow()
	}
	t.Logf("[pass] mosn server stared")

	// start client side mosn
	clientCmd := startMosn(mosnConfTaskBClient)
	defer clientCmd.Process.Kill()
	if !assert.NotNil(t, serverCmd) {
		t.Errorf("[failed] run client mosn failed")
		t.FailNow()
	}
	t.Logf("[pass] mosn client stared")

	// call
	resp, err := http.Get("http://localhost:12045")
	defer resp.Body.Close()
	if !assert.Nil(t, err) {
		t.Errorf("[failed] call mosn client err: %+v", err)
		t.FailNow()
	}

	if !assert.Equal(t, http.StatusOK, resp.StatusCode) {
		t.Errorf("[failed] call mosn client err: %+v", resp.StatusCode)
		t.FailNow()
	}
	body, err := ioutil.ReadAll(resp.Body)
	if !assert.Nil(t, err) {
		t.Errorf("[failed] read mosn client resp err: %+v", err)
		t.FailNow()
	}
	got, err := strconv.Atoi(string(body))
	if !assert.Nil(t, err) {
		t.Errorf("[failed] read mosn client resp err: %+v, body: %+v", err, body)
		t.FailNow()
	}
	if !assert.Equal(t, random, got) {
		t.Errorf("[failed] unexcepted mosn client resp: %+v, wanted: %+v", got, random)
		t.FailNow()
	}

	t.Logf("[pass] correct response")
}

func Test_TaskC(t *testing.T) {
	_, err := os.Stat(mosnConfTaskCServer)
	if !assert.Nil(t, err) {
		t.Errorf("[failed] mosn server conf err: %+v", err)
		t.FailNow()
	}

	t.Logf("[pass] mosn server config exists")

	_, err = os.Stat(mosnConfTaskCClient)
	if !assert.Nil(t, err) {
		t.Errorf("[failed] mosn client conf err: %+v", err)
		t.FailNow()
	}

	t.Logf("[pass] mosn client config exists")

	// start demo server
	go serveDemo()

	// start server side mosn
	serverCmd := startMosn(mosnConfTaskCServer)
	defer serverCmd.Process.Kill()
	if !assert.NotNil(t, serverCmd) {
		t.Errorf("[failed] run server mosn failed")
		t.FailNow()
	}
	t.Logf("[pass] mosn server stared")

	// start client side mosn
	clientCmd := startMosn(mosnConfTaskCClient)
	defer clientCmd.Process.Kill()
	if !assert.NotNil(t, serverCmd) {
		t.Errorf("[failed] run client mosn failed")
		t.FailNow()
	}
	t.Logf("[pass] mosn client stared")

	// call
	conn, err := net.Dial("tcp", "127.0.0.1:2045")
	if !assert.Nil(t, err) {
		t.Errorf("[failed] failed to conn to 127.0.0.1:2045: %+v", err)
		t.FailNow()
	}
	bytes := []byte(reqMessage)
	buf := make([]byte, 0)

	buf = append(buf, codec.Magic)
	buf = append(buf, codec.TypeMessage)
	buf = append(buf, codec.DirRequest)
	tempBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(tempBytes, requestId)
	tempBytesSec := make([]byte, 4)
	binary.BigEndian.PutUint32(tempBytesSec, uint32(len(bytes)))
	buf = append(buf, tempBytes...)
	buf = append(buf, tempBytesSec...)
	buf = append(buf, bytes...)

	//2.send message
	_, err = conn.Write(buf)
	if !assert.Nil(t, err) {
		t.Errorf("[failed] write failed to 127.0.0.1:2045: %+v", err)
		t.FailNow()
	}

	respBuff := make([]byte, 1024)

	//3.read response
	read, err := conn.Read(respBuff)
	if !assert.Nil(t, err) {
		t.Errorf("[failed] read failed from 127.0.0.1:2045: %+v", err)
		t.FailNow()
	}
	resp := respBuff[:read]

	//4.decodeResponse
	response, err := decodeResponse(nil, resp)
	if !assert.Nil(t, err) {
		t.Errorf("[failed] decodeResponse failed: %+v", err)
		t.FailNow()
	}
	fmt.Println(string(response.(*Response).Payload[:]))

	if !assert.Equal(t, respMessage, string(response.(*Response).Payload[:])) {
		t.Errorf("[failed] incorrect response: %+v", string(response.(*Response).Payload[:]))
		t.FailNow()
	}

	t.Logf("[pass] correct response")
}

func httpServe() {
	rand.Seed(time.Now().UnixNano())

	hdl := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		random = rand.Int()
		fmt.Fprintf(w, "%d", random)
	}

	http.HandleFunc("/", hdl)
	go func() {
		if err := http.ListenAndServe("127.0.0.1:1080", nil); err != nil {
			panic(err)
		}
	}()

	time.Sleep(time.Second)
}

func startMosn(path string) *exec.Cmd {
	cmd := exec.Command("/bin/bash", "-c", fmt.Sprintf("%s start -c %s", mosnExecPath, path))

	go cmd.Output()

	b, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var config v2.MOSNConfig
	if err := json.Unmarshal(b, &config); err != nil {
		panic(err)
	}

	deadline := time.NewTimer(5 * time.Second)
	ticker := time.NewTicker(500 * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			resp, err := http.Get(fmt.Sprintf("http://%s:%d/api/v1/states", config.GetAdmin().GetAddress(), config.GetAdmin().GetPortValue()))
			if err != nil {
				continue
			}
			if resp.StatusCode == http.StatusOK {
				resp.Body.Close()
				return cmd
			}
		case <-deadline.C:
			return nil
		}
	}

	// should not be here
	return nil
}

const respMessage = "Hello, I am server"

type Request struct {
	Type       byte
	RequestId  uint32
	PayloadLen uint32
	Payload    []byte
}

func decodeRequest(ctx context.Context, data []byte) (cmd interface{}, err error) {
	bytesLen := len(data)

	// 1. least bytes to decode header is RequestHeaderLen
	if bytesLen < codec.RequestHeaderLen {
		return nil, errors.New("short bytesLen")
	}

	// 2. least bytes to decode whole frame
	payloadLen := binary.BigEndian.Uint32(data[codec.RequestPayloadIndex:codec.RequestHeaderLen])
	frameLen := codec.RequestHeaderLen + int(payloadLen)
	if bytesLen < frameLen {
		return nil, errors.New("not whole bytesLen")
	}

	// 3. Request
	request := &Request{
		Type:       data[codec.TypeIndex],
		RequestId:  binary.BigEndian.Uint32(data[codec.RequestIdIndex : codec.RequestIdEnd+1]),
		PayloadLen: payloadLen,
	}

	//4. copy data for io multiplexing
	request.Payload = data[codec.RequestHeaderLen:]
	return request, err
}

func serve(c net.Conn) {
	reqBuff := make([]byte, 64)

	readLength, err := c.Read(reqBuff)
	if err == nil {
		req := reqBuff[:readLength]

		request, err := decodeRequest(nil, req)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println(string((request.(*Request).Payload)[:]))

		bytes := []byte(respMessage)
		buf := make([]byte, 0)
		buf = append(buf, codec.Magic)
		buf = append(buf, codec.TypeMessage)
		buf = append(buf, codec.DirResponse)

		tempBytes := make([]byte, 4)

		binary.BigEndian.PutUint32(tempBytes, request.(*Request).RequestId)
		tempBytesSec := make([]byte, 2)

		binary.BigEndian.PutUint16(tempBytesSec, codec.ResponseStatusSuccess)
		tempBytesThr := make([]byte, 4)

		binary.BigEndian.PutUint32(tempBytesThr, uint32(len(bytes)))
		buf = append(buf, tempBytes...)
		buf = append(buf, tempBytesSec...)
		buf = append(buf, tempBytesThr...)
		buf = append(buf, bytes...)
		c.Write(buf)
		_ = c.Close()

	}
}

func serveDemo() {
	//1.create server
	conn, err := net.Listen("tcp", "127.0.0.1:8086")
	if err != nil {
		panic("conn failed")
	}
	for {
		accept := conn.Accept
		c, err := accept()
		if err != nil {
			fmt.Println("accept closed")
			continue
		}
		//let serve do accept
		go serve(c)
	}
}

const reqMessage = "Hello World"
const requestId = 1

type Response struct {
	Type       byte
	RequestId  uint32
	PayloadLen uint32
	Payload    []byte
	Status     uint16
}

func decodeResponse(ctx context.Context, bytes []byte) (cmd interface{}, err error) {
	bytesLen := len(bytes)

	// 1. least bytes to decode header is ResponseHeaderLen
	if bytesLen < codec.ResponseHeaderLen {
		return nil, errors.New("bytesLen<ResponseHeaderLen")
	}

	payloadLen := binary.BigEndian.Uint32(bytes[codec.ResponsePayloadIndex:codec.ResponseHeaderLen])

	//2.total protocol length
	frameLen := codec.ResponseHeaderLen + int(payloadLen)
	if bytesLen < frameLen {
		return nil, errors.New("short bytesLen")
	}

	// 3.  response
	response := &Response{

		Type:       bytes[codec.DirIndex],
		RequestId:  binary.BigEndian.Uint32(bytes[codec.RequestIdIndex : codec.RequestIdEnd+1]),
		PayloadLen: payloadLen,
		Status:     codec.ResponseStatusSuccess,
	}

	//4. copy data for io multiplexing
	response.Payload = bytes[codec.ResponseHeaderLen:]
	return response, nil
}

func main() {
	//1.create client
	conn, err := net.Dial("tcp", "127.0.0.1:2045")
	if err != nil {
		panic("conn failed")
	}
	bytes := []byte(reqMessage)
	buf := make([]byte, 0)

	buf = append(buf, codec.Magic)
	buf = append(buf, codec.TypeMessage)
	buf = append(buf, codec.DirRequest)
	tempBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(tempBytes, requestId)
	tempBytesSec := make([]byte, 4)
	binary.BigEndian.PutUint32(tempBytesSec, uint32(len(bytes)))
	buf = append(buf, tempBytes...)
	buf = append(buf, tempBytesSec...)
	buf = append(buf, bytes...)

	//2.send message
	_, err = conn.Write(buf)
	if err != nil {
		panic("write failed")
	}

	respBuff := make([]byte, 1024)

	//3.read response
	read, err := conn.Read(respBuff)
	if err != nil {
		panic(err)
	}
	resp := respBuff[:read]

	//4.decodeResponse
	response, err := decodeResponse(nil, resp)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(string(response.(*Response).Payload[:]))
}
