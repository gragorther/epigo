package tokens_test

import (
	"testing"

	"github.com/gragorther/epigo/tokens"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var jwtSecret = []byte("testsecret")

const audience = "https://testserver.com"
const issuer = "https://testserver.com"

func TestCreateEmailVerification(t *testing.T) {
	table := map[string]struct {
		Email string
	}{
		"valid input": {
			Email: "testemail@testing.com",
		},
	}

	for name, test := range table {
		t.Run(name, func(t *testing.T) {
			require := require.New(t)
			gotToken, err := tokens.CreateEmailVerification(jwtSecret, test.Email, audience, issuer)
			require.NoError(err, "creating email verification token shouldn't fail")
			userEmail, err := tokens.ParseEmailVerification(jwtSecret, gotToken, issuer, audience)
			require.NoError(err, "parsing token shouldn't fail")

			assert.Equal(t, test.Email, userEmail)
		})
	}
}

const testIssuer = "https://isszuer.com"
const testEmail = "testemail@testing.com"

const testAudience = "https://server.com"

func TestParseEmailVerification(t *testing.T) {

	type want struct {
		Email string
		Error error
	}
	table := []struct {
		Name  string
		Token string
		Want  want
	}{
		{Name: "valid input", Want: want{Error: nil, Email: testEmail}},
	}

	{
		var err error
		table[0].Token, err = tokens.CreateEmailVerification(jwtSecret, testEmail, testAudience, testIssuer)
		if err != nil {
			panic(err)
		}
	}

	for _, test := range table {
		t.Run(test.Name, func(t *testing.T) {
			email, err := tokens.ParseEmailVerification(jwtSecret, test.Token, testAudience, testIssuer)
			require.NoError(t, err)
			assert.Equal(t, test.Want.Email, email)
		})
	}
}
