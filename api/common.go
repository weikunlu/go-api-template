package api

import "github.com/gin-gonic/gin"

// A helper function to get error response
func GetErrorResponse(data interface{}, desc string) map[string]interface{} {
	return gin.H{
		"metadata": gin.H{
			"status": "9999",
			"desc":   desc,
		},
		"data": data,
	}
}

// A helper function to get success response
func GetSuccessResponse(data interface{}) map[string]interface{} {
	return gin.H{
		"metadata": gin.H{
			"status": "0000",
			"desc":   "",
		},
		"data": data,
	}
}
