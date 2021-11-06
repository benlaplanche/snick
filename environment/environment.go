package environment

import "os"

func DetectENV() (response string) {

	switch {
	case os.Getenv("GITHUB_ACTIONS") != "":
		response = "gh-actions"
	default:
		response = "local"
	}

	return
}
