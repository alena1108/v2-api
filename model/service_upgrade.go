package model

import "github.com/rancher/go-rancher/client"

type ServiceUpgradeCommon struct {
	ToServiceStrategy *client.ToServiceUpgradeStrategy `json:"toServiceStrategy,omitempty" yaml:"to_service_strategy,omitempty"`
}

type ServiceUpgrade struct {
	ServiceUpgradeCommon
	InServiceStrategy *InServiceUpgradeStrategy `json:"inServiceStrategy,omitempty" yaml:"in_service_strategy,omitempty"`
}

type ServiceUpgradeDBProxy struct {
	ServiceUpgradeCommon
	InServiceStrategy *InServiceUpgradeStrategyDBProxy `json:"inServiceStrategy,omitempty" yaml:"in_service_strategy,omitempty"`
}
