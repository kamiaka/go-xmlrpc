package xmlrpc

import (
	"encoding/xml"
	"fmt"
	"io"
)

// Response ...
type Response interface {
	Value() (*Value, error)
	Write(io.Writer) error
}

type response struct {
	XMLName xml.Name `xml:"methodResponse"`
	Param   *param   `xml:"params>param,omitempty"`
	Fault   *Value   `xml:"fault>value,omitempty"`
}

// NewResponse ...
func NewResponse(val *Value) Response {
	var p *param
	if val != nil {
		p = &param{
			Value: val,
		}
	}
	return &response{
		Param: p,
	}
}

// FaultResponse represents Fault
func FaultResponse(code int, str string) Response {
	f := &Fault{
		code: code,
		str:  str,
	}
	return &response{
		Fault: f.Value(),
	}
}

func (r *response) Value() (*Value, error) {
	if r.Fault != nil {
		f := r.Fault.Struct()
		return nil, &Fault{
			code: f["faultCode"].(int),
			str:  f["faultString"].(string),
		}
	}
	if r.Param == nil {
		return nil, nil
	}
	return r.Param.Value, nil
}
func (r *response) Write(w io.Writer) error {
	_, err := w.Write(header)
	if err != nil {
		return err
	}
	return xml.NewEncoder(w).Encode(r)
}

// Fault error
type Fault struct {
	code int
	str  string
}

func (f *Fault) Code() int {
	return f.code
}
func (f *Fault) String() string {
	return f.str
}

func (f *Fault) Error() string {
	return fmt.Sprintf("Fault: %s (code: %d)", f.str, f.code)
}

func (f *Fault) Value() *Value {
	return Struct(
		Member("faultCode", Int(f.code)),
		Member("faultString", String(f.str)),
	)
}

// ParseResponse ...
func ParseResponse(r io.Reader) (*Value, error) {
	var resp response
	err := xml.NewDecoder(r).Decode(&resp)
	if err != nil {
		return nil, err
	}
	return resp.Value()
}
