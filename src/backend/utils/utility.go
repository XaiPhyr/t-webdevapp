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

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"gopkg.in/gomail.v2"
	"gopkg.in/yaml.v3"

	"github.com/uptrace/bun/extra/bundebug"
)

type (
	Config struct {
		Server   ServerConfig   `yaml:"server"`
		Database DatabaseConfig `yaml:"database"`
		Frontend FrontendConfig `yaml:"frontend"`
		Env      string         `yaml:"env"`
		SMTP     SMTPConfig     `yaml:"smtp"`
	}

	ServerConfig struct {
		Endpoint string `yaml:"endpoint"`
		JwtKey   string `yaml:"jwt_key"`
	}

	DatabaseConfig struct {
		Host    string `yaml:"host"`
		User    string `yaml:"user"`
		Pass    string `yaml:"pass"`
		Port    string `yaml:"port"`
		DB      string `yaml:"db"`
		SSLMode string `yaml:"sslmode"`
	}

	FrontendConfig struct {
		Source string `yaml:"src"`
	}

	SMTPConfig struct {
		Host string `yaml:"host"`
		User string `yaml:"user"`
		Pass string `yaml:"pass"`
		Port int    `yaml:"port"`
	}
)

func InitConfig() Config {
	var cfg Config

	env := os.Getenv("APP_ENVIRONMENT")

	rootFile := "./conf/config.yml"

	if env == "test" {
		rootFile = "../conf/config.yml"
	}

	f, err := os.Open(rootFile)
	if err != nil {
		log.Printf("Error: %s", err)
	}

	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)

	if err != nil {
		log.Printf("Error: %s", err)
	}

	return cfg
}

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
	cfg := InitConfig()

	username := cfg.Database.User
	password := cfg.Database.Pass
	host := cfg.Database.Host
	port := cfg.Database.Port
	database := cfg.Database.DB
	sslmode := cfg.Database.SSLMode

	dsn := "postgres://" + username + ":" + password + "@" + host + ":" + port + "/" + database + "?sslmode=" + sslmode
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New(), bun.WithDiscardUnknownColumns())

	if cfg.Env == "dev" {
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
	cfg := InitConfig()
	host := cfg.SMTP.Host
	port := cfg.SMTP.Port
	username := cfg.SMTP.User
	password := cfg.SMTP.Pass

	m := gomail.NewMessage(gomail.SetCharset("UTF-8"))
	m.SetHeader("From", "noreply-rdev@local")
	m.SetHeader("To", "arcadia.initiative+t_webdevapp@gmail.com")
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	m.AddAlternative("text/html", body)
	// m.Attach("")

	d := gomail.NewDialer(host, port, username, password)

	if err := d.DialAndSend(m); err != nil {
		log.Printf("Error: %s", err)
	}
}
