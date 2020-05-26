module github.com/keilerkonzept/dockerfile-json

go 1.13

require (
	github.com/moby/buildkit v0.7.1
	github.com/yalp/jsonpath v0.0.0-20180802001716-5cc68e5049a0
)


replace github.com/containerd/containerd => github.com/containerd/containerd v1.3.4
replace github.com/docker/docker => github.com/docker/docker v1.4.2-0.20200227233006-38f52c9fec82
