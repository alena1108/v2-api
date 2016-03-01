package model

type InServiceUpgradeStrategyCommon struct {
	BatchSize int64 `json:"batchSize" yaml:"batch_size" schema:"create=true" schema:"create=true"`

	IntervalMillis int64 `json:"intervalMillis" yaml:"interval_millis" schema:"create=true" schema:"create=true"`

	StartFirst bool `json:"startFirst" yaml:"start_first" schema:"create=true" schema:"create=true"`
}

type InServiceUpgradeStrategy struct {
	InServiceUpgradeStrategyCommon
	ContainerTemplates []*Container `json:"containerTemplates" schema:"create=true, type=array[container]"`

	PreviousContainerTemplates []*Container `json:"previousContainerTemplates" schema:"create=true, type=array[container]"`
}

type InServiceUpgradeStrategyDBProxy struct {
	InServiceUpgradeStrategyCommon

	LaunchConfig *ContainerDBProxy `json:"launchConfig"`

	SecondaryLaunchConfigs []*ContainerDBProxy `json:"secondaryLaunchConfigs"`

	PreviousSecondaryLaunchConfigs []*ContainerDBProxy `json:"previousSecondaryLaunchConfigs" yaml:"previous_secondary_launch_configs" schema:"create=true"`

	PreviousLaunchConfig *ContainerDBProxy `json:"secondaryLaunchConfigs" yaml:"secondary_launch_configs" schema:"create=true"`
}
