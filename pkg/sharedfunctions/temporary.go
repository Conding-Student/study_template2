package sharedfunctions

import (
	"chatbot/pkg/utils"
	"chatbot/pkg/utils/go-utils/database"
	"chatbot/pkg/utils/go-utils/encryptDecrypt"
	"fmt"
	"log"
	"net/smtp"
	"strconv"
	"time"

	"math/rand"
)

var secretKey = utils.GetEnv("SECRET_KEY")

//////////////////////////////////////////////////////////////////////////////////////////////////////
// 			        	     Ito ay para sa paggenerate ng Random Password	     		   			//
//////////////////////////////////////////////////////////////////////////////////////////////////////

func GenerateRandomPassword(length int) string {
	// Characters to choose from when generating the password
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// Seed the random number generator
	src := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(src)

	// Generate the random password
	password := make([]byte, length)
	for i := range password {
		password[i] = charset[rng.Intn(len(charset))]
	}

	return string(password)
}

//////////////////////////////////////////////////////////////////////////////////////////////////////
// 			        	  Ito ay para sa pagsend ng Temporary Password sa email	     	   			//
//////////////////////////////////////////////////////////////////////////////////////////////////////

func SendEmail(toEmail, recoveryPasswordOrOTP, username string, code int) error {
	db := database.DB

	var result map[string]any
	if err := db.Raw("SELECT * FROM get_emailsendercreds()").Scan(&result).Error; err != nil {
		return err
	}

	// Email configuration
	smtpServer, err := encryptDecrypt.Decrypt(GetStringFromMap(result, "smtp_server"), secretKey)
	if err != nil {
		return err
	}
	smtpPort := GetIntFromMap(result, "smtp_port")
	senderEmail, err := encryptDecrypt.Decrypt(GetStringFromMap(result, "sender_email"), secretKey)
	if err != nil {
		return err
	}
	senderPassword, err := encryptDecrypt.Decrypt(GetStringFromMap(result, "password"), secretKey)
	if err != nil {
		return err
	}
	headerImage := GetStringFromMap(result, "headerimage")
	headerImageName := GetStringFromMap(result, "headerimagename")

	// Compose the email message
	to := []string{toEmail}
	var body string
	var subject string

	switch code {
	case 0:
		subject = "Temporary Password"
		body = `
		<html>
			<head>
			<style>
			/* Define your CSS styles here */
			body {
			font-family: Arial, sans-serif;
			background-color: #f5f5f5;
			}
			.container {
			max-width: 600px;
			margin: 0 auto;
			padding: 20px;
			background-color: #ffffff;
			border-radius: 10px;
			box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
			}
			.header {
			text-align: center;
			margin-bottom: 20px;
			}
			.message {
			padding: 20px;
			background-color: #f0f0f0;
			border-radius: 5px;
			}
			</style>
		</head>
		<body>
			<div class="container">
				<div class="header">
					<img src="` + headerImage + `" alt="` + headerImageName + `" style="width:100px;"/>
					<h2>Temporary Password</h2>
				</div>
				<div class="message">
					<p>Attention: Your temporary password has been provided for account recovery purposes. Please note that temporary passwords are susceptible to security risks and should be changed immediately upon receipt.</p>
					<p>Your account's security is of paramount importance to us. Temporary passwords, if left unchanged, can pose vulnerabilities to unauthorized access and compromise your personal information. Hackers and malicious actors often target accounts with temporary passwords due to their temporary nature and predictable patterns.</p>
					<p>To ensure the utmost security of your account, we strongly advise changing your temporary password as soon as possible. When selecting a new password, please adhere to the following best practices for password security:</p>
					<ul>
						<li><strong>Complexity:</strong> Use a combination of uppercase letters, lowercase letters, numbers, and special characters to create a strong and unique password.</li>
						<li><strong>Length:</strong> Opt for a password that is at least 8 characters long to enhance its complexity and resilience against brute-force attacks.</li>
						<li><strong>Avoid Common Patterns:</strong> Refrain from using easily guessable information such as names, birthdays, or common phrases as your password.</li>
						<li><strong>Regular Updates:</strong> Routinely change your password every few months to mitigate the risk of unauthorized access.</li>
					</ul>
					<p>Remember, your password serves as the first line of defense against cyber threats. By following these guidelines and promptly updating your temporary password, you play a crucial role in safeguarding your account from potential security breaches.</p>
					<p>This is your CA-GABAY temporary password: <strong>` + recoveryPasswordOrOTP + `</strong></p>
				</div>
				<div>
					<p><em style="color: red;">Please note: This is an automated email. Please do not reply to this email address.</em></p>
				</div>
			</div>
		</body>
		</html>
		`
	case 1:
		subject = "One-Time Password (OTP)"
		body = `
		<html>
			<head>
			<style>
			/* Define your CSS styles here */
			body {
			font-family: Arial, sans-serif;
			background-color: #f5f5f5;
			}
			.container {
			max-width: 600px;
			margin: 0 auto;
			padding: 20px;
			background-color: #ffffff;
			border-radius: 10px;
			box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
			}
			.header {
			text-align: center;
			margin-bottom: 20px;
			}
			.message {
			padding: 20px;
			background-color: #f0f0f0;
			border-radius: 5px;
			}
			</style>
		</head>
		<body>
			<div class="container">
				<div class="header">
					<img src="` + headerImage + `" alt="` + headerImageName + `" style="width:100px;"/>
					<h2>Email Verification</h2>
				</div>
				<div class="message">
					<p>Dear User,</p>
					<p>You have requested a One-Time Password (OTP) for email verification. Please use the OTP provided below to verify your email address before proceeding to Account Creation page:</p>
					<p>Your CA-GABAY One-Time Password (OTP): <strong>` + recoveryPasswordOrOTP + `</strong></p>
					<p><strong>What to do next:</strong></p>
						<ul>
							<li>Enter the OTP on the email verification page to complete the verification process.</li>
							<li>The OTP will expire in 5 minutes. If it expires, please request a new OTP.</li>
							<li>If you did not request this OTP, please ignore this email. No further action is required.</li>
							<li>If you suspect any unauthorized activity, please contact our support team immediately.</li>
						</ul>
					<p>Your account's security is very important to us. By verifying your email address, you help protect your account from unauthorized access and ensure you receive important notifications.</p>
				</div>
				<div>
					<p><em style="color: red;">Please note: This is an automated email. Please do not reply to this email address.</em></p>
				</div>
			</div>
		</body>
		</html>
		`
	case 2:
		subject = "Temporary Credentials"
		body = `
		<html>
			<head>
			<style>
			/* Define your CSS styles here */
			body {
			font-family: Arial, sans-serif;
			background-color: #f5f5f5;
			}
			.container {
			max-width: 600px;
			margin: 0 auto;
			padding: 20px;
			background-color: #ffffff;
			border-radius: 10px;
			box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
			}
			.header {
			text-align: center;
			margin-bottom: 20px;
			}
			.message {
			padding: 20px;
			background-color: #f0f0f0;
			border-radius: 5px;
			}
			</style>
		</head>
		<body>
			<div class="container">
				<div class="header">
					<img src="` + headerImage + `" alt="` + headerImageName + `" style="width:100px;"/>
					<h2>Temporary Credentials</h2>
				</div>
				<div class="message">
					<p>Attention: You have been provided with temporary credentials for your new account setup. Please be aware that these temporary credentials, especially the password, are vulnerable to security risks and should be changed immediately upon receipt.</p>
					<p>Your account's security is our top priority. If temporary passwords are not updated promptly, they can leave your account susceptible to unauthorized access, potentially compromising your personal information. Hackers often target accounts with temporary passwords due to their temporary nature and predictable patterns.</p>
					<p>To secure your new account, we strongly recommend changing your temporary password immediately. When creating a new password, please follow these best practices for password security:</p>
					<ul>
						<li><strong>Complexity:</strong> Use a combination of uppercase and lowercase letters, numbers, and special characters to create a strong and unique password.</li>
						<li><strong>Length:</strong> Choose a password that is at least 8 characters long to increase its complexity and resistance to brute-force attacks.</li>
						<li><strong>Avoid Common Patterns:</strong> Do not use easily guessable information such as names, birthdays, or common phrases in your password.</li>
						<li><strong>Regular Updates:</strong> Change your password regularly, ideally every few months, to minimize the risk of unauthorized access.</li>
					</ul>
					<p>Your password is the first line of defense against cyber threats. By following these guidelines and updating your temporary password promptly, you play a crucial role in protecting your new account from potential security breaches.</p>
					<p>This is your CA-GABAY temporary credentials:</p>
					<p style="margin-left: 25px;">
						Username: <strong>` + username + `</strong><br>
						Temporary Password: <strong>` + recoveryPasswordOrOTP + `</strong>
					</p>
				</div>
				<div>
					<p><em style="color: red;">Please note: This is an automated email. Do not reply to this email address.</em></p>
				</div>
			</div>
		</body>
		</html>
		`
	case 3:
		subject = "One-Time Password (OTP)"
		body = `
		<html>
			<head>
			<style>
			/* Define your CSS styles here */
			body {
			font-family: Arial, sans-serif;
			background-color: #f5f5f5;
			}
			.container {
			max-width: 600px;
			margin: 0 auto;
			padding: 20px;
			background-color: #ffffff;
			border-radius: 10px;
			box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
			}
			.header {
			text-align: center;
			margin-bottom: 20px;
			}
			.message {
			padding: 20px;
			background-color: #f0f0f0;
			border-radius: 5px;
			}
			</style>
			</head>
			<body>
				<div class="container">
					<div class="header">
						<img src="` + headerImage + `" alt="` + headerImageName + `" style="width:100px;"/>
						<h2>Device Info Update</h2>
					</div>
					<div class="message">
						<p>Dear User,</p>
						<p>You have requested a One-Time Password (OTP) to update the device associated with your account. Please use the OTP provided below to proceed with the device update process:</p>
						<p>Your CA-GABAY One-Time Password (OTP): <strong>` + recoveryPasswordOrOTP + `</strong></p>
						<p><strong>What to do next:</strong></p>
							<ul>
								<li>Enter the OTP on the device verification page to complete the update process.</li>
								<li>The OTP will expire in 5 minutes. If it expires, please request a new OTP.</li>
								<li>If you did not request this OTP, please ignore this email. No further action is required.</li>
								<li>If you suspect any unauthorized activity, please contact our support team immediately.</li>
							</ul>
						<p>Your account's security is our priority. Verifying your device helps protect your account from unauthorized access and ensures a secure user experience.</p>
					</div>
					<div>
						<p><em style="color: red;">Please note: This is an automated email. Please do not reply to this email address.</em></p>
					</div>
				</div>
			</body>
		</html>
		`
	case 4:
		subject = "One-Time Password (OTP)"
		body = `
		<html>
			<head>
			<style>
			/* Define your CSS styles here */
			body {
			font-family: Arial, sans-serif;
			background-color: #f5f5f5;
			}
			.container {
			max-width: 600px;
			margin: 0 auto;
			padding: 20px;
			background-color: #ffffff;
			border-radius: 10px;
			box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
			}
			.header {
			text-align: center;
			margin-bottom: 20px;
			}
			.message {
			padding: 20px;
			background-color: #f0f0f0;
			border-radius: 5px;
			}
			</style>
			</head>
			<body>
				<div class="container">
					<div class="header">
						<img src="` + headerImage + `" alt="` + headerImageName + `" style="width:100px;"/>
						<h2>Device Info Update</h2>
					</div>
					<div class="message">
						<p>Dear User,</p>
						<p>You have requested a One-Time Password (OTP) to update CA-GABAY PIN. Please use the OTP provided below to proceed:</p>
						<p>Your CA-GABAY One-Time Password (OTP): <strong>` + recoveryPasswordOrOTP + `</strong></p>
						<p><strong>What to do next:</strong></p>
							<ul>
								<li>Enter the OTP on the verification page to complete the update process.</li>
								<li>The OTP will expire in 5 minutes. If it expires, please request a new OTP.</li>
								<li>If you did not request this OTP, please ignore this email. No further action is required.</li>
								<li>If you suspect any unauthorized activity, please contact our support team immediately.</li>
							</ul>
						<p>Your account's security is our priority. Verifying your device helps protect your account from unauthorized access and ensures a secure user experience.</p>
					</div>
					<div>
						<p><em style="color: red;">Please note: This is an automated email. Please do not reply to this email address.</em></p>
					</div>
				</div>
			</body>
		</html>
		`
	}

	message := "Subject: " + subject + "\r\n" +
		"To: " + toEmail + "\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n" +
		// "Content-Type: text/plain; charset=UTF-8\r\n" +
		"Reply-To: noreply@ca-gabay.com\r\n" +
		"\r\n" + body

	// Connect to the SMTP server
	auth := smtp.PlainAuth("", senderEmail, senderPassword, smtpServer)
	err = smtp.SendMail(smtpServer+":"+strconv.Itoa(smtpPort), auth, senderEmail, to, []byte(message))
	if err != nil {
		log.Println("Error sending email:", err)
		return err
	}

	fmt.Println("----------------------------------------------------------------------------------------------")
	fmt.Println("\n"+subject+" sent successfully to:", toEmail)
	fmt.Println("\n----------------------------------------------------------------------------------------------")
	return nil
}

