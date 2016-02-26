package model

import (
	"github.com/rancher/go-rancher/client"
)

func NewSchema() *client.Schemas {
	schemas := &client.Schemas{}

	apiVersion := schemas.AddType("apiVersion", client.Resource{})
	apiVersion.CollectionMethods = []string{}
	schemas.AddType("schema", client.Schema{})
	schemas.AddType("service", Service{})
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

	return schemas
}
