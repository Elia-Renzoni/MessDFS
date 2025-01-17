package internal_test

import (
	"net/url"
	"storageservice/internal"
	"testing"
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

	var (
		query2 = []string{
			"12",
			"Paolo",
			"Rossi",
		}
		query3 = []string{
			"34",
			"Marta",
			"Grossi",
		}
	)

	if err := controller.CreateNewDir("test_dir3"); err != nil {
		t.Errorf("%v", err)
	}

	if err := controller.WriteRemoteCSV("test_dir3", "test_file.csv", "insert", query); err != nil {
		t.Errorf("%v", err)
	}

	if err := controller.WriteRemoteCSV("test_dir3", "test_file.csv", "insert", query2); err != nil {
		t.Errorf("%v", err)
	}

	if err := controller.WriteRemoteCSV("test_dir3", "test_file.csv", "insert", query3); err != nil {
		t.Errorf("%v", err)
	}
}

func TestReadInRemoteCSV(t *testing.T) {
	controller := internal.ResourceController{}
	v := url.Values{}
	v.Set("id", "12")

	_, err := controller.ReadInRemoteCSV("test_dir3", "test_file.csv", "read", v)
	if err != nil {
		t.Errorf("%v", err)
	}
}

func TestUpdateInRemoteCSV(t *testing.T) {
	controller := internal.ResourceController{}
	v := map[string][]string{
		"Name": []string{"12", "Paolo", "Pippo"},
	}

	if err := controller.UpdateRemoteCSV("test_dir3", "test_file.csv", "update", v); err != nil {
		t.Errorf("%v", err)
	}
}

func TestCreateNewDir2(t *testing.T) {
	controller := &internal.ResourceController{}
	if err := controller.CreateNewDir("test_dir4"); err != nil {
		t.Errorf("%v", err)
	}
}
