
#!/bin/bash

set -e

APP_NAME="hpi"
PREFIX="/usr"
BINDIR="$PREFIX/bin"
DESKTOPDIR="$PREFIX/share/applications"
MIMEDIR="$PREFIX/share/mime/packages"
ICONDIR="$PREFIX/share/icons/hicolor/scalable/apps"

echo "==> Checking Go dependencies..."
go mod tidy

echo "==> Building Helwan Package Installer..."
go build -o "$APP_NAME"

echo "==> Installing HPI..."
sudo install -Dm755 "$APP_NAME" "$BINDIR/$APP_NAME"

echo "==> Installing desktop entry..."
sudo install -Dm644 \
    resources/hpi.desktop \
    "$DESKTOPDIR/hpi.desktop"

echo "==> Installing application icon..."
sudo install -Dm644 \
    resources/icons/hpi.svg \
    "$ICONDIR/hpi.svg"

echo "==> Installing MIME definition..."
sudo install -Dm644 \
    resources/hpi-mime.xml \
    "$MIMEDIR/hpi.xml"

echo "==> Updating desktop database..."
sudo update-desktop-database "$DESKTOPDIR"

echo "==> Updating MIME database..."
sudo update-mime-database "$PREFIX/share/mime"

echo "==> Registering HPI as the default package installer..."
xdg-mime default hpi.desktop application/x-arch-package

echo
echo "=========================================="
echo " Helwan Package Installer installed!"
echo "=========================================="
echo
echo "Binary:"
echo "  $BINDIR/$APP_NAME"
echo
echo "Desktop integration:"
echo "  Double-click *.pkg.tar.zst -> HPI"
echo "  Right-click -> Install Package (zst)"
echo
echo "Done."
