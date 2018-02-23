package xmlrpc

import (
	"bytes"
	"reflect"
	"testing"
)

func TestParseResponse(t *testing.T) {
	cases := []struct {
		resp      []byte
		wantParam *Value
		wantError error
	}{
		{
			resp: []byte(`
				<?xml version="1.0"?>
				<methodResponse></methodResponse>
			`),
		},
		{
			resp: []byte(`
				<?xml version="1.0"?>
				<methodResponse>
					<fault>
						<value>
							<struct>
									<member>
										<name>faultCode</name>
										<value><int>4</int></value>
									</member>
									<member>
										<name>faultString</name>
										<value><string>Too many parameters.</string></value>
									</member>
							</struct>
						</value>
					</fault>
				</methodResponse>
			`),
			wantError: &Fault{code: 4, str: "Too many parameters."},
		},
		{
			resp: []byte(`
				<?xml version="1.0"?>
				<methodResponse>
					<params>
						<param><value><string>South Dakota</string></value></param>
					</params>
				</methodResponse>
			`),
			wantParam: String("South Dakota"),
		},
	}
	for i, tc := range cases {
		got, err := ParseResponse(bytes.NewReader(tc.resp))
		if err != nil {
			if tc.wantError == nil {
				t.Fatalf("#%d: ParseResponse(%s) returns error: %s", i, tc.resp, err)
			}
			if !reflect.DeepEqual(tc.wantError, err) {
				t.Fatalf("#%d: ParseResponse(%s) returns error: %+v, want: %+v", i, tc.resp, err, tc.wantError)
			}
			if !reflect.DeepEqual(tc.wantParam, got) {
				t.Errorf("#%d: ParseResponse(%s) == %+v, want: %+v", i, tc.resp, got, tc.wantParam)
			}
		}
	}
}

func TestNewResponse(t *testing.T) {
	cases := []struct {
		param *Value
		want  []byte
	}{
		{
			param: nil,
			want:  []byte(`<?xml version="1.0"?><methodResponse></methodResponse>`),
		},
		{
			param: String("South Dakota"),
			want:  []byte(`<?xml version="1.0"?><methodResponse><params><param><value><string>South Dakota</string></value></param></params></methodResponse>`),
		},
	}

	for i, tc := range cases {
		resp := NewResponse(tc.param)
		buf := new(bytes.Buffer)
		err := resp.Write(buf)
		if err != nil {
			t.Fatalf("#%d: NewResponse(%+v).Write(buf) returns error: %s", i, tc.param, err)
		}
		got := buf.Bytes()
		if !bytes.Equal(got, tc.want) {
			t.Errorf("#%d: NewResponse(%+v).Write() == %s, want: %s", i, tc.param, got, tc.want)
		}
	}
}

func TestFaultResponse(t *testing.T) {
	cases := []struct {
		code int
		str  string
		want []byte
	}{
		{
			code: 4,
			str:  "Too many parameters.",
			want: []byte(`<?xml version="1.0"?><methodResponse><fault><value><struct><member><name>faultCode</name><value><int>4</int></value></member><member><name>faultString</name><value><string>Too many parameters.</string></value></member></struct></value></fault></methodResponse>`),
		},
	}

	for i, tc := range cases {
		resp := FaultResponse(tc.code, tc.str)
		buf := new(bytes.Buffer)
		err := resp.Write(buf)
		if err != nil {
			t.Fatalf("#%d: FaultResponse(%+v, %+v).Write(buf) returns error: %s", i, tc.code, tc.str, err)
		}
		got := buf.Bytes()
		if !bytes.Equal(got, tc.want) {
			t.Errorf("#%d: FaultResponse(%+v, %+v).Write() == %s, want: %s", i, tc.code, tc.str, got, tc.want)
		}
	}
}
