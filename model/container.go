package model

import (
	//"github.com/Sirupsen/logrus"
	"github.com/rancher/go-rancher/client"
)

type ContainerCommon struct {
	Common
	AllocationState    string                      `json:"allocationState" yaml:"allocation_state"`
	Build              *client.DockerBuild         `json:"build" yaml:"build" schema:"create=true,type=dockerBuild"`
	CapAdd             []string                    `json:"capAdd" yaml:"cap_add" schema:"create=true"`
	CapDrop            []string                    `json:"capDrop" yaml:"cap_drop" schema:"create=true"`
	Command            []string                    `json:"command" yaml:"command" schema:"create=true"schema:"create=true"`
	ContainerLinks     map[string]interface{}      `json:"containerLinks" yaml:"container_links"`
	CPUSet             string                      `json:"cpuSet" yaml:"cpu_set"`
	CPUShares          int64                       `json:"cpuShares" yaml:"cpu_shares" schema:"create=true"`
	CreateIndex        int64                       `json:"createIndex" yaml:"create_index"`
	DataVolumeMounts   map[string]interface{}      `json:"dataVolumeMounts" yaml:"data_volume_mounts" schema:"create=true"`
	DataVolumes        []string                    `json:"dataVolumes" yaml:"data_volumes" schema:"create=true"`
	DataVolumesFrom    []interface{}               `json:"dataVolumesFrom" yaml:"data_volumes_from" schema:"create=true,type=array[reference[instance]]"`
	DeploymentUnitUUID string                      `json:"deploymentUnitUuid" yaml:"deployment_unit_uuid"`
	Devices            []string                    `json:"devices" yaml:"devices" schema:"create=true"`
	DNS                []string                    `json:"dns" yaml:"dns" schema:"create=true"`
	DNSSearch          []string                    `json:"dnsSearch" yaml:"dns_search" schema:"create=true"`
	DomainName         string                      `json:"domainName" yaml:"domain_name" schema:"create=true"`
	EntryPoint         []string                    `json:"entryPoint" yaml:"entry_point" schema:"create=true"`
	Environment        map[string]interface{}      `json:"environment" yaml:"environment" schema:"create=true"`
	Expose             []string                    `json:"expose" yaml:"expose" schema:"create=true"`
	ExternalID         string                      `json:"externalId" yaml:"external_id"`
	ExtraHosts         []string                    `json:"extraHosts" yaml:"extra_hosts"`
	FirstRunning       string                      `json:"firstRunning" yaml:"first_running"`
	HealthCheck        *client.InstanceHealthCheck `json:"healthCheck" yaml:"health_check" schema:"create=true,type=instanceHealthCheck,nullable=true"`
	HealthState        string                      `json:"healthState" yaml:"health_state"`
	Hostname           string                      `json:"hostname" yaml:"hostname" schema:"create=true"`
	Labels             map[string]interface{}      `json:"labels" yaml:"labels" schema:"create=true"`
	Memory             int64                       `json:"memory" yaml:"memory" schema:"create=true,type=int"`
	NetworkMode        string                      `json:"networkMode" yaml:"network_mode"`
	PidMode            string                      `json:"pidMode" yaml:"pid_mode" schema:"create=true"`
	Ports              []string                    `json:"ports" yaml:"ports" schema:"create=true"`
	Privileged         bool                        `json:"privileged" yaml:"privileged" schema:"create=true"`
	PublishAllPorts    bool                        `json:"publishAllPorts" yaml:"publish_all_ports"`
	ReadOnly           bool                        `json:"readOnly" yaml:"read_only" schema:"create=true"`
	RequestedIPAddress string                      `json:"requestedIpAddress" yaml:"requested_ip_address"`
	RestartPolicy      *client.RestartPolicy       `json:"restartPolicy" yaml:"restart_policy" schema:"create=true,type=restartPolicy,nullable=true"`
	SecurityOpt        []string                    `json:"securityOpt" yaml:"security_opt" schema:"create=true"`
	StartCount         int64                       `json:"startCount" yaml:"start_count" `
	StartOnCreate      bool                        `json:"startOnCreate" yaml:"start_on_create" schema:"create=true"`
	StdinOpen          bool                        `json:"stdinOpen" yaml:"stdin_open" schema:"create=true"`
	Token              string                      `json:"token" yaml:"token"`
	Tty                bool                        `json:"tty" yaml:"tty" schema:"create=true"`
	User               string                      `json:"user" yaml:"user" schema:"create=true"`
	VolumeDriver       string                      `json:"volumeDriver" yaml:"volume_driver" schema:"create=true"`
	RequestedHostID    interface{}                 `json:"requestedHostId" yaml:"requested_host_id" schema:"create=true,type=reference[host]"`
}

type Container struct {
	client.Resource
	ContainerCommon
	MemSwap         int64             `json:"memSwap" yaml:"mem_swap" schema:"create=true,type=int"`
	Image           string            `json:"image" yaml:"image" schema:"create=true,update=true,nullable=true"`
	WorkDir         string            `json:"workDir" yaml:"work_dir"`
	Logging         *client.LogConfig `json:"logging" yaml:"logging" schema:"create=true,type=logConfig"`
	Revision        string            `json:"revision" yaml:"revision"`
	IPAddress       string            `json:"ipAddress" yaml:"ip_address"`
	NativeContainer bool              `json:"nativeContainer" yaml:"native_container"`
}

type ContainerDBProxy struct {
	client.Resource
	ContainerCommon
	MemorySwap         int64             `json:"memorySwap" yaml:"memory_swap"`
	ImageUUID          string            `json:"imageUuid" yaml:"image_uuid"`
	WorkingDir         string            `json:"workingDir" yaml:"working_dir" schema:"create=true"`
	LogConfig          *client.LogConfig `json:"logConfig" yaml:"log_config"`
	Version            string            `json:"version" yaml:"version"`
	PrimaryIPAddress   string            `json:"primaryIpAddress" yaml:"primary_ip_address"`
	NativeContainer    []uint8           `json:"nativeContainer" yaml:"native_container"`
	NetworkContainerID interface{}       `json:"networkContainerId" yaml:"network_container_id"`
}

type ContainerList struct {
	client.Collection
	Data []Container `json:"data"`
}

func getContainerSchema(schemas *client.Schemas) {
	container := AddType(schemas, "container", Container{})
	container.ResourceActions = map[string]client.Action{
		"create": client.Action{
			Output: "container",
		},
	}
	container.CollectionMethods = []string{"GET", "POST"}
	container.ResourceMethods = []string{"GET"}
	//TODO populate links with actions
	container.Actions = map[string]string{}
}
