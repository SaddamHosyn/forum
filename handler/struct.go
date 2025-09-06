package handler

import (
	"database/sql"
	"text/template"
)

type Text struct {
	ErrorNum int
	ErrorMes string
}

var Templates *template.Template

var DB *sql.DB
