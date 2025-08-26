package encryptDecrypt

import (
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"
	"chatbot/pkg/utils"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

type DecryptRequest struct {
	Text string `json:"toDecrypt"`
}

type EncryptRequest struct {
	Text string `json:"toEncrypt"`
}

var iv = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 13, 05}
var secretKey = utils.GetEnv("SECRET_KEY")

// const secretKey string = "abc&1*~#^2^#s0^=)^^7%b34"

func encodeBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func decodeBase64(s string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		// panic(err)
		log.Println(err)
		return nil, err
	}
	return data, nil
}

// Encrypt method is to encrypt or hide any classified text
func Encrypt(text, secretKey string) (string, error) {
	block, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		return "", err
	}
	plainText := []byte(text)
	cfb := cipher.NewCFBEncrypter(block, iv)
	cipherText := make([]byte, len(plainText))
	cfb.XORKeyStream(cipherText, plainText)
	return encodeBase64(cipherText), nil
}

// Decrypt method is to extract back the encrypted text
func Decrypt(text, secretKey string) (string, error) {
	block, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		return "", err
	}
	cipherText, err := decodeBase64(text)
	if err != nil {
		return "", err
	}
	cfb := cipher.NewCFBDecrypter(block, iv)
	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, cipherText)
	return string(plainText), nil
}

func EncryptHandler(c *fiber.Ctx) error {
	// Parse the request body into a User struct
	encryptRequest := new(EncryptRequest)

	//secretKey := utils.GetEnv("SECRET_KEY")
	// if len(secretKey) < 16 {
	// 	secretKey = fmt.Sprintf("%-16s", secretKey)
	// } else if len(secretKey) > 16 {
	// 	secretKey = secretKey[:16]
	// }

	if err := c.BodyParser(&encryptRequest); err != nil {
		fmt.Println("retCode 401")
		fmt.Println("Invalid Request")
		fmt.Println("Failed to parse request", err.Error())
		return c.Status(401).JSON(response.ResponseModel{
			RetCode: "401",
			Message: "Invalid Request",
			Data: errors.ErrorModel{
				Message:   "Failed to parse request",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	encryptedText, err := Encrypt(encryptRequest.Text, secretKey)
	if err != nil {
		fmt.Println("retCode 404")
		fmt.Println("Bad Request")
		fmt.Println("Failed to encrypt request", err.Error())
		return c.Status(403).JSON(response.ResponseModel{
			RetCode: "403",
			Message: "Forbidden Request",
			Data: errors.ErrorModel{
				Message:   "Failed to encrypt request",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	return c.SendString(encryptedText)
}

func DecryptHandler(c *fiber.Ctx) error {
	// Parse the request body into a User struct
	decryptRequest := new(DecryptRequest)

	// secretKey := utils.GetEnv("SECRET_KEY")
	// if len(secretKey) < 16 {
	// 	secretKey = fmt.Sprintf("%-16s", secretKey)
	// } else if len(secretKey) > 16 {
	// 	secretKey = secretKey[:16]
	// }

	if err := c.BodyParser(&decryptRequest); err != nil {
		fmt.Println("retCode 401")
		fmt.Println("Invalid Request")
		fmt.Println("Failed to parse request", err.Error())
		return c.Status(401).JSON(response.ResponseModel{
			RetCode: "401",
			Message: "Invalid Request",
			Data: errors.ErrorModel{
				Message:   "Failed to parse request",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	decryptedText, err := Decrypt(decryptRequest.Text, secretKey)
	if err != nil {
		fmt.Println("retCode 404")
		fmt.Println("Bad Request")
		fmt.Println("Failed to encrypt request", err.Error())
		return c.Status(403).JSON(response.ResponseModel{
			RetCode: "403",
			Message: "Forbidden Request",
			Data: errors.ErrorModel{
				Message:   "Failed to decrypt request",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	return c.SendString(decryptedText)
}
