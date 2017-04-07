#/bin/sh!

# su -
# iptables -A FORWARD -j NFQUEUE --queue-num 10
# sudo ./Agent 1338 10
# sudo hping3 -c 10000 -d 120 -S -w 64 -p 6853 --flood 10.5.0.145 --spoof 192.168.42.1
# sudo service iptables restart
