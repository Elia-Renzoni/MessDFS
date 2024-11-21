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

func TestDeleteDir(t *testing.T) {
	controller := &internal.ResourceController{}

	controller.CreateNewDir("test_dir2")

	if err := controller.DeleteDir("test_dir2"); err != nil {
		t.Errorf("%v", err)
	}
}

func TestWriteRemoteCSV(t *testing.T) {
	controller := &internal.ResourceController{}
	var query = []string{
		"ID",
		"Name",
		"Surname",
	}

	controller.CreateNewDir("test_dir2")

	if err := controller.WriteRemoteCSV("test_dir2", "test_file.csv", "insert", query); err != nil {
		t.Errorf("%v", err) // fail
	}
} 

