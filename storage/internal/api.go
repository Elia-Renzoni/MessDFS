package internal

import (
	"net"
	"net/http"
	"storage/middleware"
	"storage/model"
	"encoding/json"
	"io"
	"os"
	"strings"
)

type MessDFSStorageAPI struct {
	address string
	jwtMiddleware *Middleware
	ResourceController

	createDir CreateDirPayload
	deleteDir DeleteDirPayload
	deleteFile DeleteFilePayload
	delete DeletePayload
	insert InsertPayload
	read ReadPayload
	update UpdatePayload
}

func NewStorage(address string) *MessDFSStorageAPI {
	return &MessDFSStorageAPI{
		address: address,
		jwtMiddleware: NewMiddleware(),
	}
}

func (m *MessDFSStorageAPI) Start() {
	mux := http.NewServerMux()

	mux.HandleFunc("POST /csvi", m.insert)
	mux.HandleFunc("UPDATE /csvu", m.update)
	mux.HandleFunc("DELETE /csvd/", m.delete)
	mux.HandleFunc("READ /csvr/", m.delete)
	mux.HandleFunc("POST /ndir", m.createDirectory)
	mux.HandleFunc("DELETE /ddir/", m.deleteDirectory)
	mux.HandleFunc("DELETE /dfile/", m.deleteFile)

	http.ListenAndServe(address, nil)
}


func (m *MessDFSStorageAPI) insert(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, _ := io.ReadAll(r.Body)
	json.Unmarshal(body, m.insert)

	if err := m.WriteRemoteCSV(m.insert.User, m.insert.FileName, m.insert.QueryType, m.insert.QueryContent); err != nil {
		http.Error(w, err.Error(), 500)
	} else {
		message := map[string]string{
			"success": "Informations Added Succesfully",
		}

		writer(w, message)
	}
}

func (m *MessDFSStorageAPI) delete(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	splittedPath := strings.Split("/")

	m.delete.User = splittedPath[2]
	m.delete.FileName = splittedPath[3]
	m.delete.Query = r.URL.Query()

	if err := m.DeleteRemoteCSV(m.delete.User, m.delete.FileName, m.delete.Query); err != nil {
		http.Error(w, err.Error(), 500)
	} else {
		message := map[string]string{
			"success": "Informations Succesfully Removed",
		}

		writer(w, message)
	}
}

func (m *MessDFSStorageAPI) read(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	splittedPath := strings.Split(path, "/")

	m.read.User = splittedPath[2]
	m.read.FileName = splittedPath[3]
	m.read.Query = r.URL.Query()

	response, err := m.ReadInRemoteCSV(m.read.User, m.read.FileName, "read", m.read.Query)
	if err != nil {
		http.Error(w, err.Error(), 500)
	} else {
		writer(w, response)
	}
}

func (m *MessDFSStorageAPI) update(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, _ := io.ReadAll(r.Body)
	json.Unmarshal(body, m.update)

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
	json.Unmarshal(body, m.createDir)

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
	m.deleteDir.DirToDelete = splittedPath[2]

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
	m.deleteFile.DirName = splittedPath[2]
	m.deleteFile.FileToDelete = splittedPath[3]

	if err := m.DeleteFile(m.deleteFile.DirName, m.deleteFile.FileToDelete); err != nil {
		http.Error(w, err.Error(), 500)
	} else {
		message := map[string]string{
			"success": "File Succesfully Deleted",
		}

		writer(w, message)
	}
}

func wirter(w http.ResponseWriter, data map[string]string) {
	encoder := json.NewEncoder(w)
	
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	encoder.Encode(data)
}