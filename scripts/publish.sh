#!/bin/bash
cd "$( dirname "${BASH_SOURCE[0]}")"
cd ..
# Refresh public folder by deleting then running hugo
rm -rf public
hugo
# deploy
gcloud app deploy ../.
