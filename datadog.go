package main

import (
	"fmt"

	"github.com/DataDog/datadog-go/statsd"
	"k8s.io/client-go/pkg/api/v1"
)

var (
	statsdClient statsd.Client
)

func configureDatadog(uri string) error {
	statsdClient, err := statsd.New("127.0.0.1:8125")
	if err != nil {
		return err
	}

	statsdClient.Namespace = "k8s-event-logger"
	// c.Tags = append(c.Tags, fmt.Sprintf("cluster=%s", event.ClusterName))
	return nil
}

func sendEventToDatadog(event *v1.Event) error {
	datadogEvent := statsd.NewEvent(event.Reason, event.Message)
	datadogEvent.Tags = []string{
		fmt.Sprintf("event-type=%s", event.Type),
		fmt.Sprintf("event-reason=%s", event.Reason),
		fmt.Sprintf("event-source=%s", event.Source),
		fmt.Sprintf("involved-object-kind=%s", event.InvolvedObject.Kind),
		fmt.Sprintf("involved-object-name=%s", event.InvolvedObject.Name),
		fmt.Sprintf("involved-object-namespace=%s", event.InvolvedObject.Namespace),
	}

	return statsdClient.Event(datadogEvent)
}
