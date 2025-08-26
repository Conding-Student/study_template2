package authentication

import (
	"chatbot/pkg/utils/go-utils/database"
	"fmt"
)

type GetOTP struct {
	Verified bool   `gorm:"verified:"`
	Message  string `gorm:"message:"`
}

func VerifyOTP(otp, staffID string, validType int) (string, bool) {
	db := database.DB
	verfyOTP := new(GetOTP)

	if err := db.Raw("SELECT * FROM public.verifyotp($1, $2, $3)", staffID, otp, validType).Scan(&verfyOTP).Error; err != nil {
		return "Oops! There was an error while verifying your OTP. Please try again.", false
	}

	fmt.Println("------------------------------------------------------------------")
	fmt.Println("\nRequested OTP is Valid: ", verfyOTP.Verified)
	fmt.Println("Message: ", verfyOTP.Message)
	fmt.Println("\n------------------------------------------------------------------")

	return verfyOTP.Message, verfyOTP.Verified
}

func VerifyOTPForAccountCreation(otp, staffID string, validType int) (string, bool) {
	db := database.DB
	verfyOTP := new(GetOTP)

	if err := db.Raw("SELECT * FROM public.verifyotp($1, $2, $3)", staffID, otp, validType).Scan(&verfyOTP).Error; err != nil {
		return "Oops! There was an error while verifying your OTP. Please try again.", false
	}

	fmt.Println("------------------------------------------------------------------")
	fmt.Println("\nOTP from account creation is Valid: ", verfyOTP.Verified)
	fmt.Println("Message: ", verfyOTP.Message)
	fmt.Println("\n------------------------------------------------------------------")

	return verfyOTP.Message, verfyOTP.Verified
}
