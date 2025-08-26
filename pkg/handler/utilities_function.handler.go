package handler

import (
	"chatbot/pkg/logs"
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"
	"chatbot/pkg/utils"
	"chatbot/pkg/utils/go-utils/database"
	"chatbot/pkg/utils/go-utils/encryptDecrypt"
	"encoding/base64"
	"fmt"
	"log"
	"net/smtp"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

var messageStatus = "Opss!"

// Your existing AuthMiddleware function
func AuthMiddleware(c *fiber.Ctx) error {
	// Retrieve the value of the  header
	id := c.Params("id")
	path := c.OriginalURL()

	// secretKey := utils.GetEnv("SECRET_KEY")
	appCode := c.Get("appCode")
	secretCode := c.Get("securityCode")

	// decryptedAppcode, err := encryptDecrypt.Decrypt(appCode, secretKey)
	// if err != nil {
	// 	logs.LOSLogs(c, "Header", id, "401", "Error decrypting (appCode) "+path+" "+appCode)
	// 	return c.Status(401).JSON(response.ResponseModel{
	// 		RetCode: "401",
	// 		Message: "Invalid Request",
	// 		Data: errors.ErrorModel{
	// 			Message:   "Unauthorized Request",
	// 			IsSuccess: false,
	// 			Error:     nil,
	// 		},
	// 	})
	// }

	// decryptedSecretCode, err := encryptDecrypt.Decrypt(secretCode, secretKey)
	// if err != nil {
	// 	logs.LOSLogs(c, "Header", id, "401", "Error decrypting (secretCode) "+path+" "+secretCode)
	// 	return c.Status(401).JSON(response.ResponseModel{
	// 		RetCode: "401",
	// 		Message: "Invalid Request",
	// 		Data: errors.ErrorModel{
	// 			Message:   "Unauthorized Request",
	// 			IsSuccess: false,
	// 			Error:     nil,
	// 		},
	// 	})
	// }

	// Check if the 'appCode' header is present and has the expected value
	if appCode != utils.GetEnv("appCode") {
		logs.LOSLogs(c, "Header", id, "401", "Unauthorized Request (missing appCode) "+path)
		return c.Status(401).JSON(response.ResponseModel{
			RetCode: "401",
			Message: "Invalid Request",
			Data: errors.ErrorModel{
				Message:   "Unauthorized Request",
				IsSuccess: false,
				Error:     nil,
			},
		})
	}

	// Check if the 'secretCodes' header is present and has the expected value
	if secretCode != utils.GetEnv("securityCode") {
		logs.LOSLogs(c, "Header", id, "401", "Unauthorized Request (missing securityCode) "+path)
		return c.Status(401).JSON(response.ResponseModel{
			RetCode: "401",
			Message: "Invalid Request",
			Data: errors.ErrorModel{
				Message:   "Unauthorized Request",
				IsSuccess: false,
				Error:     nil,
			},
		})
	}
	return c.Next()
}

//////////////////////////////////////////////////////////////////////////////////////////////////////
// 			        	Ito ay para sa pagsend ng Password Reset Link sa email	        			//
//////////////////////////////////////////////////////////////////////////////////////////////////////

func SendPasswordResetLink(toEmail string) error {
	// Email configuration
	smtpServer := "smtp.gmail.com"
	smtpPort := 587
	senderEmail := "CA-GABAY@fortress-asya.com"
	senderPassword := "dqiynvnialknyost"

	// Compose the email message
	to := []string{toEmail}
	body := "Please don not share this link.\nClick this link to reset your password " + utils.GetEnv("passwordResetLink")
	subject := "Password Reset Link"
	message := "Subject: " + subject + "\r\n" +
		"To: " + toEmail + "\r\n" +
		"Content-Type: text/plain; charset=UTF-8\r\n" +
		"Reply-To: noreply@ca-gabay.com\r\n" + // Set the Reply-To header
		"\r\n" + body

	// Connect to the SMTP server
	auth := smtp.PlainAuth("", senderEmail, senderPassword, smtpServer)
	err := smtp.SendMail(smtpServer+":"+strconv.Itoa(smtpPort), auth, senderEmail, to, []byte(message))
	if err != nil {
		log.Println("Error sending email:", err)
		return err
	}

	log.Println("Recovery password email sent successfully to:", toEmail)
	return nil
}

//////////////////////////////////////////////////////////////////////////////////////////////////////
// 			        	  Ito ay para server switch kung kailangan ioff ang server	     	   		//
//////////////////////////////////////////////////////////////////////////////////////////////////////

type ServerSwitchCredentials struct {
	Switch  bool   `json:"switch"`
	Message string `json:"message"`
	Title   string `json:"title"`
}

func ServerSwitchMain(c *fiber.Ctx) error {
	db := database.DB.Begin()

	id := c.Params("id")
	// Query the server switch status from the database
	var serverSwitch ServerSwitchCredentials
	if err := db.Raw("SELECT * FROM public.getswitch(1)").Scan(&serverSwitch).Error; err != nil {
		logs.LOSLogs(c, "ServerSwitch", id, "500", err.Error())
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Database Error!",
			Data: errors.ErrorModel{
				Message:   "An error occured while connecting to database.",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	if !serverSwitch.Switch {
		logs.LOSLogs(c, "ServerSwitch", id, "503", serverSwitch.Message)
		return c.Status(503).JSON(response.ResponseModel{
			RetCode: "503",
			Message: serverSwitch.Title,
			Data: errors.ErrorModel{
				Message:   serverSwitch.Message,
				IsSuccess: false,
				Error:     nil,
			},
		})
	}
	db.Commit()
	return c.Next()
}

func ServerSwitchRegistration(c *fiber.Ctx) error {
	db := database.DB.Begin()

	id := c.Params("id")
	// Query the server switch status from the database
	var serverSwitch ServerSwitchCredentials
	if err := db.Raw("SELECT * FROM public.getswitch(2)").Scan(&serverSwitch).Error; err != nil {
		logs.LOSLogs(c, "ServerSwitch", id, "500", err.Error())
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Database Error!",
			Data: errors.ErrorModel{
				Message:   "An error occured while connecting to database.",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	if !serverSwitch.Switch {
		logs.LOSLogs(c, "ServerSwitch", id, "503", serverSwitch.Message)
		return c.Status(503).JSON(response.ResponseModel{
			RetCode: "503",
			Message: serverSwitch.Title,
			Data: errors.ErrorModel{
				Message:   serverSwitch.Message,
				IsSuccess: false,
				Error:     nil,
			},
		})
	}
	db.Commit()
	return c.Next()
}

func ServerSwitchEMPCSOA(c *fiber.Ctx) error {
	db := database.DB.Begin()

	id := c.Params("id")
	// Query the server switch status from the database
	var serverSwitch ServerSwitchCredentials
	if err := db.Raw("SELECT * FROM public.getswitch(3)").Scan(&serverSwitch).Error; err != nil {
		logs.LOSLogs(c, "ServerSwitch", id, "500", err.Error())
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Database Error!",
			Data: errors.ErrorModel{
				Message:   "An error occured while connecting to database.",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	if !serverSwitch.Switch {
		logs.LOSLogs(c, "ServerSwitch", id, "503", serverSwitch.Message)
		return c.Status(503).JSON(response.ResponseModel{
			RetCode: "503",
			Message: serverSwitch.Title,
			Data: errors.ErrorModel{
				Message:   serverSwitch.Message,
				IsSuccess: false,
				Error:     nil,
			},
		})
	}
	db.Commit()
	return c.Next()
}

func ServerSwitchChatbot(c *fiber.Ctx) error {
	db := database.DB.Begin()

	id := c.Params("id")
	// Query the server switch status from the database
	var serverSwitch ServerSwitchCredentials
	if err := db.Raw("SELECT * FROM public.getswitch(4)").Scan(&serverSwitch).Error; err != nil {
		logs.LOSLogs(c, "ServerSwitch", id, "500", err.Error())
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Database Error!",
			Data: errors.ErrorModel{
				Message:   "An error occured while connecting to database.",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	if !serverSwitch.Switch {
		logs.LOSLogs(c, "ServerSwitch", id, "503", serverSwitch.Message)
		return c.Status(503).JSON(response.ResponseModel{
			RetCode: "503",
			Message: serverSwitch.Title,
			Data: errors.ErrorModel{
				Message:   serverSwitch.Message,
				IsSuccess: false,
				Error:     nil,
			},
		})
	}
	db.Commit()
	return c.Next()
}

func ServerSwitchLoanCalculator(c *fiber.Ctx) error {
	db := database.DB.Begin()

	id := c.Params("id")
	// Query the server switch status from the database
	var serverSwitch ServerSwitchCredentials
	if err := db.Raw("SELECT * FROM public.getswitch(5)").Scan(&serverSwitch).Error; err != nil {
		logs.LOSLogs(c, "ServerSwitch", id, "500", err.Error())
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Database Error!",
			Data: errors.ErrorModel{
				Message:   "An error occured while connecting to database.",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	if !serverSwitch.Switch {
		logs.LOSLogs(c, "ServerSwitch", id, "503", serverSwitch.Message)
		return c.Status(503).JSON(response.ResponseModel{
			RetCode: "503",
			Message: serverSwitch.Title,
			Data: errors.ErrorModel{
				Message:   serverSwitch.Message,
				IsSuccess: false,
				Error:     nil,
			},
		})
	}
	db.Commit()
	return c.Next()
}

func ServerSwitchGabayKonek(c *fiber.Ctx) error {
	db := database.DB.Begin()

	id := c.Params("id")
	// Query the server switch status from the database
	var serverSwitch ServerSwitchCredentials
	if err := db.Raw("SELECT * FROM public.getswitch(6)").Scan(&serverSwitch).Error; err != nil {
		logs.LOSLogs(c, "ServerSwitch", id, "500", err.Error())
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Database Error!",
			Data: errors.ErrorModel{
				Message:   "An error occured while connecting to database.",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	if !serverSwitch.Switch {
		logs.LOSLogs(c, "ServerSwitch", id, "503", serverSwitch.Message)
		return c.Status(503).JSON(response.ResponseModel{
			RetCode: "503",
			Message: serverSwitch.Title,
			Data: errors.ErrorModel{
				Message:   serverSwitch.Message,
				IsSuccess: false,
				Error:     nil,
			},
		})
	}
	db.Commit()
	return c.Next()
}

// BasicAuth middleware function
func BasicAuth(c *fiber.Ctx) error {
	secretKey := utils.GetEnv("SECRET_KEY")

	id := c.Params("id")
	path := c.OriginalURL()
	fmt.Println(path)
	basicAuthFeature := "BasicAuth"
	userAgent := c.Get("User-Agent")

	// if strings.HasPrefix(userAgent, "Prometheus/") && c.Path() == "/metrics" {
	// 	// Allow the request without authentication
	// 	return c.Next()
	// }

	// fmt.Println(userAgent)
	decryptedValidUsername, err := encryptDecrypt.Decrypt(utils.GetEnv("username"), secretKey)
	if err != nil {
		logs.LOSLogs(c, basicAuthFeature, id, "500", err.Error()+" "+path+" "+userAgent)
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Internal Server Error",
			Data: errors.ErrorModel{
				Message:   "Failed to decrypt valid username",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	decryptedValidPassword, err := encryptDecrypt.Decrypt(utils.GetEnv("password"), secretKey)
	if err != nil {
		logs.LOSLogs(c, basicAuthFeature, id, "500", err.Error()+" "+path+" "+userAgent)
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Internal Server Error",
			Data: errors.ErrorModel{
				Message:   "Failed to decrypt valid password",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	authHeader := c.Get("Authorization")

	// fmt.Println(authHeader)
	// fmt.Println(decryptedAuthHeader)
	if authHeader == "" {
		logs.LOSLogs(c, basicAuthFeature, id, "401", "Unauthorized Request "+path+" "+authHeader+" "+userAgent)
		return c.Status(401).JSON(response.ResponseModel{
			RetCode: "401",
			Message: "Invalid Request",
			Data: errors.ErrorModel{
				Message:   "Unauthorized Request",
				IsSuccess: false,
				Error:     nil,
			},
		})
	}

	authParts := strings.Split(authHeader, " ")
	authPartsString := strings.Join(authParts, " ")
	// fmt.Println(authPartsString)
	if len(authParts) != 2 || authParts[0] != "Basic" {
		logs.LOSLogs(c, basicAuthFeature, id, "401", "Unauthorized Request "+path+" "+authPartsString+" "+userAgent)
		return c.Status(401).JSON(response.ResponseModel{
			RetCode: "401",
			Message: "Invalid Request",
			Data: errors.ErrorModel{
				Message:   "Unauthorized Request",
				IsSuccess: false,
				Error:     nil,
			},
		})
	}

	payload, err := base64.StdEncoding.DecodeString(authParts[1])
	if err != nil {
		logs.LOSLogs(c, basicAuthFeature, id, "401", err.Error()+" "+path+" "+userAgent)
		return c.Status(401).JSON(response.ResponseModel{
			RetCode: "401",
			Message: "Invalid Request",
			Data: errors.ErrorModel{
				Message:   "Failed to decode Basic Auth",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	creds := strings.Split(string(payload), ":")
	if len(creds) != 2 {
		logs.LOSLogs(c, basicAuthFeature, id, "401", "Invalid Basic Auth Username or Password "+path+" "+userAgent)
		return c.Status(401).JSON(response.ResponseModel{
			RetCode: "401",
			Message: "Invalid Request",
			Data: errors.ErrorModel{
				Message:   "Invalid Basic Auth Username or Password",
				IsSuccess: false,
				Error:     nil,
			},
		})
	}

	basicAuthUsername, err := encryptDecrypt.Decrypt(creds[0], secretKey)
	if err != nil {
		logs.LOSLogs(c, basicAuthFeature, id, "500", err.Error()+" "+path+" "+basicAuthUsername+" "+userAgent)
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Internal Server Error",
			Data: errors.ErrorModel{
				Message:   "Failed to decrypt username",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	basicAuthPassword, err := encryptDecrypt.Decrypt(creds[1], secretKey)
	if err != nil {
		logs.LOSLogs(c, basicAuthFeature, id, "500", err.Error()+" "+path+" "+basicAuthPassword+" "+userAgent)
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Internal Server Error",
			Data: errors.ErrorModel{
				Message:   "Failed to decrypt password",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	// fmt.Println(basicAuthUsername)
	// fmt.Println(basicAuthPassword)
	if basicAuthUsername != decryptedValidUsername || basicAuthPassword != decryptedValidPassword {
		logs.LOSLogs(c, basicAuthFeature, id, "401", "Invalid Basic Auth Username or Password "+path+" "+basicAuthUsername+" "+basicAuthPassword+" "+userAgent)
		return c.Status(401).JSON(response.ResponseModel{
			RetCode: "401",
			Message: "Invalid Request",
			Data: errors.ErrorModel{
				Message:   "Invalid Basic Auth Username or Password",
				IsSuccess: false,
				Error:     nil,
			},
		})
	}

	return c.Next()
}

func FormatToValidBirthday(Birthday string) (string, error) {
	parsedTime, err := time.Parse("2006-01-02 00:00:00.0", Birthday)
	if err != nil {
		log.Printf("Error parsing birthday: %v", err)
		return "", err
	}

	return parsedTime.Format("2006-01-02"), nil
}
