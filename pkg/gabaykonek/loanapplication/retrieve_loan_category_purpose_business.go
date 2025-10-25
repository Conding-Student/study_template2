package loans

import (
	"chatbot/pkg/logs"
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"
	"chatbot/pkg/models/status"
	"chatbot/pkg/sharedfunctions"
	"chatbot/pkg/utils/go-utils/database"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

type LoanCategoryResponse struct {
	Lccode       string `json:"lccode"`
	LoanCategory string `json:"loan_category"`
}

type LoanPurposeCredentials struct {
	Lccode string `json:"lccode"`
}

type LoanPurposeResponse struct {
	Lpcode       string `json:"lpcode"`
	LoanCategory string `json:"loan_purpose"`
}

type BusinessTypeCredentials struct {
	ListType int    `json:"list_type"`
	Lpcode   string `json:"lpcode"`
}

type BusinessTypeResponse struct {
	Btcode       string `json:"btcode"`
	BusinessType string `json:"business_type"`
}

type BusinessTypesList struct {
	BusinessType string `json:"business_type"`
}

// For CA-GABAY and kPlus // No params
func GetLoanCategory(c *fiber.Ctx) error {
	db := database.DB
	id := c.Params("id")
	LOSFeature := "LOS - Retrieve Loan Category" + id

	var result map[string]any
	if err := db.Raw("SELECT * FROM gabaykonekfunc.getloan_category()").Scan(&result).Error; err != nil {
		logs.LOSLogs(c, LOSFeature, id, "500", err.Error())
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: status.RetCode500,
			Data: errors.ErrorModel{
				Message:   "Problem connecting to database",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	sharedfunctions.ConvertStringToJSONMap(result)
	results := sharedfunctions.GetMap(result, "getloan_category")

	fmt.Println("--------------------------------------------------------------------------------------------------")
	fmt.Println("Fetched Loan Category For CA-GABAY and kPlus")
	fmt.Println("------------------------------------------------------------------------------------------------")

	logs.LOSLogs(c, LOSFeature, id, "200", "Successfully fetch Loan Categories.")
	return c.JSON(results)
}

// For CA-GABAY // No params
func GetLoanPurpose(c *fiber.Ctx) error {
	db := database.DB
	id := c.Params("id")
	LOSFeature := "LOS - Retrieve Loan Purpose" + id

	var result map[string]any
	if err := db.Raw("SELECT * FROM gabaykonekfunc.getloan_purpose()").Scan(&result).Error; err != nil {
		logs.LOSLogs(c, LOSFeature, id, "500", err.Error())
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: status.RetCode500,
			Data: errors.ErrorModel{
				Message:   "Problem connecting to database",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	sharedfunctions.ConvertStringToJSONMap(result)
	results := sharedfunctions.GetMap(result, "getloan_purpose")

	fmt.Println("--------------------------------------------------------------------------------------------------")
	fmt.Println("Fetched Loan Purpose For CA-GABAY")
	fmt.Println("------------------------------------------------------------------------------------------------")

	return c.JSON(results)
}

// For CA-GABAY and kPlus // No params
func GetBusinessType(c *fiber.Ctx) error {
	db := database.DB
	id := c.Params("id")
	LOSFeature := "LOS - Retrieve Loan Purpose" + id

	var result map[string]any
	if err := db.Raw("SELECT * FROM gabaykonekfunc.getbusinesstype()").Scan(&result).Error; err != nil {
		logs.LOSLogs(c, LOSFeature, id, "500", err.Error())
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: status.RetCode500,
			Data: errors.ErrorModel{
				Message:   "Problem connecting to database",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	sharedfunctions.ConvertStringToJSONMap(result)
	businessType := sharedfunctions.GetList(result, "response")

	fmt.Println("--------------------------------------------------------------------------------------------------")
	fmt.Println("Fetched Business Type For CA-GABAY and kPlus")
	fmt.Println("------------------------------------------------------------------------------------------------")

	return c.Status(200).JSON(response.ResponseModel{
		RetCode: "200",
		Message: "Successful",
		Data:    businessType,
	})
}

// For kPlus
func LoanPurpose(c *fiber.Ctx) error {
	id := c.Params("id")
	LOSFeature := "LOS - Retrieve Loan Purpose" + id

	lccode := new(LoanPurposeCredentials)
	if err := c.BodyParser(lccode); err != nil {
		return c.Status(401).JSON(response.ResponseModel{
			RetCode: "401",
			Message: status.RetCode401,
			Data: errors.ErrorModel{
				Message:   "Failed to parse request",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	// Define a slice to hold the results
	var purposes []LoanPurposeResponse

	// Execute the raw SQL query with the loan category code parameter
	query := "SELECT lpcode, loan_purpose FROM business_loan_purpose.loan_purpose WHERE lccode = ?"
	rows, err := database.DB.Raw(query, lccode.Lccode).Rows()
	if err != nil {
		logs.LOSLogs(c, LOSFeature, id, "500", err.Error())
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Internal Server Error",
			Data: errors.ErrorModel{
				Message:   "Internal Server Error",
				IsSuccess: false,
				Error:     err,
			},
		})
	}
	defer rows.Close()

	// Iterate over the rows and scan the values into the struct
	for rows.Next() {
		var purpose LoanPurposeResponse
		if err := rows.Scan(&purpose.Lpcode, &purpose.LoanCategory); err != nil {
			logs.LOSLogs(c, LOSFeature, id, "500", err.Error())
			return c.Status(500).JSON(response.ResponseModel{
				RetCode: "500",
				Message: "Internal Server Error",
				Data: errors.ErrorModel{
					Message:   "Internal Server Error",
					IsSuccess: false,
					Error:     err,
				},
			})
		}
		purposes = append(purposes, purpose)
	}

	// Check for any errors encountered during iteration
	if err := rows.Err(); err != nil {
		logs.LOSLogs(c, LOSFeature, id, "500", err.Error())
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Internal Server Error",
			Data: errors.ErrorModel{
				Message:   "Internal Server Error",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	// Return the loan purposes as a JSON response
	return c.Status(200).JSON(response.ResponseModel{
		RetCode: "200",
		Message: "Successful",
		Data:    purposes,
	})
}

// For CA-GABAY
func BusinessType(c *fiber.Ctx) error {
	db := database.DB
	// Parse the request body to get the loan purpose code
	typeOfBusiness := new(BusinessTypeCredentials)
	if err := c.BodyParser(typeOfBusiness); err != nil {
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

	var businessTypes any

	switch typeOfBusiness.ListType {
	case 1:
		// Define a slice to hold the results
		businessTypes = []BusinessTypeResponse{}
		// Execute the raw SQL query with the loan purpose code parameter
		query := "SELECT btcode, business_type FROM business_loan_purpose.business_type WHERE lpcode = ?"
		rows, err := db.Raw(query, typeOfBusiness.Lpcode).Rows()
		if err != nil {
			log.Println(err)
			return c.Status(500).JSON(response.ResponseModel{
				RetCode: "500",
				Message: "Internal Server Error",
				Data: errors.ErrorModel{
					Message:   "Internal Server Error",
					IsSuccess: false,
					Error:     err,
				},
			})
		}
		defer rows.Close()

		// Iterate over the rows and scan the values into the struct
		for rows.Next() {
			var types BusinessTypeResponse
			if err := rows.Scan(&types.Btcode, &types.BusinessType); err != nil {
				log.Println(err)
				return c.Status(500).JSON(response.ResponseModel{
					RetCode: "500",
					Message: "Internal Server Error",
					Data: errors.ErrorModel{
						Message:   "Internal Server Error",
						IsSuccess: false,
						Error:     err,
					},
				})
			}
			businessTypes = append(businessTypes.([]BusinessTypeResponse), types)
		}

		// Check for any errors encountered during iteration
		if err := rows.Err(); err != nil {
			log.Println(err)
			return c.Status(500).JSON(response.ResponseModel{
				RetCode: "500",
				Message: "Internal Server Error",
				Data: errors.ErrorModel{
					Message:   "Internal Server Error",
					IsSuccess: false,
					Error:     err,
				},
			})
		}
	case 0:
		// Define a slice to hold the results
		businessTypes = []BusinessTypesList{}
		// Execute the raw SQL query with the loan purpose code parameter
		rows, err := db.Raw("SELECT business_type FROM business_loan_purpose.business_type ORDER BY business_type ASC").Rows()
		if err != nil {
			log.Println(err)
			return c.Status(500).JSON(response.ResponseModel{
				RetCode: "500",
				Message: "Internal Server Error",
				Data: errors.ErrorModel{
					Message:   "Internal Server Error",
					IsSuccess: false,
					Error:     err,
				},
			})
		}
		defer rows.Close()

		// Iterate over the rows and scan the values into the struct
		for rows.Next() {
			var types BusinessTypesList
			if err := rows.Scan(&types.BusinessType); err != nil {
				log.Println(err)
				return c.Status(500).JSON(response.ResponseModel{
					RetCode: "500",
					Message: "Internal Server Error",
					Data: errors.ErrorModel{
						Message:   "Internal Server Error",
						IsSuccess: false,
						Error:     err,
					},
				})
			}
			businessTypes = append(businessTypes.([]BusinessTypesList), types)
		}

		// Check for any errors encountered during iteration
		if err := rows.Err(); err != nil {
			log.Println(err)
			return c.Status(500).JSON(response.ResponseModel{
				RetCode: "500",
				Message: "Internal Server Error",
				Data: errors.ErrorModel{
					Message:   "Internal Server Error",
					IsSuccess: false,
					Error:     err,
				},
			})
		}
	}

	// Return the business types as a JSON response
	return c.Status(200).JSON(response.ResponseModel{
		RetCode: "200",
		Message: "Successful",
		Data:    businessTypes,
	})
}
