// Package null provides a two-state nullable value container with no
// serialization codecs.
//
// Null[T] holds a value and a validity flag, nothing else. It
// intentionally implements neither json.Marshaler/Unmarshaler nor
// sql.Scanner/driver.Valuer: encoding/json sees only the raw struct fields,
// database drivers reject the type, and conversion stays an explicit call
// (see FromPtr and Ptr).
//
// The zero value is null. Null[T] is comparable when T is, and IsZero
// supports `json:",omitzero"`.
package null

// Null holds a value of type T and reports whether it is set.
// The zero value is null (not set). There is no third "absent vs explicit
// null" state.
type Null[T any] struct {
	Val   T
	Valid bool
}

// From returns a set Null holding v.
func From[T any](v T) Null[T] {
	return Null[T]{Val: v, Valid: true}
}

// FromPtr returns null when p is nil, otherwise a set Null holding a
// copy of *p.
func FromPtr[T any](p *T) Null[T] {
	if p == nil {
		return Null[T]{} //nolint:exhaustruct // the zero value is null by design
	}

	return From(*p)
}

// Get returns the value and whether it is set.
func (n Null[T]) Get() (T, bool) {
	return n.Val, n.Valid
}

// Or returns the value when set, fallback otherwise.
func (n Null[T]) Or(fallback T) T {
	if n.Valid {
		return n.Val
	}

	return fallback
}

// Ptr returns nil when null, otherwise a pointer to a copy of the value.
func (n Null[T]) Ptr() *T {
	if !n.Valid {
		return nil
	}

	p := n.Val

	return &p
}

// IsZero reports whether the value is null. The omitzero option of
// encoding/json uses it to omit unset fields.
func (n Null[T]) IsZero() bool {
	return !n.Valid
}
