package types

import (
	"io/ioutil"
	"net/http"
	"time"
)

type Message struct {
	Data      *[]byte
	Timestamp time.Time
	Topic     *string
}

func (mb *MessageBinding) Name() string {
	return "json"
}

//MessageBinding implements a Gin-Gonic interface to parse bodies in http requests
type MessageBinding struct{}

func (mb *MessageBinding) Bind(r *http.Request, obj interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	if err := r.Body.Close(); err != nil {
		return err
	}

	m := obj.(*Message)

	m.Data = &body

	obj = m

	return nil
}
