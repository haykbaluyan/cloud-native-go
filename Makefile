PROJ_ROOT := $(shell pwd)

build_circuitbreaker:
	echo "*** Building circuitbreaker"
	cd ${PROJ_ROOT}/circuitbreaker && go build -o ./bin/circuitbreaker ./

build_throttle:
	echo "*** Building throttle"
	cd ${PROJ_ROOT}/throttle && go build -o ./bin/throttle ./

build_retry:
	echo "*** Building retry"
	cd ${PROJ_ROOT}/retry && go build -o ./bin/retry ./

build_timeout:
	echo "*** Building timeout"
	cd ${PROJ_ROOT}/timeout && go build -o ./bin/timeout ./

build: build_circuitbreaker build_throttle build_retry build_timeout