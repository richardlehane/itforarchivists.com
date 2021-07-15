#!/bin/bash
# Refresh public folder by deleting then running hugo
rm -rf ../public
hugo
# deploy
gcloud app deploy ../.
