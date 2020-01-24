package uuid

import (
	"encoding/hex"
	"github.com/segmentio/ksuid"
)

func CreateUuidV4() string {

	uuid := ksuid.New()
	payload := formatUUID(uuid.Bytes())

	return payload
}

func formatUUID(u []byte) string {

	buf := make([]byte, 36)

	u[6] = (u[6] & 0x0f) | 0x40 // Version 4
	u[8] = (u[8] & 0x3f) | 0x80 // Variant is 10

	hex.Encode(buf[0:8], u[0:4])
	buf[8] = '-'
	hex.Encode(buf[9:13], u[4:6])
	buf[13] = '-'
	hex.Encode(buf[14:18], u[6:8])
	buf[18] = '-'
	hex.Encode(buf[19:23], u[8:10])
	buf[23] = '-'
	hex.Encode(buf[24:], u[10:16])

	return string(buf)
}
