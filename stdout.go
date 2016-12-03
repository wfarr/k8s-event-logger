package main

import (
	log "github.com/Sirupsen/logrus"

	"os"

	"k8s.io/client-go/pkg/api/v1"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.JSONFormatter{})
}

func sendEventToSTDOUT(event *v1.Event) error {
	log.WithFields(log.Fields{
		"application":               "k8s-event-logger",
		"cluster":                   event.GetClusterName(),
		"event-uid":                 event.UID,
		"event-type":                event.Type,
		"event-reason":              event.Reason,
		"event-source-component":    event.Source.Component,
		"event-source-host":         event.Source.Host,
		"involved-object-kind":      event.InvolvedObject.Kind,
		"involved-object-name":      event.InvolvedObject.Name,
		"involved-object-namespace": event.InvolvedObject.Namespace,
	}).Info("Processing Event")

	return nil
}
