package datamodel

import (
	"errors"
	"io"
	"net/url"

	cbg "github.com/whyrusleeping/cbor-gen"
)

var _ cbg.CBORMarshaler = (*CborURL)(nil)
var _ cbg.CBORUnmarshaler = (*CborURL)(nil)

type CborURL url.URL

func (cu CborURL) URL() *url.URL {
	u := url.URL(cu)
	return &u
}

func (cu CborURL) MarshalCBOR(w io.Writer) error {
	urlStr := cu.URL().String()

	if len(urlStr) > 8192 {
		return errors.New("value in field cu.URL was too long")
	}

	cw := cbg.NewCborWriter(w)

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len(urlStr))); err != nil {
		return err
	}

	if _, err := cw.WriteString(urlStr); err != nil {
		return err
	}

	return nil
}

func (cu *CborURL) UnmarshalCBOR(r io.Reader) error {
	cr := cbg.NewCborReader(r)
	sval, err := cbg.ReadStringWithMax(cr, 8192)
	if err != nil {
		return err
	}

	parsed, err := url.Parse(sval)
	if err != nil {
		return err
	}

	*(*url.URL)(cu) = *parsed

	return nil
}
