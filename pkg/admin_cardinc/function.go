package admincardinc

import (
	"chatbot/pkg/sharedfunctions"
	"chatbot/pkg/utils/go-utils/database"
)

func GetCardIncStaffInfo() (map[string]any, error) {
	db := database.DB
	var result map[string]any

	err := db.Raw("SELECT gabaykonekfunc.getcardinc_staff()").Scan(&result).Error
	if err != nil {
		return nil, err
	}

	sharedfunctions.ConvertStringToJSONMap(result)
	result = sharedfunctions.GetMap(result, "getcardinc_staff")
	return result, nil
}

// App Authority
func Get_ApprovingAuthority() (map[string]any, error) {
	db := database.DB
	var result map[string]any

	if err := db.Raw("SELECT public.get_approving_authority()").Scan(&result).Error; err != nil {
		return nil, err
	}

	// Convert JSON string to map if needed
	sharedfunctions.ConvertStringToJSONMap(result)
	result = sharedfunctions.GetMap(result, "get_approving_authority")

	return result, nil
}
func Delete_ApprovingAuthority(params ApprovingAuthority) (map[string]any, error) {
	db := database.DB
	var result map[string]any

	if err := db.Raw("SELECT public.delete_approving_authority(?)", params).Scan(&result).Error; err != nil {
		return nil, err
	}

	sharedfunctions.ConvertStringToJSONMap(result)
	result = sharedfunctions.GetMap(result, "delete_approving_authority")

	return result, nil
}
func Add_ApprovingAuthority(params ApprovingAuthority) (map[string]any, error) {
	db := database.DB
	var result map[string]any

	// Call Postgres function using JSONB
	if err := db.Raw(
		"SELECT public.add_approving_authority(?::jsonb) AS add_approving_authority",
		params,
	).Scan(&result).Error; err != nil {
		return nil, err
	}

	// Convert nested JSON string to map if needed
	sharedfunctions.ConvertStringToJSONMap(result)
	result = sharedfunctions.GetMap(result, "add_approving_authority")

	return result, nil
}

// EPN
func Getloandisbursed(params map[string]any) (map[string]any, error) {
	db := database.DB

	var result map[string]any
	// Call Postgres function using JSONB
	if err := db.Raw("SELECT * FROM gabaykonekfunc.getloandisbursed($1)", params).Scan(&result).Error; err != nil {
		return nil, err
	}

	// Convert nested JSON string to map if needed
	sharedfunctions.ConvertStringToJSONMap(result)
	result = sharedfunctions.GetMap(result, "getloandisbursed")
	return result, nil
}

// user role
func Upsertgkroles(params map[string]any) (map[string]any, error) {
	db := database.DB

	var result map[string]any
	// Call Postgres function using JSONB
	if err := db.Raw("SELECT * FROM gabaykonekfunc.upsert_gkrole($1)", params).Scan(&result).Error; err != nil {
		return nil, err
	}

	// Convert nested JSON string to map if needed
	sharedfunctions.ConvertStringToJSONMap(result)
	result = sharedfunctions.GetMap(result, "result")
	return result, nil
}

func View_roles(params map[string]any) (map[string]any, error) {
	db := database.DB

	var result map[string]any
	// Call Postgres function using JSONB
	if err := db.Raw("SELECT * FROM gabaykonekfunc.roles($1)", params).Scan(&result).Error; err != nil {
		return nil, err
	}

	// Convert nested JSON string to map if needed
	sharedfunctions.ConvertStringToJSONMap(result)
	result = sharedfunctions.GetMap(result, "result")
	return result, nil
}

func Delete_gkrole(params map[string]any) (map[string]any, error) {
	db := database.DB

	var result map[string]any
	// Call Postgres function using JSONB
	if err := db.Raw("SELECT * FROM gabaykonekfunc.delete_gkrole($1)", params).Scan(&result).Error; err != nil {
		return nil, err
	}

	// Convert nested JSON string to map if needed

	sharedfunctions.ConvertStringToJSONMap(result)
	result = sharedfunctions.GetMap(result, "result")
	return result, nil
}
