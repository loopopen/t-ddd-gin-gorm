package infra

import (
	"database/sql"
	"time"

	"github.com/loopopen/pkg/timex"
)

type Mapper struct{}

func (Mapper) SQLNullToString(x sql.Null[string]) string {
	if x.Valid {
		return x.V
	}
	return ""
}

func (Mapper) StringToSQLNull(x string) sql.Null[string] {
	return sql.Null[string]{
		V:     x,
		Valid: x != "",
	}
}

func (Mapper) SQLNullToUintPtr(x sql.Null[uint]) *uint {
	if x.Valid {
		return &x.V
	}
	return nil
}

func (Mapper) UintPtrToSQLNull(x *uint) sql.Null[uint] {
	if x == nil {
		return sql.Null[uint]{}
	}
	return sql.Null[uint]{
		V:     *x,
		Valid: true,
	}
}

func (Mapper) DateTimeToTime(dt timex.DateTime) time.Time {
	return dt.Time
}

func (Mapper) TimeToDateTime(t time.Time) timex.DateTime {
	return timex.TimeToDateTime(t)
}

func (Mapper) DateToTime(d timex.Date) time.Time {
	return d.Time
}

func (Mapper) TimeToDate(t time.Time) timex.Date {
	return timex.TimeToDate(t)
}
