#!/bin/sh
set -xe

SCRIPT_DIR="$(dirname "$(readlink -f "$0")")"

readonly VERSION="${1:-0.0.1}"
readonly APPDIR="build/bin/TorPlayer.app"

mkdir -p $APPDIR/Contents/{MacOS,Resources}

go build -o $APPDIR/Contents/MacOS/TorPlayer
chmod +x $APPDIR/Contents/MacOS/*

cat > $APPDIR/Contents/Info.plist << EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
    <dict>
        <key>CFBundlePackageType</key>
        <string>APPL</string>
        <key>CFBundleName</key>
        <string>TorPlayer</string>
        <key>CFBundleExecutable</key>
        <string>TorPlayer</string>
        <key>CFBundleIdentifier</key>
        <string>com.TorPlayer</string>
        <key>CFBundleVersion</key>
        <string>$VERSION</string>
        <key>CFBundleGetInfoString</key>
        <string>Built using Go</string>
        <key>CFBundleShortVersionString</key>
        <string>$VERSION</string>
        <key>CFBundleIconFile</key>
        <string>icon.icns</string>
        <key>LSMinimumSystemVersion</key>
        <string>10.13.0</string>
        <key>NSHighResolutionCapable</key>
        <string>true</string>
        <key>NSHumanReadableCopyright</key>
        <string>Copyright.........</string>
    </dict>
</plist>
EOF

cp $SCRIPT_DIR/icons/icon.icns $APPDIR/Contents/Resources/icon.icns
find $APPDIR

productbuild --component $APPDIR ./build/bin/TorPlayer_${VERSION}.pkg