# k8s-event-logger

Subscribes to Kubernetes events, logging them to:

* STDOUT (default)
* Bugsnag (optional)
* Statsd (optional, dogstatsd)

## Getting Started

```
# make sure you have a Secret "bugsnag" with a key "api-key" whose value is your bugsnag API Key
make deploy
```

## Configuring

To enable Bugsnag reporting, set `BUGSNAG_API_KEY` in your ENV.
The included sample deployment YAML expects this to come from a Kubernetes secret, which you can create like so:

```
$ kubectl create secret generic bugsnag --from-literal=api-key=MYSECRETAPIKEY
```

To enable Statsd reporting, set `STATSD_URL` in your ENV.
The included sample deployment YAML expects this to be `statsd:8125` (you have a service named statsd running on UDP 8125). 