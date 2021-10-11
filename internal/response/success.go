package response

// ErrorResponse is the response that represents an error.
type Response struct {
	Status string `json:"status"`
}

func SuccessResponse() Response {
	return Response{Status: "Success"}
}
