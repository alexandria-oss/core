package httputil

type GenericResponse struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}
