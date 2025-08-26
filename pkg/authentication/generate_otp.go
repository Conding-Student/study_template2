package authentication

import (
	"chatbot/pkg/sharedfunctions"
	"fmt"
)

func GenerateOtp(email, staffID string, code int) (string, error) {

	isSuccess, message, otp, err := sharedfunctions.GenOTP(staffID, code)

	if !isSuccess {
		return message, err
	}

	switch code {
	case 0:
		if err := sharedfunctions.SendEmail(email, otp, "", 1); err != nil {
			fmt.Println(err)
			return "We encounter an error while sending a One Time Password (OTP) in your email. Please try again later.", err
		}
	case 1:
		if err := sharedfunctions.SendEmail(email, otp, "", 3); err != nil {
			fmt.Println(err)
			return "We encounter an error while sending a One Time Password (OTP) in your email. Please try again later.", err
		}
	case 2:
		if err := sharedfunctions.SendEmail(email, otp, "", 4); err != nil {
			fmt.Println(err)
			return "We encounter an error while sending a One Time Password (OTP) in your email. Please try again later.", err
		}
	}

	return message, nil
}
