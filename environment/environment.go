package environment

import "os"

func DetectENV() (response string) {

	switch {
	case os.Getenv("GITHUB_ACTIONS") != "":
		response = "GitHub Actions"
	default:
		response = "Local Development"
	}

	return
}
