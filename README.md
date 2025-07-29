# 💸 Finanzas Personales - GoCardless + RabbitMQ

Esta aplicación en Go conecta con la API de GoCardless para recuperar transacciones bancarias y enviarlas a una cola rabbt.

## 🚀 Cómo arrancar con Docker (sin Compose)

### 1. Construye la imagen

Desde la raíz del proyecto, ejecuta:

```bash
docker build -t gocardless-ingest .
```

### 2. Ejecuta el contenedor

```bash
docker run --rm \
  -e GC_CLIENT_ID=your_gocardless_client_id \
  -e GC_SECRET_KEY=your_gocardless_secret_key \
  -e ACCOUNT_ID=your_account_id \
  -e MYSQL_DSN="user:password@tcp(host.docker.internal:3306)/dbname" \
  -e RABBITMQ_CONN_STR="amqp://user:password@host.docker.internal:5672/" \
  gocardless-ingest
```