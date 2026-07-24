FROM rockylinux:9.3

RUN dnf install -y rpm-build ruby ruby-devel gcc make tar gzip findutils which \
    && dnf clean all

RUN gem install --no-document fpm

WORKDIR /workspace

ENTRYPOINT ["/bin/bash"]
