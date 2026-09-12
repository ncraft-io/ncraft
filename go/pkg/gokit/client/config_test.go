package client

import (
	"testing"

	"github.com/ncraft-io/ncraft/go/pkg/ncraft/config"
	"github.com/ncraft-io/ncraft/go/pkg/ncraft/config/source/memory"
	"github.com/stretchr/testify/require"
)

func useConfig(t *testing.T, yaml string) {
	t.Helper()
	loaded, err := config.NewConfig(config.WithSource(memory.NewSource(memory.WithYAML([]byte(yaml)))))
	require.NoError(t, err)
	previous := config.DefaultConfig
	config.DefaultConfig = loaded
	t.Cleanup(func() {
		config.DefaultConfig = previous
		require.NoError(t, loaded.Close())
	})
}

func TestConfigPrecedence(t *testing.T) {
	useConfig(t, `
ncraft:
  sd:
    mode: direct
    transport: grpc
    retry: {enable: true, timeout: 1000, max: 3}
    direct:
      books: {name: original, urls: ["localhost:9001", "localhost:9002"]}
      users: {urls: ["localhost:9003"]}
  client:
    serviceName: default-service
    http: {timeout: 1000, headers: [Authorization, X-Request-ID], enveloped: true}
    grpc: {dialTimeout: 3000, waitForReady: true, maxRecvMsgSize: 4096}
    sd:
      transport: http
      retry: {timeout: 2000}
      direct:
        books: {name: common}
    book:
      serviceName: books
      http:
        headers: [X-Book-ID]
      sd:
        retry: {max: 5}
        direct:
          books: {urls: ["localhost:8000"]}
    unrelated:
      sd: {mode: invalid-for-books}
`)
	common, err := LoadConfig()
	require.NoError(t, err)
	require.Equal(t, "default-service", common.ServiceName)
	require.Equal(t, "direct", common.SD.Mode)
	require.Equal(t, "http", common.SD.Transport)
	require.Equal(t, 2000, common.SD.Retry.Timeout)
	require.Equal(t, 3, common.SD.Retry.Max)
	require.Equal(t, []string{"localhost:9001", "localhost:9002"}, common.SD.Direct["books"].Urls)

	book, err := LoadConfig("book")
	require.NoError(t, err)
	require.Equal(t, "books", book.ServiceName)
	require.Equal(t, []string{"X-Book-ID"}, book.HTTP.Headers)
	require.True(t, book.HTTP.Enveloped)
	require.Equal(t, 1000, book.HTTP.Timeout)
	require.Equal(t, &GRPCConfig{DialTimeout: 3000, WaitForReady: true, MaxRecvMsgSize: 4096}, book.GRPC)
	require.Equal(t, "direct", book.SD.Mode)
	require.Equal(t, "http", book.SD.Transport)
	require.True(t, book.SD.Retry.Enable)
	require.Equal(t, 2000, book.SD.Retry.Timeout)
	require.Equal(t, 5, book.SD.Retry.Max)
	require.Equal(t, "common", book.SD.Direct["books"].Name)
	require.Equal(t, []string{"localhost:8000"}, book.SD.Direct["books"].Urls)
	require.Equal(t, []string{"localhost:9003"}, book.SD.Direct["users"].Urls)

	missing, err := LoadConfig("missing")
	require.NoError(t, err)
	require.Equal(t, common, missing)
	require.Equal(t, book, NewConfig("book"))

	// A caller's mutations and service overrides cannot modify another load or
	// the global config store, including nested maps, pointers and slices.
	book.SD.Direct["books"].Urls[0] = "changed"
	book.SD.Retry.Max = 99
	book.HTTP.Headers[0] = "changed"
	again, err := LoadConfig()
	require.NoError(t, err)
	require.Equal(t, common, again)
	fresh, err := LoadConfig("book")
	require.NoError(t, err)
	require.Equal(t, []string{"localhost:8000"}, fresh.SD.Direct["books"].Urls)
	require.Equal(t, 5, fresh.SD.Retry.Max)
	require.Equal(t, []string{"X-Book-ID"}, fresh.HTTP.Headers)
}

