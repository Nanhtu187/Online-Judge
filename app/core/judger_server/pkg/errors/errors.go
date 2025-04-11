package errors

import (
	stderrors "errors"
	"fmt"
	"runtime"
	"sort"
	"strconv"
	"sync/atomic"
)

type domainError struct {
	rpcStatus uint32
	code      int
	message   string
}

type registeredErrorKey struct {
	rpcStatus int
	codeValue int
}

type registeredErrorValue struct {
	code     int
	location string
}

var disableNew uint32
var registeredErrors = map[registeredErrorKey]registeredErrorValue{}

func (e domainError) Error() string {
	return fmt.Sprintf("code: %d, message: %s", e.code, e.message)
}

// New creates a new error
func New(code int, httpCode int, message string) error {
	if atomic.LoadUint32(&disableNew) != 0 {
		panic("ONLY use errors.New for global variables")
	}

	if code >= 600 || code <= 200 {
		panic("code must between 200 and 600")
	}

	_, file, line, ok := runtime.Caller(1)
	if !ok {
		panic("can't get caller location")
	}

	rpcCode := httpCodeToRpcCode(httpCode)
	key := registeredErrorKey{
		rpcStatus: rpcCode,
		codeValue: code,
	}

	previous, existed := registeredErrors[key]
	if existed {
		fmt.Println("Existed at:", previous.location)
		fmt.Println(fmt.Sprintf("error code %d already existed", code))
		panic(fmt.Sprintf("error code %d already existed", code))
	}
	registeredErrors[key] = registeredErrorValue{
		code:     code,
		location: fmt.Sprintf("%s:%d", file, line),
	}

	return domainError{
		rpcStatus: uint32(rpcCode),
		code:      code,
		message:   message,
	}
}

// FinishNewErrors use to prevent New being called inside functions
func FinishNewErrors() {
	atomic.StoreUint32(&disableNew, 1)
}

// GetNextCode find the next error code
func GetNextCode(status string, httpCode string) int {
	httpCodeInt, _ := strconv.Atoi(httpCode)
	code, _ := strconv.Atoi(status)

	values := []int{400}

	httpCodeNum := httpCodeToRpcCode(httpCodeInt)
	for key := range registeredErrors {
		if key.rpcStatus == httpCodeNum && key.codeValue != 0 {
			values = append(values, code)
		}
	}
	sort.Ints(values)

	for i := httpCodeInt; i < (httpCodeInt/400+1)*100; i++ {
		fmt.Println(i, httpCodeInt)
		if i != code {
			return i
		}
	}

	nextValue := values[len(values)-1] + 1
	return nextValue
}

func WithMessage(err error, msg string) error {
	var domainNewErr domainError
	if stderrors.As(err, &domainNewErr) {
		return domainError{
			rpcStatus: domainNewErr.rpcStatus,
			code:      domainNewErr.code,
			message:   msg,
		}
	}

	return err
}

func httpCodeToRpcCode(httpCode int) int {
	switch httpCode {
	case 499:
		return 1
	case 500:
		return 13
	case 400:
		return 3
	case 504:
		return 4
	case 404:
		return 5
	case 409:
		return 6
	case 403:
		return 7
	case 429:
		return 8
	case 501:
		return 11
	case 401:
		return 16
	default:
		panic("not assign http code")
	}
}
