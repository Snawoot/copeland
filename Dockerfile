FROM --platform=$BUILDPLATFORM golang AS build

ARG GIT_DESC=undefined

WORKDIR /go/src/github.com/Snawoot/copeland
COPY . .
ARG TARGETOS TARGETARCH
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH CGO_ENABLED=0 go build -a -tags netgo -ldflags '-s -w -extldflags "-static" -X main.version='"$GIT_DESC" ./cmd/copeland

FROM scratch
COPY --from=build /go/src/github.com/Snawoot/copeland/copeland /
USER 9999:9999
ENTRYPOINT ["/copeland"]
