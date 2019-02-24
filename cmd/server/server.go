package main

import (
	"flag"

	"github.com/didil/volusnap/pkg/api"
	"github.com/sirupsen/logrus"
)

func main() {
	port := flag.Int("p", 8080, "server port")
	flag.Parse()

	logrus.Fatal(api.StartServer(*port))
}
