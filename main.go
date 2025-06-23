package main

import (
	"context"
	"flag"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	//gr := &rbacv1alpha1.GridOSGroup{}
	fmt.Println("Creating a client for rbac CRDs")
	kubeconfig := flag.String("kubeconfig", "/home/aminebt/.kube/config", "location to your kubeconfig file")

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		fmt.Printf("error %s building config from flags\n", err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)

	if err != nil {
		fmt.Printf("error $s, creating clientset\n", err.Error())
	}

	ctx := context.Background()
	namespace := "kube-system"
	pods, err := clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		fmt.Printf("error $s, while listing pods\n", err.Error())
	}

	fmt.Println("List of pods from specified namespace")
	for _, pod := range pods.Items {
		fmt.Printf("%+v\n", pod.Name)
	}
}
