package handler

import (
	"chatbot/pkg/models/model"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func ClientInformationForm(c *fiber.Ctx) error {
	// Initialize a ClientInformationForm model
	cif := new(model.ClientInformationForm)

	// Parse the request body into the ClientInformationForm struct
	if err := c.BodyParser(cif); err != nil {
		fmt.Println("Error parsing request body:", err)
		return c.Status(401).JSON(fiber.Map{
			"status":  messageStatus,
			"message": "Invalid input. Please review your data.",
		})
	}

	// Simulating a successful registration
	successfulRegistration := true

	if successfulRegistration {
		// Print the received data only when the registration is successful
		fmt.Printf("Received Data: %+v\n", cif)

		// TODO: Save the user information to the local database (if necessary)

		// Send the response back to the client
		return c.Status(201).JSON(fiber.Map{
			"status":  "Success",
			"message": "Client information successfully sent to Unit Manager.",
		})
	}

	// If registration is not successful, return an error response
	return c.Status(500).JSON(fiber.Map{
		"status":  messageStatus,
		"message": "Failed to sent client information.",
	})
}
