package internal

//go:generate go run github.com/ogen-go/ogen/cmd/ogen --no-server --debug.ignoreNotImplemented "mutualTLS security" --target custommetrics --clean ../../../../../../api/custom-metrics-api.yaml
//go:generate go run github.com/ogen-go/ogen/cmd/ogen  --no-server --debug.ignoreNotImplemented "mutualTLS security" --target scalinghistory --clean ../../../../../../api/scaling-history-api.yaml
//go:generate go run github.com/ogen-go/ogen/cmd/ogen  --no-server --debug.ignoreNotImplemented "mutualTLS security" --target policy --clean ../../../../../../api/policy-api.yaml
//go:generate go run github.com/ogen-go/ogen/cmd/ogen  --no-server --debug.ignoreNotImplemented "mutualTLS security" --target applicationmetric --clean ../../../../../../api/application-metric-api.yaml