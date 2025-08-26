package handler

import (
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"
	"chatbot/pkg/utils/go-utils/database"

	"github.com/gofiber/fiber/v2"
)

// Educational Attainment List
type EducationalAttainment struct {
	EducationalID         string `gorm:"column:educational_id;primaryKey"`
	EducationalAttainment string `gorm:"column:educational_attainment"`
}

func (EducationalAttainment) TableName() string {
	return "educational_attainment"
}

func EducationalAttainments(c *fiber.Ctx) error {
	db := database.DB

	var educations []EducationalAttainment
	if err := db.Order("CAST(educational_id AS INTEGER) ASC").Find(&educations).Error; err != nil {
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Internal server error",
			Data: errors.ErrorModel{
				Message:   "Failed to fetch educational attainment values",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	// Extract values from the result
	var educationalAttainment []string
	for _, education := range educations {
		educationalAttainment = append(educationalAttainment, education.EducationalAttainment)
	}

	return c.Status(200).JSON(response.ResponseModel{
		RetCode: "200",
		Message: "Successful!",
		Data:    educationalAttainment,
	})
}

// Gender
type Gender struct {
	GenderID string `gorm:"column:gender_id;primaryKey"`
	Gender   string `gorm:"column:gender"`
}

func (Gender) TableName() string {
	return "gender"
}

func Genders(c *fiber.Ctx) error {

	db := database.DB

	var genders []Gender
	if err := db.Find(&genders).Error; err != nil {
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Internal server error",
			Data: errors.ErrorModel{
				Message:   "Failed to fetch gender list",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	// Extract values from the result
	var listOfGenders []string
	for _, gendersList := range genders {
		listOfGenders = append(listOfGenders, gendersList.Gender)
	}

	return c.Status(200).JSON(response.ResponseModel{
		RetCode: "200",
		Message: "Successful!",
		Data:    listOfGenders,
	})
}

type MaritalStatus struct {
	MaritalStatusID string `gorm:"column:marital_status_id;primaryKey"`
	MaritalStatus   string `gorm:"column:marital_status"`
}

func (MaritalStatus) TableName() string {
	return "marital_status"
}

func MaritalStats(c *fiber.Ctx) error {
	db := database.DB

	var maritalStatus []MaritalStatus
	if err := db.Order("CAST(marital_status_id AS INTEGER) ASC").Find(&maritalStatus).Error; err != nil {
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Internal server error",
			Data: errors.ErrorModel{
				Message:   "Failed to fetch marital status values",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	// Extract values from the result
	var listOfMaritalStatus []string
	for _, institutionsEstablishedList := range maritalStatus {
		listOfMaritalStatus = append(listOfMaritalStatus, institutionsEstablishedList.MaritalStatus)
	}

	return c.Status(200).JSON(response.ResponseModel{
		RetCode: "200",
		Message: "Successful!",
		Data:    listOfMaritalStatus,
	})
}

type Occupation struct {
	OccupationID string `gorm:"column:occupation_id;primaryKey"`
	Occupation   string `gorm:"column:occupation"`
}

func (Occupation) TableName() string {
	return "occupation"
}

func Occupations(c *fiber.Ctx) error {
	db := database.DB

	var Occupation []Occupation
	if err := db.Order("CAST(occupation_id AS INTEGER) ASC").Find(&Occupation).Error; err != nil {
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Internal server error",
			Data: errors.ErrorModel{
				Message:   "Failed to fetch occupations values",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	// Extract values from the result
	var listOfOccupations []string
	for _, OccupationsList := range Occupation {
		listOfOccupations = append(listOfOccupations, OccupationsList.Occupation)
	}

	return c.Status(200).JSON(response.ResponseModel{
		RetCode: "200",
		Message: "Successful!",
		Data:    listOfOccupations,
	})
}

type Religion struct {
	OccupationID string `gorm:"column:religion_id;primaryKey"`
	Occupation   string `gorm:"column:religion"`
}

func (Religion) TableName() string {
	return "religion"
}

func Religions(c *fiber.Ctx) error {
	db := database.DB

	var Religion []Religion
	if err := db.Order("CAST(religion_id AS INTEGER) ASC").Find(&Religion).Error; err != nil {
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Internal server error",
			Data: errors.ErrorModel{
				Message:   "Failed to fetch religion values",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	// Extract values from the result
	var listOfReligion []string
	for _, ReligionList := range Religion {
		listOfReligion = append(listOfReligion, ReligionList.Occupation)
	}

	return c.Status(200).JSON(response.ResponseModel{
		RetCode: "200",
		Message: "Successful!",
		Data:    listOfReligion,
	})
}

type LegalID struct {
	ID      string `gorm:"column:id;primaryKey"`
	LegalID string `gorm:"column:legal_identification"`
}

func (LegalID) TableName() string {
	return "legal_id"
}

func LegalIDs(c *fiber.Ctx) error {
	db := database.DB

	var LegalID []LegalID
	if err := db.Order("CAST(id AS INTEGER) ASC").Find(&LegalID).Error; err != nil {
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Internal server error",
			Data: errors.ErrorModel{
				Message:   "Failed to fetch legal id's values",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	// Extract values from the result
	var listOfLegalID []string
	for _, LegalIdList := range LegalID {
		listOfLegalID = append(listOfLegalID, LegalIdList.LegalID)
	}

	return c.Status(200).JSON(response.ResponseModel{
		RetCode: "200",
		Message: "Successful!",
		Data:    listOfLegalID,
	})
}

type MonIncome struct {
	ID            string `gorm:"column:id;primaryKey"`
	MonthlyIncome string `gorm:"column:monthly_income"`
}

func (MonIncome) TableName() string {
	return "monthly_income"
}

func MonthlyIncome(c *fiber.Ctx) error {
	db := database.DB

	var MonIncome []MonIncome
	if err := db.Order("CAST(id AS INTEGER) ASC").Find(&MonIncome).Error; err != nil {
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Internal server error",
			Data: errors.ErrorModel{
				Message:   "Failed to fetch monthly income values",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	// Extract values from the result
	var listOfMonthlyIncome []string
	for _, MonthlyIncomeList := range MonIncome {
		listOfMonthlyIncome = append(listOfMonthlyIncome, MonthlyIncomeList.MonthlyIncome)
	}

	return c.Status(200).JSON(response.ResponseModel{
		RetCode: "200",
		Message: "Successful!",
		Data:    listOfMonthlyIncome,
	})
}

type NatOfBus struct {
	ID       string `gorm:"column:id;primaryKey"`
	NatOfBus string `gorm:"column:nature_of_business"`
}

func (NatOfBus) TableName() string {
	return "nature_of_business"
}

func NatureOfBusiness(c *fiber.Ctx) error {
	db := database.DB

	var NatOfBus []NatOfBus
	if err := db.Order("CAST(id AS INTEGER) ASC").Find(&NatOfBus).Error; err != nil {
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Internal server error",
			Data: errors.ErrorModel{
				Message:   "Failed to fetch nature of business values",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	// Extract values from the result
	var listOfNatureOfBusiness []string
	for _, NatureOfBusinessList := range NatOfBus {
		listOfNatureOfBusiness = append(listOfNatureOfBusiness, NatureOfBusinessList.NatOfBus)
	}

	return c.Status(200).JSON(response.ResponseModel{
		RetCode: "200",
		Message: "Successful!",
		Data:    listOfNatureOfBusiness,
	})
}

type CagabayVersion struct {
	Vesion string `json:"version"`
}

type CopyrightAndPoweredBy struct {
	Copyright     string `gorm:"column:copyright"`
	PoweredBy     string `gorm:"column:powered_by"`
	PoweredByLink string `gorm:"column:powered_by_link"`
}
