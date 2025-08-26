package logs

import (
	"chatbot/pkg/sharedfunctions"
	"chatbot/pkg/utils/go-utils/database"
	"fmt"
	"log"
)

func CardIncAuditTrail(staffid, module, activity string) error {
	db := database.DB

	currentTime, err := sharedfunctions.LocalTime()
	if err != nil {
		return err
	}

	insertQuery := `SELECT logs.insertaudtrail($1, $2, $3, $4)`

	if err := db.Exec(insertQuery, currentTime, staffid, module, activity).Error; err != nil {
		log.Println(err)
		return err
	}

	fmt.Println("Logs:", currentTime, staffid, module, activity)

	return nil
}

func ErrorLogs(staffid, module, activity string) error {
	db := database.DB

	currentTime, err := sharedfunctions.LocalTime()
	if err != nil {
		return err
	}

	insertQuery := `SELECT logs.inserterrorlogs($1, $2, $3, $4)`

	if err := db.Exec(insertQuery, currentTime, staffid, module, activity).Error; err != nil {
		log.Println(err)
		return err
	}

	fmt.Println("Logs:", currentTime, staffid, module, activity)

	return nil
}
