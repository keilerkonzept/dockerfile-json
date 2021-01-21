# ${APP}

Prints Dockerfiles as JSON to stdout, optionally evaluates build args. Uses the [official Dockerfile parser](https://github.com/moby/buildkit/blob/master/frontend/dockerfile/) from buildkit. Plays well with `jq`.

## Contents

- [Contents](#contents)
- [Get it](#get-it)
- [Usage](#usage)
- [Examples](#examples)
  - [JSON output](#json-output)
  - [Extract build stage names](#extract-build-stage-names)
  - [Extract base images](#extract-base-images)
    - [Expand build args, omit stage aliases and `scratch`](#expand-build-args-omit-stage-aliases-and-scratch)
    - [Set build args, omit stage aliases and `scratch`](#set-build-args-omit-stage-aliases-and-scratch)
    - [Expand build args, include all base names](#expand-build-args-include-all-base-names)
    - [Ignore build args, include all base names](#ignore-build-args-include-all-base-names)

## Get it

Using go get:

```bash
go get -u github.com/keilerkonzept/${APP}
```

Or [download the binary for your platform](https://github.com/keilerkonzept/${APP}/releases/latest) from the releases page.

## Usage

### CLI

```text
${APP} [PATHS...]

${USAGE}
```

## Examples

### JSON output

`Dockerfile`
```Dockerfile
ARG ALPINE_TAG=3.10

FROM alpine:${ALPINE_TAG} AS build
RUN echo "Hello world" > abc

FROM build AS test
RUN echo "foo" > bar

FROM scratch
COPY --from=build --chown=nobody:nobody abc .
CMD ["echo"]
```

```sh
$ dockerfile-json Dockerfile | jq .
```
```json
${EXAMPLE_1}
```

### Extract build stage names

`Dockerfile`
```Dockerfile
FROM maven:alpine AS build
# ...

FROM build AS test
# ...

FROM openjdk:jre-alpine
# ...
```

```sh
$ dockerfile-json --jsonpath=..As Dockerfile
```
```json
${EXAMPLE_2}
```

### Extract base images

`Dockerfile`
```Dockerfile
ARG ALPINE_TAG=3.10
ARG APP_BASE=scratch

FROM alpine:${ALPINE_TAG} AS build
# ...

FROM build
# ...

FROM ${APP_BASE}
# ...
```

#### Expand build args, omit stage aliases and `scratch`

Using `jq`:
```sh
$ dockerfile-json Dockerfile |
    jq '.Stages[] | select(.From | .Stage or .Scratch | not) | .BaseName'
```
```json
${EXAMPLE_3A}
```

Using `--jsonpath`:
```sh

$ dockerfile-json --jsonpath=..Image Dockerfile
```
```json
${EXAMPLE_3A}
```

Using `--jsonpath`, `--jsonpath-raw` output:
```sh
$ dockerfile-json --jsonpath=..Image --jsonpath-raw Dockerfile
```
```json
${EXAMPLE_3B}
```

#### Set build args, omit stage aliases and `scratch`

```sh
$ dockerfile-json --build-arg ALPINE_TAG=hello-world --jsonpath=..Image Dockerfile
```
```json
${EXAMPLE_4}
```

#### Expand build args, include all base names

```sh
$  dockerfile-json --jsonpath=..BaseName Dockerfile
```
```json
${EXAMPLE_5}
```

#### Ignore build args, include all base names

```sh
$ dockerfile-json --expand-build-args=false --jsonpath=..BaseName Dockerfile
```
```json
${EXAMPLE_6}
```
