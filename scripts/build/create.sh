#!/bin/bash

set -e

PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$PROJECT_DIR"
DIST_DIR="$PROJECT_DIR/dist"

BINARY_NAME="sato"
BUILDER_IMAGE="sato:builder"
GO_IMAGE="golang:1.26.5-alpine"

VERSION="0.0.1"

mkdir -p "$DIST_DIR"

show_env() {
    echo "[Environment]"
    echo "OS:      $(uname -srm)"

    if command -v docker >/dev/null 2>&1; then
        echo "Docker:  $(docker --version)"
    else
        echo "Docker:  NOT INSTALLED" >&2
        exit 1
    fi

    if command -v git >/dev/null 2>&1; then
        echo "Git:     $(git --version | awk '{print $3}')"
    fi

    echo "Go:      $GO_IMAGE"
    echo ""
}

build_bin() {

    show_env

    echo "[Building sato binary]"
    echo "Project directory: $PROJECT_DIR"
    echo ""

    GIT_COMMIT=$(
        cd "$PROJECT_DIR" &&
        git rev-parse --short HEAD 2>/dev/null ||
        echo "unknown"
    )

    echo "Building Docker image..."

    docker build \
        --build-arg GIT_COMMIT="$GIT_COMMIT" \
        --target builder \
        -t "$BUILDER_IMAGE" \
        "$PROJECT_DIR"

    echo "Extracting binary..."

    CONTAINER_ID=$(docker create "$BUILDER_IMAGE")

    docker cp \
        "$CONTAINER_ID:/out/$BINARY_NAME" \
        "$DIST_DIR/$BINARY_NAME"

    docker rm "$CONTAINER_ID" >/dev/null

    if [ ! -f "$DIST_DIR/$BINARY_NAME" ]; then
        echo "Build failed"
        exit 1
    fi

    SIZE=$(du -h "$DIST_DIR/$BINARY_NAME" | cut -f1)

    echo ""
    echo "Build successful"
    echo "Binary : $DIST_DIR/$BINARY_NAME"
    echo "Size   : $SIZE"
    echo "Commit : $GIT_COMMIT"
    echo ""
}

build_deb() {

    build_bin

    echo ""
    echo "[Building DEB package]"

    PKG_DIR="$DIST_DIR/deb"

    rm -rf "$PKG_DIR"

    mkdir -p \
        "$PKG_DIR/usr/local/bin" \
        "$PKG_DIR/DEBIAN"

    cp \
        "$DIST_DIR/$BINARY_NAME" \
        "$PKG_DIR/usr/local/bin/sato"

    cat > "$PKG_DIR/DEBIAN/control" <<EOF
Package: sato
Version: ${VERSION}
Section: utils
Priority: optional
Architecture: amd64
Maintainer: SATO Contributors
Description: Secure Access Task Operator
 Like sudo, but for secrets.
EOF

    docker run --rm \
        -v "$DIST_DIR:/dist" \
        ubuntu:24.04 \
        bash -c '
            apt-get update >/dev/null &&
            apt-get install -y dpkg >/dev/null &&
            dpkg-deb --build /dist/deb /dist/sato_'${VERSION}'_amd64.deb
        '

    rm -rf "$PKG_DIR"

    echo ""
    echo "Created:"
    echo "  $DIST_DIR/sato_${VERSION}_amd64.deb"
}

build_rpm() {

    build_bin

    echo ""
    echo "[Building RPM package]"

    PKG_DIR="$DIST_DIR/rpm"

    rm -rf "$PKG_DIR"

    mkdir -p \
        "$PKG_DIR/BUILD" \
        "$PKG_DIR/RPMS" \
        "$PKG_DIR/SOURCES" \
        "$PKG_DIR/SPECS" \
        "$PKG_DIR/SRPMS"

    cp "$DIST_DIR/$BINARY_NAME" \
       "$PKG_DIR/SOURCES/sato"

    cat > "$PKG_DIR/SPECS/sato.spec" <<EOF
Name: sato
Version: ${VERSION}
Release: 1
Summary: Secure Access Task Operator
License: Apache-2.0

BuildArch: x86_64

%description
Like sudo, but for secrets.

%install
mkdir -p %{buildroot}/usr/local/bin
install -m 755 %{_sourcedir}/sato %{buildroot}/usr/local/bin/sato

%files
/usr/local/bin/sato
EOF

    docker run --rm \
        -v "$DIST_DIR/rpm:/root/rpmbuild" \
        rockylinux:9 \
        bash -c '
            dnf install -y rpm-build >/dev/null &&
            rpmbuild -bb /root/rpmbuild/SPECS/sato.spec
        '

    RPM_FILE=$(find "$PKG_DIR/RPMS" -name '*.rpm' | head -1)

    if [ -n "$RPM_FILE" ]; then
        cp "$RPM_FILE" "$DIST_DIR/"
    fi

    rm -rf "$PKG_DIR"

    echo ""
    echo "Created RPM package:"
    ls -1 "$DIST_DIR"/*.rpm
}

case "$1" in

    bin)
        build_bin
        ;;

    deb)
        build_deb
        ;;

    rpm)
        build_rpm
        ;;

    *)
        echo "Usage:"
        echo "  $0 bin"
        echo "  $0 deb"
        echo "  $0 rpm"
        exit 1
        ;;

esac
