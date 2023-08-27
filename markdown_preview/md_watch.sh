#! /bin/bash

# Watches every 5 sec for change of md5sum of file.
# If file changed, then generate MD preview.

FHASH=`md5 $1`
while true; do
NHASH=`md5 $1`
echo "🔍 Checking Hash"
if [ "$NHASH" != "$FHASH" ]; then
echo "❗️ File Changed, generating preview"
./bin/markdown_preview -file $1
FHASH=$NHASH
fi
sleep 5
done