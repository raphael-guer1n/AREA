# OAuth2 Setup Guide

This guide explains how to set up and use the OAuth2 authentication system in the Auth Service.

## Overview

The OAuth2 system allows users to authenticate using external providers (Google, GitHub, Microsoft, Discord, etc.) without using any third-party OAuth2 libraries. The implementation is generic and configurable through a JSON configuration file.

## Architecture

The OAuth2 system consists of:

1. **Configuration** (`oauth2/config.go`) - Loads provider configurations from JSON
2. **Provider** (`oauth2/provider.go`) - Handles OAuth2 flow for individual providers
3. **Manager** (`oauth2/manager.go`) - Manages multiple providers and state validation
4. **Routes** - Three endpoints for OAuth2 authentication

## Configuration

### 1. Create OAuth2 Config File

Copy the example configuration:
```bash
cp oauth2.config.example.json oauth2.config.json
```

### 2. Configure Providers

Edit `oauth2.config.json` with your provider credentials:

```json
{
  "providers": {
    "google": {
      "name": "google",
      "client_id": "YOUR_GOOGLE_CLIENT_ID",
      "client_secret": "YOUR_GOOGLE_CLIENT_SECRET",
      "auth_url": "https://accounts.google.com/o/oauth2/v2/auth",
      "token_url": "https://oauth2.googleapis.com/token",
      "redirect_uri": "http://localhost:8080/auth/oauth2/callback",
      "scopes": ["openid", "email", "profile"],
      "user_info_url": "https://www.googleapis.com/oauth2/v2/userinfo"
    }
  }
}
```

### 3. Set Environment Variable (Optional)

Set the config file path in `.env`:
```bash
OAUTH2_CONFIG_PATH=oauth2.config.json
```

## Obtaining Provider Credentials

### Google OAuth2

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select existing
3. Enable "Google+ API"
4. Go to "Credentials" → "Create Credentials" → "OAuth 2.0 Client ID"
5. Set authorized redirect URI: `http://localhost:8080/auth/oauth2/callback`
6. Copy Client ID and Client Secret

### GitHub OAuth2

1. Go to GitHub Settings → Developer settings → OAuth Apps
2. Click "New OAuth App"
3. Set Authorization callback URL: `http://localhost:8080/auth/oauth2/callback`
4. Copy Client ID and generate Client Secret

### Microsoft OAuth2

