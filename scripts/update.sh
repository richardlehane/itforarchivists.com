#!/bin/bash
if [[ -z "${WIN_HOME}" ]]; then
  MY_HOME=$HOME
else
  MY_HOME=$WIN_HOME
fi
TARGET=$MY_HOME/Dropbox/programming/sites/itforarchivists.com/static/
ROY=$MY_HOME/Dropbox/programming/go/src/github.com/richardlehane/siegfried/cmd/roy/data
cd "$( dirname "${BASH_SOURCE[0]}")"
# Refresh sets dir 
rm -rf $TARGET/sets
cp -rf $ROY/sets $TARGET
# Refresh latest sigs
rm -rf $TARGET/latest
mkdir -p $TARGET/latest
cp $ROY/*.sig $TARGET/latest
# Refresh sf data
cp run_go run.go
go run run.go
rm run.go