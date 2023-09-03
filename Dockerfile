# Builder Image
# ---------------------------------------------------
FROM dimaskiddo/alpine:go-1.20 AS go-builder

WORKDIR /usr/src/app

COPY . ./

RUN go mod download \
    && CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -a -o object-storage-proxy cmd/main/main.go


# Final Image
# ---------------------------------------------------
FROM dimaskiddo/alpine:base
MAINTAINER Dimas Restu Hidayanto <drh.dimasrestu@gmail.com>

ARG SERVICE_NAME="object-storage-proxy"
ENV PATH="$PATH:/usr/app/${SERVICE_NAME}"

WORKDIR /usr/app/${SERVICE_NAME}

COPY --from=go-builder /usr/src/app/object-storage-proxy ./object-storage-proxy

CMD ["object-storage-proxy", "proxy"]
