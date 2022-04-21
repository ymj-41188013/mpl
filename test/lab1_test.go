package test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/fdingiit/mpl/pkg/sdbs/server/pkg"
	"github.com/fdingiit/mpl/pkg/simple"
	"github.com/stretchr/testify/assert"
)

func Test_Lab1_TaskA(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		header  simple.Header
		args    args
		want    []byte
		wantErr error
	}{
		{
			name: "",
			header: simple.Header{
				TotalLength: 322,
				Type:        "RQ",
				PageMark:    0,
				Checksum:    "tPK6UhVeIHb2hrsedxXMJHw         ",
				ServiceCode: 1000501,
				Reserved:    0,
			},
			args:    args{ctx: context.TODO()},
			want:    []byte("00000322RQ0tPK6UhVeIHb2hrsedxXMJHw         010005010"),
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := tt.header
			enc, err := h.Encode(tt.args.ctx)
			if !assert.Equal(t, tt.wantErr, err) {
				t.Errorf("[failed] Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, enc, tt.want) {
				t.Errorf("[failed] Encode() got = %v, want %v", enc, tt.want)
				return
			}
			hp := simple.Header{}
			err = hp.Decode(tt.args.ctx, enc)
			if !assert.Equal(t, tt.wantErr, err) {
				t.Errorf("[failed] Decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, h, hp) {
				t.Errorf("[failed] Decode() got = %v, want %v", hp, h)
			}
		})
	}
}

func Test_Lab1_TaskB_Request(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		request simple.Request
		args    args
		want    []byte
		wantErr error
	}{
		{
			name: "",
			request: simple.Request{
				Header: simple.Header{
					TotalLength: 322,
					Type:        "RQ",
					PageMark:    0,
					Checksum:    "tPK6UhVeIHb2hrsedxXMJHw         ",
					ServiceCode: 1000501,
					Reserved:    0,
				},
				UnixTimestamp: 1648811583,
				SerialNo:      12345,
				Currency:      2,
				Amount:        100,
				Unit:          0,
				OutBankId:     2,
				OutAccountId:  1234567899321,
				InBankId:      2,
				InAccountId:   3211541298661,
				Notes:         "",
			},
			args:    args{ctx: context.TODO()},
			want:    []byte("00000328RQ0tPK6UhVeIHb2hrsedxXMJHw         010005010<timestamp>1648811583</timestamp><serial_no>12345</serial_no><currency>2</currency><amount>100</amount><unit>0</unit><out_bank_id>2</out_bank_id><out_account_id>1234567899321</out_account_id><in_bank_id>2</in_bank_id><in_account_id>3211541298661</in_account_id><notes></notes>"),
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := tt.request
			enc, err := req.Encode(tt.args.ctx)
			if !assert.Equal(t, tt.wantErr, err) {
				t.Errorf("[failed] Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, enc, tt.want) {
				t.Errorf("[failed] Encode() got = %v, want %v", enc, tt.want)
				return
			}
			reqp := simple.Request{}
			err = reqp.Decode(tt.args.ctx, enc)
			if !assert.Equal(t, tt.wantErr, err) {
				t.Errorf("[failed] Decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, req, reqp) {
				t.Errorf("[failed] Decode() got = %v, want %v", reqp, req)
			}
		})
	}
}

func Test_Lab1_TaskB_Response(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name     string
		response simple.Response
		args     args
		want     []byte
		wantErr  error
	}{
		{
			name: "",
			response: simple.Response{
				Header: simple.Header{
					TotalLength: 322,
					Type:        "RS",
					PageMark:    0,
					Checksum:    "665db818fa5ef08e9f10ec77d76b9a0e",
					ServiceCode: 1000501,
					Reserved:    0,
				},
				UnixTimestamp: 1648811583,
				SerialNo:      12345,
				ErrCode:       0,
				Message:       "ok",
			},
			args:    args{ctx: context.TODO()},
			want:    []byte("00000156RS0665db818fa5ef08e9f10ec77d76b9a0e010005010<timestamp>1648811583</timestamp><serial_no>12345</serial_no><err_code>0</err_code><message>ok</message>"),
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rsp := tt.response
			enc, err := rsp.Encode(tt.args.ctx)
			if !assert.Equal(t, tt.wantErr, err) {
				t.Errorf("[failed] Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, enc, tt.want) {
				t.Errorf("[failed] Encode() got = %v, want %v", enc, tt.want)
				return
			}
			rspp := simple.Response{}
			err = rspp.Decode(tt.args.ctx, enc)
			if !assert.Equal(t, tt.wantErr, err) {
				t.Errorf("[failed] Decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, rsp, rspp) {
				t.Errorf("[failed] Decode() got = %v, want %v", rspp, rsp)
			}
		})
	}
}

