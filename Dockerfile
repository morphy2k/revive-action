FROM golang:1.23.2 AS build-env

ARG ACTION_VERSION=unknown
ARG REVIVE_VERSION=v1.4.0

ENV CGO_ENABLED=0

RUN go install -v -ldflags="-X 'github.com/mgechev/revive/cli.version=${REVIVE_VERSION}'" \
    github.com/mgechev/revive@${REVIVE_VERSION}

WORKDIR /tmp/github.com/morphy2k/revive-action
COPY . .

RUN go install -ldflags="-X 'main.version=${ACTION_VERSION}'"

FROM alpine:3.20.3

LABEL repository="https://github.com/morphy2k/revive-action"
LABEL homepage="https://github.com/morphy2k/revive-action"
LABEL maintainer="Markus Wiegand <mail@morphy.dev>"

LABEL org.opencontainers.image.source = "https://github.com/morphy2k/revive-action"
LABEL org.opencontainers.image.description="GitHub Action that runs Revive on your Go code"
LABEL org.opencontainers.image.licenses=MIT

LABEL com.github.actions.name="Revive Action"
LABEL com.github.actions.description="GitHub Action that runs Revive on your Go code"
LABEL com.github.actions.icon="code"
LABEL com.github.actions.color="blue"

COPY --from=build-env ["/go/bin/revive", "/go/bin/revive-action", "/bin/"]
COPY --from=build-env /tmp/github.com/morphy2k/revive-action/entrypoint.sh /

RUN apk add --no-cache bash gawk

ENTRYPOINT ["/entrypoint.sh"]
