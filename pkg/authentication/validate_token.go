package authentication

import (
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"
	"chatbot/pkg/models/status"
	"chatbot/pkg/sharedfunctions"
	"chatbot/pkg/utils"
	"chatbot/pkg/utils/go-utils/database"
	"chatbot/pkg/utils/go-utils/encryptDecrypt"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"slices"

	"github.com/gofiber/fiber/v2"
)

type JWTClaims struct {
	Exp      int64  `json:"exp"`
	Rolename string `json:"rolename"`
	StaffID  string `json:"staffID"`
	Username string `json:"username"`
}

type AdminRoles struct {
	SuperAdmin string
	Admin      string
}

func ValidateSuperAdminToken(c *fiber.Ctx) error {
	db := database.DB
	authHeader := c.Get("Authorization")

	if authHeader == "" {
		fmt.Println("Missing or invalid token: ", authHeader)
		return c.Status(419).JSON(response.ResponseModel{
			RetCode: "419",
			Message: status.RetCode419,
			Data: errors.ErrorModel{
				Message:   "Invalid Authentication. You are not authorized to access this feature.",
				IsSuccess: false,
				Error:     nil,
			},
		})
	}

	// Remove "Bearer " prefix if it exists
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		fmt.Println("Invalid token format: ", authHeader)
		return c.Status(419).JSON(response.ResponseModel{
			RetCode: "419",
			Message: status.RetCode419,
			Data: errors.ErrorModel{
				Message:   "Invalid Authentication. You are not authorized to access this feature.",
				IsSuccess: false,
				Error:     nil,
			},
		})
	}

	isSuccess, retCodeInt, retCode, tstatus, tmessage, err := sharedfunctions.ValidateToken(tokenString)
	if err != nil {
		return c.Status(retCodeInt).JSON(response.ResponseModel{
			RetCode: retCode,
			Message: tstatus,
			Data: errors.ErrorModel{
				Message:   tmessage,
				IsSuccess: isSuccess,
				Error:     err,
			},
		})
	}

	decryptedToken, err := encryptDecrypt.Decrypt(tokenString, utils.GetEnv("SECRET_KEY"))
	if err != nil {
		fmt.Println("Can't decrypt token: ", tokenString)
		return c.Status(419).JSON(response.ResponseModel{
			RetCode: "419",
			Message: status.RetCode419,
			Data: errors.ErrorModel{
				Message:   "Invalid Authentication. You are not authorized to access this feature.",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	// Split the token into its parts
	parts := strings.Split(decryptedToken, ".")
	if len(parts) != 3 {
		fmt.Println("Invalid token format: ", parts)
		return c.Status(419).JSON(response.ResponseModel{
			RetCode: "419",
			Message: status.RetCode419,
			Data: errors.ErrorModel{
				Message:   "Invalid Authentication. You are not authorized to access this feature.",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	// Base64 decode the payload (second part)
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		fmt.Println("Error decoding payload:", payload, err)
		return c.Status(419).JSON(response.ResponseModel{
			RetCode: "419",
			Message: status.RetCode419,
			Data: errors.ErrorModel{
				Message:   "Invalid Authentication. You are not authorized to access this feature.",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	// Unmarshal the JSON payload into a struct
	var claims JWTClaims
	if err := json.Unmarshal(payload, &claims); err != nil {
		fmt.Println("Error unmarshaling claims:", err)
		return c.Status(419).JSON(response.ResponseModel{
			RetCode: "419",
			Message: status.RetCode419,
			Data: errors.ErrorModel{
				Message:   "Invalid Authentication. You are not authorized to access this feature.",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	// Print the extracted claims
	// fmt.Println("Rolename:", claims.Rolename)
	// fmt.Println("StaffID:", claims.StaffID)
	// fmt.Println("Username:", claims.Username)
	// fmt.Println("Expiration Time:", claims.Exp)

	var superAdminRoles string
	if err := db.Raw("SELECT * FROM public.getroles(0)").Scan(&superAdminRoles).Error; err != nil {
		fmt.Println("Error fetching token duration:", err)
		return c.Status(419).JSON(response.ResponseModel{
			RetCode: "419",
			Message: status.RetCode419,
			Data: errors.ErrorModel{
				Message:   "Invalid Authentication. You are not authorized to access this feature.",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	// fmt.Println("Token Role Name: " + claims.Rolename)
	// fmt.Println("SuperAdmin Role Name: " + superAdminRoles)

	if claims.Rolename != superAdminRoles {
		fmt.Println("Not a SuperAdmin")
		return c.Status(419).JSON(response.ResponseModel{
			RetCode: "419",
			Message: status.RetCode419,
			Data: errors.ErrorModel{
				Message:   "Invalid Authentication. You are not a SuperAdmin. You are not authorized to access this feature.",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	currentTime := time.Now().Unix()
	tokenCurrentDuration := currentTime - claims.Exp

	var tokenDuration int64
	if err := db.Raw("SELECT seconds FROM authentication.expiration WHERE id = 2").Scan(&tokenDuration).Error; err != nil {
		fmt.Println("Error fetching token duration:", err)
		return c.Status(419).JSON(response.ResponseModel{
			RetCode: "419",
			Message: status.RetCode419,
			Data: errors.ErrorModel{
				Message:   "Invalid Authentication. You are not authorized to access this feature.",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	fmt.Println("------------------------------------------------------------------------------------------------------------------------")
	fmt.Println("\nToken Role Name: " + claims.Rolename)
	fmt.Println("Token Validity Duration: ", tokenDuration)
	fmt.Println("Token Current Duration: ", tokenCurrentDuration)
	fmt.Println("\n------------------------------------------------------------------------------------------------------------------------")

	if tokenCurrentDuration >= tokenDuration {
		fmt.Println("Token has expired:", claims)
		return c.Status(419).JSON(response.ResponseModel{
			RetCode: "419",
			Message: status.RetCode419,
			Data: errors.ErrorModel{
				Message:   "Token has expired. Sign in required.",
				IsSuccess: false,
				Error:     nil,
			},
		})
	}

	return c.Next()
}

func ValidateAdminToken(c *fiber.Ctx) error {
	db := database.DB
	authHeader := c.Get("Authorization")

	if authHeader == "" {
		fmt.Println("Missing or invalid token: ", authHeader)
		return c.Status(419).JSON(response.ResponseModel{
			RetCode: "419",
			Message: status.RetCode419,
			Data: errors.ErrorModel{
				Message:   "Invalid Authentication. You are not authorized to access this feature.",
				IsSuccess: false,
				Error:     nil,
			},
		})
	}

	// Remove "Bearer " prefix if it exists
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		fmt.Println("Invalid token format: ", authHeader)
		return c.Status(419).JSON(response.ResponseModel{
			RetCode: "419",
			Message: status.RetCode419,
			Data: errors.ErrorModel{
				Message:   "Invalid Authentication. You are not authorized to access this feature.",
				IsSuccess: false,
				Error:     nil,
			},
		})
	}

	isSuccess, retCodeInt, retCode, tstatus, tmessage, err := sharedfunctions.ValidateToken(tokenString)
	if err != nil {
		return c.Status(retCodeInt).JSON(response.ResponseModel{
			RetCode: retCode,
			Message: tstatus,
			Data: errors.ErrorModel{
				Message:   tmessage,
				IsSuccess: isSuccess,
				Error:     err,
			},
		})
	}

	decryptedToken, err := encryptDecrypt.Decrypt(tokenString, utils.GetEnv("SECRET_KEY"))
	if err != nil {
		fmt.Println("Can't decrypt token: ", tokenString)
		return c.Status(419).JSON(response.ResponseModel{
			RetCode: "419",
			Message: status.RetCode419,
			Data: errors.ErrorModel{
				Message:   "Invalid Authentication. You are not authorized to access this feature.",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	// Split the token into its parts
	parts := strings.Split(decryptedToken, ".")
	if len(parts) != 3 {
		fmt.Println("Invalid token format: ", parts)
		return c.Status(419).JSON(response.ResponseModel{
			RetCode: "419",
			Message: status.RetCode419,
			Data: errors.ErrorModel{
				Message:   "Invalid Authentication. You are not authorized to access this feature.",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	// Base64 decode the payload (second part)
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		fmt.Println("Error decoding payload:", payload, err)
		return c.Status(419).JSON(response.ResponseModel{
			RetCode: "419",
			Message: status.RetCode419,
			Data: errors.ErrorModel{
				Message:   "Invalid Authentication. You are not authorized to access this feature.",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	// Unmarshal the JSON payload into a struct
	var claims JWTClaims
	if err := json.Unmarshal(payload, &claims); err != nil {
		fmt.Println("Error unmarshaling claims:", err)
		return c.Status(419).JSON(response.ResponseModel{
			RetCode: "419",
			Message: status.RetCode419,
			Data: errors.ErrorModel{
				Message:   "Invalid Authentication. You are not authorized to access this feature.",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	// Print the extracted claims
	// fmt.Println("Rolename:", claims.Rolename)
	// fmt.Println("StaffID:", claims.StaffID)
	// fmt.Println("Username:", claims.Username)
	// fmt.Println("Expiration Time:", claims.Exp)

	var superAdminRoles []string
	if err := db.Raw("SELECT * FROM public.getroles(1)").Scan(&superAdminRoles).Error; err != nil {
		fmt.Println("Error fetching token duration:", err)
		return c.Status(419).JSON(response.ResponseModel{
			RetCode: "419",
			Message: status.RetCode419,
			Data: errors.ErrorModel{
				Message:   "Invalid Authentication. You are not authorized to access this feature.",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	// fmt.Println(superAdminRoles)

	validRole := slices.Contains(superAdminRoles, claims.Rolename)

	if !validRole {
		return c.Status(419).JSON(response.ResponseModel{
			RetCode: "419",
			Message: status.RetCode419,
			Data: errors.ErrorModel{
				Message:   "Invalid Authentication. You are not a SuperAdmin or Admin. You are not authorized to access this feature.",
				IsSuccess: false,
				Error:     nil,
			},
		})
	}

	currentTime := time.Now().Unix()
	tokenCurrentDuration := currentTime - claims.Exp

	var tokenDuration int64
	if err := db.Raw("SELECT seconds FROM authentication.expiration WHERE id = 2").Scan(&tokenDuration).Error; err != nil {
		fmt.Println("Error fetching token duration:", err)
		return c.Status(419).JSON(response.ResponseModel{
			RetCode: "419",
			Message: status.RetCode419,
			Data: errors.ErrorModel{
				Message:   "Invalid Authentication. You are not authorized to access this feature.",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	fmt.Println("------------------------------------------------------------------------------------------------------------------------")
	fmt.Println("\nToken Role Name: " + claims.Rolename)
	fmt.Println("Token Validity Duration: ", tokenDuration)
	fmt.Println("Token Current Duration: ", tokenCurrentDuration)
	fmt.Println("\n------------------------------------------------------------------------------------------------------------------------")

	if tokenCurrentDuration >= tokenDuration {
		fmt.Println("Token has expired:", claims)
		return c.Status(419).JSON(response.ResponseModel{
			RetCode: "419",
			Message: status.RetCode419,
			Data: errors.ErrorModel{
				Message:   "Token has expired. Sign in required.",
				IsSuccess: false,
				Error:     nil,
			},
		})
	}

	return c.Next()
}

func ValidateUserToken(c *fiber.Ctx) error {
	db := database.DB
	authHeader := c.Get("Authorization")

	if authHeader == "" {
		fmt.Println("Missing or invalid token: ", authHeader)
		return c.Status(419).JSON(response.ResponseModel{
			RetCode: "419",
			Message: status.RetCode419,
			Data: errors.ErrorModel{
				Message:   "Invalid Authentication. You are not authorized to access this feature.",
				IsSuccess: false,
				Error:     nil,
			},
		})
	}

	// Remove "Bearer " prefix if it exists
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		fmt.Println("Invalid token format: ", authHeader)
		return c.Status(419).JSON(response.ResponseModel{
			RetCode: "419",
			Message: status.RetCode419,
			Data: errors.ErrorModel{
				Message:   "Invalid Authentication. You are not authorized to access this feature.",
				IsSuccess: false,
				Error:     nil,
			},
		})
	}

	isSuccess, retCodeInt, retCode, tstatus, tmessage, err := sharedfunctions.ValidateToken(tokenString)
	if err != nil {
		return c.Status(retCodeInt).JSON(response.ResponseModel{
			RetCode: retCode,
			Message: tstatus,
			Data: errors.ErrorModel{
				Message:   tmessage,
				IsSuccess: isSuccess,
				Error:     err,
			},
		})
	}

	decryptedToken, err := encryptDecrypt.Decrypt(tokenString, utils.GetEnv("SECRET_KEY"))
	if err != nil {
		fmt.Println("Can't decrypt token: ", tokenString)
		return c.Status(419).JSON(response.ResponseModel{
			RetCode: "419",
			Message: status.RetCode419,
			Data: errors.ErrorModel{
				Message:   "Invalid Authentication. You are not authorized to access this feature.",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	// Split the token into its parts
	parts := strings.Split(decryptedToken, ".")
	if len(parts) != 3 {
		fmt.Println("Invalid token format: ", parts)
		return c.Status(419).JSON(response.ResponseModel{
			RetCode: "419",
			Message: status.RetCode419,
			Data: errors.ErrorModel{
				Message:   "Invalid Authentication. You are not authorized to access this feature.",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	// Base64 decode the payload (second part)
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		fmt.Println("Error decoding payload:", payload, err)
		return c.Status(419).JSON(response.ResponseModel{
			RetCode: "419",
			Message: status.RetCode419,
			Data: errors.ErrorModel{
				Message:   "Invalid Authentication. You are not authorized to access this feature.",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	// Unmarshal the JSON payload into a struct
	var claims JWTClaims
	if err := json.Unmarshal(payload, &claims); err != nil {
		fmt.Println("Error unmarshaling claims:", err)
		return c.Status(419).JSON(response.ResponseModel{
			RetCode: "419",
			Message: status.RetCode419,
			Data: errors.ErrorModel{
				Message:   "Invalid Authentication. You are not authorized to access this feature.",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	// Print the extracted claims
	// fmt.Println("Rolename:", claims.Rolename)
	// fmt.Println("StaffID:", claims.StaffID)
	// fmt.Println("Username:", claims.Username)
	// fmt.Println("Expiration Time:", claims.Exp)

	var roles []string
	if err := db.Raw("SELECT * FROM public.getroles(2)").Scan(&roles).Error; err != nil {
		fmt.Println("Error fetching token duration:", err)
		return c.Status(419).JSON(response.ResponseModel{
			RetCode: "419",
			Message: status.RetCode419,
			Data: errors.ErrorModel{
				Message:   "Invalid Authentication. You are not authorized to access this feature.",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	// fmt.Println(roles)

	validRole := slices.Contains(roles, claims.Rolename)

	if !validRole {
		return c.Status(419).JSON(response.ResponseModel{
			RetCode: "419",
			Message: status.RetCode419,
			Data: errors.ErrorModel{
				Message:   "Invalid Authentication. You are not authorized to access this feature.",
				IsSuccess: false,
				Error:     nil,
			},
		})
	}

	currentTime := time.Now().Unix()
	tokenCurrentDuration := currentTime - claims.Exp

	var tokenDuration int64
	if err := db.Raw("SELECT seconds FROM authentication.expiration WHERE id = 3").Scan(&tokenDuration).Error; err != nil {
		fmt.Println("Error fetching token duration:", err)
		return c.Status(419).JSON(response.ResponseModel{
			RetCode: "419",
			Message: status.RetCode419,
			Data: errors.ErrorModel{
				Message:   "Invalid Authentication. You are not authorized to access this feature.",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	fmt.Println("------------------------------------------------------------------------------------------------------------------------")
	fmt.Println("\nToken Role Name: " + claims.Rolename)
	fmt.Println("Token Validity Duration: ", tokenDuration)
	fmt.Println("Token Current Duration: ", tokenCurrentDuration)
	fmt.Println("\n------------------------------------------------------------------------------------------------------------------------")

	if tokenCurrentDuration >= tokenDuration {
		fmt.Println("Token has expired:", claims)
		return c.Status(419).JSON(response.ResponseModel{
			RetCode: "419",
			Message: status.RetCode419,
			Data: errors.ErrorModel{
				Message:   "Token has expired. Sign in required.",
				IsSuccess: false,
				Error:     nil,
			},
		})
	}

	return c.Next()
}
