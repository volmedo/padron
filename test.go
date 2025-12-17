package main

import (
	"context"
	"fmt"
	"net/url"

	spaceblobcap "github.com/alanshaw/1up-service/pkg/capabilities/space/blob"
	"github.com/alanshaw/ucantone/client"
	"github.com/alanshaw/ucantone/did"
	"github.com/alanshaw/ucantone/execution"
	"github.com/alanshaw/ucantone/ipld"
	"github.com/alanshaw/ucantone/ipld/datamodel"
	"github.com/alanshaw/ucantone/principal/ed25519"
	"github.com/alanshaw/ucantone/result"
	"github.com/alanshaw/ucantone/ucan/invocation"
	mh "github.com/multiformats/go-multihash"

	blobcap "github.com/volmedo/padron/pkg/capabilities/blob"
)

const (
	serviceID  = "did:key:z6MkoznxjrCCpQwFAD1BJP2uFiAccpo6cPLHDGqPtdjahajj"
	serviceURL = "http://localhost:3000"
)

func main() {
	alice, err := ed25519.Generate()
	if err != nil {
		panic(err)
	}

	space, err := ed25519.Generate()
	if err != nil {
		panic(err)
	}

	service, err := did.Parse(serviceID)
	if err != nil {
		panic(err)
	}

	digest, err := mh.Sum([]byte("testing 1, 2, 3"), mh.SHA2_256, -1)
	if err != nil {
		panic(err)
	}

	spaceBlobAddInv, err := spaceblobcap.Add.Invoke(
		alice,
		alice,
		&spaceblobcap.AddArguments{
			Blob: spaceblobcap.Blob{
				Digest: digest,
				Size:   12345,
			},
		},
		invocation.WithAudience(service),
	)
	if err != nil {
		panic(err)
	}

	inv, err := blobcap.Allocate.Invoke(
		alice,
		alice,
		&blobcap.AllocateArguments{
			Space: space.DID(),
			Blob: blobcap.Blob{
				Digest: digest,
				Size:   12345,
			},
			Cause: spaceBlobAddInv.Link(),
		},
		invocation.WithAudience(service),
	)
	if err != nil {
		panic(err)
	}

	url, err := url.Parse(serviceURL)
	if err != nil {
		panic(err)
	}

	client, err := client.NewHTTP(url)
	if err != nil {
		panic(err)
	}

	res, err := client.Execute(execution.NewRequest(context.Background(), inv))
	if err != nil {
		panic(err)
	}

	result.MatchResultR0(
		res.Result(),
		func(o ipld.Any) {
			ok := blobcap.AllocateOK{}
			err := datamodel.Rebind(datamodel.NewAny(o), &ok)
			if err != nil {
				panic(err)
			}
			fmt.Printf("/blob/allocate response: %+v\n\n", ok)
		},
		func(x ipld.Any) {
			fmt.Printf("Invocation failed: %v\n\n", x)
		},
	)
}
