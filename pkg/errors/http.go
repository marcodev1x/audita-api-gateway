package errors

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/audita-bids/private-kit/decode"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type HTTPError struct {
	Status  int    `json:"-"`
	Code    string `json:"code,omitempty"`
	Message string `json:"message"`
}

func (e *HTTPError) Error() string {
	return e.Message
}

func New(status int, code, message string) error {
	return &HTTPError{
		Status:  status,
		Code:    code,
		Message: message,
	}
}

func BadRequest(code, message string) error {
	return New(http.StatusBadRequest, code, message)
}

func Unauthorized(message string) error {
	return New(http.StatusUnauthorized, "unauthorized", message)
}

func Forbidden(message string) error {
	return New(http.StatusForbidden, "forbidden", message)
}

func NotFound(resource string) error {
	return New(http.StatusNotFound, "not_found", resource+" not found")
}

func Conflict(message string) error {
	return New(http.StatusConflict, "conflict", message)
}

func Internal(message string) error {
	return New(http.StatusInternalServerError, "internal_error", message)
}

func FromGRPC(err error) *HTTPError {
	st, ok := status.FromError(err)
	if !ok {
		return nil
	}

	code := st.Code()
	message := cleanGRPCMessage(st.Message())

	// Services answer with Unknown/Internal plus a plain error string, so the
	// real intent has to be read from the message itself.
	if code == codes.Unknown || code == codes.Internal {
		if classified := classify(message); classified != nil {
			return classified
		}
	}

	httpStatus := grpcToHTTPStatus(code)

	return &HTTPError{
		Status:  httpStatus,
		Code:    statusCode(httpStatus),
		Message: safeMessage(httpStatus, message),
	}
}

func ParseError(err error) *HTTPError {
	if err == nil {
		return nil
	}

	var httpErr *HTTPError
	if errors.As(err, &httpErr) {
		return httpErr
	}

	// Field validation describes the caller's own payload, so it is safe to
	// echo back as is.
	var svcErr *decode.SvcError
	if errors.As(err, &svcErr) {
		return &HTTPError{
			Status:  http.StatusBadRequest,
			Code:    "validation_error",
			Message: svcErr.Message,
		}
	}

	// A body the decoder could not read is the caller's mistake, not ours.
	if malformedPayload(err) {
		return &HTTPError{
			Status:  http.StatusBadRequest,
			Code:    "invalid_request",
			Message: "invalid request payload",
		}
	}

	if grpcErr := FromGRPC(err); grpcErr != nil {
		return grpcErr
	}

	if classified := classify(err.Error()); classified != nil {
		return classified
	}

	return &HTTPError{
		Status:  http.StatusInternalServerError,
		Code:    statusCode(http.StatusInternalServerError),
		Message: statusMessage(http.StatusInternalServerError),
	}
}

