package sntp

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
