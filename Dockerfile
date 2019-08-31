FROM golang:1.12

ENV GOPROXY https://proxy.golang.org

RUN go get -v github.com/mgechev/revive
RUN go get -v github.com/morphy2k/revive-action

COPY entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
