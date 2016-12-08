package main

import (
	"fmt"

	"github.com/bugsnag/bugsnag-go"
	"k8s.io/client-go/pkg/api/v1"
)

func configureBugsnag(apiKey, releaseStage string) error {
	bugsnag.Configure(bugsnag.Configuration{
		APIKey:       apiKey,
		ReleaseStage: releaseStage,
	})
	return nil
}

func sendEventToBugsnag(event *v1.Event) error {
	if event.Type != "Normal" {
		bugsnag.Notify(
			fmt.Errorf("Type: %s, Reason: %s, Message: %s", event.Type, event.Reason, event.Message),
			bugsnag.MetaData{
				"Event": {
					"Type":    event.Type,
					"Reason":  event.Reason,
					"Message": event.Message,
				},
				"InvolvedObject": {
					"Kind":      event.InvolvedObject.Kind,
					"Name":      event.InvolvedObject.Name,
					"Namespace": event.InvolvedObject.Namespace,
					"UID":       event.InvolvedObject.UID,
				},
				"Source": {
					"Component": event.Source.Component,
					"Host":      event.Source.Host,
				},
			})
	}
	return nil
}
