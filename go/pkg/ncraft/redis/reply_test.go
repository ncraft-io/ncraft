package redis

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReplyConversions(t *testing.T) {
	for _, v := range []any{int64(42), "42", []byte("42")} {
		n, err := Int(v, nil)
		require.NoError(t, err)
		require.Equal(t, 42, n)
		n64, err := Int64(v, nil)
		require.NoError(t, err)
		require.Equal(t, int64(42), n64)
	}
	for _, v := range []any{1.25, "1.25", []byte("1.25")} {
		n, err := Float64(v, nil)
		require.NoError(t, err)
		require.Equal(t, 1.25, n)
	}
	for _, v := range []any{true, int64(1), int64(-1), "true", []byte("1")} {
		b, err := Bool(v, nil)
		require.NoError(t, err)
		require.True(t, b)
	}
	for _, v := range []any{false, int64(0), "false", []byte("0")} {
		b, err := Bool(v, nil)
		require.NoError(t, err)
		require.False(t, b)
	}
	for _, v := range []any{"text\x00", []byte("text\x00")} {
		s, err := String(v, nil)
		require.NoError(t, err)
		require.Equal(t, "text\x00", s)
		b, err := Bytes(v, nil)
		require.NoError(t, err)
		require.Equal(t, []byte("text\x00"), b)
	}
	ints, err := Ints([]any{int64(1), "2", []byte("3"), nil}, nil)
	require.NoError(t, err)
	require.Equal(t, []int{1, 2, 3, 0}, ints)
	strings, err := Strings([]any{"a", []byte("b"), nil}, nil)
	require.NoError(t, err)
	require.Equal(t, []string{"a", "b", ""}, strings)
	m, err := StringMap([]any{"a", []byte("b"), []byte("c"), "d"}, nil)
	require.NoError(t, err)
	require.Equal(t, map[string]string{"a": "b", "c": "d"}, m)
}

func TestReplyErrors(t *testing.T) {
	converters := []func(any, error) error{
		func(v any, e error) error { _, e = Int(v, e); return e },
		func(v any, e error) error { _, e = Int64(v, e); return e },
		func(v any, e error) error { _, e = Float64(v, e); return e },
		func(v any, e error) error { _, e = Bool(v, e); return e },
		func(v any, e error) error { _, e = String(v, e); return e },
		func(v any, e error) error { _, e = Bytes(v, e); return e },
		func(v any, e error) error { _, e = Values(v, e); return e },
		func(v any, e error) error { _, e = Ints(v, e); return e },
		func(v any, e error) error { _, e = Strings(v, e); return e },
		func(v any, e error) error { _, e = StringMap(v, e); return e },
	}
	upstream := errors.New("connection failed")
	for _, convert := range converters {
		require.ErrorIs(t, convert(nil, upstream), upstream)
		require.ErrorIs(t, convert(nil, nil), ErrNil)
		require.Error(t, convert(struct{}{}, nil))
	}
	_, err := Int64("9223372036854775808", nil)
	require.Error(t, err)
	_, err = Float64("bad", nil)
	require.Error(t, err)
	_, err = Bool("bad", nil)
	require.Error(t, err)
	_, err = Ints([]any{"bad"}, nil)
	require.Error(t, err)
	_, err = Strings([]any{int64(1)}, nil)
	require.Error(t, err)
	for _, v := range []any{[]any{"unpaired"}, []any{int64(1), "v"}, []any{"k", true}} {
		_, err := StringMap(v, nil)
		require.Error(t, err)
	}
}
