package offices

import (
	// "chatbot/pkg/models/errors"
	// "chatbot/pkg/models/response"
	// "chatbot/pkg/models/status"
	"chatbot/pkg/realtime"
	"chatbot/pkg/sharedfunctions"
	"chatbot/pkg/utils/go-utils/database"
	"strings"
)

// center
func Get_Center(params *SelectCentersParams) (map[string]any, error) {
	db := database.DB
	var result map[string]any

	// Now driver will marshal params -> JSON automatically
	if err := db.Raw(`SELECT cardincoffices.get_centers(?)`, params).Scan(&result).Error; err != nil {
		return nil, err
	}

	// Convert JSON string fields to proper JSON
	sharedfunctions.ConvertStringToJSONMap(result)
	result = sharedfunctions.GetMap(result, "get_centers")

	return result, nil
}

func Upsert_Center(params *UpsertCentersParams) (map[string]any, error) {
	db := database.DB
	var result map[string]any

	// Pass the struct directly to Postgres JSONB function
	if err := db.Raw(`SELECT cardincoffices.upsert_centers(?)`, params).Scan(&result).Error; err != nil {
		return nil, err
	}

	// Convert JSON string fields to map if needed
	sharedfunctions.ConvertStringToJSONMap(result)
	result = sharedfunctions.GetMap(result, "upsert_centers") // same as Get_Center
	message := sharedfunctions.GetStringFromMap(result, "message")

	//broadcasting
	handleMessage(message, result)
	return result, nil
}

func handleMessage(message string, result map[string]any) {
	hubs := map[string]func(map[string]any){
		"Center":  realtime.UpsertCentersHub.Publish,
		"Cluster": realtime.UpsertClusterHub.Publish,
		"Region":  realtime.UpsertRegionHub.Publish,
		"Unit":    realtime.UpsertUnitsHub.Publish,
	}

	for prefix, publish := range hubs {
		if strings.HasPrefix(message, prefix+" successfully") {
			publish(result)
			return
		}
	}
}

// cluster
func Get_Clusters() (map[string]any, error) {
	db := database.DB
	var result map[string]any

	// Call the updated SQL function that returns JSONB
	if err := db.Raw(`SELECT cardincoffices.getall_cluster()`).Scan(&result).Error; err != nil {
		return nil, err
	}

	// Convert JSON string fields to proper JSON if necessary
	sharedfunctions.ConvertStringToJSONMap(result)

	// Extract the JSONB field returned by the function
	result = sharedfunctions.GetMap(result, "getall_cluster")

	return result, nil
}
func Upsert_Cluster(params *UpsertClusterParams) (map[string]any, error) {
	db := database.DB
	var result map[string]any

	if err := db.Raw(`SELECT cardincoffices.upsert_cluster(?)`, params).Scan(&result).Error; err != nil {
		return nil, err
	}

	sharedfunctions.ConvertStringToJSONMap(result)
	result = sharedfunctions.GetMap(result, "upsert_cluster")
	message := sharedfunctions.GetStringFromMap(result, "message")
	//broadcasting
	handleMessage(message, result)

	return result, nil
}

// Regions
func Get_Region(params *SelectRegionsParams) (map[string]any, error) {
	db := database.DB
	var result map[string]any

	// Now driver will marshal params -> JSON automatically
	if err := db.Raw(`SELECT cardincoffices.get_regions(?)`, params).Scan(&result).Error; err != nil {
		return nil, err
	}

	// Convert JSON string fields to proper JSON
	sharedfunctions.ConvertStringToJSONMap(result)
	result = sharedfunctions.GetMap(result, "get_regions")

	return result, nil
}
func Upsert_Region(params *UpsertRegionParams) (map[string]any, error) {
	db := database.DB
	var result map[string]any

	// Call PostgreSQL function with JSONB
	if err := db.Raw(`SELECT cardincoffices.upsert_region(?)`, params).Scan(&result).Error; err != nil {
		return nil, err
	}

	// Convert and unwrap JSON result
	sharedfunctions.ConvertStringToJSONMap(result)
	result = sharedfunctions.GetMap(result, "upsert_region")
	message := sharedfunctions.GetStringFromMap(result, "message")
	//broadcasting
	handleMessage(message, result)
	return result, nil
}

// Get units
func Get_Units(params map[string]any) (map[string]any, error) {
	db := database.DB

	// Now driver will marshal params -> JSON automatically
	if err := db.Raw("SELECT * FROM gabaykonekfunc.officesgetunits($1)", params).Scan(&params).Error; err != nil {
		return nil, err
	}

	sharedfunctions.ConvertStringToJSONMap(params)
	params = sharedfunctions.GetMap(params, "officesgetunits")

	return params, nil
}

func Upsert_Units(params *UpsertUnitsParams) (map[string]any, error) {
	db := database.DB
	var result map[string]any

	// Call PostgreSQL function with JSONB payload
	if err := db.Raw(`SELECT cardincoffices.upsert_units(?)`, params).Scan(&result).Error; err != nil {
		return nil, err
	}

	// Convert and unwrap JSON result
	sharedfunctions.ConvertStringToJSONMap(result)
	result = sharedfunctions.GetMap(result, "upsert_units")
	message := sharedfunctions.GetStringFromMap(result, "message")
	//broadcasting
	handleMessage(message, result)
	return result, nil
}

// get staff name
func Get_fullname(params map[string]any) (map[string]any, error) {
	db := database.DB
	var result map[string]any

	// Call PostgreSQL function with JSONB payload
	if err := db.Raw(`SELECT userprofile.get_fullname(?)`, params).Scan(&result).Error; err != nil {
		return nil, err
	}

	// Convert and unwrap JSON result
	sharedfunctions.ConvertStringToJSONMap(result)
	result = sharedfunctions.GetMap(result, "get_fullname")

	return result, nil
}

// designation
func GetStaffByDesignationDB(params map[string]any) (map[string]any, error) {
	db := database.DB
	var result map[string]any

	if err := db.Raw(`SELECT cardincoffices.getstaffby_designation_jsonb(?)`, params).Scan(&result).Error; err != nil {
		return nil, err
	}

	// Convert JSON string to map
	sharedfunctions.ConvertStringToJSONMap(result)
	result = sharedfunctions.GetMap(result, "getstaffby_designation_jsonb")

	return result, nil
}

// center staff
func GetCenterByStaffIDDB(params map[string]any) (map[string]any, error) {
	db := database.DB
	var result map[string]any

	if err := db.Raw(`SELECT cardincoffices.get_center_by_staff_id_jsonb(?)`, params).Scan(&result).Error; err != nil {
		return nil, err
	}

	// Convert string JSON to map
	sharedfunctions.ConvertStringToJSONMap(result)
	result = sharedfunctions.GetMap(result, "get_center_by_staff_id_jsonb")

	return result, nil
}

// update center staff
func UpdateCenterStaffDB(params map[string]any) (map[string]any, error) {
	db := database.DB
	var result map[string]any

	if err := db.Raw(`SELECT cardincoffices.update_center_staffid_jsonb(?)`, params).Scan(&result).Error; err != nil {
		return nil, err
	}

	sharedfunctions.ConvertStringToJSONMap(result)
	result = sharedfunctions.GetMap(result, "update_center_staffid_jsonb")

	return result, nil
}
