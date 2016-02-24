package model

import (
	//"github.com/Sirupsen/logrus"
	"github.com/rancher/go-rancher/client"
)

type ContainerCommon struct {
	Common
	AllocationState    string                      `json:"allocationState" yaml:"allocation_state"`
	Build              *client.DockerBuild         `json:"build" yaml:"build"`
	CapAdd             []string                    `json:"capAdd" yaml:"cap_add"`
	CapDrop            []string                    `json:"capDrop" yaml:"cap_drop"`
	Command            []string                    `json:"command" yaml:"command"`
	ContainerLinks     map[string]interface{}      `json:"containerLinks" yaml:"container_links"`
	CPUSet             string                      `json:"cpuSet" yaml:"cpu_set"`
	CPUShares          int64                       `json:"cpuShares" yaml:"cpu_shares"`
	CreateIndex        int64                       `json:"createIndex" yaml:"create_index"`
	DataVolumeMounts   map[string]interface{}      `json:"dataVolumeMounts" yaml:"data_volume_mounts" schema:",create=true"`
	DataVolumes        []string                    `json:"dataVolumes" yaml:"data_volumes"`
	DataVolumesFrom    []string                    `json:"dataVolumesFrom" yaml:"data_volumes_from"`
	DeploymentUnitUUID string                      `json:"deploymentUnitUuid" yaml:"deployment_unit_uuid"`
	Devices            []string                    `json:"devices" yaml:"devices"`
	DNS                []string                    `json:"dns" yaml:"dns"`
	DNSSearch          []string                    `json:"dnsSearch" yaml:"dns_search"`
	DomainName         string                      `json:"domainName" yaml:"domain_name"`
	EntryPoint         []string                    `json:"entryPoint" yaml:"entry_point"`
	Environment        map[string]interface{}      `json:"environment" yaml:"environment"`
	Expose             []string                    `json:"expose" yaml:"expose"`
	ExternalID         string                      `json:"externalId" yaml:"external_id"`
	ExtraHosts         []string                    `json:"extraHosts" yaml:"extra_hosts"`
	FirstRunning       string                      `json:"firstRunning" yaml:"first_running"`
	HealthCheck        *client.InstanceHealthCheck `json:"healthCheck" yaml:"health_check" schema:",create=true,type=instanceHealthCheck"`
	HealthState        string                      `json:"healthState" yaml:"health_state"`
	Hostname           string                      `json:"hostname" yaml:"hostname"`
	Labels             map[string]interface{}      `json:"labels" yaml:"labels"`
	Memory             int64                       `json:"memory" yaml:"memory" `
	NativeContainer    bool                        `json:"nativeContainer" yaml:"native_container"`
	NetworkMode        string                      `json:"networkMode" yaml:"network_mode"`
	PidMode            string                      `json:"pidMode" yaml:"pid_mode"`
	Ports              []string                    `json:"ports" yaml:"ports"`
	Privileged         bool                        `json:"privileged" yaml:"privileged"`
	PublishAllPorts    bool                        `json:"publishAllPorts" yaml:"publish_all_ports"`
	ReadOnly           bool                        `json:"readOnly" yaml:"read_only"`
	//RequestedHostId    []interface{}          `json:"requestedHostId" yaml:"requested_host_id"`
	RequestedIPAddress string                `json:"requestedIpAddress" yaml:"requested_ip_address"`
	RestartPolicy      *client.RestartPolicy `json:"restartPolicy" yaml:"restart_policy" schema:",create=true,type=restartPolicy"`
	SecurityOpt        []string              `json:"securityOpt" yaml:"security_opt"`
	StartCount         int64                 `json:"startCount" yaml:"start_count"`
	StartOnCreate      bool                  `json:"startOnCreate" yaml:"start_on_create" schema:",default=true"`
	StdinOpen          bool                  `json:"stdinOpen" yaml:"stdin_open"`
	Token              string                `json:"token" yaml:"token"`
	Tty                bool                  `json:"tty" yaml:"tty"`
	User               string                `json:"user" yaml:"user"`
	VolumeDriver       string                `json:"volumeDriver" yaml:"volume_driver"`
}

type ContainerV2 struct {
	client.Resource
	ContainerCommon
	MemSwap   int64             `json:"memSwap" yaml:"mem_swap"`
	Image     string            `json:"image" yaml:"image" schema:",create=true,update=true,nullable=true"`
	WorkDir   string            `json:"workDir" yaml:"work_dir"`
	Logging   *client.LogConfig `json:"logging" yaml:"logging"`
	Revision  string            `json:"revision" yaml:"revision"`
	IPAddress string            `json:"ipAddress" yaml:"ip_address"`
}

type ContainerV1 struct {
	client.Resource
	ContainerCommon
	MemorySwap       int64             `json:"memorySwap" yaml:"memory_swap"`
	ImageUUID        string            `json:"imageUuid" yaml:"image_uuid"`
	WorkingDir       string            `json:"workingDir" yaml:"working_dir"`
	LogConfig        *client.LogConfig `json:"logConfig" yaml:"log_config"`
	Version          string            `json:"version" yaml:"version"`
	PrimaryIPAddress string            `json:"primaryIpAddress" yaml:"primary_ip_address"`
}

type ContainerList struct {
	client.Collection
	Data []ContainerV2 `json:"data"`
}

func getContainerSchema(schemas *client.Schemas) {
	container := AddType(schemas, "container", ContainerV2{})
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
