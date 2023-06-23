/*
Copyright © 2023 Thomas von Dein

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.

Commentary:

   This little program implements:
https://aws.amazon.com/blogs/database/performance-impact-of-idle-postgresql-connections/

   It  opens 1k  concurrent  postgres connections,  does  a couple  of
   selects  and  then  stops  doing anything  further,  thus  creating
   hanging idle sessions.

*/

package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	_ "github.com/lib/pq"
	flag "github.com/spf13/pflag"
)

type Tableschema struct {
	Name string
}

const Maxloop int = 100
const Getschema string = "SELECT table_schema||'.'||table_name as relname from information_schema.tables WHERE table_schema='information_schema';"
const Gettable string = "SELECT * FROM information_schema.columns LIMIT 1;"
const StartTransaction string = "BEGIN TRANSACTION"
const HostConnection string = "user=%s dbname=%s password=%s"
const NetConnection string = "user=%s dbname=%s password=%s host=%s port=%d"

var wg sync.WaitGroup

func main() {
	var optPasswd string
	var optUser string
	var optDatabase string
	var optServer string
	var optPort int
	var optMaxconnections int
	var optTimeout int
	var optIdleTransaction bool
	var ctx context.Context
	var conn string

	flag.StringVarP(&optPasswd, "password", "p", "", "Password of the database user")
	flag.StringVarP(&optUser, "user", "u", "postgres", "Database user")
	flag.StringVarP(&optDatabase, "database", "d", "postgres", "Database")
	flag.StringVarP(&optServer, "server", "s", "localhost", "Server")
	flag.IntVarP(&optMaxconnections, "client", "c", 500, "Number of concurrent users")
	flag.IntVarP(&optPort, "port", "P", 5432, "TCP Port")
	flag.IntVarP(&optTimeout, "timeout", "t", 0, "Wether to stop the clients after N seconds")
	flag.BoolVarP(&optIdleTransaction, "idletransaction", "i", false, "Wether to stay in idle in transaction state")

	flag.Parse()

	if optServer != "" {
		conn = fmt.Sprintf(NetConnection, optUser, optDatabase, optPasswd, optServer, optPort)
	} else {
		conn = fmt.Sprintf(HostConnection, optUser, optDatabase, optPasswd)
	}

	// first do a connection test
	db, err := sql.Open("postgres", conn)
	if err != nil {
		log.Fatal(err)
	}
	db.Close()

	log.Printf("DB Connection works, firing up %d clients\n", optMaxconnections)

	if optTimeout > 0 {
		ctx, _ = context.WithTimeout(context.Background(), time.Duration(optTimeout)*time.Second)
		log.Printf("Clients will be killed after %d seconds", optTimeout)
	} else {
		ctx = context.TODO()
		log.Println("Clients will run endlessly, abort with C-c")
	}

	wg.Add(optMaxconnections)

	for m := 0; m < optMaxconnections; m++ {
		go dbClient(ctx, conn, optIdleTransaction)
	}

	wg.Wait()
}

func dbClient(ctx context.Context, conn string, idle bool) {
	defer wg.Done()

	db, err := sql.Open("postgres", conn)
	if err != nil {
		log.Fatal(err)
	}

	if idle {
		_, err := db.Exec(StartTransaction)
		if err != nil {
			log.Fatal(err)
		}
	}
	rows, err := db.Query(Getschema)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	//log.Println("Got rows")

	for rows.Next() {
		var T Tableschema
		if err := rows.Scan(&T.Name); err != nil {
			log.Fatal(err)
		}
		//log.Printf("Got table %s\n", T.Name)

		for i := 0; i < Maxloop; i++ {
			rows, err := db.Query(Gettable)
			if err != nil {
				log.Fatal(err)
			}
			rows.Close() // ignore result
		}
	}

	//log.Println("Got tables")

	// block this thread forever or timeout
	//select {}
	select {
	case <-time.After(1 * time.Second):
	case <-ctx.Done():
	}
}
