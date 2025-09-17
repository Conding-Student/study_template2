package admincardinc

import (
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"

	"github.com/gofiber/fiber/v2"
)

type ApprovingAuthority struct {
	Amount        string `json:"amount"`
	Role          string `json:"role"`
	RecommendedBy string `json:"recommendedBy"`
	ApprovedBy    string `json:"approvedBy"`
}

type ApprovingAuthorityResponse struct {
	ID            uint   `json:"id"`
	Amount        string `json:"amount"`
	Role          string `json:"role"`
	RecommendedBy string `json:"recommended_by"`
	ApprovedBy    string `json:"approved_by"`
}

func AddApprovingAuthority(c *fiber.Ctx) error {
	// Parse request body
	var appAuth ApprovingAuthority
	if err := c.BodyParser(&appAuth); err != nil {
		return c.Status(400).JSON(response.ResponseModel{
			RetCode: "400",
			Message: "Bad Request",
			Data: errors.ErrorModel{
				Message:   "Failed to parse request",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	// Call helper
	result, err := Add_ApprovingAuthority(appAuth)
	if err != nil {
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Internal Server Error",
			Data: errors.ErrorModel{
				Message:   "Problem connecting to database",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	return c.JSON(result)
}

func ViewApprovingAuthority(c *fiber.Ctx) error {
	result, err := Get_ApprovingAuthority()
	if err != nil {
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Internal Server Error",
			Data: errors.ErrorModel{
				Message:   "Problem connecting to database",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	return c.JSON(result)
}
func DeleteApprovingAuthority(c *fiber.Ctx) error {
	var appAuth ApprovingAuthority

	// Parse request body
	if err := c.BodyParser(&appAuth); err != nil {
		return c.Status(400).JSON(response.ResponseModel{
			RetCode: "400",
			Message: "Bad Request",
			Data: errors.ErrorModel{
				Message:   "Failed to parse request",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	// Call helper
	result, err := Delete_ApprovingAuthority(appAuth)
	if err != nil {
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Internal Server Error",
			Data: errors.ErrorModel{
				Message:   "Problem connecting to database",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	return c.JSON(result)
}

// func AddApprovingAuthority(c *fiber.Ctx) error {
// 	// Parse the JSON request body into an ApprovingAuthority struct
// 	var appAuth ApprovingAuthority
// 	if err := c.BodyParser(&appAuth); err != nil {
// 		fmt.Println("Bad Request")
// 		fmt.Println("Failed to parse request")
// 		return c.Status(400).JSON(response.ResponseModel{
// 			RetCode: "400",
// 			Message: "Bad Request",
// 			Data: errors.ErrorModel{
// 				Message:   "Failed to parse request",
// 				IsSuccess: false,
// 				Error:     err,
// 			},
// 		})
// 	}

// 	// Begin transaction
// 	db := database.DB.Begin()

// 	// Check if the role already exists
// 	var count int
// 	err := db.Raw("SELECT COUNT(*) FROM public.approving_authority WHERE role = ? AND recommended_by = ? AND approved_by = ?", appAuth.Role, appAuth.RecommendedBy, appAuth.ApprovedBy).Row().Scan(&count)
// 	if err != nil {
// 		db.Rollback()
// 		log.Println(err)
// 		fmt.Println("retCode 500")
// 		fmt.Println("Internal Server Error")
// 		fmt.Println("Problem connecting to database")
// 		return c.Status(500).JSON(response.ResponseModel{
// 			RetCode: "500",
// 			Message: "Internal Server Error",
// 			Data: errors.ErrorModel{
// 				Message:   "Problem connecting to database",
// 				IsSuccess: false,
// 				Error:     err,
// 			},
// 		})
// 	}

// 	var message string

// 	if count > 0 {
// 		// If the role exists, update the existing record
// 		updateStatement := `
// 			UPDATE public.approving_authority
// 			SET amount = ?, role = ?, recommended_by = ?, approved_by = ?
// 			WHERE role = ?
// 		`
// 		result := db.Exec(updateStatement, appAuth.Amount, appAuth.Role, appAuth.RecommendedBy, appAuth.ApprovedBy, appAuth.Role)
// 		if result.Error != nil {
// 			db.Rollback()
// 			log.Println(result.Error)
// 			fmt.Println("retCode 500")
// 			fmt.Println("Internal Server Error")
// 			fmt.Println("Problem connecting to database")
// 			return c.Status(500).JSON(response.ResponseModel{
// 				RetCode: "500",
// 				Message: "Internal Server Error",
// 				Data: errors.ErrorModel{
// 					Message:   "Problem connecting to database",
// 					IsSuccess: false,
// 					Error:     result.Error,
// 				},
// 			})
// 		}
// 		message = "Approving authority updated successfully"
// 	} else {
// 		// If the role doesn't exist, insert a new record
// 		insertStatement := `
// 			INSERT INTO public.approving_authority (amount, role, recommended_by, approved_by)
// 			VALUES (?, ?, ?, ?)
// 		`
// 		result := db.Exec(insertStatement, appAuth.Amount, appAuth.Role, appAuth.RecommendedBy, appAuth.ApprovedBy)
// 		if result.Error != nil {
// 			db.Rollback()
// 			log.Println(result.Error)
// 			fmt.Println("retCode 500")
// 			fmt.Println("Internal Server Error")
// 			fmt.Println("Problem connecting to database")
// 			return c.Status(500).JSON(response.ResponseModel{
// 				RetCode: "500",
// 				Message: "Internal Server Error",
// 				Data: errors.ErrorModel{
// 					Message:   "Problem connecting to database",
// 					IsSuccess: false,
// 					Error:     result.Error,
// 				},
// 			})
// 		}
// 		message = "Approving authority added successfully"
// 	}

// 	// Commit transaction
// 	db.Commit()

// 	// Return a success response
// 	return c.Status(200).JSON(response.ResponseModel{
// 		RetCode: "200",
// 		Message: message,
// 		Data:    appAuth,
// 	})
// }

// func ViewApprovingAuthority(c *fiber.Ctx) error {
// 	// Query for approving authority records
// 	rows, err := database.DB.Raw("SELECT * FROM public.approving_authority order by amount asc").Rows()
// 	if err != nil {
// 		database.DB.Rollback()
// 		log.Println(err)
// 		fmt.Println("retCode 500")
// 		fmt.Println("Internal Server Error")
// 		fmt.Println("Problem connecting to database")
// 		return c.Status(500).JSON(response.ResponseModel{
// 			RetCode: "500",
// 			Message: "Internal Server Error",
// 			Data: errors.ErrorModel{
// 				Message:   "Problem connecting to database",
// 				IsSuccess: false,
// 				Error:     err,
// 			},
// 		})
// 	}
// 	defer rows.Close()

// 	// Iterate through the rows and collect approving authority data
// 	var approvingAuthorities []ApprovingAuthorityResponse
// 	for rows.Next() {
// 		var aaResponse ApprovingAuthorityResponse
// 		if err := rows.Scan(&aaResponse.ID, &aaResponse.Amount, &aaResponse.Role, &aaResponse.RecommendedBy, &aaResponse.ApprovedBy); err != nil {
// 			database.DB.Rollback()
// 			log.Println(err)
// 			fmt.Println("retCode 500")
// 			fmt.Println("Internal Server Error")
// 			fmt.Println("Failed to scan row")
// 			return c.Status(500).JSON(response.ResponseModel{
// 				RetCode: "500",
// 				Message: "Internal Server Error",
// 				Data: errors.ErrorModel{
// 					Message:   "Failed to scan row",
// 					IsSuccess: false,
// 					Error:     err,
// 				},
// 			})
// 		}

// 		approvingAuthorities = append(approvingAuthorities, aaResponse)
// 	}

// 	// Check for any iteration errors
// 	if err := rows.Err(); err != nil {
// 		database.DB.Rollback()
// 		log.Println(err)
// 		fmt.Println("retCode 500")
// 		fmt.Println("Internal Server Error")
// 		fmt.Println("Iteration error")
// 		return c.Status(500).JSON(response.ResponseModel{
// 			RetCode: "500",
// 			Message: "Internal Server Error",
// 			Data: errors.ErrorModel{
// 				Message:   "Iteration error",
// 				IsSuccess: false,
// 				Error:     err,
// 			},
// 		})
// 	}

// 	// Return a success response with the approving authority data
// 	return c.Status(200).JSON(response.ResponseModel{
// 		RetCode: "200",
// 		Message: "Approving authorities retrieved successfully",
// 		Data:    approvingAuthorities,
// 	})
// }

// func DeleteApprovingAuthority(c *fiber.Ctx) error {
// 	// Parse the JSON request body into an ApprovingAuthority struct
// 	var appAuth ApprovingAuthority
// 	if err := c.BodyParser(&appAuth); err != nil {
// 		fmt.Println("Bad Request")
// 		fmt.Println("Failed to parse request")
// 		return c.Status(400).JSON(response.ResponseModel{
// 			RetCode: "400",
// 			Message: "Bad Request",
// 			Data: errors.ErrorModel{
// 				Message:   "Failed to parse request",
// 				IsSuccess: false,
// 				Error:     err,
// 			},
// 		})
// 	}

// 	// Begin transaction
// 	db := database.DB.Begin()

// 	// Construct the raw SQL select statement to check if the record exists
// 	selectStatement := `
// 		SELECT COUNT(*) FROM public.approving_authority
// 		WHERE amount = ? AND role = ? AND recommended_by = ? AND approved_by = ?
// 	`

// 	// Execute the raw SQL select statement to check if the record exists
// 	var count int
// 	if err := db.Raw(selectStatement, appAuth.Amount, appAuth.Role, appAuth.RecommendedBy, appAuth.ApprovedBy).Row().Scan(&count); err != nil {
// 		db.Rollback()
// 		log.Println(err)
// 		fmt.Println("retCode 500")
// 		fmt.Println("Internal Server Error")
// 		fmt.Println("Problem connecting to database")
// 		return c.Status(500).JSON(response.ResponseModel{
// 			RetCode: "500",
// 			Message: "Internal Server Error",
// 			Data: errors.ErrorModel{
// 				Message:   "Problem connecting to database",
// 				IsSuccess: false,
// 				Error:     err,
// 			},
// 		})
// 	}

// 	// If the record does not exist, return a message
// 	if count == 0 {
// 		db.Rollback()
// 		fmt.Println("retCode 404")
// 		fmt.Println("Bad Request")
// 		fmt.Println("Record does not exist")
// 		return c.Status(400).JSON(response.ResponseModel{
// 			RetCode: "404",
// 			Message: "Bad Request",
// 			Data: errors.ErrorModel{
// 				Message:   "Record does not exist",
// 				IsSuccess: false,
// 				Error:     db.Error,
// 			},
// 		})
// 	}

// 	// Construct the raw SQL delete statement
// 	deleteStatement := `
// 		DELETE FROM public.approving_authority
// 		WHERE amount = ? AND role = ? AND recommended_by = ? AND approved_by = ?
// 	`

// 	// Execute the raw SQL delete statement
// 	result := db.Exec(deleteStatement, appAuth.Amount, appAuth.Role, appAuth.RecommendedBy, appAuth.ApprovedBy)
// 	if result.Error != nil {
// 		db.Rollback()
// 		log.Println(result.Error)
// 		fmt.Println("retCode 500")
// 		fmt.Println("Internal Server Error")
// 		fmt.Println("Problem connecting to database")
// 		return c.Status(500).JSON(response.ResponseModel{
// 			RetCode: "500",
// 			Message: "Internal Server Error",
// 			Data: errors.ErrorModel{
// 				Message:   "Problem connecting to database",
// 				IsSuccess: false,
// 				Error:     result.Error,
// 			},
// 		})
// 	}

// 	// Commit transaction
// 	db.Commit()

// 	// Return a success response
// 	return c.Status(200).JSON(response.ResponseModel{
// 		RetCode: "200",
// 		Message: "Approving authority deleted successfully",
// 		Data:    nil,
// 	})
// }
