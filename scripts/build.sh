#!/usr/bin/env bash

if [[ "$0" != *scripts/*.sh ]]
then
	echo "Please, execute from project's root directory"
	exit 1
fi

source scripts/build_common.inc.sh

mkdir -p bin
rm -f bin/$EXE_NAME*

IFS='/'
for b in "${BUILD_LIST[@]}"
do
    echo "Building ${b}"
	read -ra THIS <<< "$b"
	OS=${THIS[0]}
	ARCH=${THIS[1]}
	EXT=""
	[[ $ARCH == "windows" ]] && EXT=".exe"
    GOOS=$OS GOARCH=$ARCH go build -ldflags "${FLAGS}" \
		-o "bin/${EXE_NAME}_${VER}-${OS}_${ARCH}${EXT}" \
		retroconvert.go
    if [[ $? -ne 0 ]]
    then
        echo "Compilation error!"
        exit 1
    fi
done
IFS=' '

mkdir -p dist
rm -f dist/*
echo -n "Compressing releases..."
IFS='/'
for b in "${BUILD_LIST[@]}"
do
	read -ra THIS <<< "$b"
	OS=${THIS[0]}
	ARCH=${THIS[1]}
	FILE="${EXE_NAME}_${VER}-${OS}_${ARCH}"
	if [[ $ARCH == "windows" ]]
	then
		zip -q dist/$FILE.zip bin/$FILE.exe
	else
		cp bin/$FILE dist/.
		gzip -9 dist/$FILE
	fi

    if [[ $? -ne 0 ]]
    then
        echo "Compression error!"
        exit 1
    fi
done
IFS=' '
echo "OK"

