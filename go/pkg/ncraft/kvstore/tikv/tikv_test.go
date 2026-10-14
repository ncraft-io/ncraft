package tikv

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"strings"
	"testing"
	"time"
)

func TestTikv(t *testing.T) {
	address := os.Getenv("NCRAFT_TIKV_PD")
	if address == "" {
		t.Skip("set NCRAFT_TIKV_PD to run the TiKV integration test")
	}
	s, err := NewStore(Options{Pds: strings.Split(address, ","), BucketName: fmt.Sprintf("ncraft-test-%d", time.Now().UnixNano())})
	require.NoError(t, err)
	defer s.Close()
	defer s.DeleteRange(context.Background(), nil, nil)

	// 1. test single key
	err = s.Put(context.TODO(), []byte("foo"), []byte("bar"))
	assert.NoError(t, err)

	v, err := s.Get(context.TODO(), []byte("foo"))
	assert.NoError(t, err)
	assert.Equal(t, "bar", string(v))

	err = s.Delete(context.TODO(), []byte("foo"))
	assert.NoError(t, err)

	v, err = s.Get(context.TODO(), []byte("foo"))
	assert.NoError(t, err)
	assert.Equal(t, []byte(nil), v)

	// 2. test batch keys
	err = s.BatchPut(context.TODO(), [][]byte{[]byte("foo1"), []byte("foo2")}, [][]byte{[]byte("bar1"), []byte("bar2")})
	assert.NoError(t, err)

	values, err := s.BatchGet(context.TODO(), [][]byte{[]byte("foo1"), []byte("foo2")})
	assert.NoError(t, err)
	assert.Equal(t, "bar1", string(values[0]))
	assert.Equal(t, "bar2", string(values[1]))

	err = s.BatchDelete(context.TODO(), [][]byte{[]byte("foo1"), []byte("foo2")})
	assert.NoError(t, err)

	values, err = s.BatchGet(context.TODO(), [][]byte{[]byte("foo1"), []byte("foo2")})
	assert.NoError(t, err)
	assert.Equal(t, [][]byte{nil, nil}, values)
}
