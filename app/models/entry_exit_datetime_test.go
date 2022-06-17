package models

import (
	"encoding/json"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type EntryExitDatetimeTestSuite struct {
	suite.Suite
}

func (suite *EntryExitDatetimeTestSuite) SetupTest() {}

func (suite *EntryExitDatetimeTestSuite) Test_Unmarshall() {
	suite.Run("When UTC rfc3339 nano", func() {
		ts := `{
					"exit": "2022-06-13T19:23:39.796Z",
					"entry": "2022-06-13T20:01:19.896Z"
				}`
		dr := EntryExitDatetime{}
		err := json.Unmarshal([]byte(ts), &dr)
		suite.NoErrorf(err, "")
		suite.Equal("2022-06-13T20:23:39+01:00", dr.Exit.Format(time.RFC3339))
		suite.Equal("2022-06-13T21:01:19+01:00", dr.Entry.Format(time.RFC3339))
	})

	suite.Run("When not rfc3339 nano", func() {
		ts := `{
					"exit": "2022-06-13 19:23:39",
					"entry": "2022-06-13 20:01:19"
				}`
		dr := EntryExitDatetime{}
		err := json.Unmarshal([]byte(ts), &dr)
		suite.Error(err, ErrInvalidDatetimeFormat{})
	})
}

func TestEntryExitDatetimeTestSuite(t *testing.T) {
	suite.Run(t, new(EntryExitDatetimeTestSuite))
}
