package main

import (
	"log"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// HTTP code used: 200, 201, 400, 401, 403, 404, 409, 500, 501
var gRPCCodesMap = map[codes.Code]int{
	codes.OK:                 http.StatusOK,
	codes.Canceled:           http.StatusInternalServerError,
	codes.Unknown:            http.StatusInternalServerError,
	codes.InvalidArgument:    http.StatusBadRequest,
	codes.DeadlineExceeded:   http.StatusInternalServerError,
	codes.NotFound:           http.StatusNotFound,
	codes.AlreadyExists:      http.StatusConflict,
	codes.PermissionDenied:   http.StatusUnauthorized,
	codes.Unauthenticated:    http.StatusUnauthorized,
	codes.ResourceExhausted:  http.StatusForbidden,
	codes.FailedPrecondition: http.StatusForbidden,
	codes.Aborted:            http.StatusInternalServerError,
	codes.OutOfRange:         http.StatusInternalServerError,
	codes.Unimplemented:      http.StatusNotImplemented,
	codes.Internal:           http.StatusInternalServerError,
	codes.Unavailable:        http.StatusNotFound,
	codes.DataLoss:           http.StatusNotFound,
}

func logAndSetHTTPErrorCode(w http.ResponseWriter, err error) {
	log.Println(err)
	HTTPCode := gRPCCodesMap[grpc.Code(err)]
	http.Error(w, "", HTTPCode)
}
