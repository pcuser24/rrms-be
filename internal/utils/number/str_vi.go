package number

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

var (
	digits    = []string{" không ", " một ", " hai ", " ba ", " bốn ", " năm ", " sáu ", " bảy ", " tám ", " chín "}
	thousands = []string{"", " nghìn", " triệu", " tỷ", " nghìn tỷ", " triệu tỷ"}
)

func read3digit(n int) string {
	var hundreds, tens, units int
	var res string

	hundreds = n / 100
	tens = (n % 100) / 10
	units = n % 10

	if hundreds == 0 && tens == 0 && units == 0 {
		return ""
	}

	if hundreds != 0 {
		res += digits[hundreds] + " trăm "
		if tens == 0 && units != 0 {
			res += " linh "
		}
	}

	if tens != 0 && tens != 1 {
		res += digits[tens] + " mươi"
		if tens == 0 && units != 0 {
			res += " linh "
		}
	}

	if tens == 1 {
		res += " mười "
	}

	switch units {
	case 1:
		if tens != 0 && tens != 1 {
			res += " mốt "
		} else {
			res += digits[units]
		}
	case 5:
		if tens == 0 {
			res += digits[units]
		} else {
			res += " lăm "
		}
	default:
		if units != 0 {
			res += digits[units]
		}
	}

	return res
}

var (
	ErrTooLarge = errors.New("số quá lớn")
)

func ToStr(n int64) (string, error) {
	var times, _n, i int
	var res, tmp string
	var pos [6]int

	if n == 0 {
		return "Không", nil
	}

	if n > 8999999999999999 {
		return "", ErrTooLarge
	}

	if n > 0 {
		_n = int(n)
	} else {
		_n = int(-n)
	}

	pos[5] = int(math.Floor(float64(_n) / 1000000000000000))
	_n -= pos[5] * 1000000000000000
	pos[4] = int(math.Floor(float64(_n) / 1000000000000))
	_n -= pos[4] * 1000000000000
	pos[3] = int(math.Floor(float64(_n) / 1000000000))
	_n -= pos[3] * 1000000000
	pos[2] = int(_n / 1000000)
	pos[1] = int((_n % 1000000) / 1000)
	pos[0] = int(_n % 1000)

	if pos[5] > 0 {
		times = 5
	} else if pos[4] > 0 {
		times = 4
	} else if pos[3] > 0 {
		times = 3
	} else if pos[2] > 0 {
		times = 2
	} else if pos[1] > 0 {
		times = 1
	} else {
		times = 0
	}

	for i = times; i >= 0; i-- {
		tmp = read3digit(pos[i])
		res += tmp
		if pos[i] > 0 {
			res += thousands[i]
			// if i > 0 && len(tmp) > 0 {
			// 	KetQua += ","
			// }
		}
	}

	res = strings.ToUpper(res[:1]) + res[1:]

	// Post processing
	res = strings.Join(strings.Fields(res), " ")
	res = fmt.Sprintf("%s%s", strings.ToUpper(res[:1]), res[1:])

	if n < 0 {
		return "Âm " + res, nil
	}
	return res, nil
}
