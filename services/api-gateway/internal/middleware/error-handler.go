package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func HandleError(ctx *gin.Context, err error) {
	st, ok := status.FromError(err)

	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
			"status":  false,
		})
		return
	}

	httpCode := grpcCodeToHTTPStatus(st.Code())

	ctx.JSON(httpCode, gin.H{
		"code":    httpCode,
		"message": st.Message(),
		"status":  false,
	})
}

func grpcCodeToHTTPStatus(code codes.Code) int {
	switch code {
	case codes.AlreadyExists:
		return http.StatusConflict // 409
	case codes.InvalidArgument:
		return http.StatusBadRequest // 400
	case codes.Unauthenticated:
		return http.StatusUnauthorized // 401
	case codes.PermissionDenied:
		return http.StatusForbidden // 403
	case codes.NotFound:
		return http.StatusNotFound // 404
	case codes.DeadlineExceeded:
		return http.StatusRequestTimeout // 408
	case codes.Unavailable:
		return http.StatusServiceUnavailable // 503
	default:
		return http.StatusInternalServerError // 500
	}
}
