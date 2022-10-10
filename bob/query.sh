#!/bin/bash

#curl -v -X POST http://localhost:8080 -H 'content-type: application/json' -d '{"name": "Bob"}'

curl -v -X POST $(gp url 8080) -H 'content-type: application/json' -d '{"name": "Bob"}'
