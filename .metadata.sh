# Application Builder Configuration File. Customized for: hello-world
# Each line must have an export clause.
# This file is parsed and sourced by the Makefile, Docker and Homebrew builds.
# Powered by Application Builder: https://github.com/golift/application-builder

# Bring in dynamic repo/pull/source info.
source $(dirname "${BASH_SOURCE[0]}")/init/buildinfo.sh

# Must match the repo name to make things easy. Otherwise, fix some other paths.
BINARY="discordnotifier-client"
REPO="discordnotifier-client"
# github username
GHUSER="Go-Lift-TV"
# Github repo containing homebrew formula repo.
HBREPO="golift/homebrew-mugs"
MAINT="David Newhall II <david at sleepers dot pro>"
VENDOR="Go Lift"
DESC="Unified Client for DiscordNotifier.com"
GOLANGCI_LINT_ARGS="--enable-all -D exhaustivestruct,nlreturn,forbidigo"
# Example must exist at examples/$CONFIG_FILE.example
CONFIG_FILE="dnclient.conf"
LICENSE="MIT"
# FORMULA is either 'service' or 'tool'. Services run as a daemon, tools do not.
# This affects the homebrew formula (launchd) and linux packages (systemd).
FORMULA="service"

# Used for source links and wiki links.
SOURCE_URL="https://github.com/${GHUSER}/${REPO}"
# Used for documentation links.
URL="${SOURCE_URL}"

# This parameter is passed in as -X to go build. Used to override the Version variable in a package.
# This makes a path like golift.io/version.Version=1.3.3
# Name the Version-containing library the same as the github repo, without dashes.
VERSION_PATH="golift.io/version"

# Used by homebrew downloads.
SOURCE_PATH=https://codeload.github.com/${GHUSER}/${REPO}/tar.gz/v${VERSION}

# Use upx to compress binaries. Must install upx. apt/yum/brew install upx
COMPRESS=true

export BINARY GHUSER HBREPO MAINT VENDOR DESC GOLANGCI_LINT_ARGS CONFIG_FILE
export LICENSE FORMULA SOURCE_URL URL VERSION_PATH SOURCE_PATH COMPRESS

export WINDOWS_LDFLAGS="-H windowsgui"
export MACAPP="DiscordNotifier-Client"