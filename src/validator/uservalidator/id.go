package uservalidator

import (
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/mohsenHa/messenger/param/userparam"
	"github.com/mohsenHa/messenger/pkg/errmsg"
	"github.com/mohsenHa/messenger/pkg/richerror"
)

func (v Validator) ValidateIDRequest(req userparam.IDRequest) (map[string]string, error) {
	const op = "uservalidator.ValidateIDRequest"

	if err := validation.ValidateStruct(&req,

		validation.Field(&req.PublicKey,
			validation.Required,
			validation.By(v.checkPublicKeyExist)),
	); err != nil {
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

func (v Validator) checkPublicKeyExist(value interface{}) error {
	publicKey, ok := value.(string)
	if !ok {
		return fmt.Errorf(errmsg.ErrorMsgInvalidInput)
	}
	id := v.keyGen.CreateUserID(publicKey)

	return v.checkIDExist(id)
}
