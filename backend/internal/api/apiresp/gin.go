package apiresp

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ginJson(c *gin.Context, resp *ApiResponse) {
	c.JSON(http.StatusOK, resp)
}

func GinError(c *gin.Context, err error) {
	ginJson(c, ParseError(err))
}

func GinSuccess(c *gin.Context, data any) {
	ginJson(c, ApiSuccess(data))
}
