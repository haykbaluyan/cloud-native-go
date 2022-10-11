PROJ_ROOT := $(shell pwd)

build_circuitbreaker:
	echo "*** Building circuitbreaker"
	cd ${PROJ_ROOT}/circuitbreaker && go build -o ./bin/circuitbreaker ./

build_throttle:
	echo "*** Building throttle"
	cd ${PROJ_ROOT}/throttle && go build -o ./bin/throttle ./