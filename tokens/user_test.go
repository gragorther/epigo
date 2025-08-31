package tokens_test

import (
	"testing"

	"github.com/gragorther/epigo/tokens"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
			gotToken, err := tokens.CreateUserAuth(jwtSecret, test.UserID, audience, []string{issuer})
			require.NoError(err, "creating email verification token shouldn't fail")
			userID, err := tokens.ParseUserAuth(jwtSecret, gotToken, issuer, []string{audience})
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
		table[0].Token, err = tokens.CreateUserAuth(jwtSecret, userID, audience, []string{issuer})
		if err != nil {
			panic(err)
		}
	}

	for _, test := range table {
		t.Run(test.Name, func(t *testing.T) {
			userID, err := tokens.ParseUserAuth(jwtSecret, test.Token, audience, []string{issuer})
			require.NoError(t, err)
			assert.Equal(t, test.Want.UserID, userID)
		})
	}
}
