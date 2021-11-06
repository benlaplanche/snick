package environment_test

import (
	. "github.com/benlaplanche/snick/environment"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Environment", func() {
	result := DetectENV()
	Expect(result).To(Equal("hello"))

})
