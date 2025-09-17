package middlewares

import (
	"log"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // Step1: Process the request first.

		// Step2: Check if any errors were added to the context
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				log.Print(err)
			}
		}

		// Any other steps if no errors are found
	}
}
