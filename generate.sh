#!/bin/bash

set -e
root_dir=$(cd "$(dirname "$0")";pwd)

#curl -L --socks5 127.0.0.1:1087 "https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-Country-CSV&license_key=JvbzLLx7qBZT&suffix=zip" -o GeoLite2-Country-CSV.zip
#yes|unzip GeoLite2-Country-CSV.zip
#rm GeoLite2-Country-CSV.zip
#mv -f GeoLite2* geoip
#ls ./geoip

#go get -u github.com/yuanmomo/geoip/...
if [[ ! $(command -v ${root_dir}/ip) ]] ; then
  export GO111MODULE=on
  cd ${root_dir}
  go build -o ${root_dir}/ip main.go
fi


${root_dir}/ip --country=./geoip/GeoLite2-Country-Locations-en.csv --ipv4=./geoip/GeoLite2-Country-Blocks-IPv4.csv --ipv6=./geoip/GeoLite2-Country-Blocks-IPv6.csv

chmod -x "${root_dir}"/geoip.dat
mv "${root_dir}"/geoip.dat $(/usr/local/bin/greadlink -f /usr/local/bin/v2ray | xargs dirname)/geoip.dat

rm -rf ${root_dir}/geoip

