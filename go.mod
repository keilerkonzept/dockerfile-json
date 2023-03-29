module github.com/keilerkonzept/dockerfile-json

go 1.19

require (
	github.com/moby/buildkit v0.11.5
	github.com/yalp/jsonpath v0.0.0-20180802001716-5cc68e5049a0
)

require (
	github.com/agext/levenshtein v1.2.3 // indirect
	github.com/containerd/typeurl v1.0.2 // indirect
	github.com/docker/docker v23.0.0-rc.1+incompatible // indirect
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/stretchr/testify v1.8.2 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
)

replace (
	github.com/containerd/containerd => github.com/containerd/containerd v1.7.0
	github.com/docker/docker/v23/v20 => github.com/docker/docker v23.0.2+incompatible
)
