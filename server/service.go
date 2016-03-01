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

	v1, err := s.ServiceV2ToV1(v2, r)
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

	templates, err := s.ContainerTemplatesDBProxyToV2(db.LaunchConfig, db.SecondaryLaunchConfigs, r)
	if err != nil {
		return nil, err
	}

	upgrade, err := s.UpgradeDBProxyToV2(db, r)
	if err != nil {
		return nil, err
	}

	return &model.Service{
		Resource:           db.Resource,
		ServiceCommon:      common,
		StackID:            s.obfuscateGenericID(r, "stack", db.EnvironmentID),
		ServiceIPAddress:   db.Vip,
		LinkSelector:       db.SelectorLink,
		ContainerSelector:  db.SelectorContainer,
		RetainIPAddress:    db.RetainIP,
		ContainerTemplates: templates,
		Upgrade:            upgrade,
	}, nil
}

func (s *Server) UpgradeDBProxyToV2(db *model.ServiceDBProxy, r *http.Request) (*model.ServiceUpgrade, error) {
	upgrade := &model.ServiceUpgrade{}

	if db.Upgrade != nil && db.Upgrade.InServiceStrategy != nil {
		dbProxy := &model.InServiceUpgradeStrategyDBProxy{}
		if err := convertObject(upgrade.InServiceStrategy, dbProxy); err != nil {
			return nil, err
		}
		conatinerTemplates, err := s.ContainerTemplatesDBProxyToV2(db.Upgrade.InServiceStrategy.LaunchConfig, dbProxy.SecondaryLaunchConfigs, r)
		if err != nil {
			return nil, err
		}
		previousContainerTemplates, err := s.ContainerTemplatesDBProxyToV2(dbProxy.PreviousLaunchConfig, dbProxy.PreviousSecondaryLaunchConfigs, r)
		if err != nil {
			return nil, err
		}

		upgrade.InServiceStrategy = &model.InServiceUpgradeStrategy{
			InServiceUpgradeStrategyCommon: dbProxy.InServiceUpgradeStrategyCommon,
			ContainerTemplates:             conatinerTemplates,
			PreviousContainerTemplates:     previousContainerTemplates,
		}
	}
	return upgrade, nil
}

func (s *Server) UpgradeV2ToV1(v2 *model.Service, r *http.Request) (*client.ServiceUpgrade, error) {

	inServiceStrategy := &client.InServiceUpgradeStrategy{}
	toServiceStrategy := &client.ToServiceUpgradeStrategy{}
	if v2.Upgrade != nil {
		if v2.Upgrade.InServiceStrategy != nil {
			if err := convertObject(v2.Upgrade.InServiceStrategy, inServiceStrategy); err != nil {
				return nil, err
			}
			lc, slc, err := s.ContainerTemplatesV2ToV1(v2, v2.Upgrade.InServiceStrategy.ContainerTemplates, r)
			if err != nil {
				return nil, err
			}
			inServiceStrategy.LaunchConfig = lc
			inServiceStrategy.SecondaryLaunchConfigs = slc

		} else if v2.Upgrade.ToServiceStrategy != nil {
			toServiceStrategy = v2.Upgrade.ToServiceStrategy
			toServiceStrategy.ToServiceId = s.obfuscate(r, "service", toServiceStrategy.ToServiceId)
		}
	}

	return &client.ServiceUpgrade{
		ToServiceStrategy: toServiceStrategy,
		InServiceStrategy: inServiceStrategy,
	}, nil
}

func (s *Server) ServiceV2ToV1(v2 *model.Service, r *http.Request) (*client.Service, error) {
	v1 := &client.Service{}

	if err := convertObject(v2, v1); err != nil {
		return nil, err
	}

	v1.EnvironmentId = v2.StackID
	v1.SelectorLink = v2.LinkSelector
	v1.SelectorContainer = v2.ContainerSelector
	v1.RetainIp = v2.RetainIPAddress
	lc, slc, err := s.ContainerTemplatesV2ToV1(v2, v2.ContainerTemplates, r)
	if err != nil {
		return nil, err
	}
	v1.LaunchConfig = lc
	v1.SecondaryLaunchConfigs = slc

	upgrade, err := s.UpgradeV2ToV1(v2, r)
	if err != nil {
		return nil, err
	}
	v1.Upgrade = upgrade
	logrus.Infof("service vip %v", v2.AssignServiceIPAddress)

	return v1, nil
}

func (s *Server) ContainerTemplatesV2ToV1(v2 *model.Service, templates []*model.Container, r *http.Request) (*client.LaunchConfig, []client.SecondaryLaunchConfig, error) {
	var plc *client.LaunchConfig
	var slc []client.SecondaryLaunchConfig
	if templates != nil {
		for _, template := range templates {
			v1, err := s.ContainerV2ToV1(template, r)
			if err != nil {
				return nil, nil, err
			}

			if template.Name == v2.Name {
				lc := &client.LaunchConfig{}
				if err = convertObject(v1, lc); err != nil {
					return nil, nil, err
				}
				plc = lc
			} else {
				lc := client.SecondaryLaunchConfig{}
				if err = convertObject(v1, &lc); err != nil {
					return nil, nil, err
				}
				slc = append(slc, lc)
			}
		}
	}

	return plc, slc, nil
}

func (s *Server) ContainerTemplatesDBProxyToV2(lc *model.ContainerDBProxy, slc []*model.ContainerDBProxy, r *http.Request) ([]*model.Container, error) {
	var templates []*model.Container
	if lc != nil {
		dbC := &model.ContainerDBProxy{}
		if err := convertObject(lc, &dbC); err != nil {
			return nil, err
		}
		template, err := s.ContainerDBProxyToV2(dbC, r)
		if err != nil {
			return nil, err
		}
		templates = append(templates, template)
	}

	if slc != nil {
		for _, sc := range slc {
			dbC := &model.ContainerDBProxy{}
			if err := convertObject(sc, &dbC); err != nil {
				return nil, err
			}
			template, err := s.ContainerDBProxyToV2(dbC, r)
			if err != nil {
				return nil, err
			}
			templates = append(templates, template)
		}
	}

	return templates, nil
}
