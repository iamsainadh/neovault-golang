package model

type CreateOrganizationInput struct {
	Id     string
	Name   string
	UserId string
}

type CreateOrganizationResponse struct {
	Id    string
	Name  string
	Token string
}
