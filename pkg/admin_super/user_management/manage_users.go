package usermanagement

import (
	admincardinc "chatbot/pkg/admin_cardinc"
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"
	"chatbot/pkg/models/status"
	"chatbot/pkg/realtime"
	"chatbot/pkg/sharedfunctions"

	"github.com/gofiber/fiber/v2"
)

type Request struct {
	RequestData map[string]any
}

func UpdateUsers(c *fiber.Ctx) error {
	adminAccess := c.Get("adminAccess")
	staffid := c.Params("id") // optional for logging
	request := new(Request)

	if err := c.BodyParser(&request); err != nil {
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
	// Add adminAccess into the request data
	request.RequestData["adminAccess"] = adminAccess

	result, err := sharedfunctions.UpdateUser(request.RequestData)
	if err != nil {
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: status.RetCode500,
			Data: errors.ErrorModel{
				Message:   "Error updating user due to a problem connecting to database!",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	message_update := sharedfunctions.GetStringFromMap(result, "retCode")

	if message_update == "200" {
		if staffInfo, err := admincardinc.GetCardIncStaffInfo(); err == nil {
			realtime.MainHub.Publish("ToAll", "get_cardincstaff", staffInfo)
		}
		if allUser, err := sharedfunctions.GetAllUsers(); err == nil {

			sharedfunctions.ConvertStringToJSONMap(allUser)
			allUsers := sharedfunctions.GetList(allUser, "getalluser")
			realtime.MainHub.Publish(staffid, "get_allusers", allUsers)
			//realtime.MainHub.Publish("ToAll", "get_allusers", allUsers)

		}
	}

	return c.JSON(result)
}
