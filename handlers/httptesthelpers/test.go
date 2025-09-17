package httptesthelpers

import (
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/gragorther/epigo/database/testhelpers"
)

func CreateTestContext() (c *gin.Context, w *httptest.ResponseRecorder) {
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	return
}

type HandlersTestSuite struct {
	testhelpers.DBTestSuite
}

func (h *HandlersTestSuite) AssertHTTPStatus(c *gin.Context, expected int, w *httptest.ResponseRecorder) {
	c.Writer.WriteHeaderNow()
	h.Equal(expected, w.Code, "http status codes should match")
}
