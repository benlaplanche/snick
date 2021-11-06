package environment_test

import (
	. "github.com/benlaplanche/snick/environment"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Environment Suite", func() {
	Context("It is Github Actions", func() {

		// os.Setenv("CI", "true")
		// defer os.Unsetenv("CI")

		// os.Setenv("GITHUB_ACTIONS", "true")
		// defer os.Unsetenv(("GITHUB_ACTIONS"))

		result := DetectENV()
		Expect(result).To(Equal("GitHub Actions"))
	})

})
