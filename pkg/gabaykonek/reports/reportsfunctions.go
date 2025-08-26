package reports

import (
	"chatbot/pkg/utils/go-utils/database"
	"fmt"
	"log"
	"time"

	"github.com/lib/pq"
)

type ReportsRequestBody struct {
	Staffid     string
	Designation string
	StartDate   string
	EndDate     string
}

func GetCancelledSummary(staffid string, designation int, startDate, endDate string) ([]map[string]any, error) {
	db := database.DB

	query := `
		SELECT * FROM loan_application.summaryofcancelledloans($1, $2, $3, $4)
	`
	var summary []map[string]any
	if err := db.Raw(query, staffid, designation, startDate, endDate).Scan(&summary).Error; err != nil {
		return nil, err
	}

	return summary, nil
}

func GetAppliedSummary(staffid string, designation int, startDate, endDate string) ([]map[string]any, error) {
	db := database.DB

	query := `
		SELECT * FROM loan_application.summaryofappliedloans($1, $2, $3, $4)
	`
	var summary []map[string]any
	if err := db.Raw(query, staffid, designation, startDate, endDate).Scan(&summary).Error; err != nil {
		return nil, err
	}

	return summary, nil
}

func GetRecommendedSummary(staffid string, designation int, startDate, endDate string) ([]map[string]any, error) {
	db := database.DB

	query := `
		SELECT * FROM loan_application.summaryofrecommendedloans($1, $2, $3, $4)
	`
	var summary []map[string]any
	if err := db.Raw(query, staffid, designation, startDate, endDate).Scan(&summary).Error; err != nil {
		return nil, err
	}

	return summary, nil
}

func GetApprovedSummary(staffid string, designation int, startDate, endDate string) ([]map[string]any, error) {
	db := database.DB

	query := `
		SELECT * FROM loan_application.summaryofapprovedloans($1, $2, $3, $4)
	`
	var summary []map[string]any
	if err := db.Raw(query, staffid, designation, startDate, endDate).Scan(&summary).Error; err != nil {
		return nil, err
	}

	return summary, nil
}

func GetPendingSummary(staffid string, designation int, startDate, endDate string) ([]map[string]any, error) {
	db := database.DB

	query := `
		SELECT * FROM loan_application.summaryofpendingloans($1, $2, $3, $4)
	`
	var summary []map[string]any
	if err := db.Raw(query, staffid, designation, startDate, endDate).Scan(&summary).Error; err != nil {
		return nil, err
	}

	return summary, nil
}

// func GetReleasedSummary(staffid string, designation int, startDate, endDate string) ([]map[string]any, error) {
// 	db := database.DB

// 	query := `
// 		SELECT * FROM loan_application.summaryofdisbursedloans($1, $2, $3, $4)
// 	`
// 	var summary []map[string]any
// 	if err := db.Raw(query, staffid, designation, startDate, endDate).Scan(&summary).Error; err != nil {
// 		return nil, err
// 	}

// 	return summary, nil
// }

