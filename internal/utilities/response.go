package utilities

import "github.com/gofiber/fiber/v2"

type BaseResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewSuccessResponse(c *fiber.Ctx, status int, message string, data interface{}) error {
	return c.Status(status).JSON(BaseResponse{
		Code:    status,
		Message: message,
		Data:    data,
	})
}

func NewErrorResponse(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(BaseResponse{
		Code:    status,
		Message: message,
		Data:    nil,
	})
}
