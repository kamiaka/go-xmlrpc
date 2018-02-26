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
	BoolPtr   *int       `xml:"boolean,omitempty"`
	IntPtr    *int       `xml:"int,omitempty"`
	StringPtr *string    `xml:"string,omitempty"`
	StructPtr *structVal `xml:"struct,omitempty"`
}

// StructMember ...
type StructMember struct {
	Name  string `xml:"name"`
	Value *Value `xml:"value"`
}

type structVal struct {
	Members []*StructMember `xml:"member"`
}

// Type ...
func (v *Value) Type() Type {
	return typeOf(v)
}

// Interface ...
func (v *Value) Interface() interface{} {
	if v == nil {
		return nil
	}
	if v.ArrayPtr != nil {
		return arrayToInterface(v)
	}
	if v.BoolPtr != nil {
		return *v.BoolPtr == 1
	}
	if v.StringPtr != nil {
		return *v.StringPtr
	}
	if v.IntPtr != nil {
		return *v.IntPtr
	}
	if v.StructPtr != nil {
		return structToInterface(v)
	}
	panic("Unsupported value type")
}

func arrayToInterface(v *Value) []interface{} {
	result := make([]interface{}, len(*v.ArrayPtr))
	for i, v := range *v.ArrayPtr {
		result[i] = v.Interface()
	}
	return result
}

func structToInterface(v *Value) map[string]interface{} {
	result := make(map[string]interface{})
	for _, m := range v.StructPtr.Members {
		result[m.Name] = m.Value.Interface()
	}
	return result
}

// Array ...
func (v *Value) Array() []*Value {
	if v == nil || v.ArrayPtr == nil {
		return []*Value{}
	}
	return *v.ArrayPtr
}

// Index value of array
func (v *Value) Index(i int) *Value {
	arr := v.Array()
	if len(arr) < i+1 {
		return nil
	}
	return arr[i]
}

// Bool ...
func (v *Value) Bool() bool {
	if v == nil || v.BoolPtr == nil {
		return false
	}
	return *v.BoolPtr == 1
}

// String ...
func (v *Value) String() string {
	if v == nil || v.StringPtr == nil {
		return ""
	}
	return *v.StringPtr
}

// Int ...
func (v *Value) Int() int {
	if v == nil || v.IntPtr == nil {
		return 0
	}
	return *v.IntPtr
}

// Struct ...
func (v *Value) Struct() []*StructMember {
	if v == nil || v.StructPtr == nil {
		return []*StructMember{}
	}
	return v.StructPtr.Members
}

// Field ...
func (v *Value) Field(i int) *Value {
	members := v.Struct()
	if len(members) < i+1 {
		return nil
	}
	return members[i].Value
}

// FieldByName ...
func (v *Value) FieldByName(f string) *Value {
	for _, m := range v.Struct() {
		if m.Name == f {
			return m.Value
		}
	}
	return nil
}

// Array value
func Array(vs ...*Value) *Value {
	return &Value{
		ArrayPtr: &vs,
	}
}

// Bool value
func Bool(b bool) *Value {
	var i int
	if b {
		i = 1
	}
	return &Value{
		BoolPtr: &i,
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
