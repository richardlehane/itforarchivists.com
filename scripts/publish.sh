#!/bin/bash
cd "$( dirname "${BASH_SOURCE[0]}")"
cd ..
# Refresh public folder by deleting then running hugo
rm -rf public
hugo
# launch appengine (when fixed: gcloud app deploy .)
if [[ "$OSTYPE" == "darwin"* ]]; then
 /usr/local/Caskroom/google-cloud-sdk/latest/google-cloud-sdk/platform/google_appengine/appcfg.py update -A itforarchivists -V 1 .
else
 appcfg.py update -A itforarchivists -V 1 .
fi