package blob

import (
	"github.com/alanshaw/ucantone/ucan/delegation/policy"
	"github.com/alanshaw/ucantone/validator/bindcap"
	"github.com/alanshaw/ucantone/validator/capability"
	bdm "github.com/volmedo/padron/pkg/capabilities/space/blob/datamodel"
)

const AddCommand = "/space/blob/add"

type (
	AddArguments = bdm.AddArgumentsModel
	Blob         = bdm.BlobModel
	AddOK        = bdm.AddOKModel
)

var Add, _ = bindcap.New[*AddArguments](
	AddCommand,
	capability.WithPolicyBuilder(
		policy.GreaterThan(".blob.size", 0),
		policy.LessThanOrEqual(".blob.size", 268_435_456),
	),
)
