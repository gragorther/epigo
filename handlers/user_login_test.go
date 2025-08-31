package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
	"github.com/gragorther/epigo/database/mock"
	"github.com/gragorther/epigo/handlers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoginUser(t *testing.T) {
	type want struct {
		Status int
	}
	table := []struct {
		Name  string
		Input handlers.LoginInput
	}{}

	for _, test := range table {
		t.Run(test.Name, func(t *testing.T) {
			require := require.New(t)
			assert := assert.New(t)
			r := gin.Default()
			mock := mock.NewMockDB()
			r.POST("/", handlers.LoginUser(mock, comparePasswordAndHash, JWT_SECRET, testIssuer, []string{testIssuer}))

			input, err := sonic.MarshalString(test.Input)
			require.NoError(err)
			req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/", strings.NewReader(input))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(http.StatusOK, w.Code, "http status code should indicate success")
		})
	}
}
