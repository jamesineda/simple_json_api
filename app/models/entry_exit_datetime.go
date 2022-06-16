package models

import (
	"encoding/json"
	"time"
)

type EntryExitDatetime struct {
	Entry *time.Time `json:"entry"`
	Exit  *time.Time `json:"exit"`
}

func (eed *EntryExitDatetime) UnmarshalJSON(data []byte) error {
	type Alias EntryExitDatetime

	aux := struct {
		Entry string `json:"entry"`
		Exit  string `json:"exit"`
		*Alias
	}{
		Alias: (*Alias)(eed),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// this is valid IANA timezone string (https://www.iana.org/time-zones), so no need to handle an error
	loc, _ := time.LoadLocation("Europe/London")
	ept, err := time.ParseInLocation(time.RFC3339Nano, aux.Entry, loc)
	if err != nil {
		return NewErrInvalidDatetimeFormat("entry")
	}
	entry := ept.In(loc)
	eed.Entry = &entry

	expt, err := time.ParseInLocation(time.RFC3339Nano, aux.Exit, loc)
	if err != nil {
		return NewErrInvalidDatetimeFormat("exit")
	}
	exit := expt.In(loc)
	eed.Exit = &exit

	return nil
}
