package null_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gokern/null"
)

func TestNull_ZeroValueIsNull(t *testing.T) {
	t.Parallel()

	var n null.Null[string]

	require.False(t, n.Valid)
	require.True(t, n.IsZero())

	got, ok := n.Get()
	require.False(t, ok)
	require.Empty(t, got)
}

func TestFrom(t *testing.T) {
	t.Parallel()

	n := null.From(42)

	require.True(t, n.Valid)
	require.False(t, n.IsZero())

	got, ok := n.Get()
	require.True(t, ok)
	require.Equal(t, 42, got)
}

func TestFrom_ZeroValueIsSet(t *testing.T) {
	t.Parallel()

	n := null.From(0)

	require.True(t, n.Valid)
	require.False(t, n.IsZero())
}

func TestFromPtr(t *testing.T) {
	t.Parallel()

	t.Run("nil pointer becomes null", func(t *testing.T) {
		t.Parallel()

		n := null.FromPtr[int](nil)

		require.False(t, n.Valid)
	})

	t.Run("pointer value is copied in", func(t *testing.T) {
		t.Parallel()

		v := "note"
		n := null.FromPtr(&v)

		require.True(t, n.Valid)
		require.Equal(t, "note", n.Val)

		v = "changed" // the container holds a copy, not the pointer

		require.Equal(t, "note", n.Val)
	})
}

func TestOr(t *testing.T) {
	t.Parallel()

	require.Equal(t, 7, null.From(7).Or(9))
	require.Equal(t, 9, null.Null[int]{}.Or(9))
}

func TestPtr(t *testing.T) {
	t.Parallel()

	t.Run("null yields nil", func(t *testing.T) {
		t.Parallel()

		require.Nil(t, null.Null[int]{}.Ptr())
	})

	t.Run("set yields a pointer to a copy", func(t *testing.T) {
		t.Parallel()

		n := null.From("v")
		p := n.Ptr()

		require.NotNil(t, p)
		require.Equal(t, "v", *p)

		*p = "changed" // mutating through the pointer must not reach the container

		require.Equal(t, "v", n.Val)
	})
}

func TestNull_Comparable(t *testing.T) {
	t.Parallel()

	one, alsoOne := null.From(1), null.From(1)
	setZero, nullA, nullB := null.From(0), null.Null[int]{}, null.Null[int]{}

	require.True(t, one == alsoOne)       //nolint:testifylint // == is the subject
	require.False(t, one == null.From(2)) //nolint:testifylint // == is the subject
	require.False(t, setZero == nullA)    //nolint:testifylint // == is the subject
	require.True(t, nullA == nullB)       //nolint:testifylint // == is the subject
}

func TestNull_OmitzeroIntegration(t *testing.T) {
	t.Parallel()

	type payload struct {
		Note null.Null[string] `json:"note,omitzero"`
	}

	t.Run("null field is omitted", func(t *testing.T) {
		t.Parallel()

		//nolint:musttag // Null is deliberately codec- and tag-free; the raw shape is the subject
		out, err := json.Marshal(payload{})
		require.NoError(t, err)
		require.JSONEq(t, `{}`, string(out))
	})

	t.Run("set field is emitted as the raw struct", func(t *testing.T) {
		t.Parallel()

		//nolint:musttag // Null is deliberately codec- and tag-free; the raw shape is the subject
		out, err := json.Marshal(payload{Note: null.From("hi")})
		require.NoError(t, err)
		require.JSONEq(t, `{"note":{"Val":"hi","Valid":true}}`, string(out))
	})
}
