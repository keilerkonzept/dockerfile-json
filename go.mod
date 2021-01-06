module github.com/keilerkonzept/dockerfile-json

go 1.13

require (
	github.com/moby/buildkit v0.8.1
	github.com/yalp/jsonpath v0.0.0-20180802001716-5cc68e5049a0
)

replace (
	github.com/containerd/containerd => github.com/containerd/containerd v1.3.1-0.20200512144102-f13ba8f2f2fd
	github.com/docker/docker => github.com/docker/docker v17.12.0-ce-rc1.0.20200310163718-4634ce647cf2+incompatible
)
