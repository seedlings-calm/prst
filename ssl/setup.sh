#!/bin/bash

#生成自签名证书，开发模式使用

# 检查openssl是否安装
if ! command -v openssl &> /dev/null
then
    echo "openssl could not be found. Please install it first."
    exit 1
fi

# 设置证书信息
COUNTRY="CN"
STATE="Beijing"
LOCALITY="Beijing"
ORGANIZATION="My Company, Inc"
ORG_UNIT="IT Department"
COMMON_NAME="localhost"
EMAIL="123456789@qq.com"
KEY="ssl/key.pem"
CERT="ssl/cert.pem"
# 生成自签名证书和私钥
openssl req -x509 -newkey rsa:4096 -keyout $KEY -out $CERT -days 365 -nodes \
    -subj "/C=$COUNTRY/ST=$STATE/L=$LOCALITY/O=$ORGANIZATION/OU=$ORG_UNIT/CN=$COMMON_NAME/emailAddress=$EMAIL"

# 设置文件权限
chmod 644 $CERT
chmod 600 $KEY

echo "Self-signed certificate and key have been generated:"
echo " - $CERT"
echo " - $KEY"
