# Access system health status
path "sys/health" {
  capabilities = ["read", "list"]
}

# Manage the transit secrets engine
path "transit/keys/*" {
  capabilities = [ "create", "read", "list" ]
}

# Encrypt engines secrets
path "transit/encrypt/walletera_apikeys_*" {
  capabilities = [ "create", "read", "update" ]
}

# Renew tokens
path "auth/token/renew" {
  capabilities = [ "update" ]
}

# Lookup tokens
path "auth/token/lookup" {
  capabilities = [ "update" ]
}

# Manage otp keys
path "totp/keys/walletera_*" {
  capabilities = ["create", "read", "update", "delete"]
}

# Verify an otp code
path "totp/code/walletera_*" {
  capabilities = ["update"]
}
