module github.com/Ankr-network/dccn-tools/zap/benchmarks

go 1.13

replace github.com/Ankr-network/dccn-tools/zap => ../

require (
	github.com/Ankr-network/dccn-tools/zap v0.0.0-00010101000000-000000000000
	github.com/apex/log v1.1.1
	github.com/go-kit/kit v0.9.0
	github.com/go-stack/stack v1.8.0 // indirect
	github.com/rs/zerolog v1.16.0
	github.com/sirupsen/logrus v1.4.2
	go.uber.org/multierr v1.3.0
	gopkg.in/inconshreveable/log15.v2 v2.0.0-20180818164646-67afb5ed74ec
)
