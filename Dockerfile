FROM golang:1.26.5-alpine AS builder

ARG GIT_COMMIT=unknown

WORKDIR /src

COPY go.mod ./
COPY . .
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build \
    -ldflags="-s -w -X sato/internal/sato.GitCommit=${GIT_COMMIT}" \
    -o /out/sato

FROM scratch

COPY --from=builder /out/sato /sato

ENTRYPOINT ["/sato"]
