package main

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
)

type Header struct {
	TRN []byte
	LEN []byte
	OR  []byte
	OT  []byte
}

type Session struct {
	ODAC []byte
	OTON []byte
	ONPI []byte
	STYP []byte
	PWD  []byte
	NPWD []byte
	VERS []byte
	LADC []byte
	LTON []byte
	LNPI []byte
	OPID []byte
	RES1 []byte
}

type AuthError struct {
}

func (a AuthError) Error() string {
	return "could not authenticate"
}

// stx <header> / <data> / <checksum> etx
func createLoginReq(refNum []byte, username string, pw string) []byte {
	packet := make([]byte, 0)
	packet = append(packet, STX)

	encodedPassword := fmt.Sprintf("%02X", pw)

	s := Session{
		ODAC: []byte(username),
		OTON: []byte(ORIGINATOR_TYPE_OF_NUMBER),
		PWD:  []byte(encodedPassword),
		VERS: []byte(VERS),
		ONPI: []byte(SMS_SPECIFIC),
		STYP: []byte(OPEN_SESSION),
	}

	t := reflect.TypeOf(s)
	values := reflect.ValueOf(s)
	num := t.NumField()

	fields := make([][]byte, num)
	for i := 0; i < num; i++ {
		value := values.Field(i)

		fields[i] = value.Bytes()
	}

	joined := bytes.Join(fields, []byte(DELIMITER))

	//HEADER + DATA
	part := [][]byte{
		refNum,
		[]byte(fmt.Sprintf("%05d", len(joined))),
		[]byte(OPERATION_TYPE),
		[]byte(SESSION_MANAGEMENT),
		joined,
	}

	withoutChecksum := append(bytes.Join(part, []byte(DELIMITER)), []byte(DELIMITER)...)
	withChecksum := append(withoutChecksum, checksum(withoutChecksum)...)

	packet = append(packet, withChecksum...)
	packet = append(packet, ETX)

	return packet
}

func checksum(b []byte) []byte {
	var sum byte

	for _, bt := range b {
		sum += bt
	}

	mask := sum & 255

	checksm := fmt.Sprintf("%02X", mask)

	return []byte(checksm)

}

// 00/00037/N/60/A/AUTHENTICATION FAILURE  /
func parseSessionResp(r string) error {
	//basically just look for a ACK or NEGATIVE_RESULT
	split := strings.Split(r, "/")

	if split[2] == "N" {
		err := AuthError{}
		return err
	}

	return nil
}
