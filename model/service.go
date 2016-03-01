package model

import "github.com/rancher/go-rancher/client"

type ServiceCommon struct {
	client.Resource
	Common
	CreateIndex            int                     `json:"createIndex"`
	Fqdn                   string                  `json:"fqdn" schema:"create=true"`
	Scale                  int                     `json:"scale" schema:"create=true"`
	AssignServiceIPAddress bool                    `json:"assignServiceIpAddress" schema:"create=true"`
	ExternalID             string                  `json:"externalId"`
	Metadata               map[string]interface{}  `json:"metadata" schema:"create=true,type=map[json]"`
	PublicEndpoints        []client.PublicEndpoint `json:"publicEndpoints" schema:"type=array[publicEndpoint]"`
	Restart                *client.ServiceRestart  `json:"restart" schema:"type=serviceRestart"`

	/*ContainerIds           []ID                `json:"containerIds"`
	HealthState            string              `json:"healthState"`
	HostnameOverride       string              `json:"hostnameOverride"`
	LinkedServiceIds       []ID                `json:"linkedServiceIds"`



	ScalePolicy            interface{}         `json:"scalePolicy"`
	*/
}

type Service struct {
	client.Resource
	ServiceCommon
	StackID            string          `json:"stackId" schema:"create=true,type=string"`
	ServiceIPAddress   string          `json:"serviceIpAddress"`
	LinkSelector       string          `json:"linkSelector" schema:"create=true"`
	ContainerSelector  string          `json:"containerSelector" schema:"create=true"`
	RetainIPAddress    bool            `json:"retainIpAddress" schema:"create=true"`
	ContainerTemplates []*Container    `json:"containerTemplates" schema:"create=true, type=array[container]"`
	Upgrade            *ServiceUpgrade `json:"upgrade" schema:"type=serviceUpgrade"`
}

type ServiceDBProxy struct {
	client.Resource
	ServiceCommon
	EnvironmentID          string                 `json:"environmentId"`
	Vip                    string                 `json:"vip"`
	SelectorLink           string                 `json:"selectorLink"`
	SelectorContainer      string                 `json:"selectorContainer"`
	RetainIP               bool                   `json:"retainIp"`
	LaunchConfig           *ContainerDBProxy      `json:"launchConfig"`
	SecondaryLaunchConfigs []*ContainerDBProxy    `json:"secondaryLaunchConfigs"`
	Upgrade                *ServiceUpgradeDBProxy `json:"upgrade" schema:"type=serviceUpgrade"`
}

type ServiceList struct {
	client.Collection
	Data []Service `json:"data,omitempty"`
}

func getServiceSchema(schemas *client.Schemas) {
	container := AddType(schemas, "service", Service{})
	container.ResourceActions = map[string]client.Action{
		"create": client.Action{
			Output: "service",
		},
	}
	container.CollectionMethods = []string{"GET", "POST"}
	container.ResourceMethods = []string{"GET"}
	//TODO populate links/actions
	container.Actions = map[string]string{}
}
