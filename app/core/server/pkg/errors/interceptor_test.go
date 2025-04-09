package errors

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/status"
	"testing"
)

func TestFromAndToRPCError(t *testing.T) {
	registeredErrors = map[registeredErrorKey]registeredErrorValue{}
	e := New(400, 400, "some error")
	fmt.Printf("error %#v\n", e)
	rpcErr := e.(domainError).toRPCError()

	st := status.Convert(rpcErr)

	err, ok := fromRPCStatus(st)
	assert.True(t, ok)
	assert.Equal(t, e, err)
}
