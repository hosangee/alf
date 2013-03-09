/*
	Package http implements the http api functionality. We expose a single HttpService type
	which contains which contains instances of the services implementing the different api
	entrypoints we have.
*/
package http

import (
	"github.com/emicklei/go-restful"
	"github.com/msurdi/alf/db"
	"net/http"
)

type HttpService struct {
	db          *db.DB
	bindAddress string
}

func NewHttpService(bindAddress string, db *db.DB) *HttpService {
	return &HttpService{bindAddress: bindAddress, db: db}
}

// HTTPService setups the handlers for the API entry points, and runs the server
func (self *HttpService) Start() {
	restful.Add(self.newPingService("/ping"))
	restful.Add(self.newHostsService("/host"))
	http.ListenAndServe(self.bindAddress, nil)
}

func newWs(rootPath string) (ws *restful.WebService) {
	ws = new(restful.WebService)
	ws.Path(rootPath).
		Consumes(restful.MIME_JSON, restful.MIME_XML).
		Produces(restful.MIME_JSON, restful.MIME_XML)
	return
}