// classify recognises the failures this gateway knows about. Each branch
// answers with a fixed message: the incoming text only picks the case, it is
// never forwarded to the client.
func classify(message string) *HTTPError {
	msg := strings.ToLower(message)

	switch {
	case containsAny(msg, []string{"invalid credentials", "invalid client login", "invalid login"}):
		return &HTTPError{
			Status:  http.StatusUnauthorized,
			Code:    "invalid_credentials",
			Message: "invalid credentials",
		}
	case containsAny(msg, []string{
		"token not found", "token is empty", "invalid token", "token expired",
		"unauthorized", "user not found",
	}):
		return &HTTPError{
			Status:  http.StatusUnauthorized,
			Code:    "unauthorized",
			Message: "unauthorized",
		}
	case containsAny(msg, []string{"insufficient scopes for payer"}):
		return &HTTPError{
			Status:  http.StatusForbidden,
			Code:    "plan_required",
			Message: "your plan does not include this feature",
		}
	case containsAny(msg, []string{"company required"}):
		return &HTTPError{
			Status:  http.StatusBadRequest,
			Code:    "company_required",
			Message: "choose which company you are acting as",
		}
	case containsAny(msg, []string{"suspended in this company"}):
		return &HTTPError{
			Status:  http.StatusForbidden,
			Code:    "company_suspended",
			Message: "your access to this company is suspended",
		}
	case containsAny(msg, []string{"access denied to this company"}):
		return &HTTPError{
			Status:  http.StatusForbidden,
			Code:    "company_forbidden",
			Message: "you do not have access to this company",
		}
	case containsAny(msg, []string{"also works for another company"}):
		return &HTTPError{
			Status:  http.StatusForbidden,
			Code:    "shared_identity",
			Message: "this person also works for another company",
		}
	case containsAny(msg, []string{"insufficient scopes", "permission denied", "access denied"}):
		return &HTTPError{
			Status:  http.StatusForbidden,
			Code:    "forbidden",
			Message: "insufficient permissions",
		}
	case containsAny(msg, []string{"payment does not belong to the user"}):
		return &HTTPError{
			Status:  http.StatusForbidden,
			Code:    "forbidden",
			Message: "insufficient permissions",
		}
	case containsAny(msg, []string{"payment is not cancellable"}):
		return &HTTPError{
			Status:  http.StatusConflict,
			Code:    "payment_not_cancellable",
			Message: "payment is not cancellable",
		}
	// Coupons are classified before the generic rules below on purpose:
	// "coupon not found" would otherwise be flattened into the catch-all 404,
	// and the checkout could not tell an unknown code from an expired campaign
	// or from one whose uses ran out. Each of those is a different sentence to
	// the buyer, so each gets its own code.
	case containsAny(msg, []string{"coupon expired"}):
		return &HTTPError{
			Status:  http.StatusGone,
			Code:    "coupon_expired",
			Message: "coupon expired",
		}
	case containsAny(msg, []string{"coupon max uses reached"}):
		return &HTTPError{
			Status:  http.StatusConflict,
			Code:    "coupon_exhausted",
			Message: "coupon max uses reached",
		}
	case containsAny(msg, []string{"coupon not found", "coupon is invalid"}):
		return &HTTPError{
			Status:  http.StatusNotFound,
			Code:    "coupon_not_found",
			Message: "coupon not found",
		}
	// A rejected upload is the caller's file, not our failure, and the three
	// reasons are three different fixes: send a smaller file, send a real one,
	// send a PDF. Flattening them into 500 would tell the user to try again,
	// which never works.
	case containsAny(msg, []string{"certificate file is too large"}):
		return &HTTPError{
			Status:  http.StatusRequestEntityTooLarge,
			Code:    "file_too_large",
			Message: "certificate file is too large",
		}
	case containsAny(msg, []string{"certificate file type is not allowed"}):
		return &HTTPError{
			Status:  http.StatusUnsupportedMediaType,
			Code:    "file_type_not_allowed",
			Message: "certificate file must be a PDF, PNG or JPEG",
		}
	case containsAny(msg, []string{"certificate file is empty"}):
		return &HTTPError{
			Status:  http.StatusBadRequest,
			Code:    "file_empty",
			Message: "certificate file is empty",
		}
	case containsAny(msg, []string{"limit reached"}):
		return &HTTPError{
			Status:  http.StatusConflict,
			Code:    "limit_reached",
			Message: "limit reached",
		}
	case containsAny(msg, []string{"already exists", "duplicate"}):
		return &HTTPError{
			Status:  http.StatusConflict,
			Code:    "already_exists",
			Message: "resource already exists",
		}
	case containsAny(msg, []string{"not found", "does not exist", "no documents"}):
		return &HTTPError{
			Status:  http.StatusNotFound,
			Code:    "not_found",
			Message: "resource not found",
		}
	}

	return nil
}

// malformedPayload reports whether err came from reading the request body or a
// query parameter, which the transport hands over as a plain decoding error.
func malformedPayload(err error) bool {
	var (
		syntaxErr *json.SyntaxError
		typeErr   *json.UnmarshalTypeError
		numErr    *strconv.NumError
	)

	return errors.As(err, &syntaxErr) ||
		errors.As(err, &typeErr) ||
		errors.As(err, &numErr) ||
		errors.Is(err, io.EOF) ||
		errors.Is(err, io.ErrUnexpectedEOF)
}

