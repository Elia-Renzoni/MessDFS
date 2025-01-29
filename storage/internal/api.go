package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"storageservice/internal/middleware"
	"storageservice/model"
	"strings"
)

type MessDFSStorageAPI struct {
	address       string
	jwtMiddleware *middleware.Middleware
	ResourceController

	createDir model.CreateDirPayload
	deleteDir model.DeleteDirPayload
	deleteF   model.DeleteFilePayload
	delete    model.DeletePayload
	insert    model.InsertPayload
	read      model.ReadPayload
	update    model.UpdatePayload
}

func NewStorage(address string) *MessDFSStorageAPI {
	return &MessDFSStorageAPI{
		address:       address,
		jwtMiddleware: middleware.NewMiddleware(),
	}
}

func (m *MessDFSStorageAPI) Start() {
	mux := http.NewServeMux()

	fmt.Printf("Server Listening...")

	mux.HandleFunc("/csvi", m.insertCSV)
	mux.HandleFunc("/csvu", m.updateCSV)
	mux.HandleFunc("/csvd/", m.deleteCSV)
	mux.HandleFunc("/csvr/", m.readCSV)
	mux.HandleFunc("/ndir", m.createDirectory)
	mux.HandleFunc("/ddir/", m.deleteDirectory)
	mux.HandleFunc("/dfile/", m.deleteFile)

	http.ListenAndServe(m.address, mux)
}

func (m *MessDFSStorageAPI) insertCSV(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, _ := io.ReadAll(r.Body)
	json.Unmarshal(body, &m.insert)

	fmt.Printf("Metodo ==> %v", m.insert.QueryType)
	fmt.Printf("File Name ==> %v", m.insert.FileName)
	if err := m.WriteRemoteCSV(m.insert.User, m.insert.FileName, m.insert.QueryType, m.insert.QueryContent); err != nil {
		http.Error(w, err.Error(), 500)
	} else {
		message := map[string]string{
			"success": "Informations Added Succesfully",
		}

		writer(w, message)
	}
}

func (m *MessDFSStorageAPI) deleteCSV(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	splittedPath := strings.Split(path, "/")

	m.delete.TransactionUser = splittedPath[2]
	m.delete.User = splittedPath[3]
	m.delete.FileName = splittedPath[4]
	m.delete.Query = r.URL.Query()

	if err := m.DeleteRemoteCSV(m.delete.User, m.delete.FileName, "delete", m.delete.Query); err != nil {
		http.Error(w, err.Error(), 500)
	} else {
		message := map[string]string{
			"success": "Informations Succesfully Removed",
		}

		writer(w, message)
	}
}

func (m *MessDFSStorageAPI) readCSV(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	splittedPath := strings.Split(path, "/")

	m.read.TransactionUser = splittedPath[2]
	m.read.User = splittedPath[3]
	m.read.FileName = splittedPath[4]
	m.read.Query = r.URL.Query()

	fmt.Printf("%s - %s - %s \n", m.read.User, m.read.FileName, m.read.Query)

	response, err := m.ReadInRemoteCSV(m.read.User, m.read.FileName, "read", m.read.Query)
	if err != nil {
		http.Error(w, err.Error(), 500)
	} else {
		writer(w, response)
	}
}

func (m *MessDFSStorageAPI) updateCSV(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, _ := io.ReadAll(r.Body)
	json.Unmarshal(body, &m.update)

	if err := m.UpdateRemoteCSV(m.update.User, m.update.FileName, m.update.QueryType, m.update.QueryContent); err != nil {
		http.Error(w, err.Error(), 500)
	} else {
		message := map[string]string{
			"success": "Informations Succesfully Updated",
		}

		writer(w, message)
	}
}

func (m *MessDFSStorageAPI) createDirectory(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, _ := io.ReadAll(r.Body)
	json.Unmarshal(body, &m.createDir)

	if err := m.CreateNewDir(m.createDir.DirToCreate); err != nil {
		http.Error(w, err.Error(), 500)
	} else {
		message := map[string]string{
			"success": "Directory Succesfully Created",
		}

		writer(w, message)
	}
}

func (m *MessDFSStorageAPI) deleteDirectory(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	splittedPath := strings.Split(path, "/")
	m.deleteDir.TransactionUser = splittedPath[2]
	m.deleteDir.DirToDelete = splittedPath[3]

	if err := m.DeleteDir(m.deleteDir.DirToDelete); err != nil {
		http.Error(w, err.Error(), 500)
	} else {
		message := map[string]string{
			"success": "Directory Succesfully Deleted",
		}

		writer(w, message)
	}
}

func (m *MessDFSStorageAPI) deleteFile(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	splittedPath := strings.Split(path, "/")
	m.deleteFile.TransactionUser = splittedPath[2]
	m.deleteF.DirName = splittedPath[3]
	m.deleteF.FileToDelete = splittedPath[4]

	if err := m.DeleteFile(m.deleteF.DirName, m.deleteF.FileToDelete); err != nil {
		http.Error(w, err.Error(), 500)
	} else {
		message := map[string]string{
			"success": "File Succesfully Deleted",
		}

		writer(w, message)
	}
}

func writer(w http.ResponseWriter, data map[string]string) {
	jsonData, _ := json.Marshal(data)

	fmt.Println(data)

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
