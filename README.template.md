# ${APP}

Prints Dockerfiles as JSON to stdout, optionally evaluates build args. Plays well with `jq`.

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
${APP}

${USAGE}
```

## Examples

### JSON output

`Dockerfile`
```Dockerfile
ARG ALPINE_TAG=3.10

FROM alpine:${ALPINE_TAG} AS build
RUN echo "Hello world" > abc

FROM scratch
COPY --from=build --chown=nobody:nobody abc .
CMD ["echo"]
```

```sh
$ dockerfile-json Dockerfile | jq .
```
```json
{
  "MetaArgs": [
    {
      "Key": "ALPINE_TAG",
      "Value": "3.10"
    }
  ],
  "Stages": [
    {
      "Name": "build",
      "BaseName": "alpine:3.10",
      "SourceCode": "FROM alpine:${ALPINE_TAG} AS build",
      "Platform": "",
      "IsDerived": false,
      "IsScratch": false,
      "Commands": [
        {
          "Name": "run",
          "Command": {
            "CmdLine": [
              "echo \"Hello world\" > abc"
            ],
            "PrependShell": true
          }
        }
      ]
    },
    {
      "Name": "",
      "BaseName": "scratch",
      "SourceCode": "FROM scratch",
      "Platform": "",
      "IsDerived": false,
      "IsScratch": true,
      "Commands": [
        {
          "Name": "copy",
          "Command": {
            "SourcesAndDest": [
              "abc",
              "."
            ],
            "From": "build",
            "Chown": "nobody:nobody"
          }
        },
        {
          "Name": "cmd",
          "Command": {
            "CmdLine": [
              "echo"
            ],
            "PrependShell": false
          }
        }
      ]
    }
  ]
}
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
$ dockerfile-json Dockerfile | jq '.Stages[] | .Name | select(. != "")'
```
```json
"build"
"test"
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

```sh
$ dockerfile-json Dockerfile |
    jq '.Stages[] | select((.IsDerived or .IsScratch)|not) | .BaseName'
```
```json
"alpine:3.10"
```

#### Set build args, omit stage aliases and `scratch`

```sh
$ dockerfile-json --build-arg ALPINE_TAG=hello-world Dockerfile |
    jq '.Stages[] | select((.IsDerived or .IsScratch)|not) | .BaseName'
```
```json
"alpine:hello-world"
```

#### Expand build args, include all base names

```sh
$  dockerfile-json Dockerfile | jq '.Stages[] | .BaseName'
```
```json
"alpine:3.10"
"build"
"scratch"
```

#### Ignore build args, include all base names

```sh
$ dockerfile-json --expand-build-args=false Dockerfile | jq '.Stages[] | .BaseName'
```
```json
"alpine:${ALPINE_TAG}"
"build"
"${APP_BASE}"
```
