package metrics

import (
	envutils "carduka/bidsvc/pkg/utils/env"
)

func IsEnabled() bool {
	return envutils.EnvBoolOr(EnvMetricsEnabled, false)
}
