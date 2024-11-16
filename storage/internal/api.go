package internal

import (
	"net"
	"net/http"
	"storage/middleware"
)

type MessDFSStorageAPI struct {
	address string
	jwtMiddleware *Middleware 
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


}

func (m *MessDFSStorageAPI) delete(w http.ResponseWriter, r *http.Request) {

}

func (m *MessDFSStorageAPI) read(w http.ResponseWriter, r *http.Request) {

}

func (m *MessDFSStorageAPI) update(w http.ResponseWriter, r *http.Request) {
	
}

func (m *MessDFSStorageAPI) createDirectory(w http.ResponseWriter, r *http.Request) {

}

func (m *MessDFSStorageAPI) deleteDirectory(w http.ResponseWriter, r *http.Request) {

}

func (m *MessDFSStorageAPI) deleteFile(w http.ResponseWriter, r *http.Request) {

}