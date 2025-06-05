package custom_err

type CustomError struct {
	Message string
	Type    string
}

func (r CustomError) Error() string {
	return r.Message
}

func (r CustomError) Extensions() map[string]interface{} {
	return map[string]interface{}{
		"message": r.Message,
		"type":    r.Type,
	}
}