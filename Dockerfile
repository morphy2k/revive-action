FROM golang:1.19.1 as build-env

ARG ACTION_VERSION=unknown
ARG REVIVE_VERSION=v1.2.4

ENV CGO_ENABLED=0

RUN go install -v -ldflags="-X 'main.version=${REVIVE_VERSION}'" \
    github.com/mgechev/revive@${REVIVE_VERSION}

WORKDIR /tmp/github.com/morphy2k/revive-action
COPY . .

RUN go install -ldflags="-X 'main.version=${ACTION_VERSION}'"

FROM alpine:3.17.0

LABEL repository="https://github.com/morphy2k/revive-action"
LABEL homepage="https://github.com/morphy2k/revive-action"
LABEL maintainer="Markus Wiegand <mail@morphy2k.dev>"

LABEL com.github.actions.name="Revive Action"
LABEL com.github.actions.description="Lint your Go code with Revive"
LABEL com.github.actions.icon="code"
LABEL com.github.actions.color="blue"

COPY --from=build-env ["/go/bin/revive", "/go/bin/revive-action", "/bin/"]
COPY --from=build-env /tmp/github.com/morphy2k/revive-action/entrypoint.sh /

RUN apk add --no-cache bash gawk

ENTRYPOINT ["/entrypoint.sh"]
