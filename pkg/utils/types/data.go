package types

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
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

func NStr(n sql.NullString) string {
	if !n.Valid {
		return ""
	}
	return n.String
}

func PNStr(n sql.NullString) *string {
	if !n.Valid {
		return nil
	}
	return &n.String
}

func NInt16(n sql.NullInt16) int16 {
	if !n.Valid {
		return 0
	}
	return n.Int16
}

func PNInt16(n sql.NullInt16) *int16 {
	if !n.Valid {
		return nil
	}
	return &n.Int16
}

func NInt32(n sql.NullInt32) int32 {
	if !n.Valid {
		return 0
	}
	return n.Int32
}

func PNInt32(n sql.NullInt32) *int32 {
	if !n.Valid {
		return nil
	}
	return &n.Int32
}

func NInt64(n sql.NullInt64) int64 {
	if !n.Valid {
		return 0
	}
	return n.Int64
}

func PNInt64(n sql.NullInt64) *int64 {
	if !n.Valid {
		return nil
	}
	return &n.Int64
}

func NFloat64(n sql.NullFloat64) float64 {
	if !n.Valid {
		return 0
	}
	return n.Float64
}

func PNFloat64(n sql.NullFloat64) *float64 {
	if !n.Valid {
		return nil
	}
	return &n.Float64
}

func NTime(n sql.NullTime) *time.Time {
	if !n.Valid {
		return nil
	}
	return &n.Time
}

func NBool(n sql.NullBool) bool {
	if !n.Valid {
		return false
	}
	return n.Bool
}

func BoolN(b *bool) sql.NullBool {
	if b == nil {
		return sql.NullBool{
			Valid: false,
		}
	}
	return sql.NullBool{
		Valid: true,
		Bool:  *b,
	}
}

func NUUID(n sql.NullString) uuid.UUID {
	if !n.Valid {
		return uuid.New()
	}
	return uuid.MustParse(n.String)
}

func NStrTime(n sql.NullString) *time.Time {
	if !n.Valid {
		return nil
	}
	t, _ := time.Parse(TIME_LAYOUT, n.String)
	return &t
}

func UUIDStrN(s uuid.UUID) sql.NullString {
	return sql.NullString{
		Valid:  true,
		String: s.String(),
	}
}

func UUIDN(s *uuid.UUID) uuid.NullUUID {
	nu := uuid.NullUUID{
		Valid: s != nil,
	}
	if s != nil {
		nu.UUID = *s
	}
	return nu
}

func PNUUID(u uuid.NullUUID) *uuid.UUID {
	if !u.Valid {
		return nil
	}
	return &u.UUID
}

func StrN(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{
			Valid: false,
		}
	}
	return sql.NullString{
		Valid:  true,
		String: *s,
	}
}

func Int16N(s *int16) sql.NullInt16 {
	if s == nil {
		return sql.NullInt16{
			Valid: false,
		}
	}
	return sql.NullInt16{
		Valid: true,
		Int16: *s,
	}
}

func Int32N(s *int32) sql.NullInt32 {
	if s == nil {
		return sql.NullInt32{
			Valid: false,
		}
	}
	return sql.NullInt32{
		Valid: true,
		Int32: *s,
	}
}

func Int64N(s *int64) sql.NullInt64 {
	if s == nil {
		return sql.NullInt64{
			Valid: false,
		}
	}
	return sql.NullInt64{
		Valid: true,
		Int64: *s,
	}
}

func TimeN(s *time.Time) sql.NullTime {
	if s == nil {
		return sql.NullTime{
			Valid: false,
		}
	}
	return sql.NullTime{
		Valid: true,
		Time:  *s,
	}
}

func TimeNStr(s *time.Time) sql.NullString {
	if s == nil {
		return sql.NullString{
			Valid: false,
		}
	}
	return sql.NullString{
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
