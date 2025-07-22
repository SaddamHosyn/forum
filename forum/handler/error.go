package handler

import (
	"log"
	"net/http"
	"text/template"
)

func ErrorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)

	t, err := template.ParseFiles("./assets/html/error.html")
	if err != nil {
		log.Printf("Failed to parse error.html template: %v", err)
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	var p Text

	switch status {
	case http.StatusInternalServerError:
		p = Text{
			ErrorNum: status,
			ErrorMes: "HTTP status 500: Internal Server Error",
		}

	case http.StatusBadRequest:
		p = Text{
			ErrorNum: status,
			ErrorMes: "HTTP status 400: Bad Request",
		}

	case http.StatusForbidden:
		p = Text{
			ErrorNum: status,
			ErrorMes: "HTTP status 403: Forbidden\nYou don't have permission to access this resource",
		}
	case http.StatusNotFound:
		p = Text{
			ErrorNum: status,
			ErrorMes: "HTTP status 404: Not Found",
		}
	case http.StatusUnauthorized:
		p = Text{
			ErrorNum: status,
			ErrorMes: "HTTP status 401: Unauthorized",
		}
	case http.StatusMethodNotAllowed:
		p = Text{
			ErrorNum: status,
			ErrorMes: "HTTP status 405: Method Not Allowed",
		}
	case http.StatusConflict:
		p = Text{
			ErrorNum: status,
			ErrorMes: "HTTP status 409: Conflict\nUser already exists",
		}

	default:
		status = http.StatusNotFound
		p = Text{
			ErrorNum: status,
			ErrorMes: "HTTP status 404: Page Not Found",
		}
	}

	if err := t.Execute(w, p); err != nil {
		log.Printf("Failed to execute error template: %v", err)
		ErrorHandler(w, r, http.StatusInternalServerError)
	}
}
