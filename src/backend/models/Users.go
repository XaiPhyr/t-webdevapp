package models

import (
	"database/sql"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/driver/pgdriver"
)

type (
	UserResults struct {
		Count int      `json:"count"`
		Data  *[]Users `json:"data"`
	}

	Users struct {
		bun.BaseModel `bun:"table:users,alias:u"`

		ID        int64     `bun:"id,pk,autoincrement" json:"id"`
		Username  string    `bun:"username" json:"username"`
		Password  string    `bun:"password" json:"password"`
		UserType  string    `bun:"user_type" json:"user_type"`
		Email     string    `bun:"email" json:"email"`
		UUID      string    `bun:"uuid" json:"uuid"`
		Status    string    `bun:"status,default:O" json:"status"`
		Active    bool      `bun:"active" json:"active"`
		LastLogin time.Time `bun:"last_login" json:"last_login"`
		CreatedAt time.Time `bun:"created_at,default:current_timestamp" json:"created_at"`
		UpdatedAt time.Time `bun:"updated_at,default:current_timestamp" json:"updated_at"`
		DeletedAt time.Time `bun:"deleted_at,soft_delete,nullzero" json:"deleted_at"`
	}
)

func (u *Users) NewRegister() *Users {
	return &Users{
		UUID:   uuid.Must(uuid.NewRandom()).String(),
		Status: "O",
		Active: false,
	}
}

func (u *Users) ReadAllUsers(limit, offset int) (data *[]Users, count int, err error) {
	data = new([]Users)

	q, ctx := Read(data)
	q = q.ExcludeColumn("password")
	q = q.Limit(limit)
	q = q.Offset(offset)
	q = q.Order("created_at DESC")

	count, err = q.ScanAndCount(ctx)
	return
}

func (u *Users) ReadOneUser(uuid string) (data *Users, err error) {
	data = new(Users)

	excludeCols := []string{"password"}

	q, ctx := Read(data)
	q = q.ExcludeColumn(excludeCols...)
	q = q.Where("uuid = ?", uuid)

	err = q.Scan(ctx)
	return
}

func (u *Users) CreateUser(w http.ResponseWriter, data *Users, fn func(w http.ResponseWriter, code int, message string)) (sql.Result, error) {
	q, ctx := Create(data)

	result, err := q.Exec(ctx)

	if err != nil {
		u.HandleUserError(w, err, fn)
		return nil, err
	}

	return result, nil
}

func (u *Users) UpdateUser(w http.ResponseWriter, data *Users, fn func(w http.ResponseWriter, code int, message string)) (result sql.Result, err error) {
	excludeCols := []string{"password", "created_at"}

	q, ctx := Update(data)
	q = q.ExcludeColumn(excludeCols...)
	q = q.WherePK()

	result, err = q.Exec(ctx)

	if err != nil {
		u.HandleUserError(w, err, fn)
		return nil, err
	}

	return
}

func (u *Users) DeleteAllUsers(w http.ResponseWriter, fn func(w http.ResponseWriter, code int, message string)) (result sql.Result, err error) {
	return
}

func (u *Users) DeleteUser(w http.ResponseWriter, uuid string, fn func(w http.ResponseWriter, code int, message string)) (result sql.Result, err error) {
	data := new(Users)

	q, ctx := Delete(data)
	q = q.Where("uuid = ?", uuid)

	result, err = q.Exec(ctx)

	if err != nil {
		u.HandleUserError(w, err, fn)
		return nil, err
	}

	return
}

func (u *Users) ReadByUsernameOrEmail(user string) (data *Users, err error) {
	data = new(Users)

	q, ctx := Read(data)
	q = q.Where("username = ?", user)
	q = q.WhereOr("email = ?", user)

	err = q.Scan(ctx)

	if err != nil {
		log.Printf("Error: %s", err)
	}

	return
}

func (u *Users) HandleUserError(w http.ResponseWriter, err error, fn func(w http.ResponseWriter, code int, message string)) {
	log.Printf("Handle User Error: %s", err)

	errObj := ErrorObject{
		Code:    400,
		Message: "Bad Request",
	}

	isPgErr := reflect.TypeOf(err).String() == "pgdriver.Error"

	if err != nil {
		if isPgErr {
			pgErr := err.(pgdriver.Error)

			if pgErr.IntegrityViolation() {
				log.Printf("pgErr integrity violation: %s", pgErr.Field('n'))

				switch pgErr.Field('n') {
				case "username_unique_idx":
					errObj = ErrorObject{
						Code:    http.StatusBadRequest,
						Message: "Username already taken",
					}

				case "username_alphanumeric_check":
					errObj = ErrorObject{
						Code:    http.StatusBadRequest,
						Message: "Username must be aplhanumeric",
					}
				}
			}
		}
	}

	fn(w, errObj.Code, errObj.Message)
}
