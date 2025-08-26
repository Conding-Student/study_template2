package qslsal

import (
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"
	"chatbot/pkg/models/status"
	"chatbot/pkg/sharedfunctions"

	"github.com/gofiber/fiber/v2"
)

type QslSalFields struct {
	HeaderFields     map[string]any   `json:"headerFields"`
	BusinessAssets   []map[string]any `json:"businessAssets"`
	Liabilities      []map[string]any `json:"liabilities"`
	Equity           []map[string]any `json:"equity"`
	BusinessPosition []map[string]any `json:"businessPosition"`
	BusinessExpenses []map[string]any `json:"businessExpenses"`
	NetIncome        []map[string]any `json:"netIncome"`
	HouseholdIncome  []map[string]any `json:"householdIncome"`
	Recommendation   map[string]any   `json:"recommendationFields"`
	UMApproval       map[string]any   `json:"umApprovalFields"`
	AMApproval       map[string]any   `json:"amApprovalFields"`
	RDApproval       map[string]any   `json:"rdApprovalFields"`
}

func GetFields(c *fiber.Ctx) error {
	fieldsLists, err := RetrieveQslSalFields()
	if err != nil {
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: status.RetCode500,
			Data: errors.ErrorModel{
				Message:   "Failed to fetch Business Assets fields",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	fields := QslSalFields{
		HeaderFields:     sharedfunctions.GetMap(fieldsLists, "headerfields"),
		BusinessAssets:   sharedfunctions.GetList(fieldsLists, "businessassets"),
		Liabilities:      sharedfunctions.GetList(fieldsLists, "liabilities"),
		Equity:           sharedfunctions.GetList(fieldsLists, "equity"),
		BusinessPosition: sharedfunctions.GetList(fieldsLists, "businessposition"),
		BusinessExpenses: sharedfunctions.GetList(fieldsLists, "businessexpenses"),
		NetIncome:        sharedfunctions.GetList(fieldsLists, "netincome"),
		HouseholdIncome:  sharedfunctions.GetList(fieldsLists, "householdincome"),
		Recommendation:   sharedfunctions.GetMap(fieldsLists, "recommendationfields"),
		UMApproval:       sharedfunctions.GetMap(fieldsLists, "umapprovalfields"),
		AMApproval:       sharedfunctions.GetMap(fieldsLists, "amapprovalfields"),
		RDApproval:       sharedfunctions.GetMap(fieldsLists, "rdapprovalfields"),
	}

	return c.Status(200).JSON(response.ResponseModel{
		RetCode: "200",
		Message: "Successfully fetch fields",
		Data: fiber.Map{
			"qslSalFields": fields,
		},
	})
}
