package finances_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	//"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"

	"github.com/ffigari/stored-strings/internal/basesuite"
	pkg "github.com/ffigari/stored-strings/internal/finances"
	//mocks "github.com/ffigari/stored-strings/internal/finances/mocks"
)

//go:generate mockgen -package=mocks -source=storage.go -destination=mocks/storage.go

//func (s *financesSuite) TestCreateCategory() {
//	s.Run("handler ok", func() {})
//	
//	s.Run("handler invalid input", func() {})
//
//	s.Run("creation ok", func() {
//		//ctrl := gomock.NewController(s.T())
//		//storage := mocks.NewMockstorageI(ctrl)
//
//		storage.EXPECT().CreateExpenseCategory()
//
//		pkg.Finances.CreateExpenseCategory(storage)
//	})
//}

func (s *financesSuite) TestFoo() {
	s.Run("ok", func() {
		holding := pkg.Holding{
			ARSAmount: 100,
			MutualFundsValuation: 1000,
			CEDEARsValuation: 10000,
			MEPUSDAmount: 500,
			MEPUSDQuote: 20.5,
		}

		s.Require().Equal(int64(21350), holding.Valuation())

		r := mux.NewRouter()
		server := httptest.NewServer(r)
		defer server.Close()

		req, err := http.NewRequest("GET", server.URL, nil)
		s.Require().NoError(err)

		res, err := (&http.Client{}).Do(req)
		s.Require().NoError(err)
		s.Require().Equal(http.StatusOK, res.StatusCode)
		body := s.GetBody(res)
		s.Require().Contains(body, "$213.50")
	})
}

type financesSuite struct {
	basesuite.BaseSuite
}

func TestFinancesSuite(t *testing.T) {
	suite.Run(t, new(financesSuite))
}
