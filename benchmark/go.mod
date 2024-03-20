module example

go 1.22.0

replace github.com/safeblock-dev/wr => ./..

require (
	github.com/alitto/pond v1.8.3
	github.com/devchat-ai/gopool v0.6.2
	github.com/panjf2000/ants/v2 v2.9.0
	github.com/safeblock-dev/wr v0.0.0-00010101000000-000000000000
	github.com/sourcegraph/conc v0.3.0
)

require (
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.9.0 // indirect
)
