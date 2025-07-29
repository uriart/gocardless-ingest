package models

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type Date struct {
	time.Time
}

func (d *Date) UnmarshalJSON(b []byte) error {
	s := string(b)
	if s == "null" {
		return nil
	}
	s = s[1 : len(s)-1]

	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return fmt.Errorf("error parsing date: %w", err)
	}
	d.Time = t
	return nil
}

func (d Date) Value() (driver.Value, error) {
	return d.Time, nil
}

func (d Date) MarshalJSON() ([]byte, error) {
	s := fmt.Sprintf("\"%s\"", d.Time.Format("2006-01-02"))
	return []byte(s), nil
}
