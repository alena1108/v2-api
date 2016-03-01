package server

import (
	//"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/rancher/go-rancher/client"
	"github.com/rancher/v2-api/model"
	"net/http"
	//"reflect"
	"strings"
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
            COALESCE(allocation_state, "") as allocation_state,
            COALESCE(network_container_id, 0) as network_container_id,
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
	v2 := &model.Container{}
	if err := s.parseInputParameters(r, v2); err != nil {
		return err
	}

	v1, err := s.ContainerV2ToV1(v2, r)
	if err != nil {
		return err
	}

	container, err := rancherClient.Container.Create(v1)

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

	id = s.deobfuscate(r, "instance", id)

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

		obj := &model.ContainerDBProxy{}
		obj.Type = resourceType

		var data string

		if err := rows.Scan(&obj.Name, &obj.Id, &obj.UUID, &obj.State,
			&obj.Version, &obj.FirstRunning, &obj.StartCount,
			&obj.NativeContainer, &obj.Token, &obj.ExternalID,
			&obj.DeploymentUnitUUID,
			&obj.Hostname,
			&obj.AllocationState,
			&obj.NetworkContainerID, &data); err != nil {
			return err
		}

		// Obfuscate Ids
		obj.Id = s.obfuscate(r, "instance", obj.Id)
		if err = s.parseData(data, obj); err != nil {
			return err
		}

		objV2, err := s.ContainerDBProxyToV2(obj, r)
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

func (s *Server) ContainerDBProxyToV2(db *model.ContainerDBProxy, r *http.Request) (*model.Container, error) {
	common := db.ContainerCommon
	common.Transitioning = model.GetTransitioning(common.State, common.Transitioning)
	nativeContainer := false
	if db.NativeContainer[0] == 1 {
		nativeContainer = true
	}

	if db.RequestedHostID != "" {
		common.RequestedHostID = s.obfuscateGenericID(r, "host", db.RequestedHostID)
	}

	if db.DataVolumesFrom != nil {
		var dataV []interface{}
		for _, v := range db.DataVolumesFrom {
			dataV = append(dataV, s.obfuscateGenericID(r, "instance", v))
		}
		common.DataVolumesFrom = dataV
	}

	if db.NetworkMode == "container" {
		cID := s.obfuscateGenericID(r, "instance", db.NetworkContainerID)
		common.NetworkMode = db.NetworkMode + ":" + cID
	}

	return &model.Container{
		Resource:        db.Resource,
		ContainerCommon: common,
		Image:           db.ImageUUID,
		WorkDir:         db.WorkingDir,
		Logging:         db.LogConfig,
		MemSwap:         db.MemorySwap,
		Revision:        db.Version,
		IPAddress:       db.PrimaryIPAddress,
		NativeContainer: nativeContainer,
	}, nil
}

func (s *Server) ContainerV2ToV1(v2 *model.Container, r *http.Request) (*client.Container, error) {
	v1 := &client.Container{}

	if err := convertObject(v2, v1); err != nil {
		return nil, err
	}

	if strings.HasPrefix(v2.NetworkMode, "container") {
		splitted := strings.SplitN(v2.NetworkMode, ":", 2)
		v1.NetworkContainerId = splitted[1]
		v1.NetworkMode = "container"
	}

	v1.ImageUuid = v2.Image
	v1.WorkingDir = v2.WorkDir
	v1.MemorySwap = v2.MemSwap
	v1.LogConfig = v2.Logging

	return v1, nil
}
