package errors

import (
	"github.com/stretchr/testify/assert"
	"sync/atomic"
	"testing"
)

func TestNewErrorString(t *testing.T) {
	e := New(411, 400, "some error")
	assert.Equal(t, "code: 411, message: some error", e.Error())
}

func TestNewError_OK(t *testing.T) {
	e := New(402, 400, "some error")
	assert.Equal(t, uint32(3), e.(domainError).rpcStatus)
	assert.Equal(t, 402, e.(domainError).code)
	assert.Equal(t, "some error", e.(domainError).message)
}

func TestNewError_Failed(t *testing.T) {
	table := []struct {
		name     string
		code     int
		expected string
	}{
		{
			name:     "too large code",
			code:     601,
			expected: "code must between 200 and 600",
		},
		{
			name:     "too small code",
			code:     111,
			expected: "code must between 200 and 600",
		},
	}

	for _, e := range table {
		t.Run(e.name, func(t *testing.T) {
			registeredErrors = map[registeredErrorKey]registeredErrorValue{}

			assert.PanicsWithValue(t, e.expected, func() {
				_ = New(e.code, 400, "some error")
			})
		})
	}
}

func TestNewError_Duplicated(t *testing.T) {
	registeredErrors = map[registeredErrorKey]registeredErrorValue{}

	_ = New(401, 400, "some error 1")
	assert.PanicsWithValue(t, "error code 401 already existed", func() {
		_ = New(401, 400, "some error 2")
	})
}

func TestNewError_After_Finish_Failed(t *testing.T) {
	registeredErrors = map[registeredErrorKey]registeredErrorValue{}
	FinishNewErrors()

	assert.PanicsWithValue(t, "ONLY use errors.New for global variables", func() {
		_ = New(400, 404, "some error")
	})

	atomic.StoreUint32(&disableNew, 0)
}

func TestGetNextCode(t *testing.T) {
	registeredErrors = map[registeredErrorKey]registeredErrorValue{}

	_ = New(400, 400, "some error 1")
	_ = New(401, 400, "some error 1")
	_ = New(402, 400, "some error 1")

	s := GetNextCode("402", "400")
	assert.Equal(t, 403, s)
}
