package main

import (
	"log"
	"net/http"
	"test-callback-merchant/handler"

	"github.com/gofiber/fiber/v2"
)

type CallbackRequest struct {
	UserID                string `json:"user_id"`
	MerchantTransactionID string `json:"merchant_transaction_id"`
	StatusCode            int    `json:"status_code"`
	PaymentMethod         string `json:"payment_method"`
	Amount                string `json:"amount"`
	Status                string `json:"status"`
	Currency              string `json:"currency"`
	ItemName              string `json:"item_name"`
	ItemID                string `json:"item_id"`
	ReferenceID           string `json:"reference_id"`
}

func MerchantCallback(c *fiber.Ctx) error {
	var request CallbackRequest

	if request.MerchantTransactionID == "36f3ad7d7ffd48c5a684a386de9ef440" {
		log.Println("callback requesteds")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}
	if err := c.BodyParser(&request); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	log.Printf("Received callback for Transaction ID: %s, Status Code: %d, Message: %s\n", request.MerchantTransactionID, request.StatusCode)

	// Respond back to the sender
	return c.JSON(fiber.Map{"success": true, "message": "Callback received successfully"})
}

func main() {
	app := fiber.New()

	// Daftarkan endpoint untuk menerima callback dari sistem
	app.Post("/merchant/callback", handler.MerchantCallback)
	app.Post("/test-callback", handler.SendCallbackXimpay)
	app.Post("/receive-callback1", handler.ReceiveCallback)
	app.Get("/", handler.Hello)

	// Jalankan server
	app.Listen(":3000")
}
