ARG GO_VERSION=1.24.0

FROM golang:${GO_VERSION}-bullseye AS builder

WORKDIR /app

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=bind,target=. \
    go build -o /bin/app ./internal/cmd/app

RUN curl -sSf https://truststore.pki.rds.amazonaws.com/global/global-bundle.pem -o /etc/ssl/certs/rds-ca.pem

FROM gcr.io/distroless/base-debian12:nonroot

COPY --from=builder /bin/app /bin/app
COPY --from=builder /etc/ssl/certs/rds-ca.pem /etc/ssl/certs/rds-ca.pem

ENTRYPOINT ["/bin/app"]
