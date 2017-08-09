#!/bin/bash
LOCAL=~/siegfried
TARGET=~/Dropbox/programming/sites/itforarchivists.com/public
ROY=~/Dropbox/programming/go/src/github.com/richardlehane/siegfried/cmd/roy/data
cd "$( dirname "${BASH_SOURCE[0]}")"
rm -rf $TARGET
mkdir -p $TARGET/latest
rm -rf $LOCAL/sets
cp -rf $ROY/sets $LOCAL
cp -rf $ROY/sets $TARGET
cp $ROY/default.sig $TARGET/latest
cp $ROY/pronom-tika-loc.sig $TARGET/latest
cp run_go run.go
go run run.go
rm run.go
hugo


