#!/bin/sh

# Встановлюємо параметри підключення до бази даних з змінних середовища
host="$DB_HOST"
port="$DB_PORT"

# Перевіряємо чи задані хост та порт
if [ -z "$host" ] || [ -z "$port" ]; then
  echo "DB_HOST або DB_PORT не визначені!"
  exit 1
fi

# Очікуємо, поки база даних стане доступною
until nc -z "$host" "$port"; do
  echo "Waiting for PostgreSQL at $host:$port..."
  sleep 2
done

echo "PostgreSQL is up and running!"
exec "$@"
