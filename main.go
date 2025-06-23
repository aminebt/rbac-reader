package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	rbacv1alpha1 "github.com/aminebt/rbac-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func main() {
	//gr := &rbacv1alpha1.GridOSGroup{}
	fmt.Println("Creating a client for rbac CRDs")
	kubeconfig := flag.String("kubeconfig", "/home/aminebt/.kube/config", "location to your kubeconfig file")

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		fmt.Printf("error %s building config from flags\n", err.Error())
	}

	////////////////////////////////////////////////////////////// simple clientset
	// clientset, err := kubernetes.NewForConfig(config)

	// if err != nil {
	// 	fmt.Printf("error $s, creating clientset\n", err.Error())
	// }

	// ctx := context.Background()
	// namespace := "kube-system"
	// pods, err := clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	// if err != nil {
	// 	fmt.Printf("error $s, while listing pods\n", err.Error())
	// }

	// fmt.Println("List of pods from specified namespace")
	// for _, pod := range pods.Items {
	// 	fmt.Printf("%+v\n", pod.Name)
	// }
	//////////////////////////////////////////////////////////////

	scheme := runtime.NewScheme()
	if err := rbacv1alpha1.AddToScheme(scheme); err != nil {
		fmt.Println("failed to add rbac types to scheme", "error", err)
		os.Exit(1)
	}
	fmt.Println("successfully added rbac types to scheme")

	groupCacheOpts := cache.Options{
		Scheme: scheme,
		// this requires us to explicitly start an informer for each object type
		// and helps avoid people mistakenly using the group client for other resources
		ReaderFailOnMissingInformer: false,
	}

	groupCache, err := cache.New(config, groupCacheOpts)

	// start an informer for groups
	// this is required because we set ReaderFailOnMissingInformer to true
	// _, err = groupCache.GetInformer(context.Background(), &rbacv1alpha1.GridOSGroup{})
	// if err != nil {
	// 	fmt.Println("failed to start an informer for groups", "error", err)
	// 	os.Exit(1)
	// }

	groupClient, err := client.New(config, client.Options{
		Scheme: scheme,
		Cache: &client.CacheOptions{
			Reader: groupCache,
		},
	})

	if err != nil {
		fmt.Println("failed to create client for groups", "error", err)
		os.Exit(1)
	}

	ctx := context.Background()

	go groupCache.Start(ctx)
	// if err != nil {
	// 	fmt.Println("failed to start cache", "error: ", err)
	// 	os.Exit(1)
	// }

	//namespace := "default"
	groups := &rbacv1alpha1.GridOSGroupList{}
	err = groupClient.List(ctx, groups)
	if err != nil {
		fmt.Println("failed to list groups", "error: ", err)
		os.Exit(1)
	}

	fmt.Println("List of groups from specified namespace")
	for _, group := range groups.Items {
		fmt.Printf("%+v\n", group.Name)
	}
}
