package test

type AuthBody struct {
	Email	string	`json:"email"`
	Password	string	`json:"password"`
}

type CreateProfile struct {
	Name             string `json:"name"`
	Age              int32  `json:"age"`
	City             string `json:"city"`
	CurrentSituation string `json:"current_situation"`
	Gender           string `json:"gender"`
}

type AddPreference struct {
    ActivityID int32 `json:"activity_id"`
}