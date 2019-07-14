package handlers

import (
	"encoding/json"
	"fmt"
	"golang-api-boilerplate-crud/helpers"
	"golang-api-boilerplate-crud/models"
	"golang-api-boilerplate-crud/usecases"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/getsentry/raven-go"
	"github.com/gorilla/mux"
)

// FruitsHandlersInterface ...
type FruitsHandlersInterface interface {
	Create(w http.ResponseWriter, r *http.Request)
	GetOneByID(w http.ResponseWriter, r *http.Request)
	GetAll(w http.ResponseWriter, r *http.Request)
	UpdateOneByID(w http.ResponseWriter, r *http.Request)
	UpdatePhotoOneByID(w http.ResponseWriter, r *http.Request)
	DeleteOneByID(w http.ResponseWriter, r *http.Request)
}

// Fruits Handlers
type Fruits struct{}

// NewFruitsHandlers ...
func NewFruitsHandlers() FruitsHandlersInterface {
	return Fruits{}
}

// Create ...
func (v Fruits) Create(w http.ResponseWriter, r *http.Request) {
	res := helpers.Response{}
	defer res.ServeJSON(w, r)
	defer func() {
		if res.Err != nil {
			raven.CaptureErrorAndWait(res.Err, nil)
		}
	}()

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		res.Err = fmt.Errorf("handlers: could not read request body: %s", err.Error())
		return
	}
	defer r.Body.Close()

	var p models.Fruits

	err = json.Unmarshal(b, &p)
	if err != nil {
		res.Err = fmt.Errorf("handlers: could not unmarshal request body to model: %s", err.Error())
		return
	}

	err = usecases.NewFruitsUsecase().Create(&p)
	if err != nil {
		res.Err = fmt.Errorf("handlers: %s", err.Error())
		return
	}

	res.Body.Payload = p
}

// GetOneByID ...
func (v Fruits) GetOneByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	res := helpers.Response{}
	defer res.ServeJSON(w, r)
	defer func() {
		if res.Err != nil {
			raven.CaptureErrorAndWait(res.Err, nil)
		}
	}()

	id, err := strconv.Atoi(idStr)
	if err != nil {
		res.Err = fmt.Errorf("handlers: could not convert string to int %s", err.Error())
		return
	}

	var p models.Fruits
	err = usecases.NewFruitsUsecase().GetOneByID(id, &p)
	if err != nil {
		res.Err = fmt.Errorf("handlers: %s", err.Error())
		return
	}
	res.Body.Payload = p
}

// GetAll ...
func (v Fruits) GetAll(w http.ResponseWriter, r *http.Request) {

	queryParam := r.URL.Query()
	// Set default query
	limit, offset, order, search := "10", "0", "DESC", ""
	if v := queryParam.Get("limit"); v != "" {
		limit = queryParam.Get("limit")
	}
	if v := queryParam.Get("offset"); v != "" {
		offset = queryParam.Get("offset")
	}
	if v := queryParam.Get("order"); v != "" {
		order = strings.ToUpper(queryParam.Get("order"))
	}
	if v := queryParam.Get("search"); v != "" {
		search = "%" + queryParam.Get("search") + "%"
	}

	res := helpers.Response{}
	defer res.ServeJSON(w, r)
	defer func() {
		if res.Err != nil {
			raven.CaptureErrorAndWait(res.Err, nil)
		}
	}()

	res.Body.Payload, res.Body.Count, res.Err = usecases.NewFruitsUsecase().GetAll(limit, offset, order, search)

}

// UpdateOneByID ...
func (v Fruits) UpdateOneByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	res := helpers.Response{}
	defer res.ServeJSON(w, r)
	defer func() {
		if res.Err != nil {
			raven.CaptureErrorAndWait(res.Err, nil)
		}
	}()

	id, err := strconv.Atoi(idStr)
	if err != nil {
		res.Err = fmt.Errorf("handlers: could not convert string to int %s", err.Error())
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		res.Err = fmt.Errorf("handlers: could not read request body: %s", err.Error())
		return
	}
	defer r.Body.Close()

	var p models.Fruits
	err = json.Unmarshal(b, &p)
	if err != nil {
		res.Err = fmt.Errorf("handlers: could not unmarshal request body to model: %s", err.Error())
		return
	}

	ra, err := usecases.NewFruitsUsecase().UpdateOneByID(id, &p)
	if err != nil {
		res.Err = fmt.Errorf("handlers: %s", err.Error())
		return
	}
	res.Body.Payload = fmt.Sprintf("row affected: %d", ra)
}

// UpdatePhotoOneByID ...
func (v Fruits) UpdatePhotoOneByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	res := helpers.Response{}
	defer res.ServeJSON(w, r)
	defer func() {
		if res.Err != nil {
			raven.CaptureErrorAndWait(res.Err, nil)
		}
	}()

	id, err := strconv.Atoi(idStr)
	if err != nil {
		res.Err = fmt.Errorf("handlers: could not convert string to int %s", err.Error())
		return
	}

	const MB = 1 << 20 // Megabytes
	// Parse multipartform with 10 MB max allocated
	err = r.ParseMultipartForm(10 * MB)
	if err != nil {
		res.Err = fmt.Errorf("handlers: could not parse form-data: %s", err.Error())
		return
	}

	f, fh, err := r.FormFile("image")
	if err != nil {
		res.Err = fmt.Errorf("handlers: could not get file from form-data: %s", err.Error())
		return
	}
	defer f.Close()

	// Check the content type file is allowed
	ct := fh.Header.Get("Content-Type")
	allowedTypeFile := func() bool {
		for _, v := range []string{"image/png", "image/jpeg", "image/bmp"} {
			if ct == v {
				return true
			}
		}
		return false
	}()

	if !allowedTypeFile {
		res.Err = fmt.Errorf("handlers: content type is not allowed: %s", ct)
		return
	}

	ra, err := usecases.NewFruitsUsecase().UpdatePhotoOneByID(id, &f)
	if err != nil {
		res.Err = fmt.Errorf("handlers: %s", err.Error())
		return
	}
	res.Body.Payload = fmt.Sprintf("row affected: %d", ra)
}

// DeleteOneByID ...
func (v Fruits) DeleteOneByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	res := helpers.Response{}
	defer res.ServeJSON(w, r)
	defer func() {
		if res.Err != nil {
			raven.CaptureErrorAndWait(res.Err, nil)
		}
	}()

	id, err := strconv.Atoi(idStr)
	if err != nil {
		res.Err = fmt.Errorf("handlers: could not convert string to int %s", err.Error())
		return
	}
	defer r.Body.Close()

	ra, err := usecases.NewFruitsUsecase().DeleteOneByID(id)
	if err != nil {
		res.Err = fmt.Errorf("handlers: %s", err.Error())
		return
	}
	res.Body.Payload = fmt.Sprintf("row affected: %d", ra)
}
