#!/bin/bash

set -e

current_dir=$(cd "$(dirname "$0")";pwd)
root_dir="${current_dir}/.."

cd ${root_dir}
echo "downoading maxmind ip files..... "
curl -L  "https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-Country-CSV&license_key=JvbzLLx7qBZT&suffix=zip" -o GeoLite2-Country-CSV.zip
yes|unzip GeoLite2-Country-CSV.zip
rm GeoLite2-Country-CSV.zip
mv -f GeoLite2* geoip
ls ./geoip

#go get -u github.com/yuanmomo/geoip/...
if [[ ! $(command -v ${current_dir}/ip) ]] ; then
  echo "go build ip command..... "
  export GO111MODULE=on
  cd ${root_dir}
  go build -o ${current_dir}/ip main.go
fi

echo "exec ip command..... "
cd ${root_dir}
${current_dir}/ip --country=./geoip/GeoLite2-Country-Locations-en.csv --ipv4=./geoip/GeoLite2-Country-Blocks-IPv4.csv --ipv6=./geoip/GeoLite2-Country-Blocks-IPv6.csv

echo "move geoip.dat to v2ray bin dir....."
chmod -x "${root_dir}"/geoip.dat
mv "${root_dir}"/geoip.dat ${root_dir}/publish/geoip.dat

echo "delete files ....."
rm -rf ${root_dir}/geoip

