package tokens_test

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gragorther/epigo/tokens"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var jwtSecret = []byte("testsecret")

const (
	audience = "https://testserver.com"
	issuer   = "https://testserver.com"
)

var parseEmailVerificationToken = tokens.ParseEmailVerification(jwtSecret, audience, issuer)

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
			gotToken, err := tokens.CreateEmailVerification(jwtSecret, audience, issuer)(test.Email)
			require.NoError(err, "creating email verification token shouldn't fail")
			userEmail, err := parseEmailVerificationToken(gotToken)
			require.NoError(err, "parsing token shouldn't fail")

			assert.Equal(t, test.Email, userEmail)
		})
	}
}

const testEmail = "testemail@testing.com"

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
		{Name: "invalid signing method", Want: want{Error: tokens.ErrInvalidToken}},
	}

	{
		var err error
		table[0].Token, err = tokens.CreateEmailVerification(jwtSecret, audience, issuer)(testEmail)
		if err != nil {
			t.Fatalf("failed to create token: %v", err)
		}
	}
	{
		token, err := jwt.NewWithClaims(jwt.SigningMethodHS512, tokens.EmailClaims{
			Email: "testemail@google.com",
			Claims: tokens.Claims{
				Type: tokens.TypeEmailVerification,
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    issuer,
					Audience:  []string{audience},
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)),
					IssuedAt:  jwt.NewNumericDate(time.Now()),
				},
			},
		}).SignedString(jwtSecret)
		require.NoError(t, err)
		table[1].Token = token
	}

	for _, test := range table {
		t.Run(test.Name, func(t *testing.T) {
			email, err := parseEmailVerificationToken(test.Token)
			if test.Want.Error != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, test.Want.Email, email)
		})
	}
}
