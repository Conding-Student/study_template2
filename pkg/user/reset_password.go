package users

import (
	"chatbot/pkg/handler"
	"chatbot/pkg/logs"
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/model"
	"chatbot/pkg/models/response"
	"chatbot/pkg/models/status"
	"chatbot/pkg/sharedfunctions"
	"chatbot/pkg/utils"
	"chatbot/pkg/utils/go-utils/database"
	"chatbot/pkg/utils/go-utils/encryptDecrypt"
	goerror "errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type PasswordResetCredentials struct {
	ID        uint      `gorm:"primaryKey"`
	Username  string    `json:"username"`
	StaffId   string    `json:"staffid"`
	Email     string    `json:"email"`
	Mobile    string    `json:"mobile"`
	Requested time.Time `json:"requested"`
}

type PasswordResetCredentialsLink struct {
	StaffId string `json:"staffid"`
	Email   string `json:"email"`
	Mobile  string `json:"mobile"`
}

//////////////////////////////////////////////////////////////////////////////////////////////////////
// 							This is for sending password reset link									//
//////////////////////////////////////////////////////////////////////////////////////////////////////

// Ito ay function para sa pagsend ng Password Reset Link sa email
func PasswordResetViaLink(c *fiber.Ctx) error {
	db := database.DB

	// Parse login credentials from request body
	requestBody := new(PasswordResetCredentialsLink)
	if err := c.BodyParser(requestBody); err != nil {
		fmt.Println("Failed to parse request", err)
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

	staffID := c.Params("id")
	// Find user with matching email in the database
	user := new(model.UserSignUp)
	if err := db.Where("staff_id = ? AND email = ? AND mobile = ?", requestBody.StaffId, requestBody.Email, requestBody.Mobile).First(user).Error; err != nil {
		if goerror.Is(err, gorm.ErrRecordNotFound) {
			logs.LOSLogs(c, "PasswordReset", staffID, "401", err.Error())
			return c.Status(401).JSON(response.ResponseModel{
				RetCode: "401",
				Message: status.RetCode401,
				Data: errors.ErrorModel{
					Message:   "It looks like there's a mismatch with your Staff ID, Email and Mobile No.. Please double-check and enter your correct Staff ID, Email and Mobile No..",
					IsSuccess: false,
					Error:     err,
				},
			})
		}
		logs.LOSLogs(c, "PasswordReset", staffID, "500", err.Error())
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Internal server error",
			Data: errors.ErrorModel{
				Message:   "Database Error",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	email := user.Email
	handler.SendPasswordResetLink(email)

	logs.LOSLogs(c, "PasswordReset", staffID, "200", ("Password reset link sent successfully. Please check your email." + email))
	return c.Status(200).JSON(response.ResponseModel{
		RetCode: "200",
		Message: "Password reset link sent successfully. Please check your email.",
		Data:    nil,
	})
}

//////////////////////////////////////////////////////////////////////////////////////////////////////
// 							This is for sending temporary password									//
//////////////////////////////////////////////////////////////////////////////////////////////////////

// Ito ay function para sa pagsend ng Temporary Password sa email
func PasswordReset(c *fiber.Ctx) error {

	requestBody := new(PasswordResetCredentialsLink)
	if err := c.BodyParser(requestBody); err != nil {
		fmt.Println("retCode 400")
		fmt.Println("Invalid Request")
		fmt.Println("Failed to parse request", err)
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

	staffID := c.Params("id")
	// 0 to generate
	isSuccess, message, tempPass, err := sharedfunctions.GenerateSaveTempPass(0, requestBody.StaffId, requestBody.Email, requestBody.Mobile, "")
	if err != nil {
		logs.LOSLogs(c, "PasswordReset", staffID, "500", err.Error())
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: status.RetCode500,
			Data: errors.ErrorModel{
				Message:   message,
				IsSuccess: isSuccess,
				Error:     err,
			},
		})
	}

	if !isSuccess {
		logs.LOSLogs(c, "PasswordReset", staffID, "401", message)
		return c.Status(401).JSON(response.ResponseModel{
			RetCode: "401",
			Message: status.RetCode401,
			Data: errors.ErrorModel{
				Message:   message,
				IsSuccess: isSuccess,
				Error:     err,
			},
		})
	}

	encryptedPassword, err := encryptDecrypt.Encrypt(tempPass, utils.GetEnv("SECRET_KEY"))
	if err != nil {
		logs.LOSLogs(c, "PasswordReset", staffID, "500", err.Error())
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Internal server error",
			Data: errors.ErrorModel{
				Message:   "Error encrypting password",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	// 1 to save
	isSuccess, message, encryptedPassword, err = sharedfunctions.GenerateSaveTempPass(1, requestBody.StaffId, requestBody.Email, requestBody.Mobile, encryptedPassword)
	if err != nil {
		logs.LOSLogs(c, "PasswordReset", staffID, "500", err.Error())
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: status.RetCode500,
			Data: errors.ErrorModel{
				Message:   message,
				IsSuccess: isSuccess,
				Error:     err,
			},
		})
	}

	if !isSuccess {
		logs.LOSLogs(c, "PasswordReset", staffID, "401", message)
		return c.Status(401).JSON(response.ResponseModel{
			RetCode: "401",
			Message: status.RetCode401,
			Data: errors.ErrorModel{
				Message:   message,
				IsSuccess: isSuccess,
				Error:     err,
			},
		})
	}

	// Send the combined password as the recovery password email
	if err := sharedfunctions.SendEmail(requestBody.Email, tempPass, "", 0); err != nil {
		logs.LOSLogs(c, "PasswordReset", staffID, "500", err.Error())
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: status.RetCode500,
			Data: errors.ErrorModel{
				Message:   "Error sending recovery email!",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	logs.LOSLogs(c, "PasswordReset", staffID, "200", "Temporary password sent successfully. Please check your email. "+encryptedPassword+" "+requestBody.Email)
	return c.Status(200).JSON(response.ResponseModel{
		RetCode: "200",
		Message: "Temporary password sent successfully. Please check your email. ",
		Data:    nil,
	})
}
