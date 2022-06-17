package responses

import "github.com/gofiber/fiber/v2"

type RequestResponse struct {
	Status   int        `json:"status"`
	Message  string     `json:"message"`
	Response *fiber.Map `json:"response"`
}
