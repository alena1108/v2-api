package model

import "time"

type ID string

type Common struct {
	Name        string     `json:"name" schema:",create=true,update=true"`
	Description string     `json:"description" schema:",create=true,update=true,nullable=true"`
	State       string     `json:"state"`
	UUID        string     `json:"uuid"`
	Created     *time.Time `json:"created" schema:",type=date"`
	Removed     *time.Time `json:"removed"`

	Transitioning        string `json:"transitioning"`
	TransitioningMessage string `json:"transitioningMessage"`
}
