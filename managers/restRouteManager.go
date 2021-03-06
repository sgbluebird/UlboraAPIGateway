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
package managers

import (
	"fmt"
	"strconv"
)

//InsertRestRoute in database
func (db *GatewayDB) InsertRestRoute(rr *RestRoute) *GatewayResponse {
	var rtn GatewayResponse
	dbConnected := db.DbConfig.ConnectionTest()
	if !dbConnected {
		fmt.Println("reconnection to closed database")
		db.DbConfig.ConnectDb()
	}
	var a []interface{}
	a = append(a, rr.Route, rr.ClientID)
	success, insID := db.DbConfig.InsertRestRoute(a...)
	if success == true {
		fmt.Println("inserted record")
	}
	rtn.ID = insID
	rtn.Success = success
	return &rtn
}

//UpdateRestRoute in database
func (db *GatewayDB) UpdateRestRoute(rr *RestRoute) *GatewayResponse {
	var rtn GatewayResponse
	dbConnected := db.DbConfig.ConnectionTest()
	if !dbConnected {
		fmt.Println("reconnection to closed database")
		db.DbConfig.ConnectDb()
	}
	var a []interface{}
	a = append(a, rr.Route, rr.ID, rr.ClientID)
	success := db.DbConfig.UpdateRestRoute(a...)
	if success == true {
		fmt.Println("update record")
		db.clearCache(rr.ClientID, rr.Route)
	}
	rtn.ID = rr.ID
	rtn.Success = success
	return &rtn
}

//GetRestRoute route from database
func (db *GatewayDB) GetRestRoute(rr *RestRoute) *RestRoute {
	var a []interface{}
	a = append(a, rr.ID, rr.ClientID)
	var rtn *RestRoute
	rowPtr := db.DbConfig.GetRestRoute(a...)
	if rowPtr != nil {
		//print("content row: ")
		//println(rowPtr.Row)
		foundRow := rowPtr.Row
		rtn = parseRestRouteRow(&foundRow)
	}
	return rtn
}

//GetRestRouteList route
func (db *GatewayDB) GetRestRouteList(rr *RestRoute) *[]RestRoute {
	var rtn []RestRoute
	var a []interface{}
	a = append(a, rr.ClientID)
	rowsPtr := db.DbConfig.GetRestRouteList(a...)
	if rowsPtr != nil {
		foundRows := rowsPtr.Rows
		for r := range foundRows {
			foundRow := foundRows[r]
			rowContent := parseRestRouteRow(&foundRow)
			rtn = append(rtn, *rowContent)
		}
	}
	return &rtn
}

//DeleteRestRoute in database
func (db *GatewayDB) DeleteRestRoute(rr *RestRoute) *GatewayResponse {
	var rtn GatewayResponse
	dbConnected := db.DbConfig.ConnectionTest()
	if !dbConnected {
		fmt.Println("reconnection to closed database")
		db.DbConfig.ConnectDb()
	}
	var a []interface{}
	a = append(a, rr.ID, rr.ClientID)
	success := db.DbConfig.DeleteRestRoute(a...)
	if success == true {
		fmt.Println("deleted record")
		db.clearCache(rr.ClientID, rr.Route)
	}
	rtn.ID = rr.ID
	rtn.Success = success
	return &rtn
}

func parseRestRouteRow(foundRow *[]string) *RestRoute {
	var rtn RestRoute
	if len(*foundRow) > 0 {
		ID, errID := strconv.ParseInt((*foundRow)[0], 10, 0)
		if errID != nil {
			fmt.Print(errID)
		} else {
			rtn.ID = ID
		}
		rtn.Route = (*foundRow)[1]
		CID, errID2 := strconv.ParseInt((*foundRow)[2], 10, 0)
		if errID2 != nil {
			fmt.Print(errID2)
		} else {
			rtn.ClientID = CID
		}
	}
	return &rtn
}
