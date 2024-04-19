package errors

import (
	"context"
	stderrors "errors"
	"fmt"
	"github.com/Nanhtu187/Online-Judge/app/iam/pkg/errors/generated"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
	"io"
	"net/http"
	"strconv"
)

func (e domainError) toRPCError() error {
	st := status.New(codes.Code(e.rpcStatus), e.message)

	st, err := st.WithDetails(
		&generated.ErrorDetailString{
			Field: "code",
			Value: strconv.Itoa(e.code),
		},
	)
	if err != nil {
		return err
	}

	return st.Err()
}

func fromRPCStatus(st *status.Status) (domainError, bool) {
	if len(st.Details()) == 0 {
		return domainError{}, false
	}
	code, ok := st.Details()[0].(*generated.ErrorDetailString)
	if !ok {
		return domainError{}, false
	}

	if code.Field != "code" {
		return domainError{}, false
	}

	codeNum, _ := strconv.Atoi(code.Value)
	return domainError{
		rpcStatus: uint32(st.Code()),
		code:      codeNum,
		message:   st.Message(),
	}, true
}

// UnaryServerInterceptor ...
func UnaryServerInterceptor(
	ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) (interface{}, error) {
	resp, err := handler(ctx, req)
	if err != nil {
		var domainErr domainError
		if stderrors.As(err, &domainErr) {
			return nil, domainErr.toRPCError()
		}

		st, ok := status.FromError(err)
		if ok {
			return nil, st.Err()
		}

		st = status.New(codes.Unknown, err.Error())
		return nil, st.Err()
	}
	return resp, nil
}

func statusToErrorBody(s *status.Status, marshaller runtime.Marshaler) ([]byte, error) {
	type errorBody struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	domainErr, ok := fromRPCStatus(s)
	if !ok {
		body := errorBody{
			Code:    99,
			Message: fmt.Sprintf("rpc: %d, message: %s", s.Code(), s.Message()),
		}
		return marshaller.Marshal(body)
	}

	body := errorBody{
		Code:    domainErr.code,
		Message: domainErr.message,
	}
	return marshaller.Marshal(body)
}

// CustomerHTTPError customizing error codes
func CustomerHTTPError(
	_ context.Context, _ *runtime.ServeMux, marshaler runtime.Marshaler,
	w http.ResponseWriter, _ *http.Request, err error,
) {
	const fallback = `{"code": "99", "message": "failed to marshal error message"}`

	s, ok := status.FromError(err)
	if !ok {
		s = status.New(codes.Unknown, err.Error())
	}

	w.Header().Del("Trailer")

	contentType := marshaler.ContentType(1)
	w.Header().Set("Content-Type", contentType)

	buf, merr := statusToErrorBody(s, marshaler)
	if merr != nil {
		grpclog.Infof("Failed to marshal error message: %v", merr)
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := io.WriteString(w, fallback); err != nil {
			grpclog.Infof("Failed to write response: %v", err)
		}
		return
	}

	st := runtime.HTTPStatusFromCode(s.Code())
	w.WriteHeader(st)
	if _, err := w.Write(buf); err != nil {
		grpclog.Infof("Failed to write response: %v", err)
	}
}
