package models

import (
	"context"

	"github.com/go-chi/chi"
	"github.com/uptrace/bun"

	u "t_webdevapp/utils"
)

type (
	MuxServer struct {
		Mux      *chi.Mux
		Endpoint string
	}

	ErrorObject struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
)

func Create(model interface{}) (q *bun.InsertQuery, ctx context.Context) {
	ctx = context.Background()

	db := u.InitDBConnect()
	q = db.NewInsert().Model(model)
	return
}

func Read(model interface{}) (q *bun.SelectQuery, ctx context.Context) {
	ctx = context.Background()

	db := u.InitDBConnect()
	q = db.NewSelect().Model(model)
	return
}

func Update(model interface{}) (q *bun.UpdateQuery, ctx context.Context) {
	ctx = context.Background()

	db := u.InitDBConnect()
	q = db.NewUpdate().Model(model)
	return
}

func Delete(model interface{}) (q *bun.DeleteQuery, ctx context.Context) {
	ctx = context.Background()

	db := u.InitDBConnect()
	q = db.NewDelete().Model(model)
	return
}
