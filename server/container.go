package server

import (
	//"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/rancher/go-rancher/client"
	"github.com/rancher/v2-api/model"
	"net/http"
)

type Container struct{}

func (s *Server) getContainersSQL(r *http.Request, id string) string {
	q := `
	  SELECT
	      COALESCE(name, '') as name, id, uuid, state, version,
           COALESCE(first_running, '') as first_running,
            COALESCE(start_count, 0) as start_count,
            native_container, COALESCE(token, "") as token, 
            COALESCE(external_id, "") as external_id,
            COALESCE(deployment_unit_uuid, "") as deployment_unit_uuid,
            COALESCE(hostname, "") as hostname,
             data
	  FROM instance
	  WHERE
	      account_id = :account_id
	      AND removed IS NULL
	      AND kind = 'container'`

	if id != "" {
		q += " AND id = :id"
	}

	return q
}

func (s *Server) ContainerCreate(rw http.ResponseWriter, r *http.Request) error {
	rancherClient, err := s.getClient(r)
	if err != nil {
		return err
	}
	v2 := &model.ContainerV2{}
	if err := s.parseInputParameters(r, v2); err != nil {
		return err
	}

	v1, err := FromV2(v2)
	if err != nil {
		return err
	}

	c := &client.Container{}
	if err := convertObject(v1, c); err != nil {
		return err
	}

	container, err := rancherClient.Container.Create(c)

	if err != nil {
		return err
	}

	return s.getContainer(rw, r, container.Id)
}

func (s *Server) ContainerByID(rw http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	return s.getContainer(rw, r, vars["id"])
}

func (s *Server) ContainerList(rw http.ResponseWriter, r *http.Request) error {
	return s.getContainer(rw, r, "")
}

func (s *Server) getContainer(rw http.ResponseWriter, r *http.Request, id string) error {
	resourceType := "container"

	id = s.deobfuscate(r, resourceType, id)

	rows, err := s.namedQuery(s.getContainersSQL(r, id), map[string]interface{}{
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

		obj := &model.ContainerV1{}
		obj.Type = resourceType

		var data string

		if err := rows.Scan(&obj.Name, &obj.Id, &obj.UUID, &obj.State,
			&obj.Version, &obj.FirstRunning, &obj.StartCount,
			&obj.NativeContainer, &obj.Token, &obj.ExternalID,
			&obj.DeploymentUnitUUID,
			&obj.Hostname, &data); err != nil {
			return err
		}

		// Obfuscate Ids
		obj.Id = s.obfuscate(r, resourceType, obj.Id)
		if err = s.parseData(data, obj); err != nil {
			return err
		}

		objV2, err := ToV2(obj)
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

func ToV2(v1 *model.ContainerV1) (*model.ContainerV2, error) {

	common := v1.ContainerCommon
	common.Transitioning = model.GetTransitioning(common.State, common.Transitioning)
	nativeContainer := false
	if v1.NativeContainer[0] == 1 {
		nativeContainer = true
	}

	return &model.ContainerV2{
		Resource:        v1.Resource,
		ContainerCommon: common,
		Image:           v1.ImageUUID,
		WorkDir:         v1.WorkingDir,
		Logging:         v1.LogConfig,
		MemSwap:         v1.MemorySwap,
		Revision:        v1.Version,
		IPAddress:       v1.PrimaryIPAddress,
		NativeContainer: nativeContainer,
	}, nil
}

func FromV2(v2 *model.ContainerV2) (*client.Container, error) {
	v1 := &client.Container{}

	if err := convertObject(v2, v1); err != nil {
		return nil, err
	}

	v1.ImageUuid = v2.Image
	v1.WorkingDir = v2.WorkDir
	v1.MemorySwap = v2.MemSwap

	if v2.StartOnCreate == false {
		v1.StartOnCreate = false
	}

	return v1, nil
}
