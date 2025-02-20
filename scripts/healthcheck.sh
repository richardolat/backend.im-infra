#!/bin/bash
curl -sS http://localhost:8080/ws -o /dev/null -w '%{http_code}'
