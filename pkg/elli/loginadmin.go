package elli

// import (
// 	"chatbot/pkg/models/status"
// 	"chatbot/pkg/sharedfunctions"
// 	"fmt"
// 	"chatbot/pkg/utils/go-utils/database"
// 	"github.com/gofiber/fiber/v2"
// )

// func Login_admin_account(loginCreds map[string]any) (map[string]any, bool, int, string, string, string, error){
// 	db := database.DB

// 	var response map[string]any
// 	if err := db.Raw("SELECT * FROM public.accountloginvalidationadmin($1)", loginCreds).Scan(&response).Error; err != nil {
// 		return nil, false, 500, "500", status.RetCode500, "An error occured while validating your credentials.", err
// 	}
// }
