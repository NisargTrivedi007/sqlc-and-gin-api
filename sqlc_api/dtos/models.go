package dtos

type TodoDTO struct {
	ID          int32  `json:"id"`
	Todo        string `json:"todo"`
	CreatedDate string `json:"created_date"`
	CreatedBy   int    `json:"created_by"`
	UpdatedDate string `json:"updated_date"`
}

type CreateTodoDTO struct {
	Todo      string `json:"todo"`
	CreatedBy int    `json:"created_by"`
}
type UpdateTodoDTO struct {
	ID          int32  `json:"id"`
	Todo        string `json:"todo"`
	UpdatedDate string `json:"updated_date"`
}
type UserDTO struct {
	ID          int32  `json:"id"`
	Username    string `json:"username"`
	EmailId     string `json:"email_id"`
	PhoneNo     string `json:"phone_no"`
	Password    string `json:"password"`
	CreatedDate string `json:"created_date"`
}
