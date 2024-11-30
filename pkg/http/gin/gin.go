package ginhttp

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nitesh237/go-server-template/pkg/errors"
)

type GinHttpRouter interface {
	gin.IRouter
	http.Handler
}

type Endpoint[req, resp any] func(ctx context.Context, req *req) (*resp, error)

func NewGinEndpoint[req, resp any](ep Endpoint[req, resp]) gin.HandlerFunc {
	return func(c *gin.Context) {
		r := new(req)
		if err := c.ShouldBind(&r); err != nil {
			c.JSON(errors.GetHttpCodeFromErrorType(errors.ErrInvalidArgumentStr), errors.NewErrorResponseWithDebug("Invalid Argument", err.Error(), errors.ErrInvalidArgumentStr))
			return
		}

		res, err := ep(c, r)
		if err == nil {
			c.JSON(http.StatusOK, res)
			return
		}

		errResp := errors.ErrorResponse{}
		switch {
		case errors.As(err, &errResp):
			if errResp.Code == 0 {
				errResp.Code = errors.GetErrorCodeForErrorType(errResp.ErrorType)
			}
		case errors.IsRecordNotFound(err):
			errResp = errors.NewErrorResponseWithDebug("Record Not Found", err.Error(), errors.ErrRecordNotFoundStr)
		case errors.Is(err, errors.ErrInvalidArgument):
			errResp = errors.NewErrorResponseWithDebug("Invalid Argument", err.Error(), errors.ErrInvalidArgumentStr)
		case errors.Is(err, errors.ErrAlreadyExists):
			errResp = errors.NewErrorResponseWithDebug("Already Exists", err.Error(), errors.ErrAlreadyExistsStr)
		case errors.Is(err, errors.ErrPermissionDenied):
			errResp = errors.NewErrorResponseWithDebug("Permission Denied", err.Error(), errors.ErrPermissionDeniedStr)
		case errors.Is(err, errors.ErrFailedPrecondition):
			errResp = errors.NewErrorResponseWithDebug("Failed Precondition", err.Error(), errors.ErrFailedPreconditionStr)
		default:
			errResp = errors.NewErrorResponseWithDebug("Internal Server Error", err.Error(), errors.ErrInternalServerStr)
		}

		c.JSON(errors.GetHttpCodeFromErrorType(errResp.ErrorType), errResp)
	}
}
