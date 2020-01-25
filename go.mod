module github.com/keilerkonzept/dockerfile-json

go 1.13

require (
	github.com/moby/buildkit v0.6.3
	github.com/yalp/jsonpath v0.0.0-20180802001716-5cc68e5049a0
)

replace github.com/docker/docker v1.14.0-0.20190319215453-e7b5f7dbe98c => github.com/docker/docker v1.4.2-0.20190319215453-e7b5f7dbe98c
