package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSnick(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Snick Suite")
}
