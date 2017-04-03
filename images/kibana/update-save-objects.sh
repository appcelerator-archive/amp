#!/bin/bash

input=${1:-saved-objects/kibana-amp-dashboard.json}
outputdir=${2:-saved-objects/}
f=$(mktemp)
cp $input $f
for i in $(jq  -r '.[] | ._type + "%" + ._source.title + "%" +._id' $input); do
  otype=$(echo $i | cut -d% -f1)
  title=$(echo $i | cut -d% -f2)
  id=$(echo $i | cut -d% -f3)
  sed -i.bak -e "s/$id/${otype}_${title}/g" $f
done

len=$(jq -r '. | length' $f)
for i in $(seq $len); do
  id=$(jq -r ".[$((i-1))]._id" $f)
  if [ -z "$id" ] || [ "$id" = "null" ]; then
    continue
  fi
  echo "saving $id"
  jq ".[$((i-1))]._source" $f > $outputdir/${id}.json
done
rm $f
