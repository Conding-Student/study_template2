package offices

import (
	"chatbot/pkg/sharedfunctions"
	"chatbot/pkg/utils/go-utils/database"
	"chatbot/pkg/websocket"
	"fmt"
)

// branches
func Get_Branch(params map[string]any) (map[string]any, error) {
	db := database.DB

	var result map[string]any
	// Now driver will marshal params -> JSON automatically
	if err := db.Raw(`SELECT cardincoffices.get_branches(?)`, params).Scan(&result).Error; err != nil {
		return nil, err
	}

	// Convert JSON string fields to proper JSON
	sharedfunctions.ConvertStringToJSONMap(result)
	result = sharedfunctions.GetMap(result, "get_branches")

	return result, nil
}

// center
func Get_Center(params map[string]any) (map[string]any, error) {
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

// Regions
func Get_Region(params map[string]any) (map[string]any, error) {
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

// Get units
func Get_Units(params map[string]any) (map[string]any, error) {
	db := database.DB

	var result map[string]any
	// Now driver will marshal params -> JSON automatically
	if err := db.Raw("SELECT * FROM gabaykonekfunc.officesgetunits($1)", params).Scan(&result).Error; err != nil {
		return nil, err
	}

	sharedfunctions.ConvertStringToJSONMap(result)
	result = sharedfunctions.GetMap(result, "officesgetunits")

	return result, nil
}

func Upsert_Branch(params map[string]any, params_select map[string]any) (map[string]any, error) {
	db := database.DB
	var result map[string]any

	// Pass the struct directly to Postgres JSONB function
	if err := db.Raw("SELECT * FROM cardincoffices.upsertbranches($1)", params).Scan(&result).Error; err != nil {
		return nil, err
	}

	// Convert JSON string fields to map if needed
	sharedfunctions.ConvertStringToJSONMap(result)
	result = sharedfunctions.GetMap(result, "upsertbranches") // same as Get_Center
	message := sharedfunctions.GetStringFromMap(result, "retCode")

	params_select["operation"] = 1 // to fetch all branches
	params_select["region"] = params["Region"]
	staffid := fmt.Sprintf("%v", params["staffid"])

	if clusters, err := Get_Branch(params_select); err == nil {
		handleMessage("Branch", staffid, message, clusters)
	}
	return result, nil
}

func Upsert_Center(params map[string]any, params_select map[string]any) (map[string]any, error) {
	db := database.DB
	var result map[string]any

	// Pass the struct directly to Postgres JSONB function
	if err := db.Raw(`SELECT cardincoffices.upsert_centers(?)`, params).Scan(&result).Error; err != nil {
		return nil, err
	}

	// Convert JSON string fields to map if needed
	sharedfunctions.ConvertStringToJSONMap(result)
	result = sharedfunctions.GetMap(result, "upsert_centers") // same as Get_Center
	message := sharedfunctions.GetStringFromMap(result, "retCode")

	params_select["operation"] = 1 // to fetch all centers
	params_select["brcode"] = params["brcode"]
	params_select["unitcode"] = params["unitcode"]
	staffid := fmt.Sprintf("%v", params["staffid"])

	if clusters, err := Get_Center(params_select); err == nil {
		handleMessage("Center", staffid, message, clusters)
	}
	return result, nil
}

func Upsert_Cluster(params map[string]any) (map[string]any, error) {
	db := database.DB
	var result map[string]any

	if err := db.Raw(`SELECT cardincoffices.upsert_cluster(?)`, params).Scan(&result).Error; err != nil {
		return nil, err
	}

	sharedfunctions.ConvertStringToJSONMap(result)
	result = sharedfunctions.GetMap(result, "upsert_cluster")
	message := sharedfunctions.GetStringFromMap(result, "retCode")
	staffid := fmt.Sprintf("%v", params["staffid"])

	// ✅ Safely re-fetch clusters for broadcasting
	if clusters, err := Get_Clusters(); err == nil {
		handleMessage("Cluster", staffid, message, clusters) // null staffid for cluster
	}

	return result, nil
}

func Upsert_Region(params map[string]any, params_select map[string]any) (map[string]any, error) {
	db := database.DB
	var result map[string]any

	// Call PostgreSQL function with JSONB
	if err := db.Raw(`SELECT cardincoffices.upsert_region(?)`, params).Scan(&result).Error; err != nil {
		return nil, err
	}

	// Convert and unwrap JSON result
	sharedfunctions.ConvertStringToJSONMap(result)
	result = sharedfunctions.GetMap(result, "upsert_region")
	message := sharedfunctions.GetStringFromMap(result, "retCode")

	params_select["operation"] = 1 // to fetch all regions
	params_select["cluster"] = params["cluster"]
	staffid := fmt.Sprintf("%v", params["staffid"])

	if clusters, err := Get_Region(params_select); err == nil {
		handleMessage("Region", staffid, message, clusters)
	}
	return result, nil
}

func Upsert_Units(params map[string]any, params_select map[string]any) (map[string]any, error) {
	db := database.DB
	var result map[string]any

	// Call PostgreSQL function with JSONB payload
	if err := db.Raw(`SELECT cardincoffices.upsert_units(?)`, params).Scan(&result).Error; err != nil {
		return nil, err
	}

	// Convert and unwrap JSON result
	sharedfunctions.ConvertStringToJSONMap(result)
	result = sharedfunctions.GetMap(result, "upsert_units")
	message := sharedfunctions.GetStringFromMap(result, "retCode")

	params_select["operation"] = 1 // to fetch all units
	params_select["brcode"] = params["brcode"]
	staffid := fmt.Sprintf("%v", params["staffid"])

	if clusters, err := Get_Units(params_select); err == nil {
		handleMessage("Unit", staffid, message, clusters)
	}
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

func handleMessage(functionName string, staffid string, message string, result any) {
	if message == "200" {
		fmt.Println("Broadcasting", functionName, "update to staffID:", staffid)
		hubs := map[string]func(any){
			"Center": func(data any) {
				websocket.MPublish(staffid, data, "CenOM") //
			},
			"Cluster": func(data any) {
				websocket.MPublish(staffid, data, "CluOM")
			},
			"Region": func(data any) {
				websocket.MPublish(staffid, data, "RegOM")
			},
			"Unit": func(data any) {
				websocket.MPublish(staffid, data, "UniOM")
			},
			"Branch": func(data any) {
				websocket.MPublish(staffid, data, "BrhOM")
			},
		}
		// Only call the function that matches functionName
		if publish, ok := hubs[functionName]; ok {
			publish(result)
		}
	}
}
