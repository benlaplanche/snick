package environment_test

import (
	"os"

	. "github.com/benlaplanche/snick/environment"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Environment", func() {
	Context("It is Github Actions", func() {

		BeforeEach(func() {
			os.Setenv("CI", "true")
			os.Setenv("GITHUB_ACTIONS", "true")
		})

		AfterEach(func() {
			os.Unsetenv("CI")
			os.Unsetenv("GITHUB_ACTIONS")
		})

		It("should say github actions", func() {
			result := DetectENV()
			Expect(result).To(Equal("GitHub Actions"))
		})
	})

	Context("It is local development", func() {
		It("should say local", func() {
			result := DetectENV()
			Expect(result).To(Equal("local"))
		})
	})

})
