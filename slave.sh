#!/bin/bash

# Определяем имя службы
SERVICE_NAME="myapp.service"

# Определяем пути
SERVICE_FILE="/etc/systemd/system/$SERVICE_NAME"
APP_DIR="/home/root/app"
LOG_DIR="/var/log"
ENV_FILE="$APP_DIR/env.conf"
LOG_FILE="$LOG_DIR/myapp.log"
ERROR_LOG_FILE="$LOG_DIR/myapp_error.log"

# Создаем необходимые директории
echo "Создаем директории..."
sudo mkdir -p "$APP_DIR"

# Создаем файлы, если они не существуют
echo "Создаем файлы..."
sudo touch "$ENV_FILE"
sudo touch "$LOG_FILE"
sudo touch "$ERROR_LOG_FILE"

# Настраиваем права доступа
echo "Настраиваем права доступа..."
sudo chown -R root:root "$APP_DIR"
sudo chmod -R 755 "$APP_DIR"
sudo chmod 644 "$ENV_FILE"
sudo chmod 644 "$LOG_FILE"
sudo chmod 644 "$ERROR_LOG_FILE"

# Проверяем, существует ли уже файл службы
if [ -f "$SERVICE_FILE" ]; then
  echo "Служба $SERVICE_NAME уже существует. Хотите перезаписать её? (y/n)"
  read -r choice
  if [[ "$choice" != "y" ]]; then
    echo "Операция отменена."
    exit 1
  fi
fi

# Создаем содержимое файла службы
echo "Создаем файл службы $SERVICE_NAME..."
cat <<EOF | sudo tee "$SERVICE_FILE" > /dev/null
[Unit]
Description=Wireguard Service
After=network.target

[Service]
User=root
WorkingDirectory=$APP_DIR
EnvironmentFile=$ENV_FILE
ExecStart=$APP_DIR/myapp
Restart=always
StandardOutput=append:$LOG_FILE
StandardError=append:$ERROR_LOG_FILE

[Install]
WantedBy=multi-user.target
EOF

# Проверяем, успешно ли создан файл
if [ $? -eq 0 ]; then
  echo "Файл службы $SERVICE_NAME успешно создан."
else
  echo "Ошибка при создании файла службы."
  exit 1
fi

# Включаем службу для автозапуска
echo "Включаем службу для автозапуска..."
sudo systemctl enable "$SERVICE_NAME"

# Запускаем службу
echo "Запускаем службу..."
sudo systemctl start "$SERVICE_NAME"

# Перезагружаем демон systemd
echo "Перезагружаем демон systemd..."
sudo systemctl daemon-reload

# Проверяем статус службы
echo "Проверяем статус службы..."
sudo systemctl status "$SERVICE_NAME"

if [ "$(id -u)" -ne 0 ]; then
    echo "Запустите скрипт с правами root: sudo $0" 
    exit 1
fi

# Обновление и установка WireGuard
apt update
apt install -y wireguard

# Создание директории для конфигов
mkdir -p /etc/wireguard
cd /etc/wireguard

# Генерация ключей
if [ ! -f "server_privatekey" ]; then
    wg genkey | tee server_privatekey | wg pubkey | tee server_publickey
else
    echo "Ключи уже существуют, используем существующие"
fi

# Определение сетевого интерфейса
INTERFACE=$(ip route show default | awk '/default/ {print $5}')
SERVER_IP="10.0.0.1/24"
PORT=51820

# Генерация конфига
PRIVATE_KEY=$(cat server_privatekey)

cat > wg0.conf <<EOF
[Interface]
PrivateKey = $PRIVATE_KEY
Address = $SERVER_IP
ListenPort = $PORT
PostUp = iptables -A FORWARD -i wg0 -j ACCEPT; iptables -t nat -A POSTROUTING -o $INTERFACE -j MASQUERADE
PostDown = iptables -D FORWARD -i wg0 -j ACCEPT; iptables -t nat -D POSTROUTING -o $INTERFACE -j MASQUERADE
EOF

# Включение IP forwarding
echo "net.ipv4.ip_forward=1" >> /etc/sysctl.conf
sysctl -p

# Запуск WireGuard
systemctl enable wg-quick@wg0.service
systemctl start wg-quick@wg0.service
systemctl status wg-quick@wg0.service

# Вывод информации
echo "Настройка завершена!"
echo "Публичный ключ сервера: $(cat server_publickey)"
echo "Интерфейс: $INTERFACE"
echo "Конфиг: /etc/wireguard/wg0.conf"