#!/bin/bash
cd "$( dirname "${BASH_SOURCE[0]}")/.."
#goapp serve --admin_port=8001 --port=8081
# /usr/local/Caskroom/google-cloud-sdk/latest/google-cloud-sdk/bin/dev_appserver.py --port=8081 .
if [[ "$OSTYPE" == "darwin"* ]]; then
 /usr/local/Caskroom/google-cloud-sdk/latest/google-cloud-sdk/platform/google_appengine/dev_appserver.py --port=8081 --default_gcs_bucket_name itforarchivists.appspot.com .
else
 dev_appserver.py --port=8081 --default_gcs_bucket_name itforarchivists.appspot.com .
fi