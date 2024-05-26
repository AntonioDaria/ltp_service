package api

type JSONError struct {
	Error string `json:"error"`
}

type JSONBadRequestError struct {
	Errors []string `json:"errors"`
}
