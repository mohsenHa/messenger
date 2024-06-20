package httpmsg

import (
	"errors"
	"net/http"

	"github.com/mohsenHa/messenger/logger"
	"github.com/mohsenHa/messenger/logger/loggerentity"
	"github.com/mohsenHa/messenger/pkg/errmsg"
	"github.com/mohsenHa/messenger/pkg/richerror"
)

func Error(err error) (message string, code int) {
	re := richerror.RichError{}
	ok := errors.As(err, &re)
	if ok {
		if !ok {
			return err.Error(), http.StatusBadRequest
		}
		msg := re.Message()

		code = mapKindToHTTPStatusCode(re.Kind())

		// we should not expose unexpected error messages
		if code >= http.StatusInternalServerError {
			logger.NewLog(msg).
				WithCategory(loggerentity.CategoryRequestResponse).
				WithSubCategory(loggerentity.SubCategoryInternalResponse).
				With(loggerentity.ExtraKeyErrorMessage, err.Error()).
				Error()
			msg = errmsg.ErrorMsgSomethingWentWrong
		}

		return msg, code
	}

	return err.Error(), http.StatusBadRequest
}

func mapKindToHTTPStatusCode(kind richerror.Kind) int {
	switch kind {
	case richerror.KindInvalid:
		return http.StatusUnprocessableEntity
	case richerror.KindNotFound:
		return http.StatusNotFound
	case richerror.KindForbidden:
		return http.StatusForbidden
	case richerror.KindUnexpected:
		return http.StatusInternalServerError
	default:
		return http.StatusBadRequest
	}
}
