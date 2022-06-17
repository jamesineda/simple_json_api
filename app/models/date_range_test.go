package models

import (
	"encoding/json"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type DateRangeTestSuite struct {
	suite.Suite
}

func (suite *DateRangeTestSuite) SetupTest() {}

func (suite *DateRangeTestSuite) Test_Unmarshall() {
	suite.Run("When UTC rfc3339 nano", func() {
		ts := `{
					"from": "2022-06-13T19:23:39.796Z",
					"to": "2022-06-13T20:01:19.896Z"
				}`
		dr := DatetimeRange{}
		err := json.Unmarshal([]byte(ts), &dr)
		suite.NoErrorf(err, "")
		suite.Equal("2022-06-13T20:23:39+01:00", dr.From.Format(time.RFC3339))
		suite.Equal("2022-06-13T21:01:19+01:00", dr.To.Format(time.RFC3339))
	})

	suite.Run("When not rfc3339 nano", func() {
		ts := `{
					"from": "2022-06-13 19:23:39",
					"to": "2022-06-13 20:01:19"
				}`
		dr := DatetimeRange{}
		err := json.Unmarshal([]byte(ts), &dr)
		suite.Error(err, ErrInvalidDatetimeFormat{})
	})
}

func TestDateRangeTestSuite(t *testing.T) {
	suite.Run(t, new(DateRangeTestSuite))
}
