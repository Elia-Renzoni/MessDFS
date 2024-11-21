package internal_test

import (
	"testing"
	"storageservice/internal"
)

func TestCreateNewDir(t *testing.T) {
	controller := &internal.ResourceController{}
	if err := controller.CreateNewDir("test_dir"); err != nil {
		t.Errorf("%v", err)
	}
}

