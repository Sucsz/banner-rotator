#!/bin/bash

set -e

API_URL="http://localhost:8080"

echo "Add banner to slot"
curl -s -X POST "$API_URL/slots/1/banners" \
  -H "Content-Type: application/json" \
  -d '{"banner_id": 1}'
echo -e "Done\n"

echo "Show banner"
curl -s -X POST "$API_URL/slots/1/show" \
  -H "Content-Type: application/json" \
  -d '{"group_id": 1}'
echo -e "Done\n"

echo "Click banner"
curl -s -X POST "$API_URL/slots/1/click" \
  -H "Content-Type: application/json" \
  -d '{"banner_id": 1, "group_id": 1}'
echo -e "Done\n"

echo "Remove banner from slot"
curl -s -X DELETE "$API_URL/slots/1/banners/1"
echo -e "Done\n"
