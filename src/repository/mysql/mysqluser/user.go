package mysqluser

import (
	"context"
	"database/sql"
	"errors"

	"github.com/mohsenHa/messenger/entity"
	"github.com/mohsenHa/messenger/pkg/errmsg"
	"github.com/mohsenHa/messenger/pkg/richerror"
	"github.com/mohsenHa/messenger/repository/mysql"
)

func (d *DB) Register(ctx context.Context, u entity.User) (entity.User, error) {
	const op = "mysql.Register"

	_, err := d.conn.Conn().ExecContext(ctx, `insert into users(id,public_key,code,status) values(?,?,?,?)`,
		u.ID, u.PublicKey, u.Code, u.Status)
	if err != nil {
		return entity.User{}, richerror.New(op).WithErr(err).
			WithMessage(errmsg.ErrorMsgSomethingWentWrong).WithKind(richerror.KindUnexpected)
	}

	return u, nil
}

func (d *DB) GetUserByID(ctx context.Context, id string) (entity.User, error) {
	const op = "mysql.GetUserByPublicKey"

	row := d.conn.Conn().QueryRowContext(ctx, `select * from users where id = ?`, id)

	user, err := scanUser(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.User{}, richerror.New(op).WithErr(err).
				WithMessage(errmsg.ErrorMsgNotFound).WithKind(richerror.KindNotFound)
		}

		return entity.User{}, richerror.New(op).WithErr(err).
			WithMessage(errmsg.ErrorMsgCantScanQueryResult).WithKind(richerror.KindUnexpected)
	}

	return user, nil
}

func (d *DB) UpdateCode(ctx context.Context, id, code string) error {
	const op = "mysql.Activate"

	_, err := d.conn.Conn().ExecContext(ctx, `update users set code=? where id=?`,
		code, id)
	if err != nil {
		return richerror.New(op).WithErr(err).
			WithMessage(errmsg.ErrorMsgSomethingWentWrong).WithKind(richerror.KindUnexpected)
	}

	return nil
}

func (d *DB) Activate(ctx context.Context, id string) error {
	const op = "mysql.Activate"

	_, err := d.conn.Conn().ExecContext(ctx, `update users set status=1 where id=?`,
		id)
	if err != nil {
		return richerror.New(op).WithErr(err).
			WithMessage(errmsg.ErrorMsgSomethingWentWrong).WithKind(richerror.KindUnexpected)
	}

	return nil
}

func (d *DB) IsIDUnique(id string) (bool, error) {
	const op = "mysql.IsPhoneNumberUnique"

	row := d.conn.Conn().QueryRow(`select * from users where id = ?`, id)

	_, err := scanUser(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return true, nil
		}

		return false, richerror.New(op).WithErr(err).
			WithMessage(errmsg.ErrorMsgCantScanQueryResult).WithKind(richerror.KindUnexpected)
	}

	return false, nil
}

func (d *DB) IsIDExist(id string) (bool, error) {
	unique, err := d.IsIDUnique(id)

	return !unique, err
}

func scanUser(scanner mysql.Scanner) (entity.User, error) {
	var user entity.User
	err := scanner.Scan(&user.ID, &user.PublicKey, &user.Code, &user.Status)

	return user, err
}
