package finances_test

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"

	"github.com/ffigari/stored-strings/internal/basesuite"
)

func (s *financesStorageSuite) SetupSuite() {
	s.dbPool = s.SetupDB(context.Background(), "finances_test_db")
}


func (s *financesStorageSuite) TearDownSuite() {
	s.dbPool.Close()
}

type financesStorageSuite struct {
	basesuite.BaseSuite
	dbPool *pgxpool.Pool
}

func TestFinancesStorageSuite(t *testing.T) {
	suite.Run(t, new(financesStorageSuite))
}
