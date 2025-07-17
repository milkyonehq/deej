package provider

func stringOrDefault(value string, defaultValue string) string {
	if value == "" {
		return defaultValue
	}

	return value
}
