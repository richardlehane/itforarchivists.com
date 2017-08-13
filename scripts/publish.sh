#!/bin/bash
# launch appengine
cd "$( dirname "${BASH_SOURCE[0]}")"
#gcloud app deploy .
cd ..
if [[ "$OSTYPE" == "darwin"* ]]; then
 /usr/local/Caskroom/google-cloud-sdk/latest/google-cloud-sdk/platform/google_appengine/appcfg.py update -A itforarchivists -V 1 .
else
 appcfg.py update -A itforarchivists -V 1 .
fi