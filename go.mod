module github.com/keilerkonzept/dockerfile-json

go 1.19

require (
	github.com/moby/buildkit v0.11.2
	github.com/yalp/jsonpath v0.0.0-20180802001716-5cc68e5049a0
)

require (
	github.com/Sirupsen/logrus v1.0.6 // indirect
	github.com/agext/levenshtein v1.2.3 // indirect
	github.com/containerd/containerd v1.6.16-0.20230124210447-1709cfe273d9 // indirect
	github.com/containerd/typeurl v1.0.2 // indirect
	github.com/coreos/go-systemd/v22 v22.4.0 // indirect
	github.com/cyphar/filepath-securejoin v0.2.3 // indirect
	github.com/docker/distribution v2.8.1+incompatible // indirect
	github.com/docker/docker v23.0.0-rc.1+incompatible // indirect
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/docker/libnetwork v0.5.6 // indirect
	github.com/godbus/dbus/v5 v5.0.6 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/moby/sys/mountinfo v0.6.2 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.0.3-0.20220303224323-02efb9a75ee1 // indirect
	github.com/opencontainers/runc v1.1.3 // indirect
	github.com/opencontainers/runtime-spec v1.0.3-0.20210326190908-1c3f411f0417 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/sirupsen/logrus v1.9.0 // indirect
	github.com/vishvananda/netlink v1.1.1-0.20210330154013-f5de75959ad5 // indirect
	github.com/vishvananda/netns v0.0.0-20210104183010-2eb08e3e575f // indirect
	golang.org/x/crypto v0.2.0 // indirect
	golang.org/x/sys v0.3.0 // indirect
	golang.org/x/term v0.2.0 // indirect
	google.golang.org/genproto v0.0.0-20220706185917-7780775163c4 // indirect
	google.golang.org/grpc v1.50.1 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
)

replace (
	github.com/containerd/containerd => github.com/containerd/containerd v1.6.17
	github.com/docker/docker/v23/v20 => github.com/docker/docker v23.0.1+incompatible
)
