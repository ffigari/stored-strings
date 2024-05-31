package webapi_test

import (
	"time"

	"github.com/ffigari/stored-strings/internal/webapi"
)

func (s *AtOfServeSuite) TestParseFormDatetime() {
	location, err := time.LoadLocation(
		"America/Argentina/Buenos_Aires",
	)
	s.Require().NoError(err)

	s.Run("ok", func() {
		ts, err := webapi.ParseFormDatetime("2024-06-20T09:02", location)
		s.Require().NoError(err)
		s.Require().Equal(2024, ts.Year())
		s.Require().Equal(time.Month(6), ts.Month())
		s.Require().Equal(20, ts.Day())
		s.Require().Equal(12, ts.Hour())
		s.Require().Equal(2, ts.Minute())
	})

	s.Run("wrong format fails to be parsed", func() {
		_, err := webapi.ParseFormDatetime("2024_06_20T09:02", location)
		s.Require().ErrorContains(err, "failed to match")
	})
}
