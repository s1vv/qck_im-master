package models

type LoginRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegistrarionRequest struct {
	Login           string `json:"login" binding:"required"`
	Email           string `json:"email" binding:"required"`
	Password        string `json:"password" binding:"required"`
	QckLink         string `json:"qck_link" binding:"required"`
	QckLinkPassword string `json:"qck_link_password" binding:"required"`
}

type UpdateDataLink struct {
	QckLink     string `json:"qck_link" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type ActivationData struct {
	QckLink  string `json:"qck_link" binding:"required"`
	Password string `json:"password" binding:"required"`
}
