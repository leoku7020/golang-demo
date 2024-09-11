package initkit

import (
	"os"

	"github.com/newrelic/go-agent/v3/newrelic"

	"demo/pkg/logger"
)

// ref: https://github.com/newrelic/go-agent/blob/master/v3/newrelic/config_options.go
// ConfigFromEnvironment populates the config based on environment variables:
//
//  NEW_RELIC_APP_NAME                                sets AppName
//  NEW_RELIC_ATTRIBUTES_EXCLUDE                      sets Attributes.Exclude using a comma-separated list, eg. "request.headers.host,request.method"
//  NEW_RELIC_ATTRIBUTES_INCLUDE                      sets Attributes.Include using a comma-separated list
//  NEW_RELIC_DISTRIBUTED_TRACING_ENABLED             sets DistributedTracer.Enabled using strconv.ParseBool
//  NEW_RELIC_ENABLED                                 sets Enabled using strconv.ParseBool
//  NEW_RELIC_HIGH_SECURITY                           sets HighSecurity using strconv.ParseBool
//  NEW_RELIC_HOST                                    sets Host
//  NEW_RELIC_INFINITE_TRACING_SPAN_EVENTS_QUEUE_SIZE sets InfiniteTracing.SpanEvents.QueueSize using strconv.Atoi
//  NEW_RELIC_INFINITE_TRACING_TRACE_OBSERVER_PORT    sets InfiniteTracing.TraceObserver.Port using strconv.Atoi
//  NEW_RELIC_INFINITE_TRACING_TRACE_OBSERVER_HOST    sets InfiniteTracing.TraceObserver.Host
//  NEW_RELIC_LABELS                                  sets Labels using a semi-colon delimited string of colon-separated pairs, eg. "Server:One;DataCenter:Primary"
//  NEW_RELIC_LICENSE_KEY                             sets License
//  NEW_RELIC_LOG                                     sets Logger to log to either "stdout" or "stderr" (filenames are not supported)
//  NEW_RELIC_LOG_LEVEL                               controls the NEW_RELIC_LOG level, must be "debug" for debug, or empty for info
//  NEW_RELIC_PROCESS_HOST_DISPLAY_NAME               sets HostDisplayName
//  NEW_RELIC_SECURITY_POLICIES_TOKEN                 sets SecurityPoliciesToken
//  NEW_RELIC_UTILIZATION_BILLING_HOSTNAME            sets Utilization.BillingHostname
//  NEW_RELIC_UTILIZATION_LOGICAL_PROCESSORS          sets Utilization.LogicalProcessors using strconv.Atoi
//  NEW_RELIC_UTILIZATION_TOTAL_RAM_MIB               sets Utilization.TotalRAMMIB using strconv.Atoi

func NewNewRelicApp() *newrelic.Application {
	logger.Info("NewRelicApp preparing", logger.WithFields(logger.Fields{
		"enable":             os.Getenv("NEW_RELIC_ENABLED"),
		"appName":            os.Getenv("NEW_RELIC_APP_NAME"),
		"licKey":             os.Getenv("NEW_RELIC_LICENSE_KEY"),
		"distributedTracing": os.Getenv("NEW_RELIC_DISTRIBUTED_TRACING_ENABLED"),
	}))

	app, err := newrelic.NewApplication(
		newrelic.ConfigFromEnvironment(),
		logger.ConfigNRLogger(),
	)
	if err != nil {
		logger.Panic("newrelic.NewApplication failed")
	}

	logger.Info("NewRelicApp done", logger.WithFields(logger.Fields{
		"enable":             os.Getenv("NEW_RELIC_ENABLED"),
		"appName":            os.Getenv("NEW_RELIC_APP_NAME"),
		"distributedTracing": os.Getenv("NEW_RELIC_DISTRIBUTED_TRACING_ENABLED"),
	}))

	return app
}
