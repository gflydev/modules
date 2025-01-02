#!/bin/sh
#
# Upload file to Wasabi S3
#

file_path='../../../storages3/examples/public/assets'
file_name='hero.png'
file_type='image/png'
bucket_name='avatars'

echo "\n # ============================== Get Pre-sign URL ============================== \n"
url='http://localhost:7789/api/v1/storage/presigned-url'

# Get Pre-sign URL content
json=$(curl "$url?filename=$file_name")

echo $json | jq

# Extract parameters
upload_url=$(echo $json | jq -r '.upload_url')
file_url=$(echo $json | jq -r '.file_url')

echo "\n # ============================== Upload to S3 ============================== \n"
json=$(curl --request PUT \
    --upload-file "$file_path/$file_name" \
    "$upload_url")

echo $json | jq

echo "\n # ============================== Legitimize Files ============================== \n"
url='http://localhost:7789/api/v1/storage/legitimize-files'

json=$(curl -X "PUT" "$url" \
     -H 'Content-Type: application/json; charset=utf-8' \
     -d $"{
  \"files\": [
    {
        \"file\": \"$file_url\",
        \"name\": \"$file_name\",
        \"dir\": \"$bucket_name\",
        \"legitimize_url\": \"\"
    }
  ]
}")

echo $json | jq
