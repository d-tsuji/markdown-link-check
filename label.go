package mlc

import "fmt"

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

func red(s string) string {
	return fmt.Sprintf("\u001B[31m%s\u001B[0m", s)
}

func yellow(s string) string {
	return fmt.Sprintf("\u001B[33m%s\u001B[0m", s)
}

func green(s string) string {
	return fmt.Sprintf("\u001B[32m%s\u001B[0m", s)
}

func cyan(s string) string {
	return fmt.Sprintf("\u001b[36m%s\u001B[0m", s)
}
