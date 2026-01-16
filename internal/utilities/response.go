package utilities

import "github.com/gofiber/fiber/v2"

type BaseResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type PaginationMeta struct {
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	Total     int `json:"total"`
	TotalPage int `json:"total_page"`
}

type PaginatedResponse struct {
	Code     int            `json:"code"`
	Message  string         `json:"message"`
	Data     interface{}    `json:"data"`
	Metadata PaginationMeta `json:"metadata"`
}

func NewSuccessResponse(c *fiber.Ctx, status int, message string, data interface{}) error {
	return c.Status(status).JSON(BaseResponse{
		Code:    status,
		Message: message,
		Data:    data,
	})
}

func NewPaginatedResponse(c *fiber.Ctx, status int, message string, data interface{}, page, limit, total int) error {
	totalPage := total / limit
	if total%limit != 0 {
		totalPage++
	}

	return c.Status(status).JSON(PaginatedResponse{
		Code:    status,
		Message: message,
		Data:    data,
		Metadata: PaginationMeta{
			Page:      page,
			Limit:     limit,
			Total:     total,
			TotalPage: totalPage,
		},
	})
}

func NewErrorResponse(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(BaseResponse{
		Code:    status,
		Message: message,
		Data:    nil,
	})
}
