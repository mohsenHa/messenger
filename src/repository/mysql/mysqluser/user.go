package mysqluser

import (
	"context"
	"database/sql"
	"github.com/mohsenHa/messenger/entity"
	"github.com/mohsenHa/messenger/pkg/errmsg"
	"github.com/mohsenHa/messenger/pkg/richerror"
	"github.com/mohsenHa/messenger/repository/mysql"
)

func (d *DB) Register(u entity.User) (entity.User, error) {
	const op = "mysql.Register"

	_, err := d.conn.Conn().Exec(`insert into users(public_key) values(?)`,
		u.PublicKey)
	if err != nil {
		return entity.User{}, richerror.New(op).WithErr(err).
			WithMessage(errmsg.ErrorMsgSomethingWentWrong).WithKind(richerror.KindUnexpected)
	}

	return u, nil
}

func (d *DB) GetUserByPublicKey(ctx context.Context, publicKey string) (entity.User, error) {
	const op = "mysql.GetUserByPublicKey"

	row := d.conn.Conn().QueryRowContext(ctx, `select * from users where public_key = ?`, publicKey)

	user, err := scanUser(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.User{}, richerror.New(op).WithErr(err).
				WithMessage(errmsg.ErrorMsgNotFound).WithKind(richerror.KindNotFound)
		}

		// TODO - log unexpected error for better observability
		return entity.User{}, richerror.New(op).WithErr(err).
			WithMessage(errmsg.ErrorMsgCantScanQueryResult).WithKind(richerror.KindUnexpected)
	}

	return user, nil
}

func scanUser(scanner mysql.Scanner) (entity.User, error) {
	var user entity.User

	err := scanner.Scan(&user.PublicKey)

	return user, err
}
