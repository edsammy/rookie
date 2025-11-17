source .env

AUTH_HEADER="$(printf '%s:%s' "$ROOK_CLIENT_UUID" "$ROOK_SECRET_KEY" | base64)"

curl --location 'https://api.rook-connect.review/api/v1/extraction_app/binding/' \
--header 'Content-Type: application/json' \
--header "Authorization: Basic ${AUTH_HEADER}" \
--data '{
  "user_id": "rafael1",
  "metadata": {
    "client_name": "Compass Health",
    "tyc_url": "https://example.com/terms",
    "support_url": "https://example.com/support",
    "complete_log_out": false
  },
  "salt": "12345"
}'
