package datamodel

import (
	"github.com/alanshaw/ucantone/ucan/promise"
	"github.com/multiformats/go-multihash"
)

type BlobModel struct {
	Digest multihash.Multihash `cborgen:"digest" dagjsongen:"digest"`
	Size   uint64              `cborgen:"size" dagjsongen:"size"`
}

type AddArgumentsModel struct {
	Blob BlobModel `cborgen:"blob" dagjsongen:"blob"`
}

type AddOKModel struct {
	Site promise.AwaitOK `cborgen:"site" dagjsongen:"site"`
}
