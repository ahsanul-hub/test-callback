package handler

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"time"

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

type CallbackDataXimpay struct {
	XimpayID     string `json:"ximpayid"`
	XimpayStatus int    `json:"ximpaystatus"`
	CbParam      string `json:"cbparam"`
	SecretKey    string `json:"secretkey"`
	FailCode     int    `json:"failcode"`
}

func MerchantCallback(c *fiber.Ctx) error {
	var request CallbackRequest
	log.Println(request.MerchantTransactionID)

	if request.MerchantTransactionID == "36f3ad7d7ffd48c5a684a386de9ef440" {
		log.Println("callback requested")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}
	// Parse JSON request body
	if err := c.BodyParser(&request); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	// Log the received callback data
	log.Printf("Received callback for Transaction ID: %s, Status Code: %d \n", request.MerchantTransactionID, request.StatusCode)

	// Respond back to the sender
	return c.JSON(fiber.Map{"success": true, "message": "Callback received successfully"})
}

func generateXimpayToken(data CallbackDataXimpay) string {
	joinedString := fmt.Sprintf("%s%d%s%s", data.XimpayID, data.XimpayStatus, data.CbParam, data.SecretKey)
	hash := md5.Sum([]byte(joinedString))
	return hex.EncodeToString(hash[:])
}

func SendCallbackXimpay(c *fiber.Ctx) error {

	callbackData := CallbackDataXimpay{
		XimpayID:     "1F12BB46435A46738ABBA4AF23BCFB9D",
		XimpayStatus: 1,
		CbParam:      "123456",
		SecretKey:    "ABCD",
		FailCode:     0,
	}

	ximpayToken := generateXimpayToken(callbackData)

	merchantURL := "https://merchant.com/payment-status"

	// Format URL untuk callback
	callbackURL := fmt.Sprintf("%s?ximpayid=%s&ximpaystatus=%d&cbparam=%s&ximpaytoken=%s&failcode=%d",
		merchantURL, callbackData.XimpayID, callbackData.XimpayStatus, callbackData.CbParam, ximpayToken, callbackData.FailCode)

	log.Println("Mengirim callback ke:", callbackURL)

	go func() {
		time.Sleep(2 * time.Second)
		log.Println("Callback terkirim:", callbackURL)
	}()

	return c.JSON(fiber.Map{
		"message":      "Callback diproses",
		"callback_url": callbackURL,
	})
}

func ReceiveCallback(c *fiber.Ctx) error {
	// Log headers
	log.Println("Headers:")
	for k, v := range c.GetReqHeaders() {
		log.Printf("%s: %s\n", k, v)
	}

	// Log query parameters
	log.Println("Query Parameters:")
	for k, v := range c.Queries() {
		log.Printf("%s: %s\n", k, v)
	}

	// Log request body
	body := c.Body()
	log.Println("Raw Request Body:", string(body))

	// Parse body ke map
	var requestData map[string]interface{}
	if err := c.BodyParser(&requestData); err != nil {
		log.Println("Error parsing body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid request body"})
	}

	log.Println("Parsed Request Body:", requestData)

	return c.JSON(fiber.Map{"success": true, "message": "Callback received successfully"})
}

func Hello(c *fiber.Ctx) error {
	fmt.Println("Hello endpoint reached") // Debug log
	return c.SendString("Hello, from API 1!")
}
