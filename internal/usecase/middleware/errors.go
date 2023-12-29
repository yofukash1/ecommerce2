package middleware

func EmptyFields(fields ...interface{}) bool {
	for _, f := range fields {
		if f.(string) == "" {
			return true
		}
	}
	return false
}
