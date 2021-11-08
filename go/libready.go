package main

import "C"
import (
	"context"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"os"
	"strconv"
	"strings"
)

func IsMaster(node *v1.Node) bool {
	if _, ok := node.Labels["node-role.kubernetes.io/master"]; ok {
		return true
	}

	return false
}

func Client() (*kubernetes.Clientset,error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil,err
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil,err
	}

	return clientset,nil
}

//export node_create_check
func node_create_check() (rc int,result *C.char,errStr *C.char) {
	clientset,err := Client()
	if err != nil {
		return -1, nil, C.CString(err.Error())
	}

	pods, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return -1, nil, C.CString(err.Error())
	}

	masterCount := 0
	workerCount := 0

	masterName := ""
	workerName := ""

	for _, item := range pods.Items {
		if strings.Contains(item.Name,os.Getenv("CLUSTER_NAME")) {
			if IsMaster(&item) {
				masterCount++
				masterName += item.Name + ","
			} else {
				workerCount++
				workerName += item.Name + ","
			}
		}
	}

	masterName = strings.Trim(masterName,",")
	workerName = strings.Trim(workerName,",")

	envMasterCount,err := strconv.Atoi(os.Getenv("CONTROL_PLANE_MACHINE_COUNT"))
	if err != nil {
		return -1, nil, C.CString(err.Error())
	}
	envWorkerCount,err := strconv.Atoi(os.Getenv("WORKER_MACHINE_COUNT"))
	if err != nil {
		return -1, nil, C.CString(err.Error())
	}

	resultGo := ""
	if masterCount == envMasterCount && workerCount == envWorkerCount {
		resultGo = "Created All Nodes. Nodes:"+strconv.Itoa(masterCount)+" Master-["+masterName+"] "+strconv.Itoa(workerCount)+" Worker-["+workerName+"]"
		return 0, C.CString(resultGo), nil
	}else {
		resultGo = "Waiting Create All Nodes. Created Nodes:"+strconv.Itoa(masterCount)+" Master-["+masterName+"] "+strconv.Itoa(workerCount)+" Worker-["+workerName+"]"
		return -1, nil, C.CString(resultGo)
	}
}

func main(){}
