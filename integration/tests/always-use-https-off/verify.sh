#!/bin/bash

cat ./out/result.json | jq '.result.value' > ./out/actual.txt