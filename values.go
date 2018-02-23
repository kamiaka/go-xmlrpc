package xmlrpc

type params struct {
	Data []*param `xml:"param"`
}

type param struct {
	Value *Value `xml:"value"`
}

// Value ...
type Value struct {
	ArrayPtr  *[]*Value  `xml:"array>data>value,omitempty"`
	BoolPtr   *bool      `xml:"boolean,omitempty"`
	IntPtr    *int       `xml:"int,omitempty"`
	StringPtr *string    `xml:"string,omitempty"`
	StructPtr *structVal `xml:"struct,omitempty"`
}

type StructMember struct {
	Name  string `xml:"name"`
	Value *Value `xml:"value"`
}

type structVal struct {
	Members []*StructMember `xml:"member"`
}

// Interface ...
func (v *Value) Interface() interface{} {
	if v.ArrayPtr != nil {
		return v.Array()
	}
	if v.BoolPtr != nil {
		return v.Bool()
	}
	if v.StringPtr != nil {
		return v.String()
	}
	if v.IntPtr != nil {
		return v.Int()
	}
	if v.StructPtr != nil {
		return v.Struct()
	}
	panic("Unsupported value type")
}

// Array ...
func (v *Value) Array() []interface{} {
	result := make([]interface{}, len(*v.ArrayPtr))
	for i, v := range *v.ArrayPtr {
		result[i] = v.Interface()
	}
	return result
}

// Bool ...
func (v *Value) Bool() bool {
	return *v.BoolPtr
}

// String ...
func (v *Value) String() string {
	return *v.StringPtr
}

// Int ...
func (v *Value) Int() int {
	return *v.IntPtr
}

// Struct ...
func (v *Value) Struct() map[string]interface{} {
	result := make(map[string]interface{})
	for _, m := range v.StructPtr.Members {
		result[m.Name] = m.Value.Interface()
	}
	return result
}

// Array value
func Array(vs ...*Value) *Value {
	return &Value{
		ArrayPtr: &vs,
	}
}

// Int value
func Int(i int) *Value {
	return &Value{
		IntPtr: &i,
	}
}

// String value
func String(s string) *Value {
	return &Value{
		StringPtr: &s,
	}
}

// Struct value
func Struct(members ...*StructMember) *Value {
	return &Value{
		StructPtr: &structVal{
			Members: members,
		},
	}
}

// Member of struct
func Member(name string, v *Value) *StructMember {
	return &StructMember{
		Name:  name,
		Value: v,
	}
}
