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

//build command "go build -o main *.go"
import (
	cb "UlboraApiGateway/circuitbreaker"
	gwerr "UlboraApiGateway/gwerrors"
	mgr "UlboraApiGateway/managers"
	gwmon "UlboraApiGateway/monitor"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type authHeader struct {
	token    string
	clientID int64
	userID   string
	hashed   bool
}

var gatewayDB mgr.GatewayDB
var errDB gwerr.GatewayErrorMonitor
var monDB gwmon.GatewayPerformanceMonitor
var cbDB cb.CircuitBreaker

//var gwr mgr.GatewayRoutes

func main() {

	if os.Getenv("MYSQL_PORT_3306_TCP_ADDR") != "" {
		gatewayDB.DbConfig.Host = os.Getenv("MYSQL_PORT_3306_TCP_ADDR")
	} else if os.Getenv("DATABASE_HOST") != "" {
		gatewayDB.DbConfig.Host = os.Getenv("DATABASE_HOST")
	} else {
		gatewayDB.DbConfig.Host = "localhost:3306"
	}

	if os.Getenv("DATABASE_USER_NAME") != "" {
		gatewayDB.DbConfig.DbUser = os.Getenv("DATABASE_USER_NAME")
	} else {
		gatewayDB.DbConfig.DbUser = "admin"
	}

	if os.Getenv("DATABASE_USER_PASSWORD") != "" {
		gatewayDB.DbConfig.DbPw = os.Getenv("DATABASE_USER_PASSWORD")
	} else {
		gatewayDB.DbConfig.DbPw = "admin"
	}

	if os.Getenv("DATABASE_NAME") != "" {
		gatewayDB.DbConfig.DatabaseName = os.Getenv("DATABASE_NAME")
	} else {
		gatewayDB.DbConfig.DatabaseName = "ulbora_api_gateway"
	}
	gatewayDB.ConnectDb()
	defer gatewayDB.CloseDb()
	//gwr.GwDB = gatewayDB
	errDB.DbConfig = gatewayDB.DbConfig
	monDB.CacheHost = getCacheHost()
	monDB.CallBatchSize = 10 //size of cache batch saved. normal should be 100
	monDB.DbConfig = gatewayDB.DbConfig
	cbDB.DbConfig = gatewayDB.DbConfig
	cbDB.CacheHost = getCacheHost()
	gatewayDB.GwCacheHost = getCacheHost()
	gatewayDB.Cb = cbDB

	fmt.Println("Api Gateway running on port 3011!")
	router := mux.NewRouter()
	//super admin client services
	router.HandleFunc("/rs/gwClient/add", handleClientChange)
	router.HandleFunc("/rs/gwClient/update", handleClientChange)
	router.HandleFunc("/rs/gwClient/get/{clientId}", handleClient)
	router.HandleFunc("/rs/gwClient/list", handleClientList)
	router.HandleFunc("/rs/gwClient/delete/{clientId}", handleClient)

	// super admin restRoute services
	router.HandleFunc("/rs/gwRestRouteSuper/add", handleRestRouteSuperChange)
	router.HandleFunc("/rs/gwRestRouteSuper/update", handleRestRouteSuperChange)
	router.HandleFunc("/rs/gwRestRouteSuper/get/{id}/{clientId}", handleRestRouteSuper)
	router.HandleFunc("/rs/gwRestRouteSuper/list/{clientId}", handleRestRouteSuperList)
	router.HandleFunc("/rs/gwRestRouteSuper/delete/{id}/{clientId}", handleRestRouteSuper)

	// super admin routeUrl services
	router.HandleFunc("/rs/gwRouteUrlSuper/add", handleRouteURLSuperChange)
	router.HandleFunc("/rs/gwRouteUrlSuper/update", handleRouteURLSuperChange)
	router.HandleFunc("/rs/gwRouteUrlSuper/get/{id}/{routeId}/{clientId}", handleRouteURLSuper)
	router.HandleFunc("/rs/gwRouteUrlSuper/list/{routeId}/{clientId}", handleRouteURLSuperList)
	router.HandleFunc("/rs/gwRouteUrlSuper/delete/{id}/{routeId}/{clientId}", handleRouteURLSuper)
	router.HandleFunc("/rs/gwRouteUrlSuper/activate", handleRouteURLActivateSuper)

	//super performance service
	router.HandleFunc("/rs/gwPerformanceSuper", handlePeformanceSuper)

	//super errors service
	router.HandleFunc("/rs/gwErrorsSuper", handleErrorsSuper)

	// super Breaker services
	router.HandleFunc("/rs/gwBreakerSuper/add", handleBreakerSuperChange)
	router.HandleFunc("/rs/gwBreakerSuper/update", handleBreakerSuperChange)
	router.HandleFunc("/rs/gwBreakerSuper/reset", handleBreakerSuperReset)
	router.HandleFunc("/rs/gwBreakerSuper/get/{urlId}/{routeId}/{clientId}", handleBreakerSuper)
	router.HandleFunc("/rs/gwBreakerSuper/status/{urlId}/{clientId}", handleBreakerStatusSuper)
	//router.HandleFunc("/rs/gwRouteUrlSuper/list/{routeId}/{clientId}", handleRouteURLSuperList)
	router.HandleFunc("/rs/gwBreakerSuper/delete/{urlId}/{routeId}/{clientId}", handleBreakerSuper)
	//router.HandleFunc("/rs/gwRouteUrlSuper/activate", handleRouteURLActivateSuper)

	// admin restRoute services
	router.HandleFunc("/rs/gwClientUser/get", handleUserClient)

	router.HandleFunc("/rs/gwRestRoute/add", handleRestRouteChange)
	router.HandleFunc("/rs/gwRestRoute/update", handleRestRouteChange)
	router.HandleFunc("/rs/gwRestRoute/get/{id}", handleRestRoute)
	router.HandleFunc("/rs/gwRestRoute/list", handleRestRouteList)
	router.HandleFunc("/rs/gwRestRoute/delete/{id}", handleRestRoute)

	// admin routeUrl services
	router.HandleFunc("/rs/gwRouteUrl/add", handleRouteURLChange)
	router.HandleFunc("/rs/gwRouteUrl/update", handleRouteURLChange)
	router.HandleFunc("/rs/gwRouteUrl/get/{id}/{routeId}", handleRouteURL)
	router.HandleFunc("/rs/gwRouteUrl/list/{routeId}", handleRouteURLList)
	router.HandleFunc("/rs/gwRouteUrl/delete/{id}/{routeId}", handleRouteURL)
	router.HandleFunc("/rs/gwRouteUrl/activate", handleRouteURLActivate)

	//admin performance service
	router.HandleFunc("/rs/gwPerformance", handlePeformance)

	//admin errors service
	router.HandleFunc("/rs/gwErrors", handleErrors)

	// admin Breaker services
	router.HandleFunc("/rs/gwBreaker/add", handleBreakerChange)
	router.HandleFunc("/rs/gwBreaker/update", handleBreakerChange)
	router.HandleFunc("/rs/gwBreaker/reset", handleBreakerReset)
	router.HandleFunc("/rs/gwBreaker/get/{urlId}/{routeId}", handleBreaker)
	router.HandleFunc("/rs/gwBreaker/status/{urlId}", handleBreakerStatus)
	//router.HandleFunc("/rs/gwRouteUrlSuper/list/{routeId}/{clientId}", handleRouteURLSuperList)
	router.HandleFunc("/rs/gwBreaker/delete/{urlId}/{routeId}", handleBreaker)
	//router.HandleFunc("/rs/gwRouteUrlSuper/activate", handleRouteURLActivateSuper)

	//gateway routes
	router.HandleFunc("/np/{route}/{rname}/{fpath:[^.]+}", handleGwRoute)
	router.HandleFunc("/{route}/{fpath:[^ ]+}", handleGwRoute)
	//disgard -- router.HandleFunc("/{route}/{fpath:[^.]+}", handleGwRoute)
	router.HandleFunc("/{route}", handleGwRoute)
	http.ListenAndServe(":3011", router)
}
