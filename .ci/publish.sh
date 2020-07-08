#!/bin/bash

set -e

echo "Publishing artifacts"

token="$DEPLOY_TOKEN"
ref="$1"
namespace="wobcom"
proj_name="cumulus-exporter"
proj="$namespace%2F$proj_name"
gitlab="https://gitlab.com"
api="$gitlab/api/v4"


echo "Uploading the binary to $gitlab"
out=$(curl -f \
	   --request POST \
           --header "PRIVATE-TOKEN: $token" \
           --form "file=@$CI_PROJECT_DIR/cumulus-exporter" \
	   "$api/projects/$proj/uploads")


echo "Response from gitlab is:"
echo "$out"
url=$(echo "$out" | jq -r '.full_path')

body=$(cat <<JSON
{
  "ref": "$ref",
  "tag_name": "$ref",
  "name": "$ref",
  "assets": {
    "links": [
      { "name": "cumulus-exporter",
        "url": "$gitlab$url",
        "filepath": "/binaries/cumulus-exporter"
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
     "$api/projects/$proj/releases"
