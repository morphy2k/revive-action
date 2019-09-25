FROM golang:1.13

LABEL repository="https://github.com/morphy2k/revive-action"
LABEL homepage="https://github.com/morphy2k/revive-action"
LABEL maintainer="Markus Wiegand <mail@morphy2k.dev>"

LABEL com.github.actions.name="Revive Action"
LABEL com.github.actions.description="Lint your Go code with Revive"
LABEL com.github.actions.icon="code"
LABEL com.github.actions.color="blue"

RUN go get -v github.com/mgechev/revive
RUN go get -v github.com/morphy2k/revive-action

COPY entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
