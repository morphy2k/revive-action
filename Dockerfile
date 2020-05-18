FROM golang:1.14.3

LABEL repository="https://github.com/morphy2k/revive-action"
LABEL homepage="https://github.com/morphy2k/revive-action"
LABEL maintainer="Markus Wiegand <mail@morphy2k.dev>"

LABEL com.github.actions.name="Revive Action"
LABEL com.github.actions.description="Lint your Go code with Revive"
LABEL com.github.actions.icon="code"
LABEL com.github.actions.color="blue"

ENV GO111MODULE=on

RUN go get -v github.com/mgechev/revive@v1
RUN go get -v github.com/morphy2k/revive-action

COPY entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
