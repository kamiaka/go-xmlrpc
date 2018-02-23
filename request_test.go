package xmlrpc

import (
	"bytes"
	"reflect"
	"testing"
)

func TestNewRequest(t *testing.T) {
	cases := []struct {
		name   string
		params []*Value
		want   []byte
	}{
		{
			name:   "emptyParam.method",
			params: []*Value{},
			want:   []byte(`<?xml version="1.0"?><methodCall><methodName>emptyParam.method</methodName></methodCall>`),
		},
		{
			name: "some.method",
			params: []*Value{
				String("foo"),
				Int(42),
			},
			want: []byte(`<?xml version="1.0"?><methodCall><methodName>some.method</methodName><params><param><value><string>foo</string></value></param><param><value><int>42</int></value></param></params></methodCall>`),
		},
	}

	for i, tc := range cases {
		req, err := NewRequest(tc.name, tc.params...)
		if err != nil {
			t.Fatalf("#%d: NewRequest(%#v, %+v...) returns error: %s", i, tc.name, tc.params, err)
		}
		buf := new(bytes.Buffer)
		err = req.Write(buf)
		if err != nil {
			t.Fatalf("#%d: NewRequest(%#v, %+v...).Write(buf) returns error: %s", i, tc.name, tc.params, err)
		}
		got := buf.Bytes()
		if !bytes.Equal(got, tc.want) {
			t.Errorf("#%d: NewRequest(%#v, %+v...).Write() == %s, want: %s", i, tc.name, tc.params, got, tc.want)
		}
	}
}

func TestParseRequest(t *testing.T) {
	cases := []struct {
		wantName   string
		wantParams []*Value
		xml        []byte
	}{
		{
			wantName:   "emptyParam.method",
			wantParams: []*Value{},
			xml:        []byte(`<?xml version="1.0"?><methodCall><methodName>emptyParam.method</methodName></methodCall>`),
		},
		{
			wantName: "some.method",
			wantParams: []*Value{
				String("foo"),
				Int(42),
			},
			xml: []byte(`<?xml version="1.0"?><methodCall><methodName>some.method</methodName><params><param><value><string>foo</string></value></param><param><value><int>42</int></value></param></params></methodCall>`),
		},
	}

	for i, tc := range cases {
		got, err := ParseRequest(bytes.NewReader(tc.xml))
		if err != nil {
			t.Fatalf("#%d: ParseRequest(...) returns error: %s", i, err)
		}
		if got.Method() != tc.wantName {
			t.Errorf("#%d: ParseRequest(...).Method() == %s, want: %s", i, got.Method(), tc.wantName)
		}
		if !reflect.DeepEqual(got.Params(), tc.wantParams) {
			t.Errorf("#%d: ParseRequest(...).Params() == %#v, want: %#v", i, got.Params(), tc.wantParams)
		}
	}
}
