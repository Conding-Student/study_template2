package handler

type CagabayVersionDetail struct {
	ID                int    `json:"id"`
	Version           int    `json:"version"`
	UpdateDescription string `json:"update_description"`
	UpdateMessage     string `json:"update_message"`
	AndroidLink       string `json:"android_link"`
	IOSLink           string `json:"ios_link"`
	ButtonName        string `json:"button_name"`
}

type InstitutionEstablished struct {
	ID              string `json:"id"`
	Institutions    string `json:"institutions"`
	DateEstablished string `json:"date_established"`
	Image           string `json:"image"`
}
