package db_test

import (
	"testing"

	"github.com/gragorther/epigo/database/testhelpers"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	testhelpers.DBTestSuite
}

func TestDB(t *testing.T) {
	suite.Run(t, new(Suite))
}
