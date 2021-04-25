package kessaku_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestKessaku(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Kessaku Suite")
}
