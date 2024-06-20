package messagevalidator

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/mohsenHa/messenger/param/messageparam"
	"github.com/mohsenHa/messenger/pkg/errmsg"
	"github.com/mohsenHa/messenger/pkg/richerror"
)

func (v Validator) ValidateReceiveRequest(req messageparam.ReceiveRequest) (map[string]string, error) {
	const op = "messagevalidator.ValidateReceiveRequest"

	if err := validation.ValidateStruct(&req); err != nil {
		fieldErrors := make(map[string]string)

		errV := validation.Errors{}
		ok := errors.As(err, &errV)
		if ok {
			for key, value := range errV {
				if value != nil {
					fieldErrors[key] = value.Error()
				}
			}
		}

		return fieldErrors, richerror.New(op).WithMessage(errmsg.ErrorMsgInvalidInput).
			WithKind(richerror.KindInvalid).
			WithMeta(map[string]interface{}{"req": req}).WithErr(err)
	}

	return map[string]string{}, nil
}
