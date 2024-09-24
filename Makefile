VERSION := $(shell jq -r .version <package.json)

APP      := $(shell jq -r .name <package.json)
PACKAGES := $(shell go list -f {{.Dir}} ./...)
GOFILES  := $(addsuffix /*.go,$(PACKAGES))
GOFILES  := $(wildcard $(GOFILES))

.PHONY: clean release binaries README.md

clean:
	rm -rf binaries/
	rm -rf release/

release: README.md zip
	git add README.template.md
	git add README.md
	git add Makefile
	git commit -am "Release $(VERSION)" || true
	git push
	hub release create $(VERSION) -m "$(VERSION)" -a release/$(APP)_$(VERSION)_osx_x86_64.tar.gz -a release/$(APP)_$(VERSION)_windows_x86_64.zip -a release/$(APP)_$(VERSION)_linux_x86_64.tar.gz -a release/$(APP)_$(VERSION)_windows_x86_32.zip -a release/$(APP)_$(VERSION)_linux_x86_32.tar.gz -a release/$(APP)_$(VERSION)_linux_arm64.tar.gz

README.md:
	go get github.com/keilerkonzept/$(APP) && <README.template.md subst \
		EXAMPLE_6="$$($(APP)  --expand-build-args=false --jsonpath=$$..BaseName examples/Dockerfile.3 | jq .)" \
		EXAMPLE_5="$$($(APP) --jsonpath=$$..BaseName examples/Dockerfile.3 | jq .)" \
		EXAMPLE_4="$$($(APP) --jsonpath=$$..Image --build-arg ALPINE_TAG=hello-world examples/Dockerfile.3 | jq .)" \
		EXAMPLE_3B="$$($(APP) --jsonpath=$$..Image --jsonpath-raw examples/Dockerfile.3)" \
		EXAMPLE_3A="$$($(APP) --jsonpath=$$..Image examples/Dockerfile.3 | jq .)" \
		EXAMPLE_2="$$($(APP) --jsonpath=$$..As examples/Dockerfile.2 | jq .)" \
		EXAMPLE_1="$$($(APP) examples/Dockerfile.1 | jq .)" \
		VERSION="$(VERSION)" APP="$(APP)" USAGE="$$($(APP) -h 2>&1)" > README.md

zip: release/$(APP)_$(VERSION)_osx_x86_64.tar.gz release/$(APP)_$(VERSION)_windows_x86_64.zip release/$(APP)_$(VERSION)_linux_x86_64.tar.gz release/$(APP)_$(VERSION)_windows_x86_32.zip release/$(APP)_$(VERSION)_linux_x86_32.tar.gz release/$(APP)_$(VERSION)_linux_arm64.tar.gz release/$(APP)_$(VERSION)_osx_arm64.tar.gz

binaries: binaries/osx_arm64/$(APP) binaries/osx_x86_64/$(APP) binaries/windows_x86_64/$(APP).exe binaries/linux_x86_64/$(APP) binaries/windows_x86_32/$(APP).exe binaries/linux_x86_32/$(APP)

release/$(APP)_$(VERSION)_osx_x86_64.tar.gz: binaries/osx_x86_64/$(APP)
	mkdir -p release
	tar cfz release/$(APP)_$(VERSION)_osx_x86_64.tar.gz -C binaries/osx_x86_64 $(APP)

binaries/osx_x86_64/$(APP): $(GOFILES)
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-X main.version=$(VERSION)" -o binaries/osx_x86_64/$(APP) .

release/$(APP)_$(VERSION)_windows_x86_64.zip: binaries/windows_x86_64/$(APP).exe
	mkdir -p release
	cd ./binaries/windows_x86_64 && zip -r -D ../../release/$(APP)_$(VERSION)_windows_x86_64.zip $(APP).exe

binaries/windows_x86_64/$(APP).exe: $(GOFILES)
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-X main.version=$(VERSION)" -o binaries/windows_x86_64/$(APP).exe .

release/$(APP)_$(VERSION)_linux_x86_64.tar.gz: binaries/linux_x86_64/$(APP)
	mkdir -p release
	tar cfz release/$(APP)_$(VERSION)_linux_x86_64.tar.gz -C binaries/linux_x86_64 $(APP)

binaries/linux_x86_64/$(APP): $(GOFILES)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-X main.version=$(VERSION)" -o binaries/linux_x86_64/$(APP) .

release/$(APP)_$(VERSION)_windows_x86_32.zip: binaries/windows_x86_32/$(APP).exe
	mkdir -p release
	cd ./binaries/windows_x86_32 && zip -r -D ../../release/$(APP)_$(VERSION)_windows_x86_32.zip $(APP).exe

binaries/windows_x86_32/$(APP).exe: $(GOFILES)
	GOOS=windows GOARCH=386 CGO_ENABLED=0 go build -ldflags "-X main.version=$(VERSION)" -o binaries/windows_x86_32/$(APP).exe .

release/$(APP)_$(VERSION)_linux_x86_32.tar.gz: binaries/linux_x86_32/$(APP)
	mkdir -p release
	tar cfz release/$(APP)_$(VERSION)_linux_x86_32.tar.gz -C binaries/linux_x86_32 $(APP)

binaries/linux_x86_32/$(APP): $(GOFILES)
	GOOS=linux GOARCH=386 CGO_ENABLED=0 go build -ldflags "-X main.version=$(VERSION)" -o binaries/linux_x86_32/$(APP) .

release/$(APP)_$(VERSION)_linux_arm64.tar.gz: binaries/linux_arm64/$(APP)
	mkdir -p release
	tar cfz release/$(APP)_$(VERSION)_linux_arm64.tar.gz -C binaries/linux_arm64 $(APP)

binaries/linux_arm64/$(APP): $(GOFILES)
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags "-X main.version=$(VERSION)" -o binaries/linux_arm64/$(APP) .

release/$(APP)_$(VERSION)_osx_arm64.tar.gz: binaries/osx_arm64/$(APP)
	mkdir -p release
	tar cfz release/$(APP)_$(VERSION)_osx_arm64.tar.gz -C binaries/osx_arm64 $(APP)

binaries/osx_arm64/$(APP): $(GOFILES)
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -ldflags "-X main.version=$(VERSION)" -o binaries/osx_arm64/$(APP) .
