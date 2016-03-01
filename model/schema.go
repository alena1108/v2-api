package model

import (
	"github.com/rancher/go-rancher/client"
)

func NewSchema() *client.Schemas {
	schemas := &client.Schemas{}

	apiVersion := schemas.AddType("apiVersion", client.Resource{})
	apiVersion.CollectionMethods = []string{}
	schemas.AddType("schema", client.Schema{})
	getServiceSchema(schemas)
	getContainerSchema(schemas)

	restartPolicy := schemas.AddType("restartPolicy", client.RestartPolicy{})
	restartPolicy.CollectionMethods = []string{}
	logConfig := schemas.AddType("logConfig", client.LogConfig{})
	logConfig.CollectionMethods = []string{}
	healthCheck := schemas.AddType("instanceHealthCheck", client.InstanceHealthCheck{})
	healthCheck.CollectionMethods = []string{}
	strategy := schemas.AddType("recreateOnQuorumStrategyConfig", client.RecreateOnQuorumStrategyConfig{})
	strategy.CollectionMethods = []string{}
	dockerBuild := schemas.AddType("dockerBuild", client.DockerBuild{})
	dockerBuild.CollectionMethods = []string{}

	restartStrategy := schemas.AddType("rollingRestartStrategy", client.RollingRestartStrategy{})
	restartStrategy.CollectionMethods = []string{}
	restart := schemas.AddType("serviceRestart", client.ServiceRestart{})
	restart.CollectionMethods = []string{}

	upgrade := schemas.AddType("serviceUpgrade", ServiceUpgrade{})
	upgrade.CollectionMethods = []string{}

	toService := schemas.AddType("toServiceUpgradeStrategy", client.InServiceUpgradeStrategy{})
	toService.CollectionMethods = []string{}

	inService := schemas.AddType("inServiceUpgradeStrategy", InServiceUpgradeStrategy{})
	inService.CollectionMethods = []string{}

	return schemas
}
