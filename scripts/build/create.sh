#!/bin/bash

set -euo pipefail

PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$PROJECT_DIR"

DIST_DIR="$PROJECT_DIR/dist"

BINARY_NAME="sato"
BUILDER_IMAGE="sato:builder"
GO_IMAGE="golang:1.26.5-alpine"

VERSION="0.0.3"

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

    echo "[Building SATO]"
    echo "Project directory: $PROJECT_DIR"
    echo ""

    GIT_COMMIT=$(
        cd "$PROJECT_DIR" &&
        git rev-parse --short=8 HEAD 2>/dev/null ||
        echo "unknown"
    )

    echo "Building Docker image..."

    docker build \
        --build-arg VERSION="$VERSION" \
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

ensure_bin() {
    if [ ! -f "$DIST_DIR/$BINARY_NAME" ]; then
        build_bin
    fi
}

build_deb_package() {
    (
        ensure_bin

        echo ""
        echo "[Building DEB package]"

        PKG_DIR="$DIST_DIR/deb"
        HOST_UID="$(id -u)"
        HOST_GID="$(id -g)"

        cleanup_deb() {
            rm -rf "$PKG_DIR" 2>/dev/null || true
        }

        trap cleanup_deb EXIT

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
                dpkg-deb --build /dist/deb /dist/sato_'${VERSION}'_amd64.deb &&
                chown -R '"$HOST_UID:$HOST_GID"' /dist/deb /dist/sato_'${VERSION}'_amd64.deb
            '

        if [ ! -f "$DIST_DIR/sato_${VERSION}_amd64.deb" ]; then
            echo "DEB build failed"
            exit 1
        fi

        echo ""
        echo "Created:"
        echo "  $DIST_DIR/sato_${VERSION}_amd64.deb"
    )
}

build_rpm_package() {
    (
        ensure_bin

        echo ""
        echo "[Building RPM package]"

        PKG_DIR="$DIST_DIR/rpm"
        HOST_UID="$(id -u)"
        HOST_GID="$(id -g)"

        cleanup_rpm() {
            rm -rf "$PKG_DIR" 2>/dev/null || true
        }

        trap cleanup_rpm EXIT

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
                rpmbuild -bb /root/rpmbuild/SPECS/sato.spec &&
                chown -R '"$HOST_UID:$HOST_GID"' /root/rpmbuild
            '

        RPM_FILE=$(find "$PKG_DIR/RPMS" -name '*.rpm' | head -1)

        if [ -z "$RPM_FILE" ]; then
            echo "RPM build failed"
            exit 1
        fi

        cp "$RPM_FILE" "$DIST_DIR/"

        echo ""
        echo "Created RPM package:"
        ls -1 "$DIST_DIR"/*.rpm
    )
}

build_all() {
    build_bin
    build_deb_package
    build_rpm_package

    echo ""
    echo "[All artifacts created]"
    echo "  $DIST_DIR/$BINARY_NAME"
    echo "  $DIST_DIR/sato_${VERSION}_amd64.deb"

    if ls "$DIST_DIR"/*.rpm >/dev/null 2>&1; then
        ls -1 "$DIST_DIR"/*.rpm | sed 's/^/  /'
    fi
}

case "${1:-}" in

    bin)
        build_bin
        ;;

    deb)
        build_deb_package
        ;;

    rpm)
        build_rpm_package
        ;;

    all)
        build_all
        ;;

    *)
        echo "Usage:"
        echo "  $0 bin"
        echo "  $0 deb"
        echo "  $0 rpm"
        echo "  $0 all"
        exit 1
        ;;

esac
