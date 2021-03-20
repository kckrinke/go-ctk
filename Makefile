#!/usr/bin/make -f

prefix ?= $(GOPATH)
bindir ?= $(prefix)/bin

.PHONY: all help build debug clean dev run profile.cpu profile.mem

all: build

help:
	@echo "usage: make {help|build|debug|clean|dev|run|profile.cpu|profile.mem}"
	@echo
	@echo "  build: compiles the cmd applications (go-charmap, go-dialog, go-ctk)"
	@echo "  debug: compiles the cmd applications with debugging features"
	@echo "  clean: cleans build files, logs and other intermediaries"
	@echo "  dev: debugging build of ctk-app (work-in-progress)"
	@echo "  run: run a dev build and reset the terminal if it fails"
	@echo "  profile.cpu: run with CPU profiling enabled, prompt to open pprof"
	@echo "  profile.mem: run with MEM profiling enabled, prompt to open pprof"
	@echo

build: build-go-ctk

build-examples: build-go-charmap build-go-dialog

build-go-charmap:
	@echo "building: go-charmap [release]"
	@cd ./cmd/go-charmap && go build -v -o ../../go-charmap .

build-go-dialog:
	@echo "building: go-dialog [release]"
	@cd ./cmd/go-dialog && go build -v -o ../../go-dialog .

build-go-ctk:
	@echo "building: go-ctk [release]"
	@cd ./cmd/go-ctk && go build -v -o ../../go-ctk .

debug: debug-go-charmap debug-go-dialog debug-go-ctk

debug-go-charmap:
	@echo "building: go-charmap [debug]"
	@cd ./cmd/go-charmap && go build \
		-ldflags="-X 'main.IncludeProfiling=true'" \
		-gcflags=all="-N -l" \
	  -v -o ../../go-charmap \
		.

debug-go-dialog:
	@echo "building: go-dialog [debug]"
	@cd ./cmd/go-dialog && go build \
		-ldflags="-X 'main.IncludeProfiling=true'" \
		-gcflags=all="-N -l" \
	  -v -o ../../go-dialog \
		.

debug-go-ctk:
	@echo "building: go-ctk [debug]"
	@cd cmd/go-ctk && go build \
		-ldflags="-X 'main.IncludeProfiling=true'" \
		-gcflags=all="-N -l" \
	  -v -o ../../go-ctk \
		.

install: build-go-ctk
	@echo "installing: go-ctk [release]"
	@install -v --target-directory=$(bindir) go-ctk

debug-install: debug-go-ctk
	@echo "installing: go-ctk [debug]"
	@install -v --target-directory=$(bindir) go-ctk

clean-logs:
	@echo "cleaning *.log and pprof.* files"
	@rm -fv *.log            || true
	@rm -fv pprof.*          || true

clean: clean-logs
	@echo "cleaning compiled outputs"
	@rm -fv go-charmap       || true
	@rm -fv go-dialog        || true
	@rm -fv go-ctk           || true
	@rm -fv hello-world      || true
	@rm -fv ctk-app          || true
	@rm -fv ctk-kitchen-sink || true
	@rm -rfv go_*            || true

dev:
	@echo "building: ctk-app [dev]"
	@go build -v \
		-ldflags="-X 'main.IncludeProfiling=true'" \
		-gcflags=all="-N -l" \
		./example/ctk-app/ctk-app.go

run: export GO_CDK_LOG_FILE=./ctk-app.log
run: export GO_CDK_LOG_LEVEL=debug
run: export GO_CDK_LOG_FULL_PATHS=true
run: dev
	@[ -f ctk-app ] && ./ctk-app || reset

profile.cpu: export GO_CDK_LOG_FILE=./ctk-app.log
profile.cpu: export GO_CDK_LOG_LEVEL=debug
profile.cpu: export GO_CDK_LOG_FULL_PATHS=true
profile.cpu: export GO_CDK_PROFILE=cpu
profile.cpu: dev
	@[ -f ctk-app ] && ./ctk-app || reset
	@read -p "Press enter to open a pprof instance" JUNK && go tool pprof -http=:8080 /tmp/cdk.pprof/cpu.pprof

profile.mem: export GO_CDK_LOG_FILE=./ctk-app.log
profile.mem: export GO_CDK_LOG_LEVEL=debug
profile.mem: export GO_CDK_LOG_FULL_PATHS=true
profile.mem: export GO_CDK_PROFILE=mem
profile.mem: dev
	@[ -f ctk-app ] && ./ctk-app || reset
	@read -p "Press enter to open a pprof instance" JUNK && go tool pprof -http=:8080 /tmp/cdk.pprof/mem.pprof
