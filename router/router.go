package router

import (
	"github.com/rancher/go-rancher/api"

	"github.com/gorilla/mux"
	"github.com/rancher/v2-api/model"
	"github.com/rancher/v2-api/server"
)

func New(s *server.Server) *mux.Router {
	schemas := model.NewSchema()
	router := mux.NewRouter().StrictSlash(true)

	// API framework routes
	router.Methods("GET").Path("/").Handler(api.VersionsHandler(schemas, "v1", "v2"))
	router.Methods("GET").Path("/v2/schemas").Handler(api.SchemasHandler(schemas))
	router.Methods("GET").Path("/v2/schemas/{id}").Handler(api.SchemaHandler(schemas))
	router.Methods("GET").Path("/v2").Handler(api.VersionHandler(schemas, "v2"))

	f := s.HandlerFunc
	router.Methods("GET").Path("/v2/services").Handler(f(schemas, s.ServiceList))
	router.Methods("GET").Path("/v2/service").Handler(f(schemas, s.ServiceList))
	router.Methods("GET").Path("/v2/services/{id}").Handler(f(schemas, s.ServiceByID))
	router.Methods("GET").Path("/v2/service/{id}").Handler(f(schemas, s.ServiceByID))

	router.Methods("GET").Path("/v2/containers").Handler(f(schemas, s.ContainerList))
	router.Methods("GET").Path("/v2/container").Handler(f(schemas, s.ContainerList))
	router.Methods("GET").Path("/v2/containers/{id}").Handler(f(schemas, s.ContainerByID))
	router.Methods("GET").Path("/v2/container/{id}").Handler(f(schemas, s.ContainerByID))
	router.Methods("POST").Path("/v2/containers").Handler(f(schemas, s.ContainerCreate))

	return router
}
