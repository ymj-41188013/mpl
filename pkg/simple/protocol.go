package simple

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"mosn.io/pkg/buffer"
	"strconv"
	"strings"
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
	var total = 0
	var err error
	totalLen := strings.TrimLeft(string(data[0:8]), "0")
	if totalLen != "" {
		total, err = strconv.Atoi(totalLen)
		if err != nil {
			return errors.New(fmt.Sprintf("failed to decode package len %d, err: %v", total, err))
		}
	}
	var serviceCode = 0
	serviceCodeLen := strings.TrimLeft(string(data[43:51]), "0")
	if serviceCodeLen != "" {
		serviceCode, err = strconv.Atoi(serviceCodeLen)
		if err != nil {
			return errors.New(fmt.Sprintf("failed to decode package len %d, err: %v", serviceCode, err))
		}
	}
	var pageMarkData int64
	binary.Read(bytes.NewBuffer(data[10:11]), binary.BigEndian, &pageMarkData)
	var reservedData int64
	binary.Read(bytes.NewBuffer(data[51:52]), binary.BigEndian, &reservedData)
	h.TotalLength = total
	h.Type = bytes.NewBuffer(data[8:10]).String()
	h.PageMark = int(pageMarkData)
	h.Checksum = bytes.NewBuffer(data[11:43]).String()
	h.ServiceCode = serviceCode
	h.Reserved = int(reservedData)
	return nil
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
type XmlRequest struct {
	UnixTimestamp int    `xml:"timestamp" json:"timestamp"`
	SerialNo      int    `xml:"serial_no" json:"serial_no"`
	Currency      int    `xml:"currency" json:"currency"`
	Amount        int    `xml:"amount" json:"amount"`
	Unit          int    `xml:"unit" json:"unit"`
	OutBankId     int    `xml:"out_bank_id" json:"out_bank_id"`
	OutAccountId  int    `xml:"out_account_id" json:"out_account_id"`
	InBankId      int    `xml:"in_bank_id" json:"in_bank_id"`
	InAccountId   int    `xml:"in_account_id" json:"in_account_id"`
	Notes         string `xml:"notes" json:"notes"`
}

func (r *Request) Encode(ctx context.Context) ([]byte, error) {
	buf := buffer.GetIoBuffer(r.TotalLength)
	prefixOfZero(buf, r.TotalLength)
	buf.WriteString(r.Header.Type)
	buf.WriteString(strconv.Itoa(r.Header.PageMark))
	buf.WriteString(r.Header.Checksum)
	buf.WriteString("0" + strconv.Itoa(r.Header.ServiceCode))
	buf.WriteString(strconv.Itoa(r.Header.Reserved))

	xmlRequest := XmlRequest{
		UnixTimestamp: r.UnixTimestamp,
		SerialNo:      r.SerialNo,
		Currency:      r.Currency,
		Amount:        r.Amount,
		Unit:          r.Unit,
		OutBankId:     r.OutBankId,
		OutAccountId:  r.OutAccountId,
		InBankId:      r.InBankId,
		InAccountId:   r.InAccountId,
		Notes:         r.Notes,
	}
	xml, _ := xml.Marshal(xmlRequest)
	xml = xml[len("<XmlRequest>") : len(xml)-len("</XmlRequest>")]
	buf.Write(xml)

	return buf.Bytes(), nil
}
func prefixOfZero(buf buffer.IoBuffer, payloadLen int) {
	rayLen := strconv.Itoa(payloadLen)
	if count := 8 - len(rayLen); count > 0 {
		for i := 0; i < count; i++ {
			buf.WriteString("0")
		}
	}
	buf.WriteString(rayLen)
}

