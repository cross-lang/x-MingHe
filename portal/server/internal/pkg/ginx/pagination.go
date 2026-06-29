package ginx

import (
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	// DefaultPageSize 默认每页数量
	DefaultPageSize = 10
	// MaxPageSize 最大每页数量
	MaxPageSize = 100
	// DefaultPage 默认页码
	DefaultPage = 1
)

// Pagination 分页信息
type Pagination struct {
	Page       int64 `json:"page"`        // 当前页码
	Size       int64 `json:"size"`        // 每页数量
	Total      int64 `json:"total"`       // 总记录数
	TotalPages int64 `json:"total_pages"` // 总页数
}

// PaginationRequest 分页请求参数
type PaginationRequest struct {
	Page int `form:"page" binding:"omitempty,min=1"`
	Size int `form:"size" binding:"omitempty,min=1,max=100"`
}

// NewPaginationRequest 从上下文创建分页请求
func NewPaginationRequest(c *gin.Context) *PaginationRequest {
	req := &PaginationRequest{}

	pageStr := c.DefaultQuery("page", strconv.Itoa(DefaultPage))
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = DefaultPage
	}
	req.Page = page

	sizeStr := c.DefaultQuery("size", strconv.Itoa(DefaultPageSize))
	size, err := strconv.Atoi(sizeStr)
	if err != nil || size < 1 {
		size = DefaultPageSize
	}
	if size > MaxPageSize {
		size = MaxPageSize
	}
	req.Size = size

	return req
}

// GetOffset 获取分页偏移量
func (p *PaginationRequest) GetOffset() int {
	return (p.Page - 1) * p.Size
}

// GetLimit 获取分页限制
func (p *PaginationRequest) GetLimit() int {
	return p.Size
}

// NewPagination 创建分页信息
func NewPagination(page, size, total int64) *Pagination {
	totalPages := int64(math.Ceil(float64(total) / float64(size)))
	if totalPages == 0 {
		totalPages = 1
	}

	return &Pagination{
		Page:       page,
		Size:       size,
		Total:      total,
		TotalPages: totalPages,
	}
}

// PagedResponse 分页响应数据
type PagedResponse struct {
	Data       interface{}  `json:"data"`
	Pagination *Pagination `json:"pagination"`
}

// SuccessPagedResponse 成功的分页响应
func SuccessPagedResponse(c *gin.Context, data interface{}, pagination *Pagination) {
	response := Response{
		Code: 0,
		Msg:  "success",
		Data: PagedResponse{
			Data:       data,
			Pagination: pagination,
		},
	}

	// 添加请求ID
	if requestID, exists := c.Get("request_id"); exists {
		if id, ok := requestID.(string); ok {
			response.RequestID = id
		}
	}

	c.JSON(200, response)
}

// ListResponse 列表响应数据（带分页信息）
type ListResponse struct {
	Items      interface{}  `json:"items"`
	Pagination *Pagination `json:"pagination"`
}

// SuccessListResponse 成功的列表响应
func SuccessListResponse(c *gin.Context, items interface{}, pagination *Pagination) {
	response := Response{
		Code: 0,
		Msg:  "success",
		Data: ListResponse{
			Items:      items,
			Pagination: pagination,
		},
	}

	// 添加请求ID
	if requestID, exists := c.Get("request_id"); exists {
		if id, ok := requestID.(string); ok {
			response.RequestID = id
		}
	}

	c.JSON(200, response)
}

// ParsePagination 解析分页参数并返回偏移量和限制
func ParsePagination(c *gin.Context) (page, size, offset, limit int) {
	paginationReq := NewPaginationRequest(c)
	return paginationReq.Page, paginationReq.Size, paginationReq.GetOffset(), paginationReq.GetLimit()
}

// BuildPagination 构建分页信息
func BuildPagination(page, size, total int) *Pagination {
	return NewPagination(int64(page), int64(size), int64(total))
}
