#!/bin/bash
BASE=/mnt/c/Users/richa/Code
TARGET=$BASE/sites/itforarchivists.com/static/
ROY=$BASE/go/github.com/richardlehane/siegfried/cmd/roy/data
# Refresh sets dir 
rm -rf $TARGET/sets
cp -rf $ROY/sets $TARGET
# Refresh latest sigs
rm -rf $TARGET/latest
mkdir -p $TARGET/latest
cp $ROY/*.sig $TARGET/latest
