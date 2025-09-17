package logs

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

func CardIncAuditTrail(staffID, module, activity string) error {
	params := map[string]any{
		"staffID":  staffID,
		"module":   module,
		"activity": activity,
	}

	result, err := InsertAuditTrail(params)
	if err != nil {
		return err
	}

	fmt.Println("✅ Audit trail inserted successfully:", result)
	return nil
}

func ErrorLogs(staffID, module, activity string) error {
	params := map[string]any{
		"staffID": staffID,
		"module":  module,
		"error":   activity,
	}

	result, err := InsertErrorLog(params)
	if err != nil {
		return err
	}

	fmt.Println("✅ Error log inserted successfully:", result)
	return nil
}

// func CardIncAuditTrail(staffid, module, activity string) error {
// 	db := database.DB

// 	currentTime, err := sharedfunctions.LocalTime()
// 	if err != nil {
// 		return err
// 	}

// 	insertQuery := `SELECT logs.insertaudtrail($1, $2, $3, $4)`

// 	if err := db.Exec(insertQuery, currentTime, staffid, module, activity).Error; err != nil {
// 		log.Println(err)
// 		return err
// 	}

// 	fmt.Println("Logs:", currentTime, staffid, module, activity)

//		return nil
//	}

// func ErrorLogs(staffid, module, activity string) error {
// 	db := database.DB

// 	currentTime, err := sharedfunctions.LocalTime()
// 	if err != nil {
// 		return err
// 	}

// 	insertQuery := `SELECT logs.inserterrorlogs($1, $2, $3, $4)`

// 	if err := db.Exec(insertQuery, currentTime, staffid, module, activity).Error; err != nil {
// 		log.Println(err)
// 		return err
// 	}

// 	fmt.Println("Logs:", currentTime, staffid, module, activity)

//		return nil
//	}
func TriggerErrorLog(c *fiber.Ctx) error {
	staffid := c.Query("staffid")
	module := c.Query("module")
	activity := c.Query("activity")

	err := ErrorLogs(staffid, module, activity)
	if err != nil {
		// log the exact error
		log.Println("Error inserting log:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to insert log: %v", err),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Log inserted successfully",
	})
}
