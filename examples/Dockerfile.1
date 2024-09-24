ARG ALPINE_TAG=3.10

FROM alpine:${ALPINE_TAG} AS build
RUN --network=host echo "Hello world" > abc

FROM build AS test
RUN --security=insecure echo "foo" > bar

FROM scratch
COPY --from=build --chown=nobody:nobody abc .
RUN --mount=type=bind,source=./abc,target=/def
CMD ["echo"]
