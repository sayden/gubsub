/*
 * Copyright (c) 2013-2014, Jeremy Bingham (<jbingham@gmail.com>)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Package serfin bundles up serf functions for goiardi.
package serfin

import (
	"encoding/json"
	"os"

	serfclient "github.com/hashicorp/serf/client"
	log "github.com/Sirupsen/logrus"
)

// Serfer is the common serf client for gubsub.
var Serfer *serfclient.RPCClient

// StartSerfin sets up the serf instance and starts listening for events and
// queries from other serf instances.
func StartSerfin() error {
	var err error
	Serfer, err = serfclient.NewRPCClient(Conf.SerfAddr)
	if err != nil {
		log.Fatal(err.Error())
		os.Exit(1)
	}

	if Conf.SerfEventAnnounce {
		err = Serfer.UserEvent("gubsub-join", []byte(Conf.Hostname), true)
		if err != nil {
			log.Fatal(err.Error())
			os.Exit(1)
		}
	}

	return nil
}

// SendEvent sends a serf event out from gubsub.
func SendEvent(eventName string, payload interface{}) {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Error(err.Error())
		return
	}
	err = Serfer.UserEvent(eventName, jsonPayload, true)
	if err != nil {
		log.Debug(err.Error())
	}
	return
}

// SendQuery sends a basic, no frills query out over serf.
func SendQuery(queryName string, payload interface{}) {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Errorf(err.Error())
		return
	}
	q := &serfclient.QueryParam{Name: queryName, Payload: jsonPayload}
	err = Serfer.Query(q)
	if err != nil {
		log.Debug(err.Error())
	}
	return
}