func TestConfigExplicitOverrides(t *testing.T) {
	useConfig(t, `
sd:
  mode: direct
  retry: {enable: true, timeout: 1000, max: 3}
  direct:
    books: {name: books, urls: ["localhost:9000"]}
client:
  serviceName: shared
  http: {timeout: 2000,enveloped: true,headers: [Authorization]}
  grpc: {waitForReady: true, dialTimeout: 1000}
  zero:
    serviceName: ""
    http: {timeout: 0,enveloped: false,headers: []}
    grpc: {waitForReady: false, dialTimeout: 0}
    sd:
      retry: {enable: false, timeout: 0, max: 0}
      direct:
        books: {name: "", urls: []}
  clear:
    http: null
    grpc: null
    headers: null
    sd: null
`)
	zero, err := LoadConfig("zero")
	require.NoError(t, err)
	require.Empty(t, zero.ServiceName)
	require.False(t, zero.HTTP.Enveloped)
	require.NotNil(t, zero.HTTP.Headers)
	require.Empty(t, zero.HTTP.Headers)
	require.Zero(t, zero.HTTP.Timeout)
	require.False(t, zero.GRPC.WaitForReady)
	require.Zero(t, zero.GRPC.DialTimeout)
	require.False(t, zero.SD.Retry.Enable)
	require.Zero(t, zero.SD.Retry.Timeout)
	require.Zero(t, zero.SD.Retry.Max)
	require.Empty(t, zero.SD.Direct["books"].Name)
	require.Empty(t, zero.SD.Direct["books"].Urls)

	clear, err := LoadConfig("clear")
	require.NoError(t, err)
	require.Nil(t, clear.HTTP)
	require.Nil(t, clear.GRPC)
	require.Empty(t, clear.SD.Mode)
	require.Nil(t, clear.SD.Retry)
	require.Nil(t, clear.SD.Direct)
}

func TestConfigSourcesAndServiceNames(t *testing.T) {
	for _, tc := range []struct {
		name, yaml string
		service    []string
		transport  string
	}{
		{"standalone", "sd: {mode: direct, transport: grpc}", nil, "grpc"},
		{"common", "client: {sd: {mode: direct, transport: http}}", nil, "http"},
		{"service only", "client: {books: {sd: {mode: direct, transport: http}}}", []string{"books"}, "http"},
		{"namespace preference", "sd: {transport: grpc}\nncraft: {sd: {transport: http}}", nil, "http"},
		{"namespace fallback", "sd: {transport: http}\nncraft: {client: {headers: [token]}}", nil, "http"},
		{"literal dotted name", "client: {'sample.v1.BookService': {sd: {transport: http}}}", []string{"sample.v1.BookService"}, "http"},
		{"nested path", "client: {sample: {v1: {BookService: {sd: {transport: grpc}}}}}", []string{"sample", "v1", "BookService"}, "grpc"},
		{"literal first", "client: {'sample.Book': {sd: {transport: http}}, sample: {Book: {sd: {transport: grpc}}}}", []string{"sample.Book"}, "http"},
		{"namespaced service", "client: {book: {sd: {transport: grpc}}}\nncraft: {client: {book: {sd: {transport: http}}}}", []string{"book"}, "http"},
		{"missing everything", "{}", []string{"missing"}, ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			useConfig(t, tc.yaml)
			cfg, err := LoadConfig(tc.service...)
			require.NoError(t, err)
			require.Equal(t, tc.transport, cfg.SD.Transport)
		})
	}
}

func TestConfigInvalidValues(t *testing.T) {
	for _, tc := range []struct{ yaml, path string }{
		{"sd: invalid", "sd"},
		{"client: invalid", "client"},
		{"client: {books: invalid}", "books"},
		{"client: {http: {headers: invalid}}", "headers"},
		{"client: {sd: {retry: {max: invalid}}}", "max"},
		{"client: {books: {http: {timeout: invalid}}}", "timeout"},
	} {
		t.Run(tc.yaml, func(t *testing.T) {
			useConfig(t, tc.yaml)
			cfg, err := LoadConfig("books")
			require.Nil(t, cfg)
			require.ErrorContains(t, err, tc.path)
			require.Nil(t, NewConfig("books"))
		})
	}
	useConfig(t, "{}")
	for _, name := range []string{"", " ", "books.", ".books", "a..books"} {
		_, err := LoadConfig(name)
		require.Error(t, err)
	}
}
