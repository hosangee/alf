package http

import (
	"github.com/emicklei/go-restful"
)

// newPingService returns a *restful.WebService for handling /ping requests
func (self *HttpService) newPingService(rootPath string) (ws *restful.WebService) {
	// Get
	getPing := func(request *restful.Request, response *restful.Response) {
		response.Write([]byte("pong"))
	}

	// Mapping
	ws = newWs(rootPath)

	ws.Route(ws.GET("/").
		// for documentation
		Doc("Test if alf is alive").To(getPing))

	return ws
}
