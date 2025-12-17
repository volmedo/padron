package main

import (
	bdm "github.com/volmedo/padron/pkg/capabilities/blob/datamodel"
	cbg "github.com/whyrusleeping/cbor-gen"
)

func main() {
	if err := cbg.WriteMapEncodersToFile("../cbor_gen.go", "datamodel",
		bdm.AllocateArgumentsModel{},
		bdm.BlobModel{},
		bdm.AllocateOKModel{},
		bdm.BlobAddressModel{},
		bdm.AcceptArgumentsModel{},
		bdm.AcceptOKModel{},
	); err != nil {
		panic(err)
	}
}
