#!/bin/bash

set -e

echo "Publishing artifacts"

if [ -z $DEPLOY_TOKEN ]; then
    echo "DEPLOY_TOKEN not set!"
    exit 42
fi

token="$DEPLOY_TOKEN"
gitlab="https://gitlab.com"
api="$gitlab/api/v4"


echo "Uploading the amd64 binary to $gitlab"
out_amd64=$(curl -f \
	   --request POST \
           --header "PRIVATE-TOKEN: $token" \
           --form "file=@$CI_PROJECT_DIR/cumulus-exporter-amd64" \
	   "$api/projects/$CI_PROJECT_ID/uploads")


echo "Response from gitlab is:"
echo "$out_amd64"

echo "Uploading the arm32 binary to $gitlab"
out_arm32=$(curl -f \
	   --request POST \
           --header "PRIVATE-TOKEN: $token" \
           --form "file=@$CI_PROJECT_DIR/cumulus-exporter-arm32" \
	   "$api/projects/$CI_PROJECT_ID/uploads")


echo "Response from gitlab is:"
echo "$out_arm32"

url_amd64=$(echo "$out_amd64" | jq -r '.full_path')
url_arm32=$(echo "$out_arm32" | jq -r '.full_path')

body=$(cat <<JSON
{
  "tag_name": "$CI_COMMIT_TAG",
  "name": "$ref",
  "assets": {
    "links": [
      { "name": "cumulus-exporter-amd64",
        "url": "$gitlab$url_amd64",
        "filepath": "/binaries/cumulus-exporter-amd64"
      },
      { "name": "cumulus-exporter-arm32",
        "url": "$gitlab$url_arm32",
        "filepath": "/binaries/cumulus-exporter-arm32"
      }
    ]
  }
}
JSON
)

echo "Using the following body..."
echo "$body" 

echo "... creating a release"
curl -f \
     -o - \
     --header 'Content-Type: application/json' \
     --header "PRIVATE-TOKEN: $token" \
     --data "$body" \
     --request POST \
     "$api/projects/$CI_PROJECT_ID/releases"

