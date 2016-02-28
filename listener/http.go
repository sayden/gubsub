package listener

import (
	"net/http"
	"net/url"

	"bytes"
	log "github.com/sayden/gubsub/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/sayden/gubsub/Godeps/_workspace/src/github.com/satori/go.uuid"
	"github.com/sayden/gubsub/dispatcher"
	"github.com/sayden/gubsub/types"
	"io/ioutil"
)

//NewHTTPGETListener must be called to create a new listener that will execute
//an HTTP method to a given URL
func NewHTTPListener(msg types.HTTPListener, endpoint *string, id *uuid.UUID) {
	_, err := url.Parse(msg.TargetURL)
	if err != nil {
		log.WithFields(log.Fields{
			"TargetURL": msg.TargetURL,
		}).Error("Couldn't parse TargetURL")
		return
	}

	//Creates new listener
	l := types.NewListener(endpoint)

	*id = l.ID

	dispatcher.AddListenerToTopic(l, endpoint)

	go launchHTTPResponder(l, &msg)

}

func launchHTTPResponder(l *types.Listener, msg *types.HTTPListener) {
	r, err := http.NewRequest(msg.Method, msg.TargetURL, nil)
	if err != nil {
		log.WithFields(log.Fields{
			"Method":    msg.Method,
			"TargetURL": msg.TargetURL,
		}).Error("Could not create HTTP client with specified parameters", err)
	}
	h := http.Header{}
	h.Add("Content-Type", "application/json")
	r.Header = h

	for {
		m := <-l.Ch
		r.Body = ioutil.NopCloser(bytes.NewReader(*m.Data))
		client := http.Client{}
		resp, err := client.Do(r)
		if err != nil || resp.StatusCode != 200 {
			log.WithFields(log.Fields{
				"TargetURL": msg.TargetURL,
				"Data":      string(*m.Data),
				"Topic":     *m.Topic,
				"Timestamp": m.Timestamp,
				"Method":    msg.Method,
			}).Error("Failed to notify listener", err)
		}
	}
}
