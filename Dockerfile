FROM --platform=$BUILDPLATFORM golang:1.23.5 AS build-env

ARG VERSION

ARG TARGETOS
ARG TARGETARCH

ENV CGO_ENABLED=0

WORKDIR /src
COPY . .

RUN GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -ldflags="-X 'main.version=${VERSION}'"

FROM ghcr.io/mgechev/revive:1.6.0

LABEL repository="https://github.com/morphy2k/revive-action"
LABEL homepage="https://github.com/morphy2k/revive-action"
LABEL maintainer="Markus Wiegand <mail@morphy.dev>"

LABEL org.opencontainers.image.title="Revive Action"
LABEL org.opencontainers.image.source="https://github.com/morphy2k/revive-action"
LABEL org.opencontainers.image.description="GitHub Action that runs Revive on your Go code"
LABEL org.opencontainers.image.licenses=MIT

LABEL com.github.actions.name="Revive Action"
LABEL com.github.actions.description="GitHub Action that runs Revive on your Go code"
LABEL com.github.actions.icon="code"
LABEL com.github.actions.color="blue"

COPY --from=build-env ["/src/revive-action", "/revive-action"]

ENTRYPOINT ["/revive-action"]
