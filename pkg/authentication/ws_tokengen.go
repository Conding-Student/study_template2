package authentication

import (
	"chatbot/pkg/models/response"
	"chatbot/pkg/utils"
	"chatbot/pkg/utils/go-utils/database"
	"chatbot/pkg/utils/go-utils/encryptDecrypt"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// WebSocketToken represents a short-lived token for WebSocket connections
type WebSocketToken struct {
	Token         string    `json:"token" gorm:"primaryKey"`
	OriginalToken string    `json:"original_token" gorm:"not null"`
	UserID        string    `json:"user_id" gorm:"not null"`
	ExpiresAt     time.Time `json:"expires_at" gorm:"not null"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// TableName specifies the table name with schema for GORM
func (WebSocketToken) TableName() string {
	return "authentication.websocket_tokens"
}

// WebSocketTokenRequest represents the request for a WebSocket token
type WebSocketTokenRequest struct {
	ConnectionType string `json:"connection_type"` // "articles", "trivia", "news"
}

// WebSocketTokenResponse represents the response with WebSocket token
type WebSocketTokenResponse struct {
	WebSocketToken string    `json:"websocket_token"`
	ExpiresAt      time.Time `json:"expires_at"`
	WebSocketURL   string    `json:"websocket_url"`
}

// GenerateWebSocketToken generates a short-lived token for WebSocket connections
func GenerateWebSocketToken(c *fiber.Ctx) error {
	db := database.DB

	// Validate admin token first
	if err := ValidateAdminToken(c); err != nil {
		return err
	}

	var req WebSocketTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(response.ResponseModel{
			RetCode: "400",
			Message: "Invalid request body",
			Data:    nil,
		})
	}

	// Get the original token from header
	authHeader := c.Get("Authorization")
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		return c.Status(401).JSON(response.ResponseModel{
			RetCode: "401",
			Message: "Invalid token format",
			Data:    nil,
		})
	}

	// Extract user info from token
	claims, err := extractClaimsFromToken(tokenString)
	if err != nil {
		return c.Status(401).JSON(response.ResponseModel{
			RetCode: "401",
			Message: "Invalid token",
			Data:    nil,
		})
	}

	// Generate a short-lived WebSocket token
	wsToken := uuid.New().String()
	expiresAt := time.Now().Add(5 * time.Minute) // Short expiration

	// Store the token in database
	wsTokenRecord := WebSocketToken{
		Token:         wsToken,
		OriginalToken: tokenString,
		UserID:        claims.StaffID,
		ExpiresAt:     expiresAt,
	}

	if err := db.Create(&wsTokenRecord).Error; err != nil {
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Failed to generate WebSocket token",
			Data:    nil,
		})
	}

	// Build WebSocket URL
	wsURL := fmt.Sprintf("ws://%s/api/public/v1/admin/ws/%s?token=%s",
		c.Hostname(), req.ConnectionType, wsToken)

	responseData := WebSocketTokenResponse{
		WebSocketToken: wsToken,
		ExpiresAt:      expiresAt,
		WebSocketURL:   wsURL,
	}

	return c.JSON(response.ResponseModel{
		RetCode: "200",
		Message: "WebSocket token generated successfully",
		Data:    responseData,
	})
}

// ValidateWebSocketToken validates the short-lived WebSocket token
func ValidateWebSocketToken(token string) (*WebSocketToken, error) {
	db := database.DB

	var wsToken WebSocketToken
	result := db.Where("token = ? AND expires_at > ?", token, time.Now()).First(&wsToken)
	if result.Error != nil {
		return nil, fmt.Errorf("invalid or expired WebSocket token")
	}

	// Optional: You can also validate the original token here if needed

	return &wsToken, nil
}

// CleanupExpiredWebSocketTokens removes expired WebSocket tokens
func CleanupExpiredWebSocketTokens() {
	db := database.DB
	for {
		time.Sleep(1 * time.Hour) // Run cleanup every hour
		db.Where("expires_at <= ?", time.Now()).Delete(&WebSocketToken{})
	}
}

// Helper function to extract claims from token
func extractClaimsFromToken(tokenString string) (*JWTClaims, error) {
	decryptedToken, err := encryptDecrypt.Decrypt(tokenString, utils.GetEnv("SECRET_KEY"))
	if err != nil {
		return nil, err
	}

	parts := strings.Split(decryptedToken, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid token format")
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}

	var claims JWTClaims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, err
	}

	return &claims, nil
}
