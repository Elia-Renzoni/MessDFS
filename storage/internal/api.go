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
	effectiveOperations ResourceController

	createDir CreateDirPayload
	deleteDir DeleteDirPayload
	deleteFile DeleteFilePayload
	deleteCSVContent DeletePayload
	insertCSVContent InsertPayload
	readCSVContent ReadPayload
	updateCSVContent UpdatePayload
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
	json.Unmarshal(body, m.insertCSVContent)

	if err := m.effectiveOperations.WriteRemoteCSV(m.insertCSVContent.User, 
												   m.insertCSVContent.FileName, 
												   m.insertCSVContent.QueryType, 
												   m.insertCSVContent.QueryContent); err != nil {
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

	m.deleteCSVContent.User = splittedPath[2]
	m.deleteCSVContent.FileName = splittedPath[3]
	m.deleteCSVContent.Query = r.URL.Query()

	if err := m.effectiveOperations.DeleteRemoteCSV(m.deleteCSVContent.User, 
													m.deleteCSVContent.FileName, 
													m.deleteCSVContent.Query); err != nil {
		http.Error(w, err.Error(), 500)
	} else {
		message := map[string]string{
			"success": "Informations Succesfully Removed",
		}

		writer(w, message)
	}
}

func (m *MessDFSStorageAPI) read(w http.ResponseWriter, r *http.Request) {

}

func (m *MessDFSStorageAPI) update(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, _ := io.ReadAll(r.Body)
	json.Unmarshal(body, m.updateCSVContent)

	if err := m.effectiveOperations.UpdateRemoteCSV(m.updateCSVContent.User, 
													m.updateCSVContent.FileName, 
													m.updateCSVContent.QueryType, 
													m.updateCSVContent.QueryContent); err != nil {
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

	if err := m.effectiveOperations.CreateNewDir(m.createDir.DirToCreate); err != nil {
		http.Error(w, err.Error(), 500)
	} else {
		message := map[string]string{
			"success": "Directory Succesfully Created",
		}

		writer(w, message)
	}
}

func (m *MessDFSStorageAPI) deleteDirectory(w http.ResponseWriter, r *http.Request) {

}

func (m *MessDFSStorageAPI) deleteFile(w http.ResponseWriter, r *http.Request) {

}

func wirter(w http.ResponseWriter, data map[string]string) {
	encoder := json.NewEncoder(w)
	
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	encoder.Encode(data)
}