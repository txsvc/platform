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
		Equal(*Key) bool

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

func DecodeKey(encoded string) (*Key, error) {
	return nil, nil
}
