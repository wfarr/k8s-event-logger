package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/fields"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"

	// Only required to authenticate against GKE clusters
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func main() {
	kubeconfig := flag.String("kubeconfig", "", "Path to a kube config. Only required if out-of-cluster.")
	bugsnagAPIKey := flag.String("bugsnag-api-key", "", "Enables reporting non-Normal events to bugsnag using the given API key.")
	statsd := flag.String("statsd", "", "URI to send statsd metrics to. Must be dogstatsd compatible.")
	flag.Parse()

	if *bugsnagAPIKey != "" {
		configureBugsnag(*bugsnagAPIKey)
	}

	if *statsd != "" {
		configureDatadog(*statsd)
	}

	config, err := buildConfig(*kubeconfig)
	if err != nil {
		panic(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	stop := make(chan struct{}, 1)
	source := cache.NewListWatchFromClient(
		clientset.Core().RESTClient(),
		"events",
		api.NamespaceAll,
		fields.Everything())

	create := func(obj interface{}) {
		event := obj.(*v1.Event)

		if *bugsnagAPIKey != "" {
			if err := sendEventToBugsnag(event); err != nil {
				log.WithError(err)
			}
		}

		if *statsd != "" {
			if err := sendEventToDatadog(event); err != nil {
				log.WithError(err)
			}
		}

		if err := sendEventToSTDOUT(event); err != nil {
			log.WithError(err)
		}
	}

	_, controller := cache.NewInformer(
		source,
		&v1.Event{},
		time.Second*0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: create,
		})

	go controller.Run(stop)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	s := <-signals
	fmt.Printf("received signal %#v, exiting...\n", s)
	close(stop)
	os.Exit(0)
}

func buildConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	return rest.InClusterConfig()
}

func create(obj interface{}) {

}
