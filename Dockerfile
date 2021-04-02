FROM golang:1.16.3 as build-env

ENV GO111MODULE=on
ENV CGO_ENABLED=0

RUN go get -v github.com/mgechev/revive@v1

WORKDIR /tmp/github.com/morphy2k/revive-action
COPY . .

RUN go install

FROM alpine:3.13.0

LABEL repository="https://github.com/morphy2k/revive-action"
LABEL homepage="https://github.com/morphy2k/revive-action"
LABEL maintainer="Markus Wiegand <mail@morphy2k.dev>"

LABEL com.github.actions.name="Revive Action"
LABEL com.github.actions.description="Lint your Go code with Revive"
LABEL com.github.actions.icon="code"
LABEL com.github.actions.color="blue"

COPY --from=build-env ["/go/bin/revive", "/go/bin/revive-action", "/bin/"]
COPY --from=build-env /tmp/github.com/morphy2k/revive-action/entrypoint.sh /

RUN apk add --no-cache bash

ENTRYPOINT ["/entrypoint.sh"]
