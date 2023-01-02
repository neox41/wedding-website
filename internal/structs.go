package internal

type AttendeesOutput struct {
	ID           int    `json:"id"`
	IDNucleo          int    `json:"nucleo"`
	NucleoName string `json:"nucleoname"`
	Name         string `json:"name"`
	Surname     string `json:"surname"`
	Email   string `json:"email"`
	Status       string `json:"status"`
	Link       string `json:"link"`
	English       bool `json:"english"`

	Vegan      bool `json:"vegan"`
	Vegetarian bool `json:"vegetarian"`
	Requirements string `json:"requirements"`
	Transport bool `json:"transport"`
	Attendance string `json:"attendance"`

	Category string `json:"category"`
	Abroad bool `json:"abroad"`
	Bambino bool `json:"bambino"`
	Beve bool `json:"beve"`
	Table string `json:"table"`
	Replied bool `json:"replied"`
}
type FamiliesOutput struct {
	Name string `json:"name"`
	IDNucleo          int    `json:"nucleo"`
	Email   string `json:"email"`
	English       int `json:"english"`
	Link       string `json:"link"`
	Replied       bool `json:"replied"`
	Status       string `json:"status"`
}
type UpdateStatus struct {
	Id int `json:"id"`
	Status  string `json:"status"`
}
type UpdateFamily struct {
	Email   string `json:"email"`
	Id       int `json:"id"`
}

type PageInfo struct{
	TitleEN string
	TitleIT string
	NameEN string
	NameIT string
	Location string
	Active bool
}
type Attendee struct{
	ID int
	Name string
	Surname string
	Vegan      bool
	Vegetarian bool
	Requirements string
	Transport bool
	Attendance string
}
type Pages struct{
	Pages []PageInfo
	Error string
	English bool
	LinkLang string
	IdInvite string
	Attendees []Attendee
}