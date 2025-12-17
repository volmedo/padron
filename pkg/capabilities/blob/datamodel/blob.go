package datamodel

import (
	"github.com/alanshaw/ucantone/did"
	"github.com/alanshaw/ucantone/ucan"
	"github.com/alanshaw/ucantone/ucan/promise"
	"github.com/multiformats/go-multihash"
	cbg "github.com/whyrusleeping/cbor-gen"
)

type AllocateArgumentsModel struct {
	Space did.DID   `cborgen:"space" dagjsongen:"space"`
	Blob  BlobModel `cborgen:"blob" dagjsongen:"blob"`
	Cause ucan.Link `cborgen:"cause" dagjsongen:"cause"`
}

type BlobModel struct {
	Digest multihash.Multihash `cborgen:"digest" dagjsongen:"digest"`
	Size   uint64              `cborgen:"size" dagjsongen:"size"`
}

type AllocateOKModel struct {
	Size    uint64            `cborgen:"size" dagjsongen:"size"`
	Address *BlobAddressModel `cborgen:"address,omitempty" dagjsongen:"address,omitempty"`
}

type BlobAddressModel struct {
	URL     CborURL           `cborgen:"url" dagjsongen:"url"`
	Headers map[string]string `cborgen:"headers" dagjsongen:"headers"`
	Expires cbg.CborTime      `cborgen:"expires" dagjsongen:"expires"`
}

type AcceptArgumentsModel struct {
	Space multihash.Multihash `cborgen:"space" dagjsongen:"space"`
	Blob  BlobModel           `cborgen:"blob" dagjsongen:"blob"`
	Put   promise.AwaitOK     `cborgen:"_put" dagjsongen:"_put"`
}

type AcceptOKModel struct {
	Site ucan.Link `cborgen:"site" dagjsongen:"site"`
}
