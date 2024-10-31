.PHONY: all

all: build

build:
	docker build -t namu-rank-archive .
