#!/usr/bin/env bash

if [ -z "$1" ]
then
    echo "[ERRO] No parameter was given!"
    echo "[INFO] You must provide dn interface name as input parameter"
    echo "Usage:"
    echo "$0 <dn_interface>"
    echo "Example:"
    echo "$0 enp0s4"
else
    echo "[INFO] Using $1 as interface name"

    echo -n "[INFO] Applying iptables rules... "
    sudo iptables -t nat -A POSTROUTING -o $1 -j MASQUERADE
    sudo iptables -I FORWARD 1 -j ACCEPT
    echo "[OK]"
    echo -n "[INFO] Setting kernel net.ipv4.ip_forward flag... "
    sudo sysctl -w net.ipv4.ip_forward=1 >/dev/null
    echo "[OK]"
    echo -n "[INFO] Stopping ufw firewall... "
    sudo systemctl stop ufw
    echo "[OK]"

    echo "[INFO] Configuration applied successfully"
fi
