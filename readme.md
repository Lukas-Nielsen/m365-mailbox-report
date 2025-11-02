# M365 Mailbox Reporter

Generates a Markdown report for BookStack.

## Azure Permissions

The following permissions are required:

- `Mail.Read`
- `MailboxSettings.Read`
- `User.Read.All`
- `Reports.Read.All`
- `Directory.Read.All`

## BookStack Configuration

### Required Information

- **Page ID** – locate via the BookStack API.
- **API Token** – generate in your account settings.

### Steps

1. **Log in** to your BookStack instance.
2. **Create a new page** (or use an existing one) that will hold the report.
3. **Retrieve the Page ID**:
   - Navigate to `https://<your-bookstack-domain>/api/pages`.
   - Find the desired page and note its `id`.
4. **Generate an API Token**:
   - Go to **My Account → Auth** (`/my-account/auth`).
   - Create a token and copy both the **Token ID** and **Token Secret**.

These values are required for the reporter to post the generated Markdown report to BookStack.

## Docker

The repository includes a default `docker‑compose.yaml`.
