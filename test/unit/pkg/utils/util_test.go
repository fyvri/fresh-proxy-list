package util_test

import (
	"io"
	"net/http"
)

var (
	expectedButGotMessage             = "Expected %v = %v, but got = %v"
	expectedErrorButGotMessage        = "Expected %v error = %v, but got = %v"
	expectedReturnNonNil              = "Expected %v to return a non-nil %v"
	expectedTypeAssertionErrorMessage = "Expected type assertion error, but got = %v"
	testScheme                        = "https"
	testHost                          = "example.com"
	testPath                          = "/path"
	testRawQuery                      = "query=1"
	testRawURL                        = testScheme + "://" + testHost
	testFullURL                       = testRawURL + testPath + "?" + testRawQuery
)

type mockTransport struct {
	response *http.Response
	err      error
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.response, m.err
}

type mockReadCloser struct {
	data     []byte
	errRead  error
	errClose error
}

func (m *mockReadCloser) Read(p []byte) (int, error) {
	if m.errRead != nil {
		return 0, m.errRead
	}
	copy(p, m.data)
	return len(m.data), io.EOF
}

func (m *mockReadCloser) Close() error {
	if m.errClose != nil {
		return m.errClose
	}
	return nil
}
