package httputil

type GenericResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}
