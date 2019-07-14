package helpers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

// ResponseBody struct
type ResponseBody struct {
	StatusCode    int         `json:"status_code"`
	StatusMessage string      `json:"status_message"`
	Description   string      `json:"description"`
	Count         int64       `json:"count"`
	Offset        int64       `json:"offset"`
	Limit         int64       `json:"limit"`
	Href          string      `json:"href"`
	Payload       interface{} `json:"payload"`
}

// Response struct
type Response struct {
	Body ResponseBody
	Err  error
}

// ServeJSON serve json with container Response
func (c *Response) ServeJSON(w http.ResponseWriter, r *http.Request) {

	defer func() {
		b, err := json.Marshal(c.Body)
		if err != nil {
			log.Printf("helpers: could not json marshal: %s", err.Error())
		}
		_, err = w.Write(b)
		if err != nil {
			log.Printf("helpers: could not write: %s", err.Error())
		}
	}()

	w.Header().Add("Content-Type", "application/json")
	c.Body.Href = r.RequestURI

	if c.Err != nil {
		c.Body.StatusMessage = "Error"
		c.Body.Description = c.Err.Error()
		c.Body.StatusCode = 400
		w.WriteHeader(400)
	} else {
		c.Body.StatusMessage = "Success"
		c.Body.StatusCode = 200
		c.Body.Limit = 10

		if v := r.URL.Query().Get("offset"); v != "" {
			vInt, err := strconv.Atoi(v)
			if err != nil {
				log.Println(err)
			}
			c.Body.Offset = int64(vInt)
		}
		if v := r.URL.Query().Get("limit"); v != "" {
			vInt, err := strconv.Atoi(v)
			if err != nil {
				log.Println(err)
			}
			c.Body.Limit = int64(vInt)
		}
	}
}
