package xmlrpc

// Type of value
type Type int

// Types
const (
	TypeUnknown Type = iota
	TypeArray
	TypeBool
	TypeInt
	TypeString
	TypeStruct
)

func typeOf(v *Value) Type {
	if v.ArrayPtr != nil {
		return TypeArray
	}
	if v.BoolPtr != nil {
		return TypeBool
	}
	if v.IntPtr != nil {
		return TypeInt
	}
	if v.StringPtr != nil {
		return TypeString
	}
	if v.StructPtr != nil {
		return TypeStruct
	}
	return TypeUnknown
}
