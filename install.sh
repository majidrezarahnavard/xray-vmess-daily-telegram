#!/bin/bash


install_dir=/root/xray-configuration
mkdir $install_dir


# wget https://raw.githubusercontent.com/majidrezarahnavard/xray-configuration/main/reality.json


#instal monitoring
apt-get update
apt-get install nload
apt-get install htop
apt-get install iftop
apt-get install vnstat
apt-get install speedtest-cli
apt-get install net-tools
apt-get install git
apt-get install cron
apt-get install curl tar unzip jq -y
apt-get install -y jq



journalctl --vacuum-time=1d


timedatectl set-timezone UTC
timedatectl
echo "UTC" | sudo tee /etc/timezone
cat /etc/timezone


bash -c "$(curl -L https://github.com/XTLS/Xray-install/raw/main/install-release.sh)" @ install --beta -u root



# Generate key pair
echo "Generating key pair..."
key_pair=$(xray x25519)
echo "Key pair generation complete."
echo $key_pair

#store public key in a file
touch $install_dir/key_pair.txt
echo $key_pair > $install_dir/key_pair.txt


# Create xray.service
cat > /etc/systemd/system/xray.service <<EOF
[Unit]
After=network.target nss-lookup.target

[Service]
User=root
WorkingDirectory=$install_dir
CapabilityBoundingSet=CAP_NET_ADMIN CAP_NET_BIND_SERVICE CAP_NET_RAW
AmbientCapabilities=CAP_NET_ADMIN CAP_NET_BIND_SERVICE CAP_NET_RAW
ExecStart=$install_dir/xray run -c $install_dir/reality.json
ExecReload=/bin/kill -HUP \$MAINPID
Restart=on-failure
RestartSec=10
LimitNOFILE=infinity
StandardOutput=null

[Install]
WantedBy=multi-user.target
EOF


# Install apache2 and clone the website
apt-get install apache2

cd /var/www/html/
git clone https://github.com/codingstella/vCard-personal-portfolio.git
cp -ar ./vCard-personal-portfolio/*  /var/www/html/
rm -rf ./vCard-personal-portfolio/


# Install cron job 
croncmd="cd $install_dir && $install_dir/xray-telegram > $install_dir/cronjob.log 2>&1"
cronjob="30 11 * * * $croncmd"
( crontab -l | grep -v -F "$croncmd" ; echo "$cronjob" ) | crontab -