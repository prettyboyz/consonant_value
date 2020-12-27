
package toml

// tomlType represents any Go type that corresponds to a TOML type.
// While the first draft of the TOML spec has a simplistic type system that
// probably doesn't need this level of sophistication, we seem to be militating
// toward adding real composite types.
type tomlType interface {
	typeString() string
}

// typeEqual accepts any two types and returns true if they are equal.
func typeEqual(t1, t2 tomlType) bool {
	if t1 == nil || t2 == nil {
		return false
	}
	return t1.typeString() == t2.typeString()
}

func typeIsHash(t tomlType) bool {
	return typeEqual(t, tomlHash) || typeEqual(t, tomlArrayHash)
}

type tomlBaseType string

func (btype tomlBaseType) typeString() string {
	return string(btype)
}

func (btype tomlBaseType) String() string {
	return btype.typeString()
}

var (
	tomlInteger   tomlBaseType = "Integer"
	tomlFloat     tomlBaseType = "Float"
	tomlDatetime  tomlBaseType = "Datetime"