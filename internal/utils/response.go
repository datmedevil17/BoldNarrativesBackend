package utils

import (
	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type PaginatedResponse struct {
	Success    bool        `json:"success"`
	Data       interface{} `json:"data,omitempty"`
	Total      uint64      `json:"total"`
	PageSizes  int         `json:"pageSizes"`
	Page       int         `json:"page"`
	TotalPages int         `json:"totalPages"`
}

func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode,Response{
		Success:true,
		Message:message,
		Data:data,
	})

}

func ErrorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode,Response{
		Success:false,
		Error:message,
	})
}

func PaginatedSuccessResponse(c *gin.Context, statusCode int, data interface{}, total uint64, pageSizes int, page int, totalPages int) {
	c.JSON(statusCode,PaginatedResponse{
		Success:true,
		Data:data,
		Total:total,
		PageSizes:pageSizes,
		Page:page,
		TotalPages:totalPages,
	})
}
