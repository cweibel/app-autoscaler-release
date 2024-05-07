package forwarder

import (
	"code.cloudfoundry.org/app-autoscaler/src/autoscaler/metricsforwarder/config"
	"code.cloudfoundry.org/app-autoscaler/src/autoscaler/models"
	"code.cloudfoundry.org/go-loggregator/v9"

	"code.cloudfoundry.org/lager/v3"
)

type Emitter interface {
	EmitMetric(*models.CustomMetric)
}

type MetronEmitter struct {
	client *loggregator.IngressClient
	logger lager.Logger
}

const METRICS_FORWARDER_ORIGIN = "autoscaler_metrics_forwarder"

func hasLoggregatorConfig(conf *config.Config) bool {
	return conf.LoggregatorConfig.MetronAddress != ""
}

func NewMetricForwarder(logger lager.Logger, conf *config.Config) (Emitter, error) {
	if hasLoggregatorConfig(conf) {
		return NewMetronEmitter(logger, conf)
	} else {
		return NewSyslogEmitter(logger, conf)
	}
}
