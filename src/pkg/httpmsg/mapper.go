package httpmsg

import (
	"github.com/mohsenHa/messenger/logger"
	"github.com/mohsenHa/messenger/logger/loggerentity"
	"github.com/mohsenHa/messenger/pkg/errmsg"
	"github.com/mohsenHa/messenger/pkg/richerror"
	"net/http"
)

func Error(err error) (message string, code int) {
	switch err.(type) {
	case richerror.RichError:
		re := err.(richerror.RichError)
		msg := re.Message()

		code := mapKindToHTTPStatusCode(re.Kind())

		// we should not expose unexpected error messages
		if code >= 500 {
			logger.NewLog(msg).
				WithCategory(loggerentity.CategoryRequestResponse).
				WithSubCategory(loggerentity.SubCategoryInternalResponse).
				With(loggerentity.ExtraKeyErrorMessage, err.Error()).
				Error()
			msg = errmsg.ErrorMsgSomethingWentWrong
		}

		return msg, code
	default:
		return err.Error(), http.StatusBadRequest
	}
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
