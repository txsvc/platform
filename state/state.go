package state

type (
	Key interface {

		// Name returns the key within the given namespace
		Name() string

		// Kind returns the kind of key
		Kind() string

		// String returns a string representation of the key
		String() string

		// Equal reports whether two keys are equal. Two keys are equal if they are both nil, or if their kinds and names are equal.
		Equal(Key) bool

		// Encode returns an opaque representation of the key suitable for use in HTML and URLs
		Encode() string

		// NativeKey returns the storage provider specific representation of the key
		NativeKey() interface{}
	}

	KeyImpl struct {
		name string
		kind string
	}

	StateProvider interface {
		DecodeKey(string) (*Key, error)
		NewKey(string, string) (*Key, error)
	}
)

var (
	// Interface guard
	_ Key = (*KeyImpl)(nil)
)

func NewKey(kind, name string) interface{} {
	return &KeyImpl{name: name, kind: kind}
}

func DecodeKey(encoded string) (*Key, error) {
	return nil, nil
}

func (k *KeyImpl) Name() string {
	return k.name
}

func (k *KeyImpl) Kind() string {
	return k.kind
}

func (k *KeyImpl) String() string {
	return k.kind + "." + k.name
}

func (k *KeyImpl) Equal(key Key) bool {
	if key == nil {
		return false
	}
	return k.kind == key.Kind() && k.name == key.Name()
}

func (k *KeyImpl) Encode() string {
	return k.kind + "." + k.name
}

func (k *KeyImpl) NativeKey() interface{} {
	return nil
}
