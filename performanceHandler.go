/*
 Copyright (C) 2017 Ulbora Labs Inc. (www.ulboralabs.com)
 All rights reserved.

 Copyright (C) 2017 Ken Williamson
 All rights reserved.

 Certain inventions and disclosures in this file may be claimed within
 patents owned or patent applications filed by Ulbora Labs Inc., or third
 parties.

 This program is free software: you can redistribute it and/or modify
 it under the terms of the GNU Affero General Public License as published
 by the Free Software Foundation, either version 3 of the License, or
 (at your option) any later version.

 This program is distributed in the hope that it will be useful,
 but WITHOUT ANY WARRANTY; without even the implied warranty of
 MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 GNU Affero General Public License for more details.

 You should have received a copy of the GNU Affero General Public License
 along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package main

import (
	gwmon "UlboraApiGateway/monitor"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	uoauth "github.com/Ulbora/go-ulbora-oauth2"
)

func handlePeformanceSuper(w http.ResponseWriter, r *http.Request) {
	auth := getAuth(r)
	me := new(uoauth.Claim)
	me.Role = "superAdmin"
	me.Scope = "read"
	w.Header().Set("Content-Type", "application/json")
	cType := r.Header.Get("Content-Type")
	if cType != "application/json" {
		http.Error(w, "json required", http.StatusUnsupportedMediaType)
	} else {
		switch r.Method {
		case "POST":
			me.URI = "/ulbora/rs/gwPerformanceSuper"
			valid := auth.Authorize(me)
			if valid != true {
				w.WriteHeader(http.StatusUnauthorized)
			} else {
				p := new(gwmon.GwPerformance)
				decoder := json.NewDecoder(r.Body)
				error := decoder.Decode(&p)
				if error != nil {
					log.Println(error.Error())
					http.Error(w, error.Error(), http.StatusBadRequest)
				} else if p.ClientID == 0 || p.RestRouteID == 0 || p.RouteURIID == 0 {
					http.Error(w, "bad request", http.StatusBadRequest)
				} else {
					resOut := monDB.GetRoutePerformance(p)
					//fmt.Print("response: ")
					//fmt.Println(resOut)
					resJSON, err := json.Marshal(resOut)
					if err != nil {
						log.Println(error.Error())
						http.Error(w, "json output failed", http.StatusInternalServerError)
					}
					w.WriteHeader(http.StatusOK)
					//fmt.Fprint(w, string(resJSON))
					if string(resJSON) == "null" {
						fmt.Fprint(w, "[]")
					} else {
						fmt.Fprint(w, string(resJSON))
					}
				}
			}
		}
	}
}

func handlePeformance(w http.ResponseWriter, r *http.Request) {
	auth := getAuth(r)
	me := new(uoauth.Claim)
	me.Role = "admin"
	me.Scope = "read"
	w.Header().Set("Content-Type", "application/json")
	cType := r.Header.Get("Content-Type")
	if cType != "application/json" {
		http.Error(w, "json required", http.StatusUnsupportedMediaType)
	} else {
		switch r.Method {
		case "POST":
			me.URI = "/ulbora/rs/gwPerformance"
			valid := auth.Authorize(me)
			if valid != true {
				w.WriteHeader(http.StatusUnauthorized)
			} else {
				p := new(gwmon.GwPerformance)
				decoder := json.NewDecoder(r.Body)
				error := decoder.Decode(&p)
				if error != nil {
					log.Println(error.Error())
					http.Error(w, error.Error(), http.StatusBadRequest)
				} else if p.RestRouteID == 0 || p.RouteURIID == 0 {
					http.Error(w, "bad request", http.StatusBadRequest)
				} else {
					p.ClientID = auth.ClientID
					resOut := monDB.GetRoutePerformance(p)
					//fmt.Print("response: ")
					//fmt.Println(resOut)
					resJSON, err := json.Marshal(resOut)
					if err != nil {
						log.Println(error.Error())
						http.Error(w, "json output failed", http.StatusInternalServerError)
					}
					w.WriteHeader(http.StatusOK)
					//fmt.Fprint(w, string(resJSON))
					if string(resJSON) == "null" {
						fmt.Fprint(w, "[]")
					} else {
						fmt.Fprint(w, string(resJSON))
					}
				}
			}
		}
	}
}
