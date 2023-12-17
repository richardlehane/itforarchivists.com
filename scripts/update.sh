#!/bin/bash
BASE=/mnt/c/Users/richa/Code
TARGET=$BASE/sites/itforarchivists.com/static
ROY=$BASE/go/github.com/richardlehane/siegfried/cmd/roy/data
# Refresh sets dir 
rm $TARGET/sets/*
cp $ROY/sets/* $TARGET/sets/
# Refresh latest sigs
rm $TARGET/latest/1_10/*
cp $ROY/*.sig $TARGET/latest/1_10/
