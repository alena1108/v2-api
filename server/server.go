package server

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/jmoiron/sqlx"
	"github.com/rancher/go-rancher/api"
	"github.com/rancher/go-rancher/client"
	"strconv"
)

type Server struct {
	DB                 *sqlx.DB
	driver, driverName string
}

type SchemaConvertor interface {
	FromSchema(obj interface{}) (interface{}, error)
	ToSchema(obj interface{}) (interface{}, error)
}

func New(driver, driverName string) (*Server, error) {
	db, err := sqlx.Open(driver, driverName)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Server{
		driver:     driver,
		driverName: driverName,
		DB:         db,
	}, err
}

func (s *Server) namedQuery(query string, args map[string]interface{}) (*sqlx.Rows, error) {
	rows, err := s.DB.NamedQuery(query, args)
	return rows, err
}

func (s *Server) handleError(rw http.ResponseWriter, r *http.Request, err error) {
	apiError := client.ServerApiError{
		Type:    "error",
		Status:  500,
		Code:    "ServerError",
		Message: err.Error(),
	}
	data, err := json.Marshal(&apiError)
	if err == nil {
		rw.Header().Add("Content-Type", "application/json")
		rw.WriteHeader(apiError.Status)
		rw.Write(data)
	} else {
		rw.WriteHeader(http.StatusInternalServerError)
		logrus.Errorf("Fail to marshall: %v", err)
	}
}

func (s *Server) HandlerFunc(schemas *client.Schemas, f func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return api.ApiHandlerFunc(schemas, func(rw http.ResponseWriter, r *http.Request) {
		if err := f(rw, r); err != nil {
			s.handleError(rw, r, err)
		}
	})
}

func (s *Server) writeResponse(err error, r *http.Request, data interface{}) error {
	if err != nil {
		return err
	}
	api.GetApiContext(r).Write(data)
	return nil
}

func (s *Server) deobfuscate(r *http.Request, typeName string, id string) string {
	return strings.TrimPrefix(id, getObfuscator(typeName))
}

func getObfuscator(typeName string) string {
	obfuscator := "1"
	return obfuscator + typeName[0:1]
}

func (s *Server) obfuscate(r *http.Request, typeName string, id string) string {
	if id == "" {
		return ""
	}
	return getObfuscator(typeName) + id
}

func (s *Server) getClient(r *http.Request) (*client.RancherClient, error) {
	return client.NewRancherClient(&client.ClientOpts{
		Url: "http://localhost:8080/v1/projects/1a5/schemas",
	})
}

func (s *Server) parseInputParameters(r *http.Request, obj interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&obj)

	return nil
}

func (s *Server) obfuscateGenericID(r *http.Request, typeName string, id interface{}) string {
	if i, ok := id.(float64); ok {
		str := strconv.FormatFloat(i, 'f', -1, 64)
		return s.obfuscate(r, typeName, str)
	}
	return ""
}
