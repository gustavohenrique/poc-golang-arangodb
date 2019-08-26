#!/bin/sh

database="elearning"
password="root"

for i in `/bin/ls -1 *.json`; do
  collection_name=`echo $i | awk -F '_' '{print $1}'`
  collection_type=`echo $i | awk -F '_' '{print $2}'`
  echo "Importing $i to $collection_name... $collection_type"
  edge=""
  if [ "$collection_type" == "edge" ]; then
    edge="--create-collection-type edge"
  fi
  arangoimport \
    --file "$i" \
    --collection "$collection_name" \
    --create-collection $edge \
    --overwrite \
    --batch-size 805306368 \
    --create-database \
    --server.database="$database" \
    --server.password="$password"
done
