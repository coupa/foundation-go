package persistence_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestPersistence(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Persistence Suite")
}

type TestEntity struct {
	ID      int64  `db:"id" json:"id"`
	StringA string `db:"string_a" json:"string_a"`
	IntA    int    `db:"int_a" json:"int_a"`
}
