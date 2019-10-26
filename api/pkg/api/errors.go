package api

import (
	"log"
	"net/http"

	"github.com/go-chi/render"
)

type errorRenderer func(error) render.Renderer

func HandleApiErr(err error, renderer errorRenderer, w http.ResponseWriter, r *http.Request) {
	if err != nil {
		log.Println(err)
		err = render.Render(w, r, renderer(err))
		if err != nil {
			log.Fatalf("render of error response failed: %s", err)
		}
	}
}

type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrTokenCreation(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 500,
		StatusText:     "Token creation error",
		ErrorText:      err.Error(),
	}
}

func ErrAuthenication(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 401,
		StatusText:     "Authentication error",
		ErrorText:      err.Error(),
	}
}

func ErrDatabase(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 500,
		StatusText:     "Database error",
		ErrorText:      err.Error(),
	}
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}

func ErrNotFound(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 404,
		StatusText:     "Resource not found.",
		ErrorText:      err.Error(),
	}
}

func ErrRender(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 422,
		StatusText:     "Error rendering response.",
		ErrorText:      err.Error(),
	}
}
