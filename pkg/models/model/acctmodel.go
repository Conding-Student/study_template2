package model

type UserRequestHeader struct {
	StaffID        string
	Status         string
	Remarks        string
	RoleID         int
	PersonalInfo   PersonalInformation
	AddressInfo    AddressInformation
	ContactInfo    ContactInformation
	EmploymentInfo EmploymentInformation
	AccountInfo    AccountInformation
	DeviceInfo     DeviceInformation
}

type PersonalInformation struct {
	FirstName   string
	MiddleName  string
	LastName    string
	NickName    string
	Birthdate   string
	Age         string
	Gender      string
	CivilStatus string
}

type AddressInformation struct {
	AddressType int64
	Sitio       string
	Barangay    string
	City        string
	Province    string
	ZipCode     string
}

type ContactInformation struct {
	EmailAddress string
	Mobile       string
}

type EmploymentInformation struct {
	Cluster            string
	Institution        string
	Region             string
	Department         string
	Area               string
	Designation        string
	Unit               string
	DateHired          string
	DateRegularization string
	EmploymentType     string
	Assignment         string
	JobLevel           string
	JobGrade           string
	EmploymentStatus   string
}

type AccountInformation struct {
	Username string
	Password string
}

type DeviceInformation struct {
	DeviceUsed  string
	DeviceModel string
	DeviceID    string
}

// -------------------------------------------------------------- //

type LoginCredentials struct {
	StaffId  string `json:"staffid"`
	Password string `json:"password"`
}
