package phpcrudapi

import (
	"strings"
	"time"
)

var dbtimeformat = "2006-01-02 15:04:05.999999"

type Time time.Time

func (t Time) Time() time.Time {
	return time.Time(t)
}

func (t Time) MarshalJSON() ([]byte, error) {
	return []byte(time.Time(t).UTC().Format(dbtimeformat)[:]), nil
}

func (t *Time) UnmarshalJSON(data []byte) error {
	input := strings.Trim(string(data), `"`)
	pt, err := time.Parse(dbtimeformat, input)
	if err != nil {
		return err
	}
	*t = Time(pt)
	return nil
}

var dbdateformat = "2006-01-02"

type Date time.Time

func (d Date) Time() time.Time {
	return time.Time(d)
}

func (d Date) MarshalJSON() ([]byte, error) {
	return []byte(time.Time(d).UTC().Format(dbdateformat)[:]), nil
}

func (d *Date) UnmarshalJSON(data []byte) error {
	input := strings.Trim(string(data), `"`)
	pd, err := time.Parse(dbdateformat, input)
	if err != nil {
		return err
	}
	*d = Date(pd)
	return nil
}
