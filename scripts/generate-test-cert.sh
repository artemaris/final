#!/bin/bash

# Скрипт для генерации самоподписанного TLS сертификата для тестирования
# ВНИМАНИЕ: Этот сертификат подходит только для тестирования!
# Для продакшена используйте сертификаты от доверенного CA (например, Let's Encrypt)

CERT_DIR="certs"
CERT_FILE="$CERT_DIR/server.crt"
KEY_FILE="$CERT_DIR/server.key"

# Создаем директорию для сертификатов
mkdir -p "$CERT_DIR"

# Генерируем самоподписанный сертификат
openssl req -x509 -newkey rsa:4096 -keyout "$KEY_FILE" -out "$CERT_FILE" \
  -days 365 -nodes -subj "/C=RU/ST=Moscow/L=Moscow/O=GophKeeper/CN=localhost"

# Устанавливаем права доступа
chmod 600 "$KEY_FILE"
chmod 644 "$CERT_FILE"

echo "✅ Тестовый TLS сертификат создан:"
echo "   Certificate: $CERT_FILE"
echo "   Private Key: $KEY_FILE"
echo ""
echo "Для запуска сервера с TLS:"
echo "   ./bin/gophkeeper-server --tls-cert=$CERT_FILE --tls-key=$KEY_FILE"
echo ""
echo "⚠️  ВНИМАНИЕ: Это самоподписанный сертификат для тестирования!"
echo "   Браузеры будут показывать предупреждение о безопасности."

