FROM golang:1.26.5-alpine

ARG USERNAME=sato
ARG USER_UID=1000
ARG USER_GID=1000

WORKDIR /workspaces/sato

RUN apk add --no-cache \
    bash bash-completion git docker-cli docker-cli-compose docker-cli-buildx \
    coreutils findutils grep sed make ca-certificates sudo acl openssh

RUN addgroup -g "${USER_GID}" "${USERNAME}" \
    && adduser -D -u "${USER_UID}" -G "${USERNAME}" -s /bin/bash "${USERNAME}" \
    && echo "${USERNAME} ALL=(ALL) NOPASSWD:ALL" > /etc/sudoers.d/${USERNAME} \
    && chmod 0440 /etc/sudoers.d/${USERNAME}

RUN mkdir -p /home/${USERNAME}/go/pkg/mod /home/${USERNAME}/.cache/go-build \
    && chown -R ${USERNAME}:${USERNAME} /home/${USERNAME}

ENV CGO_ENABLED=0
ENV GOPATH=/home/${USERNAME}/go
ENV GOMODCACHE=/home/${USERNAME}/go/pkg/mod
ENV GOCACHE=/home/${USERNAME}/.cache/go-build
