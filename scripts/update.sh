#!/bin/bash
TARGET=$HOME/Dropbox/programming/sites/itforarchivists.com/static/
ROY=$HOME/Dropbox/programming/go/github.com/richardlehane/siegfried/cmd/roy/data
cd "$( dirname "${BASH_SOURCE[0]}")"
# Refresh sets dir 
rm -rf $TARGET/sets
cp -rf $ROY/sets $TARGET
# Refresh latest sigs
rm -rf $TARGET/latest
mkdir -p $TARGET/latest
cp $ROY/*.sig $TARGET/latest
