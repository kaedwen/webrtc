#!/bin/sh

set +x

lego --path "/etc/ldap/cert" --email "${ACME_EMAIL}" --server "https://acme.zerossl.com/v2/DV90" --accept-tos --eab \
--kid "${ACME_KID}" \
--hmac "${ACME_HMAC}" \
--key-type "rsa4096" \
--dns "httpreq" \
--domains "${ACME_DOMAINS}" \
--pem \
renew \
--renew-hook="systemctl restart ring"