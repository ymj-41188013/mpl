package test

import (
	"context"
	"testing"

	"github.com/fdingiit/mpl/pkg/plugin/simple"
	simple2 "github.com/fdingiit/mpl/pkg/simple"
	"github.com/stretchr/testify/assert"
	"mosn.io/api"
	"mosn.io/pkg/buffer"
)

func Test_Lab2_TaskA_Interface(t *testing.T) {
	var c api.XProtocolCodec
	var x api.XProtocol

	c = simple.CodecSimple{}
	x = c.NewXProtocol(context.TODO())

	t.Logf("codec: %+v", c)
	t.Logf("xprotocol: %+v", x)
}

func Test_Lab2_TaskA_LoadCodec(t *testing.T) {
	var c api.XProtocolCodec

	c = simple.LoadCodec()
	t.Logf("codec: %+v", c)
}

func Test_Lab2_TaskB_ProtocolName(t *testing.T) {
	tests := []struct {
		name string
		want api.ProtocolName
	}{
		{
			name: "",
			want: simple.ProtocolName,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			codec := simple.LoadCodec()
			xp := codec.NewXProtocol(context.TODO())

			if !assert.Equal(t, tt.want, xp.Name()) {
				t.Errorf("[failed] incorrect protocol name: %+v", xp.Name())
				t.FailNow()
			}

			if !assert.Equal(t, tt.want, codec.ProtocolName()) {
				t.Errorf("[failed] incorrect protocol name: %+v", codec.ProtocolName())
				t.FailNow()
			}
		})
	}
}

func Test_Lab2_TaskB_ProtocolMatch(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		want api.MatchResult
	}{
		{
			name: "",
			data: []byte{},
			want: api.MatchAgain,
		},
		{
			name: "notenough",
			data: []byte{},
			want: api.MatchAgain,
		},
		{
			name: "",
			data: []byte("12345shangshandalaohulaohumeidazhaodadaoxiaosongshu"),
			want: api.MatchFailed,
		},
		{
			name: "",
			data: []byte("00000328RQ0tPK6UhVeIHb2hrsedxXMJHw         010005010<timestamp>1648811583</timestamp><serial_no>12345</serial_no><currency>2</currency><amount>100</amount><unit>0</unit><out_bank_id>2</out_bank_id><out_account_id>1234567899321</out_account_id><in_bank_id>2</in_bank_id><in_account_id>3211541298661</in_account_id><notes></notes>"),
			want: api.MatchSuccess,
		},
		{
			name: "",
			data: []byte("00000156RS0665db818fa5ef08e9f10ec77d76b9a0e010005010<timestamp>1648811583</timestamp><serial_no>12345</serial_no><err_code>0</err_code><message>ok</message>"),
			want: api.MatchSuccess,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := simple.CodecSimple{}
			got := c.ProtocolMatch()(tt.data)
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("[failed] incorrect protocol match")
				t.FailNow()
			}
		})
	}
}

func Test_Lab2_TaskC_Encode(t *testing.T) {
	type args struct {
		ctx   context.Context
		model interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    api.IoBuffer
		wantErr bool
	}{
		{
			name: "",
			args: args{
				ctx: context.TODO(),
				model: simple2.Request{
					Header: simple2.Header{
						TotalLength: 328,
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
			},
			want: buffer.NewIoBufferBytes([]byte("00000328RQ0tPK6UhVeIHb2hrsedxXMJHw         010005010<timestamp>1648811583</timestamp><serial_no>12345</serial_no><currency>2</currency><amount>100</amount><unit>0</unit><out_bank_id>2</out_bank_id><out_account_id>1234567899321</out_account_id><in_bank_id>2</in_bank_id><in_account_id>3211541298661</in_account_id><notes></notes>")),
		},
		{
			name: "",
			args: args{
				ctx: context.TODO(),
				model: simple2.Response{
					Header: simple2.Header{
						TotalLength: 156,
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
			},
			want: buffer.NewIoBufferBytes([]byte("00000156RS0665db818fa5ef08e9f10ec77d76b9a0e010005010<timestamp>1648811583</timestamp><serial_no>12345</serial_no><err_code>0</err_code><message>ok</message>")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			X := simple.XProtocolSimple{}
			got, err := X.Encode(tt.args.ctx, tt.args.model)
			if !assert.Nil(t, err) {
				t.Errorf("[failed] Encode() error = %v", err)
				t.FailNow()
			}
			if !assert.Equal(t, tt.want.String(), got.String()) {
				t.Errorf("[failed] value mismatch, got: %+v, wanted: %+v", tt.want.String(), got.String())
				t.FailNow()
			}
		})
	}
}

func Test_Lab2_TaskC_Decode(t *testing.T) {
	type args struct {
		ctx  context.Context
		data api.IoBuffer
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "",
			args: args{
				ctx:  context.TODO(),
				data: buffer.NewIoBufferBytes([]byte("00000328RQ0tPK6UhVeIHb2hrsedxXMJHw         010005010<timestamp>1648811583</timestamp><serial_no>12345</serial_no><currency>2</currency><amount>100</amount><unit>0</unit><out_bank_id>2</out_bank_id><out_account_id>1234567899321</out_account_id><in_bank_id>2</in_bank_id><in_account_id>3211541298661</in_account_id><notes></notes>")),
			},
			want: simple2.Request{
				Header: simple2.Header{
					TotalLength: 328,
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
			wantErr: false,
		},
		{
			name: "",
			args: args{
				ctx:  context.TODO(),
				data: buffer.NewIoBufferBytes([]byte("00000156RS0665db818fa5ef08e9f10ec77d76b9a0e010005010<timestamp>1648811583</timestamp><serial_no>12345</serial_no><err_code>0</err_code><message>ok</message>")),
			},
			want: simple2.Response{
				Header: simple2.Header{
					TotalLength: 156,
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
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			X := simple.XProtocolSimple{}
			got, err := X.Decode(tt.args.ctx, tt.args.data)
			if !assert.Nil(t, err) {
				t.Errorf("[failed] Decode() error = %v", err)
				t.FailNow()
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("[failed] value mismatch, got: %+v, wanted: %+v", tt.want, got)
				t.FailNow()
			}
		})
	}
}
