package random

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

var r *rand.Rand

const (
	ALPHABET     = "abcdefghijklmnopqrstuvwxyz"
	ALPHANUMERIC = "abcdefghijklmnopqrstuvwxyz0123456789"
	NUMERIC      = "0123456789"
)

func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func RandomInt32(min, max int32) int32 {
	return min + r.Int31n(max-min+1)
}

func RandomInt64(min, max int64) int64 {
	return min + r.Int63n(max-min+1)
}

func RandomFloat32(min, max float32) float32 {
	return min + r.Float32()*(max-min)
}

func RandomFloat64(min, max float64) float64 {
	return min + r.Float64()*(max-min)
}

func randomStr(n int, alphabet string) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[r.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomAlphabetStr(n int) string {
	return randomStr(n, ALPHABET)
}

func RandomAlphanumericStr(n int) string {
	return randomStr(n, ALPHANUMERIC)
}

func RandomNumericStr(n int) string {
	return randomStr(n, NUMERIC)
}

func RandomIPv4Address() string {
	return fmt.Sprintf("%d.%d.%d.%d", RandomInt64(0, 255), RandomInt64(0, 255), RandomInt64(0, 255), RandomInt64(0, 255))
}

func RandomEmail() string {
	return fmt.Sprintf("%s@gmail.com", RandomAlphabetStr(8))
}

func RandomPhoneNumber() string {
	return fmt.Sprintf("0%d", RandomInt64(100000000, 999999999))
}

func RandomDate() time.Time {
	return time.Date(
		int(RandomInt64(1970, 2020)),
		time.Month(RandomInt64(1, 12)),
		int(RandomInt64(1, 28)),
		int(RandomInt64(0, 23)),
		int(RandomInt64(0, 59)),
		int(RandomInt64(0, 59)),
		int(RandomInt64(0, 999999999)),
		time.UTC,
	)
}

func RandomURL() string {
	return fmt.Sprintf("https://%s.com", RandomAlphabetStr(8))
}

// Shuffle slice using Fisher-Yates algorithm. Modifies the original slice.
func ShuffleSlice[T any](slice []T) {
	for i := len(slice) - 1; i > 0; i-- {
		j := r.Intn(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
}

func RandomlyPickNFromSlice[T any](slice []T, n int) []T {
	s := make([]T, len(slice))
	copy(s, slice)
	ShuffleSlice(s)
	return s[:n]
}
