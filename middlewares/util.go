package middlewares

import (
	"fmt"
	"net/http"
)

func SetHttpAuthHeaderToken(header *http.Header, token string) {
	if header == nil {
		header = new(http.Header)
	}
	header.Set("Authorization", fmt.Sprint("Bearer ", token))
}
