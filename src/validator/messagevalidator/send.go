package messagevalidator

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/mohsenHa/messenger/param/messageparam"
	"github.com/mohsenHa/messenger/pkg/errmsg"
	"github.com/mohsenHa/messenger/pkg/richerror"
)

func (v Validator) ValidateSendRequest(req messageparam.SendRequest) (map[string]string, error) {
	const op = "messagevalidator.ValidateSendRequest"

	if err := validation.ValidateStruct(&req,

		validation.Field(&req.ToId,
			validation.Required,
			validation.By(v.checkUserExist),
		),
		validation.Field(&req.Message,
			validation.Required,
		),
	); err != nil {
		fieldErrors := make(map[string]string)

		errV, ok := err.(validation.Errors)
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

	return nil, nil
}
func (v Validator) checkUserExist(value interface{}) error {
	id := value.(string)

	if isExist, err := v.urepo.IsIdExist(id); err != nil || !isExist {
		if err != nil {
			return err
		}

		if !isExist {
			return fmt.Errorf(errmsg.ErrorMsgNotFound)
		}
	}

	return nil
}
