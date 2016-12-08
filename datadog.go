package main

import (
	"fmt"
	"os"

	dogstatsd "github.com/Shopify/go-dogstatsd"
	"k8s.io/client-go/pkg/api/v1"
)

func sendEventToDatadog(event *v1.Event) error {
	tags := []string{
		fmt.Sprintf("event-type:%s", event.Type),
		fmt.Sprintf("event-reason:%s", event.Reason),
		fmt.Sprintf("event-source:%s", event.Source),
		fmt.Sprintf("involved-object-kind:%s", event.InvolvedObject.Kind),
		fmt.Sprintf("involved-object-name:%s", event.InvolvedObject.Name),
		fmt.Sprintf("involved-object-namespace:%s", event.InvolvedObject.Namespace),
	}

	statsdClient, err := dogstatsd.New(os.Getenv("STATSD_URL"), &dogstatsd.Context{
		Namespace: "k8s-event-logger.",
	})
	defer statsdClient.Close()

	if err != nil {
		return err
	}

	return statsdClient.Event(event.Reason, event.Message, tags)
}
