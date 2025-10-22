package responses

type SuccessResponse struct {
	Message string   `json:"message"`
	Results []Result `json:"results"`
}

type Result struct {
	Operation string `json:"operation"`
	Success   bool   `json:"success"`
}
