package render

import (
	"os"
	"testing"
)

func TestInitTemplates(t *testing.T) {
	if _, err := os.Stat("../templates"); err == nil {
		_ = os.Chdir("..")
	}

	InitTemplates()

	if Templates == nil {
		t.Fatal("Expected Templates to be initialized, got nil")
	}
}
