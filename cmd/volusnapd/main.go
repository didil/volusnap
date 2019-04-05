package main

import "github.com/sirupsen/logrus"

func main() {
	cmd := buildRootCmd()
	err := cmd.Execute()
	if err != nil {
		logrus.Fatalf("volusnapd cmd err: %v", err)
	}
}
