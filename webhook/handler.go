package webhook

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// Handler processes incoming UQPay webhook requests.
func Handler(c *gin.Context) {
	// 1. Read body
	// Gin's BindJSON or ShouldBindJSON consumes the body, so we need to read it first
	// to verify the signature. c.GetRawData() reads the body and restores it for binding.
	body, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read request body"})
		return
	}

	// 2. Get headers
	signature := c.GetHeader("x-wk-signature")
	timestamp := c.GetHeader("x-wk-timestamp")

	if signature == "" || timestamp == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing signature or timestamp headers"})
		return
	}

	// 3. Verify signature
	secret := os.Getenv("UQPAY_WEBHOOK_SECRET")
	if secret == "" {
		log.Println("Error: UQPAY_WEBHOOK_SECRET environment variable not set")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if !VerifySignature(body, signature, timestamp, secret) {
		log.Printf("Invalid signature. Signature: %s, Timestamp: %s", signature, timestamp)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid signature"})
		return
	}

	// 4. Unmarshal JSON
	// Since we already read the body, we can unmarshal it manually or use ShouldBindBodyWith if we needed binding.
	// But since we have 'body' bytes, json.Unmarshal is efficient.
	var event UQPayEvent
	if err := json.Unmarshal(body, &event); err != nil {
		log.Printf("Failed to unmarshal JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
		return
	}

	// 5. Log the event (Business Logic)
	log.Printf("Received UQPay Event: Type=%s, Name=%s, ID=%s", event.EventType, event.EventName, event.EventID)
	log.Printf("Event Data: %+v", event.Data)

	// 6. Return 200 OK
	c.String(http.StatusOK, "OK")
}
