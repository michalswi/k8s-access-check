package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"./pkg/checker"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// go run main.go --dir <json> --ns <namespace> --run-outside-k-cluster true

func newClientSet(runOutsideKcluster bool) (*kubernetes.Clientset, error) {

	kubeConfigLocation := ""

	if runOutsideKcluster == true {
		homeDir := os.Getenv("HOME")
		kubeConfigLocation = filepath.Join(homeDir, ".kube", "config")
	}

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigLocation)
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(config)
}

func decodeJson(dir string) (map[string]interface{}, error) {

	byteValue, err := ioutil.ReadFile(dir)
	if err != nil {
		log.Printf("%v", err)
	}

	var results map[string]interface{}
	// fmt.Println(results["action"])

	if err := json.Unmarshal([]byte(byteValue), &results); err != nil {
		log.Printf("%v", err)
	}

	takeAction := results["action"].(map[string]interface{})

	return takeAction, err
}

func main() {

	// import json
	var dir string
	flag.StringVar(&dir, "dir", "", "Set this flag when passing a json file, e.g. '--dir /tmp/your.json'.")

	// namespace
	var ns string
	flag.StringVar(&ns, "ns", "default", "Set this flag when changing default namespace.")

	// cloudconfig
	runOutsideKcluster := flag.Bool("run-outside-k-cluster", false, "Set this flag when running outside of the cluster.")
	flag.Parse()

	// 'dir' must be
	if dir == "" {
		log.Fatal("Flag '--dir' not provided")
	}

	decjson, err := decodeJson(dir)
	if err != nil {
		log.Fatal(err)
	}

	clientset, err := newClientSet(*runOutsideKcluster)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Init namespace: %s\n", ns)

	// struct 'Params' defined in 'checker.go'
	params := &checker.Params{Kclient: clientset, Decjson: decjson, Ns: ns}

	// one by one
	// params.GetKubeVersion()
	// params.WhatCanIdo()
	// params.WhatCanIdoList()

	// run specific one
	params.Runner()
}
