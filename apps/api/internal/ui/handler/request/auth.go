package request

type User struct {
	Username string `json:"username"`
}

type FinishUserRegister struct {
	Username string `json:"username"`
}
