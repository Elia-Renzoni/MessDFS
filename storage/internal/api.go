package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"storageservice/model"
	"strings"
)

type MessDFSStorageAPI struct {
	// IP address + listen port
	address string

	// embedded structure
	ResourceController
	serviceConn *AuthServiceTrigger

	// data plane structures
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
		address:     address,
		serviceConn: NewAuthServiceTrigger(),
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

	if isTransactionOk := m.serviceConn.CheckTransactionOwner(m.insert.TransactionUser, m.insert.User); !isTransactionOk {
		writer(w, map[string]string{"err": "Transaction Not Allowed"}, 400)
		return
	}

	fmt.Printf("Metodo ==> %v", m.insert.QueryType)
	fmt.Printf("File Name ==> %v", m.insert.FileName)
	if err := m.WriteRemoteCSV(m.insert.User, m.insert.FileName, m.insert.QueryType, m.insert.QueryContent); err != nil {
		http.Error(w, err.Error(), 500)
	} else {
		message := map[string]string{
			"success": "Informations Added Succesfully",
		}

		writer(w, message, 201)
	}
}

func (m *MessDFSStorageAPI) deleteCSV(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	splittedPath := strings.Split(path, "/")

	m.delete.TransactionUser = splittedPath[2]
	m.delete.User = splittedPath[3]
	m.delete.FileName = splittedPath[4]
	m.delete.Query = r.URL.Query()

	if isTransactionOk := m.serviceConn.CheckTransactionOwner(m.delete.TransactionUser, m.delete.User); !isTransactionOk {
		writer(w, map[string]string{"err": "Transaction Not Allowed"}, 400)
		return
	}

	if err := m.DeleteRemoteCSV(m.delete.User, m.delete.FileName, "delete", m.delete.Query); err != nil {
		http.Error(w, err.Error(), 500)
	} else {
		message := map[string]string{
			"success": "Informations Succesfully Removed",
		}

		writer(w, message, 200)
	}
}

func (m *MessDFSStorageAPI) readCSV(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	splittedPath := strings.Split(path, "/")

	m.read.TransactionUser = splittedPath[2]
	m.read.Friend = splittedPath[3]
	m.read.User = splittedPath[4]
	m.read.FileName = splittedPath[5]
	m.read.Query = r.URL.Query()

	if isTransactionOk := m.serviceConn.CheckFriendship(m.read.TransactionUser, m.read.Friend); !isTransactionOk {
		writer(w, map[string]string{"err": "Transaction Not Allowed"}, 400)
		return
	}

	fmt.Printf("%s - %s - %s \n", m.read.User, m.read.FileName, m.read.Query)

	response, err := m.ReadInRemoteCSV(m.read.User, m.read.FileName, "read", m.read.Query)
	if err != nil {
		http.Error(w, err.Error(), 500)
	} else {
		writer(w, response, 200)
	}
}

func (m *MessDFSStorageAPI) updateCSV(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, _ := io.ReadAll(r.Body)
	json.Unmarshal(body, &m.update)

	if isTransactionOk := m.serviceConn.CheckTransactionOwner(m.update.TransactionUser, m.update.User); !isTransactionOk {
		writer(w, map[string]string{"err": "Transaction Not Allowed"}, 400)
		return
	}

	if err := m.UpdateRemoteCSV(m.update.User, m.update.FileName, m.update.QueryType, m.update.QueryContent); err != nil {
		http.Error(w, err.Error(), 500)
	} else {
		message := map[string]string{
			"success": "Informations Succesfully Updated",
		}

		writer(w, message, 201)
	}
}

func (m *MessDFSStorageAPI) createDirectory(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, _ := io.ReadAll(r.Body)
	json.Unmarshal(body, &m.createDir)

	fmt.Printf("%s OK", m.createDir.TransactionUser)
	if isTransactionOk := m.serviceConn.CheckTransactionOwner(m.createDir.TransactionUser, m.createDir.DirToCreate); !isTransactionOk {
		writer(w, map[string]string{"err": "Transaction Not Allowed"}, 400)
		return
	}

	if err := m.CreateNewDir(m.createDir.DirToCreate); err != nil {
		http.Error(w, err.Error(), 500)
	} else {
		message := map[string]string{
			"success": "Directory Succesfully Created",
		}

		writer(w, message, 201)
	}
}

func (m *MessDFSStorageAPI) deleteDirectory(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	splittedPath := strings.Split(path, "/")
	m.deleteDir.TransactionUser = splittedPath[2]
	m.deleteDir.DirToDelete = splittedPath[3]

	if isTransactionOk := m.serviceConn.CheckTransactionOwner(m.deleteDir.TransactionUser, m.deleteDir.DirToDelete); !isTransactionOk {
		writer(w, map[string]string{"err": "Transaction Not Allowed"}, 400)
		return
	}

	if okDelete := m.serviceConn.DeleteDirectoryInAuth(m.deleteDir.DirToDelete); !okDelete {
		http.Error(w, "Operation Not Allowed", 500)
		return
	}

	if err := m.DeleteDir(m.deleteDir.DirToDelete); err != nil {
		http.Error(w, err.Error(), 500)
	} else {
		message := map[string]string{
			"success": "Directory Succesfully Deleted",
		}

		writer(w, message, 200)
	}
}

func (m *MessDFSStorageAPI) deleteFile(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	splittedPath := strings.Split(path, "/")
	m.deleteF.TransactionUser = splittedPath[2]
	m.deleteF.DirName = splittedPath[3]
	m.deleteF.FileToDelete = splittedPath[4]

	if isTransactionOk := m.serviceConn.CheckTransactionOwner(m.deleteF.TransactionUser, m.deleteF.DirName); !isTransactionOk {
		writer(w, map[string]string{"err": "Transaction Not Allowed"}, 400)
		return
	}

	if err := m.DeleteFile(m.deleteF.DirName, m.deleteF.FileToDelete); err != nil {
		http.Error(w, err.Error(), 500)
	} else {
		message := map[string]string{
			"success": "File Succesfully Deleted",
		}

		writer(w, message, 200)
	}
}

func writer(w http.ResponseWriter, data map[string]string, statusCode int) {
	jsonData, _ := json.Marshal(data)

	fmt.Println(data)

	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
