package types

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const TIME_LAYOUT = "2006-01-02 15:04:05"

func Ptr[T interface{}](t T) *T {
	return &t
}

type ExString string
type ExTime time.Time

func (e *ExString) ToString() string {
	return string(*e)
}

func (e *ExString) ToTime() *time.Time {
	t, err := time.Parse(TIME_LAYOUT, string(*e))
	if err != nil {
		return nil
	}
	return &t
}

func (e *ExTime) ToTime() time.Time {
	return time.Time(*e)
}

func (e *ExTime) ToString() string {
	return time.Time(*e).Format(TIME_LAYOUT)
}

func NStr(n pgtype.Text) string {
	if !n.Valid {
		return ""
	}
	return n.String
}

func PNStr(n pgtype.Text) *string {
	if !n.Valid {
		return nil
	}
	return &n.String
}

func NInt16(n pgtype.Int2) int16 {
	if !n.Valid {
		return 0
	}
	return n.Int16
}

func PNInt16(n pgtype.Int2) *int16 {
	if !n.Valid {
		return nil
	}
	return &n.Int16
}

func NInt32(n pgtype.Int4) int32 {
	if !n.Valid {
		return 0
	}
	return n.Int32
}

func PNInt32(n pgtype.Int4) *int32 {
	if !n.Valid {
		return nil
	}
	return &n.Int32
}

func NInt64(n pgtype.Int8) int64 {
	if !n.Valid {
		return 0
	}
	return n.Int64
}

func PNInt64(n pgtype.Int8) *int64 {
	if !n.Valid {
		return nil
	}
	return &n.Int64
}

func NFloat32(n pgtype.Float4) float32 {
	if !n.Valid {
		return 0
	}
	return n.Float32
}

func PNFloat32(n pgtype.Float4) *float32 {
	if !n.Valid {
		return nil
	}
	return &n.Float32
}

func Float32N(f *float32) pgtype.Float4 {
	if f == nil {
		return pgtype.Float4{
			Valid: false,
		}
	}
	return pgtype.Float4{
		Valid:   true,
		Float32: *f,
	}
}

func NFloat64(n pgtype.Float8) float64 {
	if !n.Valid {
		return 0
	}
	return n.Float64
}

func PNFloat64(n pgtype.Float8) *float64 {
	if !n.Valid {
		return nil
	}
	return &n.Float64
}

func Float64N(f *float64) pgtype.Float8 {
	if f == nil {
		return pgtype.Float8{
			Valid: false,
		}
	}
	return pgtype.Float8{
		Valid:   true,
		Float64: *f,
	}
}

func NTime(n pgtype.Time) *time.Time {
	if !n.Valid {
		return nil
	}
	// get time.Time from microseconds
	t := time.Unix(0, n.Microseconds*1000)
	return &t
}

func NBool(n pgtype.Bool) bool {
	if !n.Valid {
		return false
	}
	return n.Bool
}

func BoolN(b *bool) pgtype.Bool {
	if b == nil {
		return pgtype.Bool{
			Valid: false,
		}
	}
	return pgtype.Bool{
		Valid: true,
		Bool:  *b,
	}
}

func PNBool(n pgtype.Bool) *bool {
	if !n.Valid {
		return nil
	}
	return &n.Bool
}

// func NUUID(n pgtype.Text) uuid.UUID {
// 	if !n.Valid {
// 		return uuid.New()
// 	}
// 	return uuid.MustParse(n.String)
// }

func NUUID(n pgtype.UUID) uuid.UUID {
	if !n.Valid {
		return uuid.Nil
	}
	return n.Bytes
}

func NStrTime(n pgtype.Text) *time.Time {
	if !n.Valid {
		return nil
	}
	t, _ := time.Parse(TIME_LAYOUT, n.String)
	return &t
}

func UUIDStrN(s uuid.UUID) pgtype.Text {
	return pgtype.Text{
		Valid:  true,
		String: s.String(),
	}
}

func UUIDN(s uuid.UUID) pgtype.UUID {
	nu := pgtype.UUID{
		Valid: s != uuid.Nil,
	}
	if s != uuid.Nil {
		nu.Bytes = s
	}
	return nu
}

func PNUUID(u pgtype.UUID) uuid.UUID {
	if !u.Valid {
		return uuid.Nil
	}
	val, err := uuid.FromBytes(u.Bytes[:])
	if err != nil {
		return uuid.Nil
	}
	return val
}

func StrN(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{
			Valid: false,
		}
	}
	return pgtype.Text{
		Valid:  true,
		String: *s,
	}
}

func Int16N(s *int16) pgtype.Int2 {
	if s == nil {
		return pgtype.Int2{
			Valid: false,
		}
	}
	return pgtype.Int2{
		Valid: true,
		Int16: *s,
	}
}

func Int32N(s *int32) pgtype.Int4 {
	if s == nil {
		return pgtype.Int4{
			Valid: false,
		}
	}
	return pgtype.Int4{
		Valid: true,
		Int32: *s,
	}
}

func Int64N(s *int64) pgtype.Int8 {
	if s == nil {
		return pgtype.Int8{
			Valid: false,
		}
	}
	return pgtype.Int8{
		Valid: true,
		Int64: *s,
	}
}

func TimeN(s *time.Time) pgtype.Time {
	if s == nil {
		return pgtype.Time{
			Valid: false,
		}
	}
	return pgtype.Time{
		Valid:        true,
		Microseconds: s.UnixNano() / 1000,
	}
}

func TimeNStr(s *time.Time) pgtype.Text {
	if s == nil {
		return pgtype.Text{
			Valid: false,
		}
	}
	return pgtype.Text{
		Valid:  true,
		String: s.Format(TIME_LAYOUT),
	}
}

func GetString(v interface{}) *string {
	if str, ok := v.(string); ok {
		return &str
	}
	return nil
}
