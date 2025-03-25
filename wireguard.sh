#!/bin/bash
apt update
apt install -y wireguard

mkdir -p /etc/wireguard
cd /etc/wireguard
wg genkey | tee /etc/wireguard/server_privatekey | wg pubkey | tee /etc/wireguard/server_publickey
SERVER_IP="10.0.0.1/24"
PORT=51820
PRIVATE_KEY=$(cat server_privatekey)
PUBLIC_KEY=$(cat server_publickey)
cat > wg0.conf <<EOF
[Interface]
PrivateKey = $PRIVATE_KEY
Address = $SERVER_IP
ListenPort = $PORT
PostUp = iptables -A FORWARD -i %i -j ACCEPT; iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE
PostDown = iptables -D FORWARD -i %i -j ACCEPT; iptables -t nat -D POSTROUTING -o eth0 -j MASQUERADE
EOF

# Включение IP forwarding
echo "net.ipv4.ip_forward=1" >> /etc/sysctl.conf
# Запуск WireGuard
systemctl enable wg-quick@wg0.service
systemctl start wg-quick@wg0.service