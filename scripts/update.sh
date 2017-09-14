#!/bin/bash
LOCAL=~/siegfried/
TARGET=~/Dropbox/programming/sites/itforarchivists.com/static/
ROY=~/Dropbox/programming/go/src/github.com/richardlehane/siegfried/cmd/roy/data
cd "$( dirname "${BASH_SOURCE[0]}")"
# Refresh sets dir for local and static
rm -rf $LOCAL/sets
rm -rf $TARGET/sets
cp -rf $ROY/sets $LOCAL
cp -rf $ROY/sets $TARGET
# Refresh latest sigs in static
rm -rf $TARGET/latest
mkdir -p $TARGET/latest
cp $ROY/*.sig $TARGET/latest
# Refresh sf data
cp run_go run.go
go run run.go
rm run.go