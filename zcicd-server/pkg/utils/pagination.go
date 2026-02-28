package utils

import (
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Pagination struct {
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
	Total    int64 `json:"total"`
}

type PageResponse struct {
	Items      interface{} `json:"items"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	Total      int64       `json:"total"`
	TotalPages int         `json:"total_pages"`
}

func ParsePagination(c *gin.Context) Pagination {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	return Pagination{Page: page, PageSize: pageSize}
}

func (p Pagination) Offset() int {
	return (p.Page - 1) * p.PageSize
}

func (p Pagination) TotalPages() int {
	if p.PageSize <= 0 {
		return 0
	}
	return int(math.Ceil(float64(p.Total) / float64(p.PageSize)))
}

func NewPageResponse(items interface{}, pagination Pagination) PageResponse {
	return PageResponse{
		Items:      items,
		Page:       pagination.Page,
		PageSize:   pagination.PageSize,
		Total:      pagination.Total,
		TotalPages: pagination.TotalPages(),
	}
}
