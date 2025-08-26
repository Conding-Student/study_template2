package model

import (
	"time"

	"gorm.io/gorm"
)

type UserSignUp struct {
	gorm.Model
	ID                string `gorm:"primaryKey;autoIncrement"`
	Institution       string `gorm:"not null"`
	Department        string `gorm:"not null"`
	Area              string `gorm:"not null"`
	Unit              string `gorm:"not null"`
	Designation       string `gorm:"not null"`
	Age               string `gorm:"not null"`
	Street            string `gorm:"not null"`
	Barangay          string `gorm:"not null"`
	Municipality      string `gorm:"not null"`
	Province          string `gorm:"not null"`
	Email             string `gorm:"not null"`
	Gender            string `gorm:"not null"`
	Picture           string `gorm:"null"`
	Username          string `gorm:"not null"`
	Mobile            string `gorm:"not null"`
	Password          string `gorm:"not null"`
	Birthday          string `gorm:"not null"`
	StaffId           string `gorm:"not null"`
	DeviceModel       string `gorm:"not null"`
	FirstName         string `gorm:"not null"`
	MiddleName        string `gorm:"not null"`
	LastName          string `gorm:"not null"`
	DeviceID          string `gorm:"not null"`
	FailedAttempts    int
	LastFailedAttempt time.Time `gorm:"column:last_failed_attempt"`
	Role              int
}

type UserData struct {
	ID                string `gorm:"primaryKey;autoIncrement"`
	Institution       string `gorm:"not null"`
	Department        string `gorm:"not null"`
	Area              string `gorm:"not null"`
	Unit              string `gorm:"not null"`
	Designation       string `gorm:"not null"`
	Age               string `gorm:"not null"`
	Street            string `gorm:"not null"`
	Barangay          string `gorm:"not null"`
	Municipality      string `gorm:"not null"`
	Province          string `gorm:"not null"`
	Email             string `gorm:"not null"`
	Gender            string `gorm:"not null"`
	Picture           string `gorm:"null"`
	Username          string `gorm:"not null"`
	Mobile            string `gorm:"not null"`
	Password          string `gorm:"not null"`
	Birthday          string `gorm:"not null"`
	StaffId           string `gorm:"not null"`
	DeviceModel       string `gorm:"not null"`
	FirstName         string `gorm:"not null"`
	MiddleName        string `gorm:"not null"`
	LastName          string `gorm:"not null"`
	DeviceID          string `gorm:"not null"`
	FailedAttempts    int
	LastFailedAttempt time.Time `gorm:"column:last_failed_attempt"`
	AccessToken       string    `json:"accessToken"`
}

type LockTimeInvalidCredential struct {
	ID              int `gorm:"primaryKey" json:"id"`
	MaxAttempts     int `json:"max_attempts"`
	LockoutDuration int `json:"lockout_duration"`
}

type ClientInformationForm struct {
	gorm.Model
	Cid                      string `gorm:"not null"`
	ClientType               string `gorm:"not null"`
	MembershipClassification string `gorm:"not null"`
	MembershipType           string `gorm:"not null"`
	Institution              string `gorm:"not null"`
	Area                     string `gorm:"not null"`
	Unit                     string `gorm:"not null"`
	Center                   string `gorm:"not null"`
	Title                    string `gorm:"not null"`
	FirstName                string `gorm:"not null"`
	MiddleName               string `gorm:"not null"`
	LastName                 string `gorm:"not null"`
	SuffixName               string `gorm:"not null"`
	MaidenName               string `gorm:"not null"`
	Gender                   string `gorm:"not null"`
	Citizenship              string `gorm:"not null"`
	MaritalStatus            string `gorm:"not null"`
	Birthday                 string `gorm:"not null"`
	Age                      string `gorm:"not null"`
	Occupation               string `gorm:"not null"`
	EducationalAttainment    string `gorm:"not null"`
	Disable                  string `gorm:"not null"`
	OtherLegalId             string `gorm:"not null"`
	Tin                      string `gorm:"not null"`
	Sss                      string `gorm:"not null"`
	NatureOfBusiness         string `gorm:"not null"`
	SourceOfFund             string `gorm:"not null"`
	NameOfCompany            string `gorm:"not null"`
	HouseholdMonthlyIncome   string `gorm:"not null"`
	NumberOfHouseholdMembers string `gorm:"not null"`
	NumberOfCollegeChildren  string `gorm:"not null"`
	Email                    string `gorm:"not null"`
	Landline                 string `gorm:"not null"`
	Mobile                   string `gorm:"not null"`
	Healthy                  string `gorm:"not null"`
	PlaceOfBirthBarangay     string `gorm:"not null"`
	PlaceOfBirthMunicipality string `gorm:"not null"`
	PlaceOfBirthProvince     string `gorm:"not null"`
	PresentAddStreet         string `gorm:"not null"`
	PresentAddBarangay       string `gorm:"not null"`
	PresentAddMunicipality   string `gorm:"not null"`
	PresentAddProvince       string `gorm:"not null"`
	HomeAddStreet            string `gorm:"not null"`
	HomeAddBarangay          string `gorm:"not null"`
	HomeAddMunicipality      string `gorm:"not null"`
	HomeAddProvince          string `gorm:"not null"`
	BusinessAddStreet        string `gorm:"not null"`
	BusinessAddBarangay      string `gorm:"not null"`
	BusinessAddMunicipality  string `gorm:"not null"`
	BusinessAddProvince      string `gorm:"not null"`
	SpouseFirstName          string `gorm:"not null"`
	SpouseMiddleName         string `gorm:"not null"`
	SpouseLastName           string `gorm:"not null"`
	SpouseSuffixName         string `gorm:"not null"`
	SpouseCitizenship        string `gorm:"not null"`
	SpouseCardMember         string `gorm:"not null"`
	SpouseBirthday           string `gorm:"not null"`
	SpouseAge                string `gorm:"not null"`
	SpousePlaceOfBirth       string `gorm:"not null"`
	SpouseGender             string `gorm:"not null"`
	SpouseOccupation         string `gorm:"not null"`
	SpouseSourceOfIncome     string `gorm:"not null"`
	SpouseBusiness           string `gorm:"not null"`
	SpouseMobile             string `gorm:"not null"`
	PPIScore                 string `gorm:"not null"`
}
