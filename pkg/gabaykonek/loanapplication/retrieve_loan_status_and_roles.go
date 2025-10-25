package loans

import (
	"chatbot/pkg/logs"
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"
	"chatbot/pkg/sharedfunctions"
	"chatbot/pkg/utils/go-utils/database"

	"github.com/gofiber/fiber/v2"
)

func LoanStatusAndRoles(c *fiber.Ctx) error {
	db := database.DB
	id := c.Params("id")
	LOSFeature := "LOS - Retrieve Loan Status and Roles" + id

	var result map[string]any
	if err := db.Raw("SELECT * FROM gabaykonekfunc.getloan_statusandroles()").Scan(&result).Error; err != nil {
		logs.LOSLogs(c, LOSFeature, id, "500", err.Error())
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "200",
			Message: "Successful!",
			Data: errors.ErrorModel{
				IsSuccess: false,
				Message:   "An error occured while connecting to database",
				Error:     err,
			},
		})
	}

	sharedfunctions.ConvertStringToJSONMap(result)
	statusAndRoles := sharedfunctions.GetMap(result, "getloan_statusandroles")

	logs.LOSLogs(c, LOSFeature, id, "200", "Fetch loan status and  user roles successfully!")
	return c.JSON(statusAndRoles)
}
