package blob

import (
	"github.com/alanshaw/ucantone/ucan/delegation/policy"
	"github.com/alanshaw/ucantone/validator/bindcap"
	"github.com/alanshaw/ucantone/validator/capability"
	cbg "github.com/whyrusleeping/cbor-gen"

	bdm "github.com/volmedo/padron/pkg/capabilities/blob/datamodel"
)

type (
	AllocateArguments = bdm.AllocateArgumentsModel
	AllocateOK        = bdm.AllocateOKModel
	BlobAddress       = bdm.BlobAddressModel
	AcceptArguments   = bdm.AcceptArgumentsModel
	AcceptOK          = bdm.AcceptOKModel
	Blob              = bdm.BlobModel
	CborURL           = bdm.CborURL
	CborTime          = cbg.CborTime
)

const AllocateCommand = "/blob/allocate"

var Allocate, _ = bindcap.New[*AllocateArguments](
	AllocateCommand,
	capability.WithPolicyBuilder(
		policy.GreaterThan(".blob.size", 0),
		policy.LessThanOrEqual(".blob.size", 268_435_456),
	),
)

const AcceptCommand = "/blob/accept"

var Accept, _ = bindcap.New[*AcceptArguments](AcceptCommand)
