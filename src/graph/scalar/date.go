package scalar

import (
	"fmt"
	"io"
	"strconv"
	"time"
)

// Maps a Date GraphQL scalar to a Go time.Time struct.
// This scalar adheres to the time.RFC3339 format.
// https://pkg.go.dev/time#pkg-constants
type Date time.Time

// UnmarshalGQL implements the graphql.Unmarshaler interface
func (t *Date) UnmarshalGQL(v interface{}) error {
	dt, ok := v.(string)
	if !ok {
		return fmt.Errorf("DateTime must be a string")
	}

	godt, err := time.Parse(time.DateOnly, dt)
	if err != nil {
		return err
	}

	*t = Date(godt)
	return nil
}

// MarshalGQL implements the graphql.Marshaler interface
func (t Date) MarshalGQL(w io.Writer) {
	_, _ = w.Write([]byte(strconv.Quote(time.Time(t).Format(time.DateOnly))))
	//w.Write([]byte(strconv.Quote(time.Time(t).UTC().Format(time.RFC3339))))
}
