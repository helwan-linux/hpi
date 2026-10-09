# Helwan Package Installer (HPI)

**Helwan Package Installer (HPI)** is a simple graphical package installer for Arch Linux and Helwan Linux.

It provides an easy way to install local `.pkg.tar.zst` packages without requiring the user to manually run `pacman` commands in a terminal.

## Features

* 📦 Install `.pkg.tar.zst` packages graphically
* 🔍 Display package information before installation
* 🔗 Handle package dependencies through Pacman
* 🌐 Resolve available repository dependencies automatically
* 🔐 Use PolicyKit for administrative privileges
* 🖱️ Double-click integration for `.pkg.tar.zst` files
* 📁 Nemo right-click integration
* 🖥️ GTK3 graphical interface
* ⚡ Lightweight and simple
* 🚫 No daemon or background service
* 🌍 Designed for Helwan Linux and Arch Linux

## How It Works

Simply open a local package:

```text
Double-click package.pkg.tar.zst
        ↓
Helwan Package Installer
        ↓
Package information
        ↓
Dependency check
        ↓
Installation through Pacman
        ↓
Done
```

HPI does not replace Pacman. It provides a graphical interface on top of the existing Arch Linux package management system.

## Installation

### From the Helwan Linux package

If an HPI package is available:

```bash
sudo pacman -U hpi-1.0.0-1-x86_64.pkg.tar.zst
```

### Build from source

Clone the repository:

```bash
git clone https://github.com/helwan-linux/hpi.git
cd hpi
```

Build and install:

```bash
makepkg -si
```

Or use the included installation script:

```bash
chmod +x install.sh
./install.sh
```

## Usage

Open a local package with:

```bash
hpi package.pkg.tar.zst
```

You can also simply double-click a `.pkg.tar.zst` file from the file manager.

In Nemo, you can right-click a package and select:

**Install Package (zst)**

## Dependencies

HPI uses the following main system components:

* GTK3
* Pacman
* PolicyKit
* Go

The actual package installation and dependency resolution are handled by Pacman.

## Project Structure

```text
hpi/
├── gui/            # Graphical interface
├── installer/      # Installation and Pacman integration
├── package/        # Package parsing and metadata
├── i18n/           # Translation support
├── resources/      # Desktop, MIME, icon and Nemo integration
├── system/         # Desktop integration
├── tests/          # Tests
├── main.go
├── install.sh
└── PKGBUILD
```

## Philosophy

HPI follows the Helwan Linux philosophy:

> **Keep things simple.**

The goal is not to replace the powerful tools already provided by Arch Linux.

Instead, HPI makes common package installation tasks easier for users who prefer a graphical workflow.

## Status

HPI is under active development.

The core installation workflow is functional and has been tested with real Arch Linux packages.

sudo pacman -S go gtk3
git clone https://github.com/helwan-linux/hpi.git
cd hpi
go build -o hpi
./hpi

## License

HPI is free and open-source software released under the **GNU General Public License v3.0 or later**.

## Links

* **Helwan Linux:** https://helwan-linux.github.io/helwanlinux/
* **Helwan Linux GitHub:** https://github.com/helwan-linux
* **HPI Repository:** https://github.com/helwan-linux/hpi
