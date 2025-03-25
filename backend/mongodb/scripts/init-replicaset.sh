#/bin/bash

# Initiate replicaset
mongosh "mongodb://mongodb:27017/" -u admin -p pass --eval 'rs.initiate()'
