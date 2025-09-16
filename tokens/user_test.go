package tokens_test

import (
	"testing"

	"github.com/gragorther/epigo/tokens"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	createUserAuth = tokens.CreateUserAuth(jwtSecret, []string{audience}, issuer)
	parseUserAuth  = tokens.ParseUserAuth(jwtSecret, []string{audience}, issuer)
)

func TestCreateUserAuth(t *testing.T) {
	table := map[string]struct {
		UserID uint
	}{
		"valid input": {
			UserID: 12,
		},
	}

	for name, test := range table {
		t.Run(name, func(t *testing.T) {
			require := require.New(t)
			gotToken, err := createUserAuth(test.UserID)
			require.NoError(err, "creating email verification token shouldn't fail")
			userID, err := parseUserAuth(gotToken)
			require.NoError(err, "parsing token shouldn't fail")

			assert.Equal(t, test.UserID, userID)
		})
	}
}

func TestParseUserAuth(t *testing.T) {
	const userID = 12
	type want struct {
		UserID uint
	}
	table := []struct {
		Name  string
		Token string
		Want  want
	}{
		{Name: "valid input", Want: want{UserID: userID}},
	}

	{
		var err error
		table[0].Token, err = createUserAuth(userID)
		if err != nil {
			panic(err)
		}
	}

	for _, test := range table {
		t.Run(test.Name, func(t *testing.T) {
			userID, err := parseUserAuth(test.Token)
			require.NoError(t, err)
			assert.Equal(t, test.Want.UserID, userID)
		})
	}
}