func Test_Lab1_TaskC(t *testing.T) {
	rand.Seed(time.Now().Unix())

	// start sdbs server
	go pkg.Serve()
	err := portCheck1("9999")
	if !assert.Nil(t, err) {
		t.Errorf("[failed] server failed to start error = %v", err.Error())
		return
	}

	// start gateway
	startCmd := exec.Command("/bin/bash", "start_gateway.sh")
	startCmd.Dir = "/Users/dingfei/Go/src/github.com/fdingiit/mpl/pkg/sdbs/gateway/"
	go func() {
		output, err := startCmd.Output()
		if err != nil {
			fmt.Println(output)
			panic(err)
		}
	}()

	// defer stop gateway
	stopCmd := exec.Command("/bin/bash", "stop_gateway.sh")
	stopCmd.Dir = "/Users/dingfei/Go/src/github.com/fdingiit/mpl/pkg/sdbs/gateway/"
	defer stopCmd.Output()

	err = portCheck1("80")
	defer startCmd.Process.Signal(os.Kill)
	if !assert.Nil(t, err) {
		t.Errorf("[failed] gateway failed to start error = %v", err)
		return
	}

	time.Sleep(time.Second)

	type args struct {
		ctx context.Context
	}
	reqs := []GWRequest{
		{
			Timestamp:    time.Now().Unix(),
			SerialNo:     rand.Int63(),
			Currency:     rand.Int31(),
			Amount:       rand.Int63(),
			Unit:         rand.Int31(),
			OutBankId:    rand.Int63(),
			OutAccountId: rand.Int63(),
			InBankId:     rand.Int63(),
			InAccountId:  rand.Int63(),
			Notes:        "random ut",
		},
	}

	tests := []struct {
		name    string
		req     GWRequest
		args    args
		want    GWResponse
		wantErr error
	}{
		{
			name: "",
			req:  reqs[0],
			args: args{ctx: context.TODO()},
			want: GWResponse{
				SerialNo: reqs[0].SerialNo,
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqData, _ := json.Marshal(tt.req)
			buf := bytes.NewBuffer(reqData)
			reqHttp, _ := http.NewRequest("POST", "http://127.0.0.1:80/transfer", buf)
			reqHttp.Header.Set("X-SDBS-PAGING-MASK", "0")
			reqHttp.Header.Set("X-SDBS-CHECKSUM", "665db818fa5ef08e9f10ec77d76b9a0e")
			cli := http.Client{}
			rspHttp, err := cli.Do(reqHttp)
			if !assert.Nil(t, err) {
				t.Errorf("[failed] http.Post() error = %+v", err.Error())
				return
			}
			defer rspHttp.Body.Close()

			rspData, err := ioutil.ReadAll(rspHttp.Body)
			if !assert.Equal(t, tt.wantErr, err) {
				t.Errorf("[failed] ioutil.ReadAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !assert.Equal(t, http.StatusOK, rspHttp.StatusCode) {
				t.Errorf("[failed] http.Post() statuscode = %v, data = %v", rspHttp.StatusCode, string(rspData))
				return
			}

			if !assert.NotNil(t, rspHttp.Body) {
				t.Errorf("[failed] http.Post() nil response")
				return
			}

			var rsp GWResponse
			err = json.Unmarshal(rspData, &rsp)
			if !assert.Nil(t, err) {
				t.Errorf("[failed] Decode() error = %v", err.Error())
				return
			}
			if !assert.Equal(t, rsp.SerialNo, tt.want.SerialNo) {
				t.Errorf("[failed] incorrect rsp got = %v, want %v", rsp, tt.want)
				return
			}
		})
	}
}

func portCheck1(port string) error {
	fmt.Printf("checking port: %+v\n", port)

	ticker := time.NewTicker(time.Second)
	timer := time.NewTimer(time.Second * 10)
	for {
		select {
		case <-ticker.C:
			conn, err := net.DialTimeout("tcp", net.JoinHostPort("127.0.0.1", port), time.Second)
			if err != nil || conn == nil {
				continue
			}
			conn.Close()
			return nil
		case <-timer.C:
			return errors.New("timeout")
		}
	}
}

type GWRequest struct {
	Timestamp int64 `json:"timestamp"`

	SerialNo int64 `json:"serial_no"`

	Currency int32 `json:"currency"`

	Amount int64 `json:"amount"`

	Unit int32 `json:"unit"`

	OutBankId int64 `json:"out_bank_id"`

	OutAccountId int64 `json:"out_account_id"`

	InBankId int64 `json:"in_bank_id"`

	InAccountId int64 `json:"in_account_id"`

	Notes string `json:"notes,omitempty"`
}

type GWResponse struct {
	Timestamp int64 `json:"timestamp"`

	SerialNo int64 `json:"serial_no"`

	ErrCode int32 `json:"err_code"`

	Message string `json:"message,omitempty"`
}
