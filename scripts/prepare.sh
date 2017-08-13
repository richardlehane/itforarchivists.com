#!/bin/bash
LOCAL=~/siegfried/
TARGET=~/Dropbox/programming/sites/itforarchivists.com/public/
ROY=~/Dropbox/programming/go/src/github.com/richardlehane/siegfried/cmd/roy/data
cd "$( dirname "${BASH_SOURCE[0]}")"
rm -rf $TARGET
rm -rf $LOCAL/sets
mkdir -p $TARGET/latest
cp -rf $ROY/sets $LOCAL
cp -rf $ROY/sets $TARGET
cp $ROY/*.sig $TARGET/latest
cp run_go run.go
go run run.go
rm run.go
cd ..
hugo


