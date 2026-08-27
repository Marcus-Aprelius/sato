FROM golang:1.26.7-alpine AS builder

ARG VERSION=0.0.0-dev
ARG GIT_COMMIT=unknown

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build \
    -ldflags="-s -w -X sato/internal/sato.Version=v${VERSION} -X sato/internal/sato.GitCommit=${GIT_COMMIT}" \
    -o /out/sato

FROM scratch

COPY --from=builder /out/sato /sato

ENTRYPOINT ["/sato"]
