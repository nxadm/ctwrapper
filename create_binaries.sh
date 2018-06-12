#!/bin/bash -e
set -xv
APP=ctwrapper
BIN_DIR=/var/tmp/$APP
BASE_NAME="$BIN_DIR/$APP"
PLATFORMS=("windows/amd64" "windows/386" "darwin/amd64" "darwin/386" "linux/amd64" "linux/386")

function build {
    GOOS=$1
    GOARCH=$2
    OUTPUT="${BASE_NAME}-${GOOS}-${GOARCH}"
    if [ $GOOS = "windows" ]; then
        OUTPUT+='.exe'
    fi
    GOOS=$GOOS GOARCH=$GOARCH go build -o $OUTPUT
    cd $BIN_DIR
    sha512sum $OUTPUT > $OUTPUT.sha512
}

mkdir -p $BIN_DIR
for i in ${PLATFORMS[@]}; do
    PLATFORMS_SPLIT=(${i//\// })
    GOOS=${PLATFORMS_SPLIT[0]}
    GOARCH=${PLATFORMS_SPLIT[1]}
    build $GOOS $GOARCH
done
