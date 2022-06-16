package models

import (
	"encoding/json"
	"fmt"
	"time"
)

type ErrInvalidDateFormat struct {
	field string
}

func NewErrInvalidDateFormat(field string) *ErrInvalidDateFormat {
	return &ErrInvalidDateFormat{field}
}

func (e *ErrInvalidDateFormat) Error() string {
	return fmt.Sprintf("%s invalid date format - must be YYYY-MM-DD", e.field)
}

var ErrTotalDueMustBeGreaterThenZero error = fmt.Errorf("total_due must be greater than 0")
var ErrReducedAmountMustBeGreaterThenZero error = fmt.Errorf("reduced_amount must be greater than 0")

type Ticket struct {
	Sref                  string             `json:"sref" binding:"required"`
	NoticeNumber          string             `json:"notice_number" binding:"required"`
	VehicleRegistration   string             `json:"vehicle_registration" binding:"required"`
	Contravention         string             `json:"contravention" binding:"required"`
	ContraventionDatetime *time.Time         `json:"contravention_datetime" binding:"required"`
	EntryExitDatetime     *EntryExitDatetime `json:"entry_exit_datetime"`
	ObservationDatetime   *DatetimeRange     `json:"observation_datetime"`
	Location              string             `json:"location" binding:"required"`
	NoticeToKeeper        NoticeToKeeper     `json:"notice_to_keeper" binding:"required"`
	Profa                 bool               `json:"profa"`
	TotalDue              int                `json:"total_due"`
	ReducedAmount         int                `json:"reduced_amount"`
	ReducePeriodEnds      *time.Time         `json:"reduce_period_ends"`
	Photos                Photos             `json:"photos" binding:"required"`
	PaymentUrl            string             `json:"payment_url"`
	AppealUrl             string             `json:"appeal_url_url"`
}

func (t *Ticket) UnmarshalJSON(data []byte) error {
	type Alias Ticket

	aux := struct {
		ContraventionDatetime string `json:"contravention_datetime" binding:"required"`
		ReducePeriodEnds      string `json:"reduce_period_ends"`

		*Alias
	}{
		Alias: (*Alias)(t),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// this is valid IANA timezone string (https://www.iana.org/time-zones), so no need to handle an error
	loc, _ := time.LoadLocation("Europe/London")
	cpt, err := time.Parse(time.RFC3339Nano, aux.ContraventionDatetime)
	if err != nil {
		return NewErrInvalidDatetimeFormat("contravention_datetime")
	}
	contraventionDatetime := cpt.In(loc)
	t.ContraventionDatetime = &contraventionDatetime

	rpept, err := time.ParseInLocation("2006-01-02", aux.ReducePeriodEnds, loc)
	if err != nil {
		return NewErrInvalidDateFormat("reduce_period_ends")
	}
	reducedPeriodEnds := rpept.In(loc)
	t.ReducePeriodEnds = &reducedPeriodEnds

	if t.TotalDue <= 0 {
		return ErrTotalDueMustBeGreaterThenZero
	}

	if t.ReducedAmount <= 0 {
		return ErrReducedAmountMustBeGreaterThenZero
	}

	return nil
}
