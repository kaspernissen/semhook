#!/bin/bash

curl -X POST -F "rule=@test/test_rule.yaml" http://localhost:8080/ask
