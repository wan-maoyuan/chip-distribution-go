# 程序名称
NAME = chip-distribution
# 版本
# VERSION = $(shell git describe --tags --always)
VERSION = v1.0.0
# git 提交 Hash
COMMIT_HASH = $(shell git show -s --format=%H)
# build 时间
BUILD_TIME ?= $(shell date +%Y%m%d%H%M%S)
# go可执行程序输出目录
DIST_FOLDER := dist
# docker镜像构建目录
RELEASE_FOLDER := release
# go文件列表
GOFILES := $(shell find . ! -path "./vendor/*" -name "*.go")
# 构建附加选项
BUILD_OPTS := -ldflags "-s -w -X 'main.Version=${VERSION}' -X 'main.CommitHash=${COMMIT_HASH}' -X 'main.BuildTime=${BUILD_TIME}'"
# 编译环境
BUILD_ENV := GOOS=linux GOARCH=amd64


.PHONY: build
# linux 平台打包编译
build: ${DIST_FOLDER}/${NAME}
# 构建目标
${DIST_FOLDER}/${NAME}: ${GOFILES}
	go mod tidy
	${BUILD_ENV} go build ${BUILD_OPTS} -o $@


.PHONY: container
# docker 镜像打包
container: ${DIST_FOLDER}/${NAME}
	docker build -t ${NAME}:${VERSION}					\
	    -f ${RELEASE_FOLDER}/Dockerfile .;


.PHONY: clean
# 清理
clean:
	-rm -rf $(DIST_FOLDER)/*
	-go clean 
	-go clean -cache