package utils

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"gopkg.in/gomail.v2"

	"github.com/uptrace/bun/extra/bundebug"
)

var (
	env         = os.Getenv("ENVIRONEMNT")
	dsn         = os.Getenv("DATABASE_URL")
	smtphost    = os.Getenv("SMTP_HOST")
	smtpport, _ = strconv.Atoi(os.Getenv("SMTP_PORT"))
	smtpuser    = os.Getenv("SMTP_USER")
	smtppass    = os.Getenv("SMTP_PASS")
)

func ParseHTML(path string, data interface{}) (body string, err error) {
	t := template.New(filepath.Base(path)).Funcs(template.FuncMap{})
	t, err = t.ParseFiles(path)

	if err != nil {
		fmt.Println("Error loading template", err.Error())
		return "", err
	} else {
		var tpl bytes.Buffer

		if err = t.Execute(&tpl, data); err == nil {
			body = tpl.String()
		}
	}

	return
}

func InitDBConnect() *bun.DB {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New(), bun.WithDiscardUnknownColumns())

	if env == "dev" {
		db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	}

	return db
}

func ListenNotify() {
	db := InitDBConnect()
	ctx := context.Background()

	ln := pgdriver.NewListener(db)

	go func() {
		if err := ln.Listen(ctx, "last_login"); err != nil {
			panic(err)
		}
	}()
}

func Mailer(subject, body string) {
	m := gomail.NewMessage(gomail.SetCharset("UTF-8"))
	m.SetHeader("From", "noreply-rdev@local")
	m.SetHeader("To", "arcadia.initiative+t_webdevapp@gmail.com")
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	m.AddAlternative("text/html", body)
	// m.Attach("")

	fmt.Println(smtphost, smtpport, smtpuser, smtppass)

	d := gomail.NewDialer(smtphost, smtpport, smtpuser, smtppass)

	if err := d.DialAndSend(m); err != nil {
		log.Printf("Error: %s", err)
	}
}
