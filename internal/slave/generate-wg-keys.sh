#!/bin/bash

# Generate private key
private_key=$(wg genkey)

# Generate public key from private key
public_key=$(echo $private_key | wg pubkey)

echo "Private Key: $private_key"
echo "Public Key: $public_key"

# Save keys to .env file
echo "WG_PRIVATE_KEY=$private_key" > .env
echo "WG_PEER_PUBLIC_KEY=$public_key" >> .env 