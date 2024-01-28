#!/bin/sh
set -xe

SCRIPT_DIR="$(dirname "$(readlink -f "$0")")"

readonly VERSION="${1:-0.0.1}"
readonly APP=TorPlayer
readonly APPDIR=build/bin/${APP}_${VERSION}


mkdir -p "$APPDIR/usr/bin"
mkdir -p "$APPDIR/usr/share/applications"
mkdir -p "$APPDIR/usr/share/icons/hicolor/1024x1024/apps"
mkdir -p "$APPDIR/usr/share/icons/hicolor/256x256/apps"
mkdir -p "$APPDIR/DEBIAN"

go build -o "$APPDIR/usr/bin/$APP"
chmod +x "$APPDIR/usr/bin/$APP"

cp "$SCRIPT_DIR/icons/icon.png" "$APPDIR/usr/share/icons/hicolor/1024x1024/apps/${APP}.png"
cp "$SCRIPT_DIR/icons/icon.png" "$APPDIR/usr/share/icons/hicolor/256x256/apps/${APP}.png"

cat > "$APPDIR/usr/share/applications/${APP}.desktop" << EOF
[Desktop Entry]
Version=1.0
Type=Application
Name=$APP
Exec=$APP
Icon=$APP
Terminal=false
StartupWMClass=TorPlayer
EOF

cat > "$APPDIR/DEBIAN/control" << EOF
Package: ${APP}
Version: 1.0-0
Section: base
Priority: optional
Architecture: amd64
Maintainer: Chien Nguyen
Description: TorPlayer is a simple video player from the torrent network.
EOF

dpkg-deb --build "$APPDIR"