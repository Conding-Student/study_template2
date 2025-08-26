package gabaykonekdashboard

import (
	"chatbot/pkg/sharedfunctions"
	"chatbot/pkg/utils/go-utils/database"
	"fmt"
	"log"

	"github.com/lib/pq"
)

func GetReleasedLoanAndTotal(staffid string, startDate, endDate string) ([]map[string]any, error) {
	db := database.DB
	eSystem := database.EsystemDB

	var designation int64
	if err := db.Raw("SELECT public.convertdesigtoint($1)", staffid).Scan(&designation).Error; err != nil {
		return nil, err
	}

	var summary []map[string]any
	if designation == 0 {

		query := `
		SELECT * FROM loan_application.getparamsforaoloanreleasessum($1)
		`

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
		area := aoparams[0]["area"].(string)
		var centers []string

		for _, param := range aoparams {
			if unitCenter, ok := param["unitcenter"].(string); ok {
				centers = append(centers, unitCenter)
			}
		}

		// fmt.Println(aoparams)
		// fmt.Println(staffid)
		// fmt.Println(area)
		// fmt.Println(centers)

		eSystemQuery := `SELECT * FROM staging.cagabay_getloansandtotal($1, $2, $3, $4, $5)`
		if err := eSystem.Raw(eSystemQuery, designation, area, pq.Array(centers), startDate, endDate).Scan(&summary).Error; err != nil {
			log.Println("esystem error", err)
			return nil, err
		}
		return summary, nil
	} else if designation == 1 {
		query := `
		SELECT * FROM loan_application.getparamsforumloanreleasessum($1)
		`

		var umparams []map[string]any
		if err := db.Raw(query, staffid).Scan(&umparams).Error; err != nil {
			log.Println("ca-gabay error", err)
			return nil, err
		}

		if len(umparams) == 0 {
			return summary, nil
		}

		area := umparams[0]["area"].(string)
		var unit []int64

		for _, param := range umparams {
			if unitCenter, ok := param["unitcenter"].(int64); ok {
				unit = append(unit, unitCenter)
			}
		}

		// fmt.Println(umparams)
		// fmt.Println(staffid)
		// fmt.Println(area)
		fmt.Println(unit)

		eSystemQuery := `SELECT * FROM staging.cagabay_getloansandtotal($1, $2, $3, $4, $5)`
		if err := eSystem.Raw(eSystemQuery, designation, area, pq.Array(unit), startDate, endDate).Scan(&summary).Error; err != nil {
			log.Println("esystem error", err)
			return nil, err
		}
		return summary, nil
	} else if designation == 2 {
		query := `
		SELECT * FROM loan_application.getparamsforamloanreleasessum($1)
		`

		var amparams []map[string]any
		if err := db.Raw(query, staffid).Scan(&amparams).Error; err != nil {
			log.Println("ca-gabay error", err)
			return nil, err
		}

		if len(amparams) == 0 {
			return summary, nil
		}

		var areaList []string
		for _, param := range amparams {
			if areas, ok := param["area"].(string); ok {
				areaList = append(areaList, areas)
			}
		}

		// fmt.Println(amparams)
		// fmt.Println(staffid)
		// fmt.Println(areaList)

		eSystemQuery := `SELECT * FROM staging.cagabay_getloansandtotal($1, $2, $3, $4, $5)`
		if err := eSystem.Raw(eSystemQuery, designation, nil, pq.Array(areaList), startDate, endDate).Scan(&summary).Error; err != nil {
			log.Println("esystem error", err)
			return nil, err
		}
		return summary, nil
	} else if designation == 3 {
		query := `
		SELECT * FROM loan_application.getparamsforrdloanreleasessum($1)
		`

		var rdparams []map[string]any
		if err := db.Raw(query, staffid).Scan(&rdparams).Error; err != nil {
			log.Println("ca-gabay error", err)
			return nil, err
		}

		if len(rdparams) == 0 {
			return summary, nil
		}

		var areaList []string
		for _, param := range rdparams {
			if areas, ok := param["area"].(string); ok {
				areaList = append(areaList, areas)
			}
		}

		// fmt.Println(rdparams)
		// fmt.Println(staffid)
		// fmt.Println(areaList)

		eSystemQuery := `SELECT * FROM staging.cagabay_getloansandtotal($1, $2, $3 $4, $5)`
		if err := eSystem.Raw(eSystemQuery, designation, nil, pq.Array(areaList), startDate, endDate).Scan(&summary).Error; err != nil {
			log.Println("esystem error", err)
			return nil, err
		}
		return summary, nil
	} else {
		return nil, fmt.Errorf("error generating summary of releases due to wrong gabaykonek roles")
	}
}

func GetStaffEfficiency(staffID string) (map[string]any, error) {
	db := database.DB
	eSystem := database.EsystemDB

	var response map[string]any
	if err := db.Raw("SELECT * FROM cardincoffices.getoffices($1)", staffID).Scan(&response).Error; err != nil {
		return nil, err
	}

	sharedfunctions.ConvertStringToJSONMap(response)

	offices := sharedfunctions.GetMap(response, "offices")
	fmt.Println(offices)

	var staffEff map[string]any
	if err := eSystem.Raw("SELECT * FROM public.cagabay_staffeff($1)", offices).Scan(&staffEff).Error; err != nil {
		return nil, err
	}

	sharedfunctions.ConvertStringToJSONMap(staffEff)

	return staffEff, nil
}