//////////////////////////////////////////////////////////////////////////////////////////////////////
// 			        	  			Ito ay para sa temporary password	 			    	   		//
//////////////////////////////////////////////////////////////////////////////////////////////////////

// func TemporaryPassword(email, username string) (string, error) {
// 	// Generate a random valid password
// 	randomPassword := GenerateRandomPassword(7)

// 	// Combine "CA-GABAY" with the random password
// 	combinedPassword := "CA-GABAY#FDSAPInc.@8=" + randomPassword

// 	secretKey := utils.GetEnv("SECRET_KEY")
// 	encryptedPassword, err := encryptDecrypt.Encrypt(combinedPassword, secretKey)
// 	if err != nil {
// 		return "Error encrypting password", err
// 	}

// 	password := encryptedPassword

// 	// Send the combined password as the recovery password email
// 	if err := SendEmail(email, combinedPassword, username, 2); err != nil {
// 		return "Error sending recovery email!", err
// 	}

// 	return password, nil
// }

func CreateTemporaryUsername(firstName, middleName, lastName string) (string, error) {
	if len(firstName) == 0 || len(middleName) == 0 || len(lastName) == 0 {
		return "", fmt.Errorf("error generating temporary username")
	}

	// Get the first letter of the first name and concatenate it with the last name
	tempUsername := firstName + "." + string(middleName[0]) + "." + lastName

	// Ensure the username is at least 8 characters long
	if len(tempUsername) < 8 {
		// Generate random characters to pad the username to at least 8 characters
		paddingLength := 8 - len(tempUsername)
		tempUsername += GenerateRandomPassword(paddingLength)
	}

	fmt.Println("--------------------------------------------------------------------------------------------------")
	fmt.Println("\nSuccessful?:", true)
	fmt.Println("Message:", tempUsername)
	fmt.Println("\n------------------------------------------------------------------------------------------------")

	return tempUsername, nil
}
