package main

import (
	"distributeKV/config"
	"distributeKV/db"
	"distributeKV/web"
	"flag"
	"fmt"
	"log"
	"net/http"
)

var (
	dbLocationFlag      = flag.String("db-location", "my.db", "The location of db file")
	httpAddrFlag        = flag.String("http-addr", "127.0.0.1:8080", "The host address")
	partitionConfigFlag = flag.String("partition-config", "static_partition.toml", "location of the partition config file")
	partitionFlag       = flag.String("partition", "", "Partition name ")
	launchFlag          = flag.String("launch-mode", "", "launch mode for db server or proxy server")
)

func parseFlag() {
	flag.Parse()

	if *launchFlag != "Proxy" && *dbLocationFlag == "" {
		log.Fatalf("DB location file cannot be empty for DB Mode!\n")
	}

	if *launchFlag == "" {
		log.Fatalf("Launch mode cannot be empty\n")
	}
}

func main() {

	parseFlag()

	fmt.Println("Welcome to DistributeKV!")

	// read the partition config
	var partitions = config.ParseCofig(*partitionConfigFlag)

	if *launchFlag == "DB" {
		var curPartition config.Partition
		var partitionCount = len(partitions.Partitions)
		for _, p := range partitions.Partitions {
			if p.Name == *partitionFlag {
				curPartition = p
			}
		}

		db, err, closeFun := db.NewDatabase(*dbLocationFlag)
		if err != nil {
			log.Fatalf("NewDatabase(%s), %v", *dbLocationFlag, err)
		}
		defer closeFun()

		webServer := web.CreateServer(db, curPartition, partitionCount)
		http.HandleFunc("/get", webServer.GetHandler)
		http.HandleFunc("/set", webServer.SetHandler)
		fmt.Printf("[%s] Web Server is listening on %s ...\n", curPartition.Name, *httpAddrFlag)
		log.Fatal(webServer.ListenAndServe())
	} else {

		proxyServer := web.CreatProxyServer(&partitions, "127.0.0.1:8888")
		http.HandleFunc("/get", proxyServer.GetHandler)
		http.HandleFunc("/set", proxyServer.SetHandler)
		fmt.Printf("Proxy Server is listening on %s ...\n", "127.0.0.1:8888")
		log.Fatal(proxyServer.ListenAndServeProxy())
	}

}
