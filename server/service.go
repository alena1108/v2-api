package server

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/rancher/go-rancher/client"
	"github.com/rancher/v2-api/model"
)

func (s *Server) ServiceByID(rw http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	return s.getService(rw, r, vars["id"])
}

func (s *Server) ServiceList(rw http.ResponseWriter, r *http.Request) error {
	return s.getService(rw, r, "")
}

func (s *Server) getService(rw http.ResponseWriter, r *http.Request, id string) error {
	resourceType := "service"

	id = s.deobfuscate(r, "service", id)

	rows, err := s.namedQuery(s.getServicesSQL(r, id), map[string]interface{}{
		"account_id": s.getAccountID(r),
		"id":         id,
	})
	if err != nil {
		return err
	}
	defer rows.Close()

	response := &client.GenericCollection{
		Collection: client.Collection{
			Type:         "collection",
			ResourceType: resourceType,
		},
	}

	for rows.Next() {

		obj := &model.ServiceDBProxy{}
		obj.Type = resourceType

		var data string

		if err := rows.Scan(&obj.Name, &obj.Id, &obj.UUID, &obj.EnvironmentID,
			&obj.State, &obj.CreateIndex, &obj.Vip, &obj.ExternalID,
			&obj.SelectorLink, &obj.SelectorContainer, &data); err != nil {
			return err
		}

		// Obfuscate Ids
		obj.Id = s.obfuscate(r, "service", obj.Id)
		if err = s.parseData(data, obj); err != nil {
			return err
		}

		objV2, err := s.ServiceDBProxyToV2(obj, r)
		if err != nil {
			return err
		}
		if id != "" {
			return s.writeResponse(rows.Err(), r, objV2)
		}
		response.Data = append(response.Data, objV2)
	}

	return s.writeResponse(rows.Err(), r, response)
}

func (s *Server) getServicesSQL(r *http.Request, id string) string {
	q := `
		SELECT
			name, id, uuid, environment_id, state, 
            COALESCE(create_index,0) as create_index,
            COALESCE(vip,"") as vip,
            COALESCE(external_id,"") as external_id,
            COALESCE(selector_link,"") as selector_link,
            COALESCE(selector_container,"") as selector_container,

            data 
		FROM service
		WHERE
			account_id = :account_id
			AND removed IS NULL`

	if id != "" {
		q += " AND id = :id"
	}

	return q
}

func (s *Server) ServiceCreate(rw http.ResponseWriter, r *http.Request) error {
	rancherClient, err := s.getClient(r)
	if err != nil {
		return err
	}
	v2 := &model.Service{}
	if err := s.parseInputParameters(r, v2); err != nil {
		return err
	}

	v1, err := s.ServiceFromV2ToV1(v2, r)
	if err != nil {
		return err
	}

	service, err := rancherClient.Service.Create(v1)

	if err != nil {
		return err
	}

	return s.getService(rw, r, service.Id)
}

func (s *Server) ServiceDBProxyToV2(db *model.ServiceDBProxy, r *http.Request) (*model.Service, error) {
	common := db.ServiceCommon
	common.Transitioning = model.GetTransitioning(common.State, common.Transitioning)

	if common.PublicEndpoints != nil {
		for key, val := range common.PublicEndpoints {
			val.InstanceId = s.obfuscateGenericID(r, "instance", val.InstanceId)
			val.HostId = s.obfuscateGenericID(r, "host", val.HostId)
			common.PublicEndpoints[key] = val
		}
	}

	return &model.Service{
		Resource:          db.Resource,
		ServiceCommon:     common,
		StackID:           s.obfuscateGenericID(r, "stack", db.EnvironmentID),
		ServiceIPAddress:  db.Vip,
		LinkSelector:      db.SelectorLink,
		ContainerSelector: db.SelectorContainer,
		RetainIPAddress:   db.RetainIP,
	}, nil
}

func (s *Server) ServiceFromV2ToV1(v2 *model.Service, r *http.Request) (*client.Service, error) {
	v1 := &client.Service{}

	if err := convertObject(v2, v1); err != nil {
		return nil, err
	}

	v1.EnvironmentId = v2.StackID
	v1.SelectorLink = v2.LinkSelector
	v1.SelectorContainer = v2.ContainerSelector
	v1.RetainIp = v2.RetainIPAddress
	logrus.Infof("service vip %v", v2.AssignServiceIPAddress)

	return v1, nil
}
