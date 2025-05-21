# Агрегатор турів

**Агрегатор турів** — це Go-сервіс для централізованого пошуку туристичних пропозицій від різних туроператорів.

---

## Основні можливості

- Пошук турів через REST API
- Асинхронна обробка задач (RabbitMQ)
- Кешування результатів (Memcached, TTL 15 хв)
- Зберігання цін у ClickHouse
- Гнучке підключення нових туроператорів через адаптери

---

## Вимоги

- Go 1.22.0
- MySQL
- ClickHouse
- Memcached
- RabbitMQ

---

## Швидкий старт

### 1. Встановлення Go

```bash
curl -O https://dl.google.com/go/go1.22.0.linux-amd64.tar.gz
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
```

### 2. Клонування та збірка

```bash
git clone https://github.com/tomazov/gonapi.git
cd gonapi
go build -o aggregator main.go
```

### 3. Налаштування .env

```dotenv
DB_HOST=localhost
DB_PORT=3306
DB_USER=touruser
DB_PASSWORD=secret
DB_NAME=tourdb
CH_HOST=localhost
CH_PORT=9000
CH_USER=default
CH_PASSWORD=
CH_DATABASE=tourdata
MEMCACHED_HOST=localhost
MEMCACHED_PORT=11211
RABBITMQ_URL=amqp://guest:guest@localhost:5672/
```

### 4. Запуск

```bash
./aggregator
```

---

## API-приклад

```http
GET /getResults?from=1831&to=115&checkIn=2025-06-01&adults=2
```

---

## Формат відповіді

```json
{
  "searchId": "abc123-def456",
  "lastResult": false,
  "results": [
    {
      "operator": "TourOperatorA",
      "tourName": "Holiday Hotel 5*",
      "departureDate": "2025-06-01",
      "nights": 7,
      "price": 1000,
      "currency": "USD",
      "persons": 2
    }
  ]
}
```

---

## Логування

- Для кожного туроператора лог-файл: `/var/log/vendor_<rec_id>.log`

---

## Статус

- ✅ Архітектура погоджена
- ✅ Реалізовано кешування
- ⏳ Йде інтеграція ClickHouse

---
