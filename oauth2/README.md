# OAuth 2.0 an authorization framework

## Questions

1. Is always the authorization code stored in `code` from the response?
2. What are CSRF attacks?
3. How many kind of grant type exists?

## OAuth Roles

- Resource Owner
- Client
- Client ID/Secret
- Resource Server
- Authorization Server
- Authorization Grant
- Redirect URL
- Access Token
- Scope
- Consent
- Back/Front Channel

## Flow

First the application (client) has to register/sign up with the Server/API

- client gets the `client_id` and `client_secret`
- client has to specify the `redirect URLs` (usually https)
- state (opaque to OAuth2)

After registration the client has to request the authorization server (`server`) for authorization grant. The
authorization grant could have 
