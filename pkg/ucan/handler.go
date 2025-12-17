package ucan

import (
	"github.com/alanshaw/ucantone/execution"
	"github.com/alanshaw/ucantone/validator"
)

type Handler struct {
	Capability validator.Capability
	Handler    execution.HandlerFunc
}
