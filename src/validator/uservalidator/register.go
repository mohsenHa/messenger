package uservalidator

import (
	"errors"
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/mohsenHa/messenger/param/userparam"
	"github.com/mohsenHa/messenger/pkg/errmsg"
	"github.com/mohsenHa/messenger/pkg/richerror"
)

func (v Validator) ValidateRegisterRequest(req userparam.RegisterRequest) (map[string]string, error) {
	const op = "uservalidator.ValidateRegisterRequest"

	if err := validation.ValidateStruct(&req,

		validation.Field(&req.PublicKey,
			validation.Required,
			validation.By(v.checkIDUniqueness)),
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

func (v Validator) checkIDUniqueness(value interface{}) error {
	publicKey, ok := value.(string)
	if !ok {
		return fmt.Errorf(errmsg.ErrorMsgInvalidInput)
	}
	id := v.keyGen.CreateUserID(publicKey)

	if isUnique, err := v.repo.IsIDUnique(id); err != nil || !isUnique {
		if err != nil {
			return err
		}

		if !isUnique {
			return fmt.Errorf(errmsg.ErrorMsgPublicKeyIsAlreadyRegisteredTryToLogin)
		}
	}

	return nil
}
