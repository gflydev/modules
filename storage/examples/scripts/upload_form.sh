#!/bin/sh
#
# Upload file to local
#

file_path='../public/assets'
file_name='hero.png'
file_type='image/png'

echo "\n # ============================== Upload to Server ============================== \n"
upload_url='http://localhost:7789/api/v1/storage/uploads'

json=$(curl "$upload_url" \
  --form "file=@$file_path/$file_name;type=$file_type")

echo $json | jq

