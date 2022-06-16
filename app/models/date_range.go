package models

import (
	"encoding/json"
	"fmt"
	"time"
)

//var ErrInvalidDatetimeFormat = fmt.Errorf("invalid datetime format, must be RFC3339 Nano")

type ErrInvalidDatetimeFormat struct {
	field string
}

func NewErrInvalidDatetimeFormat(field string) *ErrInvalidDatetimeFormat {
	return &ErrInvalidDatetimeFormat{field}
}

func (e *ErrInvalidDatetimeFormat) Error() string {
	return fmt.Sprintf("%s invalid datetime format - must be RFC3339 Nano", e.field)
}

type DatetimeRange struct {
	To   *time.Time `json:"to"`
	From *time.Time `json:"from"`
}

func (dr *DatetimeRange) UnmarshalJSON(data []byte) error {
	type Alias DatetimeRange

	aux := struct {
		To   string `json:"to"`
		From string `json:"from"`
		*Alias
	}{
		Alias: (*Alias)(dr),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// this is valid IANA timezone string (https://www.iana.org/time-zones), so no need to handle an error
	loc, _ := time.LoadLocation("Europe/London")
	tpt, err := time.ParseInLocation(time.RFC3339Nano, aux.To, loc)
	if err != nil {
		return NewErrInvalidDatetimeFormat("to")
	}
	to := tpt.In(loc)
	dr.To = &to

	fpt, err := time.ParseInLocation(time.RFC3339Nano, aux.From, loc)
	if err != nil {
		return NewErrInvalidDatetimeFormat("from")
	}
	from := fpt.In(loc)
	dr.From = &from

	return nil
}
