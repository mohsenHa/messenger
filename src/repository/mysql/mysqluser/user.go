package mysqluser

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/mohsenHa/messenger/entity"
	"github.com/mohsenHa/messenger/pkg/errmsg"
	"github.com/mohsenHa/messenger/pkg/richerror"
	"github.com/mohsenHa/messenger/repository/mysql"
)

func (d *DB) Register(ctx context.Context, u entity.User) (entity.User, error) {
	const op = "mysql.Register"

	_, err := d.conn.Conn().ExecContext(ctx, `insert into users(id,public_key,active_code,status) values(?,?,?,?)`,
		u.Id, u.PublicKey, u.ActiveCode, u.Status)
	if err != nil {
		fmt.Println(err)
		return entity.User{}, richerror.New(op).WithErr(err).
			WithMessage(errmsg.ErrorMsgSomethingWentWrong).WithKind(richerror.KindUnexpected)
	}

	return u, nil
}

func (d *DB) GetUserById(ctx context.Context, id string) (entity.User, error) {
	const op = "mysql.GetUserByPublicKey"

	row := d.conn.Conn().QueryRowContext(ctx, `select * from users where id = ?`, id)

	user, err := scanUser(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.User{}, richerror.New(op).WithErr(err).
				WithMessage(errmsg.ErrorMsgNotFound).WithKind(richerror.KindNotFound)
		}

		return entity.User{}, richerror.New(op).WithErr(err).
			WithMessage(errmsg.ErrorMsgCantScanQueryResult).WithKind(richerror.KindUnexpected)
	}

	return user, nil
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

func (d *DB) IsIdUnique(id string) (bool, error) {
	const op = "mysql.IsPhoneNumberUnique"

	row := d.conn.Conn().QueryRow(`select * from users where id = ?`, id)

	_, err := scanUser(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return true, nil
		}

		return false, richerror.New(op).WithErr(err).
			WithMessage(errmsg.ErrorMsgCantScanQueryResult).WithKind(richerror.KindUnexpected)
	}

	return false, nil
}

func scanUser(scanner mysql.Scanner) (entity.User, error) {
	var user entity.User
	err := scanner.Scan(&user.Id, &user.PublicKey, &user.ActiveCode, &user.Status)

	return user, err
}
