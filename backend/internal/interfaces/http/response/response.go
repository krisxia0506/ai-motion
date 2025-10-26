package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type PaginationData struct {
	Items      interface{}       `json:"items"`
	Pagination *PaginationInfo   `json:"pagination,omitempty"`
}

type PaginationInfo struct {
	Page       int  `json:"page"`
	PageSize   int  `json:"pageSize"`
	Total      int  `json:"total"`
	TotalPages int  `json:"totalPages"`
	HasNext    bool `json:"hasNext"`
	HasPrev    bool `json:"hasPrev"`
}

const (
	CodeSuccess = 0

	CodeInvalidParams     = 10001
	CodeResourceNotFound  = 10002
	CodeFileParseError    = 30002
	CodeAIServiceError    = 40001
	CodeGenerationError   = 40003
	CodeDatabaseError     = 50001
	CodeInternalError     = 50002
)

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: "success",
		Data:    data,
	})
}

func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: message,
		Data:    data,
	})
}

func SuccessList(c *gin.Context, items interface{}, page, pageSize, total int) {
	totalPages := (total + pageSize - 1) / pageSize
	if totalPages < 1 {
		totalPages = 1
	}

	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: "success",
		Data: PaginationData{
			Items: items,
			Pagination: &PaginationInfo{
				Page:       page,
				PageSize:   pageSize,
				Total:      total,
				TotalPages: totalPages,
				HasNext:    page < totalPages,
				HasPrev:    page > 1,
			},
		},
	})
}

func Error(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

func ErrorWithData(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

func InvalidParams(c *gin.Context, message string) {
	Error(c, CodeInvalidParams, message)
}

func ResourceNotFound(c *gin.Context, message string) {
	Error(c, CodeResourceNotFound, message)
}

func FileParseError(c *gin.Context, message string) {
	Error(c, CodeFileParseError, message)
}

func AIServiceError(c *gin.Context, message string) {
	Error(c, CodeAIServiceError, message)
}

func GenerationError(c *gin.Context, message string) {
	Error(c, CodeGenerationError, message)
}

func DatabaseError(c *gin.Context, message string) {
	Error(c, CodeDatabaseError, message)
}

func InternalError(c *gin.Context, message string) {
	Error(c, CodeInternalError, message)
}
