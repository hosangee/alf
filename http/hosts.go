package http

import (
	log "code.google.com/p/log4go"
	"github.com/emicklei/go-restful"
	"github.com/msurdi/alf/db"
	"net/http"
)

// newHostService returns a *restful.WebService for handling /hosts requests
func (self *HttpService) newHostsService(rootPath string) (ws *restful.WebService) {
	ws = newWs(rootPath)

	// Get all hosts
	findAllHosts := func(request *restful.Request, response *restful.Response) {
		var hosts []db.Host
		err := self.db.Hosts.FindAll(&hosts)
		if err != nil {
			log.Info("Error processing request: " + err.Error())
			response.WriteError(http.StatusInternalServerError, err)
		} else {
			response.WriteEntity(hosts)
		}
	}
	ws.Route(ws.GET("/").
		Doc("Get all the hosts").To(findAllHosts))

	// Get one host by id
	findOneHost := func(request *restful.Request, response *restful.Response) {
		id := request.PathParameter("id")
		var host db.Host
		err := self.db.Hosts.FindById(id, &host)
		if err != nil {
			log.Info("Error processing request: " + err.Error())
			response.WriteError(http.StatusInternalServerError, err)
		} else {
			response.WriteEntity(host)
		}
	}
	ws.Route(ws.GET("/{id}").
		Doc("Get all the hosts").
		Param(ws.PathParameter("id", "the identifier for a host")).
		To(findOneHost))

	// Create a new host
	createHost := func(request *restful.Request, response *restful.Response) {
		host := &db.Host{}
		err := request.ReadEntity(host)
		if err != nil {
			response.WriteError(http.StatusInternalServerError, err)
		} else {
			err := self.db.Hosts.Save(host)
			if err != nil {
				response.WriteError(http.StatusInternalServerError, err)

			} else {
				response.WriteHeader(http.StatusCreated)
				response.WriteEntity(host)
			}
		}
	}
	ws.Route(ws.POST("/").
		Doc("Create a new host").
		To(createHost))

	return
}
