package main

import (
	"distributeKV/db"
	"distributeKV/web"
	"flag"
	"fmt"
	"log"
	"net/http"
)

var (
	dbLocationFlag = flag.String("db-location", "", "The location of db file")
	httpAddrFlag   = flag.String("http-addr", "127.0.0.1:8080", "The host address")
)

func parseFlag() {
	flag.Parse()

	if *dbLocationFlag == "" {
		log.Fatalf("DB location file cannot be empty!\n")
	}
}

func main() {

	parseFlag()

	fmt.Println("Welcome to DistributeKV!")
	db, err, closeFun := db.NewDatabase(*dbLocationFlag)
	if err != nil {
		log.Fatalf("NewDatabase(%s), %v", *dbLocationFlag, err)
	}
	defer closeFun()

	webServer := web.CreateServer(db)
	http.HandleFunc("/get", webServer.GetHandler)
	http.HandleFunc("/set", webServer.SetHandler)
	fmt.Printf("Http Server is listening on %s ...\n", *httpAddrFlag)
	log.Fatal(webServer.ListenAndServe(*httpAddrFlag))
}
