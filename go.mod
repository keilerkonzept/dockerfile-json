module github.com/keilerkonzept/dockerfile-json

go 1.21

toolchain go1.22.1

require (
	github.com/moby/buildkit v0.13.1
	github.com/yalp/jsonpath v0.0.0-20180802001716-5cc68e5049a0
)

require (
	github.com/Sirupsen/logrus v1.0.6 // indirect
	github.com/agext/levenshtein v1.2.3 // indirect
	github.com/containerd/containerd v1.7.13 // indirect
	github.com/containerd/typeurl v1.0.2 // indirect
	github.com/containerd/typeurl/v2 v2.1.1 // indirect
	github.com/coreos/go-systemd/v22 v22.5.0 // indirect
	github.com/cyphar/filepath-securejoin v0.2.4 // indirect
	github.com/docker/distribution v2.8.2+incompatible // indirect
	github.com/docker/docker v25.0.3+incompatible // indirect
	github.com/docker/go-connections v0.5.0 // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/docker/libnetwork v0.5.6 // indirect
	github.com/godbus/dbus/v5 v5.1.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/moby/sys/mountinfo v0.7.1 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.0-rc5 // indirect
	github.com/opencontainers/runc v1.1.7 // indirect
	github.com/opencontainers/runtime-spec v1.1.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/vishvananda/netlink v1.2.1-beta.2 // indirect
	github.com/vishvananda/netns v0.0.4 // indirect
	golang.org/x/crypto v0.17.0 // indirect
	golang.org/x/sys v0.16.0 // indirect
	golang.org/x/term v0.15.0 // indirect
	google.golang.org/genproto v0.0.0-20231016165738-49dd2c1f3d0b // indirect
	google.golang.org/grpc v1.59.0 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
)

replace (
	github.com/containerd/containerd => github.com/containerd/containerd v1.7.14
	github.com/docker/docker/v23/v20 => github.com/docker/docker v26.0.0+incompatible
)
