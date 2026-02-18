module github.com/Gattancha-Computer-Services/rubrik-exporter

go 1.25

require (
	github.com/machinebox/graphql v0.2.2
	github.com/prometheus/client_golang v1.21.0
)

require github.com/pkg/errors v0.9.1 // indirect

replace (
	github.com/Gattancha-Computer-Services/rubrik-exporter => ./
	github.com/Sirupsen/logrus => github.com/sirupsen/logrus v1.9.3
)
