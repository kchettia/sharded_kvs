package routes

type KeyCountResp struct {
	Message  string `json:"message"`
	KeyCount int    `json:"key-count"`
}

type PutResp struct {
	Message  string `json:"message"`
	Replaced bool   `json:"replaced"`
	Address  string `json:"address"`
}

type PutErrorResp struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

type GetResp struct {
	Message   string `json:"message"`
	DoesExist bool   `json:"doesExist"`
	Value     string `json:"value"`
}

type GetErrorResp struct {
	Message   string `json:"message"`
	DoesExist bool   `json:"doesExist"`
	Error     string `json:"error"`
}

type DelResp struct {
	Message   string `json:"message"`
	DoesExist bool   `json:"doesExist"`
	Address   string `json:"address"`
}

type DelErrorResp struct {
	Message   string `json:"message"`
	DoesExist bool   `json:"doesExist"`
	Error     string `json:"error"`
}

type ViewChangeErrorResp struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}
type ViewChangeResp struct {
	Message string `json:"message"`
}
