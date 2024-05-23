package uservalidator

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/mohsenHa/messenger/param/userparam"
	"github.com/mohsenHa/messenger/pkg/errmsg"
	"github.com/mohsenHa/messenger/pkg/richerror"
)

func (v Validator) ValidateIdRequest(req userparam.IdRequest) (map[string]string, error) {
	const op = "uservalidator.ValidateIdRequest"

	if err := validation.ValidateStruct(&req,

		validation.Field(&req.PublicKey,
			validation.Required,
			validation.By(v.checkPublicKeyExist)),
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

func (v Validator) checkPublicKeyExist(value interface{}) error {
	publicKey := value.(string)
	id := v.keyGen.CreateUserId(publicKey)

	return v.checkIdExist(id)
}
