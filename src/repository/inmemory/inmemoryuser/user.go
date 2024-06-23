package inmemoryuser

import (
	"context"
	"github.com/mohsenHa/messenger/entity"
	"github.com/mohsenHa/messenger/pkg/errmsg"
	"github.com/mohsenHa/messenger/pkg/richerror"
	"time"
)

func (d *DB) Register(_ context.Context, u entity.User) (entity.User, error) {
	d.table.Store(u.ID, u)

	return u, nil
}

func (d *DB) GetUserByID(_ context.Context, id string) (entity.User, error) {
	const op = "mysql.GetUserByPublicKey"

	u, ok := d.table.Load(id)
	if !ok {
		return entity.User{}, richerror.New(op).
			WithMessage(errmsg.ErrorMsgNotFound).WithKind(richerror.KindNotFound)
	}
	if user, ok := u.(entity.User); ok {
		return user, nil
	}

	return entity.User{}, richerror.New(op).
		WithMessage(errmsg.ErrorMsgNotFound).WithKind(richerror.KindNotFound)
}

func (d *DB) UpdateCode(ctx context.Context, id, code string) error {
	const op = "mysql.Activate"
	u, err := d.GetUserByID(ctx, id)
	if err != nil {
		return richerror.New(op).WithErr(err).
			WithMessage(errmsg.ErrorMsgSomethingWentWrong).WithKind(richerror.KindUnexpected)
	}
	u.Code = code
	_, err = d.Register(ctx, u)
	if err != nil {
		return richerror.New(op).WithErr(err).
			WithMessage(errmsg.ErrorMsgSomethingWentWrong).WithKind(richerror.KindUnexpected)
	}

	return nil
}

func (d *DB) Activate(ctx context.Context, id string) error {
	const op = "mysql.Activate"

	u, err := d.GetUserByID(ctx, id)
	if err != nil {
		return richerror.New(op).WithErr(err).
			WithMessage(errmsg.ErrorMsgSomethingWentWrong).WithKind(richerror.KindUnexpected)
	}
	u.Status = 1
	_, err = d.Register(ctx, u)
	if err != nil {
		return richerror.New(op).WithErr(err).
			WithMessage(errmsg.ErrorMsgSomethingWentWrong).WithKind(richerror.KindUnexpected)
	}

	return nil
}

func (d *DB) IsIDUnique(id string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	cancel()
	_, err := d.GetUserByID(ctx, id)
	if err != nil {
		return true, nil
	}

	return false, nil
}

func (d *DB) IsIDExist(id string) (bool, error) {
	unique, err := d.IsIDUnique(id)

	return !unique, err
}