func GetReleasedSummary(staffid string, designation int, startDate, endDate string) ([]map[string]any, error) {
	db := database.DB
	eSystem := database.EsystemDB

	var summary []map[string]any
	if designation == 0 {

		query := `
		SELECT * FROM loan_application.getparamsforaoloanreleasessum($1)
		`
		// Fetch the area and unitcenter
		var aoparams []map[string]any
		if err := db.Raw(query, staffid).Scan(&aoparams).Error; err != nil {
			log.Println("ca-gabay error", err)
			return nil, err
		}

		if len(aoparams) == 0 {
			log.Println("No params found for", staffid, designation)
			return summary, nil
		}

		if len(aoparams) == 0 {
			return summary, nil
		}
		// Extract the area code (assuming it's consistent as "U5") and unitcenter values
		area := aoparams[0]["area"].(string) // assuming area is consistent
		var centers []string

		// Loop over the aoparams and collect the unitcenter values into the slice
		for _, param := range aoparams {
			if unitCenter, ok := param["unitcenter"].(string); ok {
				centers = append(centers, unitCenter)
			}
		}

		// fmt.Println(aoparams)
		// fmt.Println(staffid)
		// fmt.Println(area)
		// fmt.Println(centers)

		// Now pass the area and centers to the second query
		eSystemQuery := `SELECT * FROM public.cagabay_getaoloanreleasesummary($1, $2, $3, $4)`
		if err := eSystem.Raw(eSystemQuery, area, pq.Array(centers), startDate, endDate).Scan(&summary).Error; err != nil {
			log.Println("esystem error", err)
			return nil, err
		}
		return summary, nil
	} else if designation == 1 {
		query := `
		SELECT * FROM loan_application.getparamsforumloanreleasessum($1)
		`
		// Fetch the area and unitcenter
		var umparams []map[string]any
		if err := db.Raw(query, staffid).Scan(&umparams).Error; err != nil {
			log.Println("ca-gabay error", err)
			return nil, err
		}

		if len(umparams) == 0 {
			return summary, nil
		}

		// Extract the area code (assuming it's consistent as "U5") and unitcenter values
		area := umparams[0]["area"].(string) // assuming area is consistent
		var unit []int64

		// Loop over the aoparams and collect the unitcenter values into the slice
		for _, param := range umparams {
			if unitCenter, ok := param["unitcenter"].(int64); ok {
				unit = append(unit, unitCenter)
			}
		}

		// fmt.Println(umparams)
		// fmt.Println(staffid)
		// fmt.Println(area)
		// fmt.Println(unit)

		// Now pass the area and centers to the second query
		eSystemQuery := `SELECT * FROM public.cagabay_getunitloanreleasesummary($1, $2, $3, $4)`
		if err := eSystem.Raw(eSystemQuery, area, pq.Array(unit), startDate, endDate).Scan(&summary).Error; err != nil {
			log.Println("esystem error", err)
			return nil, err
		}
		return summary, nil
	} else if designation == 2 {
		query := `
		SELECT * FROM loan_application.getparamsforamloanreleasessum($1)
		`
		// Fetch the area and unitcenter
		var amparams []map[string]any
		if err := db.Raw(query, staffid).Scan(&amparams).Error; err != nil {
			log.Println("ca-gabay error", err)
			return nil, err
		}

		if len(amparams) == 0 {
			return summary, nil
		}

		var areaList []string
		// Loop over the aoparams and collect the unitcenter values into the slice
		for _, param := range amparams {
			if areas, ok := param["area"].(string); ok {
				areaList = append(areaList, areas)
			}
		}

		// fmt.Println(amparams)
		// fmt.Println(staffid)
		// fmt.Println(areaList)

		// Now pass the area and centers to the second query
		eSystemQuery := `SELECT * FROM public.cagabay_getarealoanreleasesummary($1, $2, $3)`
		if err := eSystem.Raw(eSystemQuery, pq.Array(areaList), startDate, endDate).Scan(&summary).Error; err != nil {
			log.Println("esystem error", err)
			return nil, err
		}
		return summary, nil
	} else if designation == 3 {
		query := `
		SELECT * FROM loan_application.getparamsforrdloanreleasessum($1)
		`
		// Fetch the area and unitcenter
		var rdparams []map[string]any
		if err := db.Raw(query, staffid).Scan(&rdparams).Error; err != nil {
			log.Println("ca-gabay error", err)
			return nil, err
		}

		if len(rdparams) == 0 {
			return summary, nil
		}

		var areaList []string
		// Loop over the aoparams and collect the unitcenter values into the slice
		for _, param := range rdparams {
			if areas, ok := param["area"].(string); ok {
				areaList = append(areaList, areas)
			}
		}

		// fmt.Println(rdparams)
		// fmt.Println(staffid)
		// fmt.Println(areaList)

		// Now pass the area and centers to the second query
		eSystemQuery := `SELECT * FROM public.cagabay_getarealoanreleasesummary($1, $2, $3)`
		if err := eSystem.Raw(eSystemQuery, pq.Array(areaList), startDate, endDate).Scan(&summary).Error; err != nil {
			log.Println("esystem error", err)
			return nil, err
		}
		return summary, nil
	}

	return nil, fmt.Errorf("error generating summary of releases")
}

func GetBranchName(staffid string, designation int) (string, error) {
	db := database.DB

	query := `
		SELECT * FROM loan_application.getbranchname($1, $2)
	`

	var branchName *string
	if err := db.Raw(query, staffid, designation).Scan(&branchName).Error; err != nil {
		return "", err
	}

	if branchName == nil {
		return "", nil // or return an error if preferred
	}

	return *branchName, nil
}

func GetGeneratedBy(staffid string) (string, error) {
	db := database.DB

	query := `
		SELECT * FROM loan_application.getgeneratedby($1)
	`

	var generatedBy string
	if err := db.Raw(query, staffid).Scan(&generatedBy).Error; err != nil {
		return "", err
	}

	return generatedBy, nil
}

func GetDateRange(startDate, endDate string) (string, error) {
	// Parse the input dates
	dateStart, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return "", err
	}
	dateEnd, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return "", err
	}

	if dateStart.After(dateEnd) {
		return "", fmt.Errorf("start date cannot be after end date")
	}

	var dateRange string
	// Check if the dates are the same
	if dateStart.Equal(dateEnd) {
		dateRange = fmt.Sprintf("as of %s", dateStart.Format("January 2, 2006"))
	} else {
		dateRange = fmt.Sprintf("from %s to %s", dateStart.Format("January 2, 2006"), dateEnd.Format("January 2, 2006"))
	}
	// Return the formatted date range
	return dateRange, nil
}
