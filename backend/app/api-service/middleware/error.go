package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/STLeee/mediation-platform/backend/app/api-service/model"
)

// ErrorHandler is a middleware for handling errors
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		err := c.Errors.Last()
		if err != nil {
			// Handle error
			httpStatusCodeError, ok := err.Err.(model.HttpStatusCodeError)
			if !ok {
				httpStatusCodeError = model.HttpStatusCodeError{
					StatusCode: 500,
					Err:        err.Err,
				}
			}

			// Print error message
			fmt.Printf("Error (%d): %s\n", httpStatusCodeError.StatusCode, httpStatusCodeError.Error())

			// Send response
			response := model.MessageResponse{
				Message: httpStatusCodeError.Error(),
			}
			c.JSON(httpStatusCodeError.StatusCode, response)
			c.Abort()
			return
		}
	}
}
