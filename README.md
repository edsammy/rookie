# [WIP]
## Rook playground
Code is messy. Playing around with go and [Rook API](https://docs.tryrook.io/docs/rookconnect/introduction/)
```
cp example.env .env
```
and update with your values from https://www.clients.portal.tryrook.io/home/settings/tokens

### Link User to Developer account
1. Generate a QR code
  (update user_id name that you want on line 9 of scripts/get_code.sh)
  ```
  bash scripts/get_code.sh
  ```

2. Copy the png text from the returned data (`data:image/png;base64, iVBORw0KG...`)

3. Paste into browser to render QR code

4. Download Rook Extraction App: https://docs.tryrook.io/docs/ROOKExtractionApp/introduction/#how-to-get-the-app

5. Scan the QR code

### Fetch user data
Update user_ID in [fetch_data.go](./scripts/fetch_data.go#L25)
