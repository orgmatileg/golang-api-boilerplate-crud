package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"golang-api-boilerplate-crud/models"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

// ImgurRepositoryInterface ...
type ImgurRepositoryInterface interface {
	Create(bs string) (url string, err error)
}

// Imgur Repository
type Imgur struct{}

// NewImgurRepository ...
func NewImgurRepository() ImgurRepositoryInterface {
	return Imgur{}
}

// Create ...
func (v Imgur) Create(bs string) (string, error) {
	var imgurlClientID = os.Getenv("IMGUR_CLIENT_ID")
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	// Write image field to form-data
	err := w.WriteField("image", bs)
	if err != nil {
		return "", fmt.Errorf("repository: could not write form field: %s", err.Error())
	}

	defer func() {
		err = w.Close()
		if err != nil {
			log.Printf("repository: could not close writer: %s", err.Error())
		}
	}()

	c := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	url := "https://api.imgur.com/3/image"
	req, err := http.NewRequest(http.MethodPost, url, &buf)
	if err != nil {
		return "", fmt.Errorf("repository: could not make new request: %s", err.Error())
	}

	// Set HTTP header
	req.Header.Add("Authorization", fmt.Sprintf("Client-ID %s", imgurlClientID))
	req.Header.Set("Content-Type", w.FormDataContentType())

	// make request to imgur API
	res, err := c.Do(req)
	if err != nil {
		return "", fmt.Errorf("repository: could not execute request to imgur API: %s", err.Error())
	}
	defer func() {
		err = res.Body.Close()
		if err != nil {
			log.Printf("repository: could not close http resp body: %s", err.Error())
		}
	}()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("repository: could not read response body: %s", err.Error())
	}

	// Unmarshal response body to model
	var i models.Imgur
	err = json.Unmarshal(body, &i)
	if err != nil {
		return "", fmt.Errorf("repository: could not unmarshall body to model: %s", err.Error())
	}

	return i.Data.Link, nil
}