// safeMessage decides whether a downstream message may reach the client: 5xx
// never carries detail, and anything shaped like infrastructure output falls
// back to the generic text for the status.
func safeMessage(status int, message string) string {
	message = strings.TrimSpace(message)

	if message == "" || status >= http.StatusInternalServerError || leaksInternals(message) {
		return statusMessage(status)
	}

	return message
}

// maxMessageLen is the length past which a message stops looking like a
// business rule and starts looking like a dump.
const maxMessageLen = 180

func leaksInternals(message string) bool {
	if len(message) > maxMessageLen || strings.ContainsAny(message, "\n\t") {
		return true
	}

	return containsAny(strings.ToLower(message), internalHints)
}

// internalHints are fragments that only ever come from storage engines,
// infrastructure or the Go runtime — never from a business rule.
var internalHints = []string{
	"mongo", "sql", "postgres", "pq:", "pgx", "redis", "kafka", "elastic",
	"collection", "database", "objectid", "duplicate key", "index out of range",
	"dial tcp", "dial udp", "connection refused", "connection reset", "no such host",
	"i/o timeout", "transport:", "rpc error", "grpc", "x509", "certificate",
	"localhost", "127.0.0.1", ".svc", "http://", "https://",
	"panic", "goroutine", "runtime error", "nil pointer", "stack", "syscall",
	".go:", "0x",
}

var statusMessages = map[int]string{
	http.StatusBadRequest:          "invalid request",
	http.StatusUnauthorized:        "unauthorized",
	http.StatusForbidden:           "insufficient permissions",
	http.StatusNotFound:            "resource not found",
	http.StatusRequestTimeout:      "request timeout",
	http.StatusConflict:            "resource already exists",
	http.StatusTooManyRequests:     "too many requests, try again later",
	http.StatusInternalServerError: "internal server error",
	http.StatusNotImplemented:      "not implemented",
	http.StatusServiceUnavailable:  "service temporarily unavailable",
	http.StatusGatewayTimeout:      "upstream service timeout",
}

var statusCodes = map[int]string{
	http.StatusBadRequest:          "bad_request",
	http.StatusUnauthorized:        "unauthorized",
	http.StatusForbidden:           "forbidden",
	http.StatusNotFound:            "not_found",
	http.StatusRequestTimeout:      "request_timeout",
	http.StatusConflict:            "conflict",
	http.StatusTooManyRequests:     "too_many_requests",
	http.StatusInternalServerError: "internal_error",
	http.StatusNotImplemented:      "not_implemented",
	http.StatusServiceUnavailable:  "service_unavailable",
	http.StatusGatewayTimeout:      "gateway_timeout",
}

func statusMessage(status int) string {
	if msg, ok := statusMessages[status]; ok {
		return msg
	}

	return statusMessages[http.StatusInternalServerError]
}

func statusCode(status int) string {
	if code, ok := statusCodes[status]; ok {
		return code
	}

	return statusCodes[http.StatusInternalServerError]
}

func grpcToHTTPStatus(code codes.Code) int {
	switch code {
	case codes.OK:
		return http.StatusOK
	case codes.Canceled:
		return http.StatusRequestTimeout
	case codes.Unknown:
		return http.StatusInternalServerError
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.DeadlineExceeded:
		return http.StatusGatewayTimeout
	case codes.NotFound:
		return http.StatusNotFound
	case codes.AlreadyExists:
		return http.StatusConflict
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.ResourceExhausted:
		return http.StatusTooManyRequests
	case codes.FailedPrecondition:
		return http.StatusBadRequest
	case codes.Aborted:
		return http.StatusConflict
	case codes.OutOfRange:
		return http.StatusBadRequest
	case codes.Unimplemented:
		return http.StatusNotImplemented
	case codes.Internal:
		return http.StatusInternalServerError
	case codes.Unavailable:
		return http.StatusServiceUnavailable
	case codes.DataLoss:
		return http.StatusInternalServerError
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}

func cleanGRPCMessage(msg string) string {
	if idx := strings.LastIndex(msg, "desc = "); idx != -1 {
		return msg[idx+7:]
	}
	return msg
}

func containsAny(s string, substrs []string) bool {
	for _, substr := range substrs {
		if strings.Contains(s, substr) {
			return true
		}
	}
	return false
}
