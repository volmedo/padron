package main

import (
	"context"
	"fmt"
	"net/url"
	"time"

	spaceblobcap "github.com/alanshaw/1up-service/pkg/capabilities/space/blob"
	"github.com/alanshaw/ucantone/client"
	"github.com/alanshaw/ucantone/did"
	"github.com/alanshaw/ucantone/execution"
	"github.com/alanshaw/ucantone/ipld/datamodel"
	"github.com/alanshaw/ucantone/principal/ed25519"
	"github.com/alanshaw/ucantone/result"
	"github.com/alanshaw/ucantone/ucan"
	"github.com/alanshaw/ucantone/ucan/delegation"
	"github.com/alanshaw/ucantone/ucan/invocation"
	mh "github.com/multiformats/go-multihash"

	blobcap "github.com/alanshaw/libracha/capabilities/blob"
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

	spaceToAliceDel, err := delegation.Delegate(
		space,
		alice,
		space,
		"/ucan/*",
		delegation.WithExpiration(ucan.UTCUnixTimestamp(time.Now().Add(1*time.Minute).Unix())), // 1 min
	)
	if err != nil {
		panic(err)
	}

	spaceBlobAddInv, err := spaceblobcap.Add.Invoke(
		alice,
		space,
		&spaceblobcap.AddArguments{
			Blob: spaceblobcap.Blob{
				Digest: digest,
				Size:   12345,
			},
		},
		invocation.WithAudience(service),
		invocation.WithProofs(spaceToAliceDel.Link()),
	)
	if err != nil {
		panic(err)
	}

	inv, err := blobcap.Allocate.Invoke(
		alice,
		space,
		&blobcap.AllocateArguments{
			Blob: blobcap.Blob{
				Digest: digest,
				Size:   12345,
			},
			Cause: spaceBlobAddInv.Link(),
		},
		invocation.WithAudience(service),
		invocation.WithProofs(spaceToAliceDel.Link()),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("spaceToAliceDel:", spaceToAliceDel.Link())
	fmt.Println("spaceBlobAddInv:", spaceBlobAddInv.Link())
	fmt.Println("inv:", inv.Link())

	svcURL, err := url.Parse(serviceURL)
	if err != nil {
		panic(err)
	}

	client, err := client.NewHTTP(svcURL)
	if err != nil {
		panic(err)
	}

	req := execution.NewRequest(
		context.Background(),
		inv,
		execution.WithDelegations(spaceToAliceDel),
		execution.WithInvocations(spaceBlobAddInv),
	)
	res, err := client.Execute(req)
	if err != nil {
		panic(err)
	}

	o, x := result.Unwrap(res.Result())
	if x != nil {
		fmt.Printf("Invocation failed: %v\n\n", x)
		return
	}

	ok := blobcap.AllocateOK{}
	if err := datamodel.Rebind(datamodel.NewAny(o), &ok); err != nil {
		panic(err)
	}

	fmt.Println("/blob/allocate response:")
	fmt.Printf("  Size: %d\n", ok.Size)
	if ok.Address != nil {
		fmt.Printf("  Upload URL: %s\n", ok.Address.URL.URL().String())
		fmt.Printf("  Headers:    %v\n", ok.Address.Headers)
		fmt.Printf("  Expires:    %s\n", ok.Address.Expires.Time().Format(time.RFC3339))
	}
	fmt.Println()

	fmt.Println("/blob/allocate successful!")
}
