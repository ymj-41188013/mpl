package simple

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
)

type Header struct {
	TotalLength int
	Type        string
	PageMark    int
	Checksum    string
	ServiceCode int
	Reserved    int
}

func (h *Header) Encode(ctx context.Context) ([]byte, error) {
	return []byte(fmt.Sprintf("%08d%s%d%32s%08d%d", h.TotalLength, h.Type, h.PageMark, h.Checksum, h.ServiceCode, h.Reserved)), nil
}

func (h *Header) Decode(ctx context.Context, data []byte) error {
	// todo lab1-task-a
	panic("implement me")
}

type Request struct {
	Header
	UnixTimestamp int    `xml:"timestamp"`
	SerialNo      int    `xml:"serial_no"`
	Currency      int    `xml:"currency"`
	Amount        int    `xml:"amount"`
	Unit          int    `xml:"unit"`
	OutBankId     int    `xml:"out_bank_id"`
	OutAccountId  int    `xml:"out_account_id"`
	InBankId      int    `xml:"in_bank_id"`
	InAccountId   int    `xml:"in_account_id"`
	Notes         string `xml:"notes"`
}

func (r *Request) Encode(ctx context.Context) ([]byte, error) {
	// todo lab1-task-b
	panic("implement me")
}

func (r *Request) Decode(ctx context.Context, data []byte) error {
	// todo lab1-task-b
	panic("implement me")
}

type Response struct {
	Header
	UnixTimestamp int64  `xml:"timestamp"`
	SerialNo      int    `xml:"serial_no"`
	ErrCode       int    `xml:"err_code"`
	Message       string `xml:"message"`
}

func (r *Response) Encode(ctx context.Context) ([]byte, error) {
	// todo lab1-task-b
	panic("implement me")
}

func (r *Response) Decode(ctx context.Context, data []byte) error {
	if len(data) < 52 {
		return errors.New("incorrect response data length")
	}

	if err := r.Header.Decode(ctx, data[:52]); err != nil {
		return err
	}

	if len(data) < r.TotalLength {
		return errors.New("incorrect response data length")
	}

	var xmlData []byte
	xmlData = append([]byte("<response>"), data[52:]...)
	xmlData = append(xmlData, []byte("</response>")...)

	if err := xml.Unmarshal(xmlData, r); err != nil {
		return err
	}

	return nil
}