func (r *Request) Decode(ctx context.Context, data []byte) error {
	var packetLen = 0
	var err error

	rawLen := strings.TrimLeft(string(data[0:8]), "0")
	if rawLen != "" {
		packetLen, err = strconv.Atoi(rawLen)
		if err != nil {
			return errors.New(fmt.Sprintf("failed to decode package len %d, err: %v", packetLen, err))
		}
	}
	var total = 0
	totalLe := strings.TrimLeft(string(data[0:8]), "0")
	if totalLe != "" {
		total, err = strconv.Atoi(totalLe)
		if err != nil {
			return errors.New(fmt.Sprintf("failed to decode package len %d, err: %v", total, err))
		}
	}
	var serviceCode = 0
	serviceCodeLen := strings.TrimLeft(string(data[43:51]), "0")
	if serviceCodeLen != "" {
		serviceCode, err = strconv.Atoi(serviceCodeLen)
		if err != nil {
			return errors.New(fmt.Sprintf("failed to decode package len %d, err: %v", serviceCode, err))
		}
	}
	var pageMarkData int64
	binary.Read(bytes.NewBuffer(data[10:11]), binary.BigEndian, &pageMarkData)
	var reservedData int64
	binary.Read(bytes.NewBuffer(data[51:52]), binary.BigEndian, &reservedData)
	header := &Header{}
	header.TotalLength = total
	header.Type = bytes.NewBuffer(data[8:10]).String()
	header.PageMark = int(pageMarkData)
	header.Checksum = bytes.NewBuffer(data[11:43]).String()
	header.ServiceCode = serviceCode
	header.Reserved = int(reservedData)
	r.Header = *header

	xh, err := parseXmlHeader(data[52:packetLen])
	if err != nil {
		return err
	}

	r.UnixTimestamp, _ = strconv.Atoi(xh["timestamp"])
	r.Currency, _ = strconv.Atoi(xh["currency"])
	r.Amount, _ = strconv.Atoi(xh["amount"])
	r.InAccountId, _ = strconv.Atoi(xh["in_account_id"])
	r.SerialNo, _ = strconv.Atoi(xh["serial_no"])
	r.InBankId, _ = strconv.Atoi(xh["in_bank_id"])
	r.OutAccountId, _ = strconv.Atoi(xh["out_account_id"])
	r.OutBankId, _ = strconv.Atoi(xh["out_bank_id"])
	r.Unit, _ = strconv.Atoi(xh["unit"])
	r.Notes, _ = xh["notes"]
	return nil
}

type XmlHeader map[string]string

type KeyValueEntry struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

func (m XmlHeader) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if len(m) == 0 {
		return nil
	}

	if err := e.EncodeToken(start); err != nil {
		return err
	}

	for k, v := range m {
		e.Encode(KeyValueEntry{XMLName: xml.Name{Local: k}, Value: v})
	}

	return e.EncodeToken(start.End())
}

func (m *XmlHeader) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	*m = XmlHeader{}
	for {
		var e KeyValueEntry

		err := d.Decode(&e)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		(*m)[e.XMLName.Local] = e.Value
	}

	return nil
}

func parseXmlHeader(data []byte) (XmlHeader, error) {
	xmlBody := string(data)
	header := XmlHeader{}
	xmlBody = "<header>" + xmlBody + "</header>"
	if xmlBody != "" {
		buf := bytes.NewBufferString(xmlBody)
		xmlDecoder := xml.NewDecoder(buf)
		if err := xmlDecoder.Decode(&header); err != nil {
			return nil, err
		}
	}
	return header, nil
}

type Response struct {
	Header
	UnixTimestamp int64  `xml:"timestamp"`
	SerialNo      int    `xml:"serial_no"`
	ErrCode       int    `xml:"err_code"`
	Message       string `xml:"message"`
}

type XmlResponse struct {
	UnixTimestamp int64  `xml:"timestamp"`
	SerialNo      int    `xml:"serial_no"`
	ErrCode       int    `xml:"err_code"`
	Message       string `xml:"message"`
}

func (r *Response) Encode(ctx context.Context) ([]byte, error) {
	// todo lab1-task-b
	buf := buffer.GetIoBuffer(r.TotalLength)
	prefixOfZero(buf, r.TotalLength)
	buf.WriteString(r.Header.Type)
	buf.WriteString(strconv.Itoa(r.Header.PageMark))
	buf.WriteString(r.Header.Checksum)
	buf.WriteString("0" + strconv.Itoa(r.Header.ServiceCode))
	buf.WriteString(strconv.Itoa(r.Header.Reserved))

	xmlResponse := XmlResponse{
		UnixTimestamp: r.UnixTimestamp,
		SerialNo:      r.SerialNo,
		ErrCode:       r.ErrCode,
		Message:       r.Message,
	}
	xml, _ := xml.Marshal(xmlResponse)
	xml = xml[len("<XmlResponse>") : len(xml)-len("</XmlResponse>")]
	buf.Write(xml)
	return buf.Bytes(), nil
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
