package markdownlink

func getStatusLabel(statusCode int) string {
	if statusCode >= 500 {
		return yellow("⚠")
	} else if statusCode >= 400 {
		return red("✖")
	} else if statusCode >= 200 {
		return green("✓")
	} else {
		return yellow("⚠")
	}
}
