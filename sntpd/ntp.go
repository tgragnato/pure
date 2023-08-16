package main

import "time"

const (
	LI_NO_WARNING      = 0
	LI_ALARM_CONDITION = 3
	VN_FIRST           = 1
	VN_LAST            = 4
	MODE_CLIENT        = 3
	FROM_1900_TO_1970  = 2208988800
)

func validFormat(request []byte) bool {
	var l = request[0] >> 6
	var v = (request[0] << 2) >> 5
	var m = (request[0] << 5) >> 5
	if (l == LI_NO_WARNING) || (l == LI_ALARM_CONDITION) {
		if (v >= VN_FIRST) && (v <= VN_LAST) {
			if m == MODE_CLIENT {
				return true
			}
		}
	}
	return false
}

func int2bytes(i int64) []byte {
	b := make([]byte, 4)
	h1 := i >> 24
	h2 := (i >> 16) - (h1 << 8)
	h3 := (i >> 8) - (h1 << 16) - (h2 << 8)
	h4 := byte(i)
	b[0] = byte(h1)
	b[1] = byte(h2)
	b[2] = byte(h3)
	b[3] = byte(h4)
	return b
}

func generate(request []byte) []byte {
	second := time.Now().Unix() + FROM_1900_TO_1970
	fraction := int64(time.Now().Nanosecond()) + FROM_1900_TO_1970
	response := make([]byte, 48)
	response[0] = (request[0] & 0x38) + 4
	response[1] = 1
	response[2] = request[2]
	response[3] = 0xEC
	response[12] = 0x4E
	response[13] = 0x49
	response[14] = 0x43
	response[15] = 0x54
	copy(response[16:20], int2bytes(second)[0:])
	copy(response[24:32], request[40:48])
	copy(response[32:36], int2bytes(second)[0:])
	copy(response[36:40], int2bytes(fraction)[0:])
	copy(response[40:48], response[32:40])
	return response
}