1. Go to [Azure Portal](https://portal.azure.com/)
2. Navigate to "Azure Active Directory" → "App registrations"
3. Click "New registration"
4. Set Redirect URI: `http://localhost:8080/auth/oauth2/callback`
5. Copy Application (client) ID
6. Go to "Certificates & secrets" → Create new client secret

### Discord OAuth2

1. Go to [Discord Developer Portal](https://discord.com/developers/applications)
2. Create "New Application"
3. Go to OAuth2 section
4. Add redirect: `http://localhost:8080/auth/oauth2/callback`
5. Copy Client ID and Client Secret

## API Endpoints

### 1. List Available Providers

**GET** `/auth/oauth2/providers`

Returns all configured OAuth2 providers.

**Response:**
```json
{
  "success": true,
  "data": {
    "providers": ["google", "github", "microsoft", "discord"]
  }
}
```

### 2. Generate Authorization URL

**GET** `/auth/oauth2/authorize?provider=google`

Generates the OAuth2 authorization URL to redirect users to the provider's login page.

**Parameters:**
- `provider` (required): Provider name (e.g., "google", "github")

**Response:**
```json
{
  "success": true,
  "data": {
    "auth_url": "https://accounts.google.com/o/oauth2/v2/auth?client_id=...&state=...",
    "provider": "google"
  }
}
```

**Frontend Usage:**
```javascript
// 1. Get auth URL
const response = await fetch('/auth/oauth2/authorize?provider=google');
const data = await response.json();

// 2. Redirect user to provider
window.location.href = data.data.auth_url;
```

### 3. Handle OAuth2 Callback

**GET** `/auth/oauth2/callback?code=...&state=...`

Handles the OAuth2 callback after user authenticates with the provider.

**Parameters:**
- `code` (required): Authorization code from provider
- `state` (required): CSRF protection state parameter

**Response:**
```json
{
  "success": true,
  "data": {
    "provider": "google",
    "user_info": {
      "id": "123456789",
      "email": "user@example.com",
      "name": "John Doe",
      "username": "johndoe"
    },
    "access_token": "ya29.a0...",
    "token_type": "Bearer",
    "expires_in": 3600
  }
}
```

## OAuth2 Flow

### Complete Authentication Flow

```
1. Frontend: GET /auth/oauth2/authorize?provider=google
   ← Returns auth_url

2. Frontend: Redirect user to auth_url
   → User authenticates with Google

3. Provider: Redirects to /auth/oauth2/callback?code=...&state=...
   → Auth service validates state
   → Exchanges code for access token
   → Fetches user info from provider
   ← Returns user info and tokens

4. Backend: Create or link user account
   → Generate JWT token
   ← Return JWT to frontend

5. Frontend: Store JWT and authenticate user
```

## Implementation Example

### Frontend Flow

```javascript
class OAuth2Client {
  async startAuth(provider) {
    // Get authorization URL
    const response = await fetch(
      `/auth/oauth2/authorize?provider=${provider}`
    );
    const data = await response.json();

    if (data.success) {
      // Redirect to provider
      window.location.href = data.data.auth_url;
    }
  }

  async handleCallback() {
    // Extract code and state from URL
    const params = new URLSearchParams(window.location.search);
    const code = params.get('code');
    const state = params.get('state');

    if (code && state) {
      // Send to backend
      const response = await fetch(
        `/auth/oauth2/callback?code=${code}&state=${state}`
      );
      const data = await response.json();

      if (data.success) {
        // User authenticated successfully
        console.log('User info:', data.data.user_info);
        // TODO: Link account or create session
      }
    }
  }
}
```

## Database Schema

The system includes a new table for storing OAuth2 account links:

```sql
CREATE TABLE oauth2_accounts (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider VARCHAR(50) NOT NULL,
    provider_user_id VARCHAR(255) NOT NULL,
    access_token TEXT NOT NULL,
    refresh_token TEXT,
    token_type VARCHAR(50),
    expires_at TIMESTAMP,
    scope TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(provider, provider_user_id)
);
```

## Security Features

1. **CSRF Protection**: State parameter validates the OAuth2 flow
2. **Single-use States**: State tokens are deleted after use
3. **Secure Token Storage**: Tokens stored in database, not exposed to frontend
4. **Provider Validation**: Only configured providers are allowed

## Adding New Providers

To add a new OAuth2 provider:

1. Obtain OAuth2 credentials from the provider
2. Add configuration to `oauth2.config.json`:

```json
{
  "providers": {
    "newprovider": {
      "name": "newprovider",
      "client_id": "YOUR_CLIENT_ID",
      "client_secret": "YOUR_CLIENT_SECRET",
      "auth_url": "https://provider.com/oauth/authorize",
      "token_url": "https://provider.com/oauth/token",
      "redirect_uri": "http://localhost:8080/auth/oauth2/callback",
      "scopes": ["email", "profile"],
      "user_info_url": "https://provider.com/api/user"
    }
  }
}
```

3. Restart the service

## Troubleshooting

### OAuth2 Endpoints Return 503

- Ensure `oauth2.config.json` exists and is valid
- Check the `OAUTH2_CONFIG_PATH` environment variable
- Verify the file contains at least one provider

### Invalid State Error

- State tokens are single-use and expire
- User must complete the flow without refreshing
- Don't reuse authorization URLs

### Token Exchange Fails

- Verify client_id and client_secret are correct
- Check that redirect_uri matches exactly in provider settings
- Ensure provider URLs (auth_url, token_url) are correct

### User Info Request Fails

- Verify user_info_url is correct for the provider
- Check that requested scopes include user profile access
- Some providers use different field names (configure mapping if needed)

## Next Steps

1. Implement user account linking in the callback handler
2. Store OAuth2 tokens in the database
3. Add token refresh functionality
4. Implement account merging logic
5. Add provider-specific user info mapping
