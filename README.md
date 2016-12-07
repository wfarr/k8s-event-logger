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

To enable Statsd reporting, set `STATSD_URL` in your ENV.