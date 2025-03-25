#/bin/bash

# List all json files in document folder
for file in /data/document/*.json
do
  # Get collection name from file name
  collection=$(basename $file .json)
  # Import json file to mongodb
  mongoimport -h rs0/mongodb:27017 -u admin -p pass --authenticationDatabase admin --ssl --tlsInsecure --db mediation-platform --collection $collection --file $file --jsonArray
done
