package xmlrpc

import (
	"encoding/xml"
	"io"
)

// Request ...
type Request interface {
	Method() string
	Params() []*Value
	Write(io.Writer) error
}

type request struct {
	XMLName    xml.Name  `xml:"methodCall"`
	MethodName string    `xml:"methodName"`
	Parameters *[]*param `xml:"params>param,omitempty"`
}

// NewRequest ...
func NewRequest(method string, args ...*Value) (Request, error) {
	var params *[]*param
	if len(args) > 0 {
		ps := make([]*param, len(args))
		for i, a := range args {
			ps[i] = &param{
				Value: a,
			}
		}
		params = &ps
	}

	return &request{
		MethodName: method,
		Parameters: params,
	}, nil
}

// ParseRequest ...
func ParseRequest(r io.Reader) (Request, error) {
	var req request
	err := xml.NewDecoder(r).Decode(&req)
	if err != nil {
		return nil, err
	}
	return &req, nil
}

func (r *request) Method() string {
	return r.MethodName
}

func (r *request) Params() []*Value {
	if r.Parameters == nil {
		return []*Value{}
	}
	values := make([]*Value, len(*r.Parameters))
	for i, p := range *r.Parameters {
		values[i] = p.Value
	}
	return values
}

var header = []byte(`<?xml version="1.0"?>`)

func (r *request) Write(w io.Writer) error {
	_, err := w.Write(header)
	if err != nil {
		return err
	}
	return xml.NewEncoder(w).Encode(r)
}
