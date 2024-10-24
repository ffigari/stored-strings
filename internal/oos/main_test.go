package oos_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/ffigari/stored-strings/internal/basesuite"
	pkg "github.com/ffigari/stored-strings/internal/oos"
)

//go:generate mockgen -package=mocks -source=interfaces.go -destination=mocks/main.go

func (s *S) TestMain() {
	s.Run("ok", func() {
		files, err := pkg.ReadFiles("/test_dir")
		s.Require().NoError(err)
		s.Require().Equal(2, len(files))
		s.Require().Equal("file1", files[0].Name())
		s.Require().Equal("content1\n", string(files[0].Content()))
		s.Require().Equal("file2", files[1].Name())
		s.Require().Equal("content2\n", string(files[1].Content()))
	})

}

type S struct {
	basesuite.BaseSuite
}

func TestFinancesSuite(t *testing.T) {
	suite.Run(t, new(S))
}
