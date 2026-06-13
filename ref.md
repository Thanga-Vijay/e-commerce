Based on your project structure, here's the **complete order** to run docker-compose files for local testing:

## 🚀 Complete Docker Compose Startup Order

### **Step 1: Start Infrastructure (Databases, Redis, Kafka)**

```bash
# Start databases and Redis
docker-compose up -d postgres redis

# Wait 10-15 seconds for databases to initialize
# Then start Kafka with Zookeeper
docker-compose -f docker-compose.kafka.yml up -d

# Wait 20-30 seconds for Kafka to be ready
# Verify Kafka is running
docker-compose -f docker-compose.kafka.yml ps
```

### **Step 2: Create Kafka Topics**

```bash
# Create all required topics
make kafka-topics-create

# Or manually:
docker-compose -f docker-compose.kafka.yml exec kafka kafka-topics.sh \
  --create --bootstrap-server localhost:9092 \
  --topic user.registered --partitions 3 --replication-factor 1

# Verify topics created
make kafka-topics-list
```

### **Step 3: Start All Microservices**

```bash
# Start all 9 backend services
docker-compose up -d \
  auth-service \
  product-service \
  cart-service \
  wishlist-service \
  order-service \
  payment-service \
  inventory-service \
  notification-service \
  reporting-service

# Wait 20-30 seconds for services to initialize
# Check logs
docker-compose logs -f auth-service product-service
```

### **Step 4: Start Frontend**

```bash
# Start React frontend
docker-compose up -d frontend

# Or if you prefer to run frontend locally (hot reload):
cd frontend
npm install
npm run dev
```

### **Step 5: Start Monitoring (Optional)**

```bash
# Start Prometheus, Grafana, Alertmanager
docker-compose -f docker-compose.monitoring.yml up -d

# Wait 10 seconds, then access:
# Grafana: http://localhost:3001 (admin/admin)
# Prometheus: http://localhost:9090
```

---

## 📝 Complete Command Sequence

**Copy-paste this complete sequence:**

```bash
# 1. Start databases and Redis
docker-compose up -d postgres redis
echo "Waiting for databases to initialize..."
sleep 15

# 2. Start Kafka infrastructure
docker-compose -f docker-compose.kafka.yml up -d
echo "Waiting for Kafka to be ready..."
sleep 30

# 3. Create Kafka topics (if you have the script)
make kafka-topics-create || echo "Skipping Kafka topics (run manually if needed)"

# 4. Start all microservices
docker-compose up -d \
  auth-service \
  product-service \
  cart-service \
  wishlist-service \
  order-service \
  payment-service \
  inventory-service \
  notification-service \
  reporting-service

echo "Waiting for services to start..."
sleep 20

# 5. Start frontend
docker-compose up -d frontend

# 6. Optional: Start monitoring
docker-compose -f docker-compose.monitoring.yml up -d

# 7. Check status
docker-compose ps
docker-compose -f docker-compose.kafka.yml ps
docker-compose -f docker-compose.monitoring.yml ps

echo "✅ All services started!"
```

---

## 🎯 Using Makefile (Easiest Way)

If you have the Makefile, use these commands in order:

```bash
# 1. Start infrastructure (databases, Redis, Kafka)
make infra-up

# 2. Start all services
make services-up

# 3. Start frontend
make frontend-up

# 4. Start monitoring (optional)
make monitoring-up

# 5. Check everything
make status
```

---

## 🌐 Access Your Services

After startup (wait 1-2 minutes total), access:

| Service | URL | Notes |
|---------|-----|-------|
| **Frontend** | http://localhost:3000 | React UI |
| **Auth API** | http://localhost:8081/health | JWT auth |
| **Product API** | http://localhost:8082/health | Catalog |
| **Cart API** | http://localhost:8083/health | Shopping cart |
| **Order API** | http://localhost:8085/health | Orders |
| **Payment API** | http://localhost:8086/health | Stripe |
| **Kafka UI** | http://localhost:8090 | Topics, consumers |
| **Grafana** | http://localhost:3001 | Dashboards (admin/admin) |
| **Prometheus** | http://localhost:9090 | Metrics |

---

## 🔍 Verify Everything is Running

```bash
# Check all containers
docker ps

# Check logs for any errors
docker-compose logs --tail=50

# Check specific service
docker-compose logs -f auth-service

# Health check all services
curl http://localhost:8081/health
curl http://localhost:8082/health
curl http://localhost:8083/health
```

---

## 🛑 Shutdown Order (When Done)

```bash
# Stop in reverse order
docker-compose down                              # Services & frontend
docker-compose -f docker-compose.monitoring.yml down
docker-compose -f docker-compose.kafka.yml down

# Or stop everything at once
make down

# Remove volumes (clean slate)
docker-compose down -v
```

---

## ⚠️ Common Issues & Solutions

**Issue: "Port already in use"**
```bash
# Find what's using the port
netstat -ano | findstr :8081  # Windows
lsof -i :8081                 # Mac/Linux

# Stop the process or change port in docker-compose.yml
```

**Issue: "Database connection refused"**
```bash
# Wait longer for databases to initialize
docker-compose logs postgres

# Or restart databases
docker-compose restart postgres
```

**Issue: "Kafka broker not available"**
```bash
# Kafka takes 20-30 seconds to start
docker-compose -f docker-compose.kafka.yml logs kafka

# Restart Kafka
docker-compose -f docker-compose.kafka.yml restart kafka
```

**Issue: Services can't connect to each other**
```bash
# Check network
docker network ls
docker network inspect e-commerce_default

# Restart services
docker-compose restart
```

---

## 💡 Pro Tips

1. **First Time Setup**: Run `docker-compose pull` to download all images first
2. **Watch Logs**: Use `docker-compose logs -f` to see real-time logs
3. **Clean Start**: Use `docker-compose down -v` to remove volumes and start fresh
4. **Resource Usage**: Monitor with `docker stats` to see CPU/memory usage
5. **Rebuild After Code Changes**: `docker-compose up -d --build <service-name>`

---

**Quick Start (All-in-One):**
```bash
make stack-up        # Starts everything
make stack-down      # Stops everything
make stack-restart   # Restarts everything
```

That's the complete order! Start with infrastructure → microservices → frontend → monitoring (optional). 🚀

```bash
docker-compose up -d
docker-compose -f docker-compose.kafka.yml up -d
docker-compose -f docker-compose.monitring.yml up -d
docker-compose ps
docker-compose -f docker-compose.kafka.yml ps
```

docker compose \
  -f docker-compose.yml \
  -f docker-compose.kafka.yml \
  -f docker-compose.monitoring.yml config | grep DATA_SOURCE_NAME \
  up -d