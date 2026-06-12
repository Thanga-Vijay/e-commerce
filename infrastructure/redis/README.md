# Redis Configuration

Redis setup for caching and session management.

## Deployment

- Redis Cluster Mode (6 nodes)
- Redis Sentinel for high availability
- Persistent storage with AOF

## Use Cases

### Product Service
- Product details cache (1 hour TTL)
- Category list cache (24 hours TTL)
- Search results cache (15 minutes TTL)

### Cart Service
- User cart cache (24 hours TTL)
- Cart session data

### Auth Service
- JWT blacklist (TTL = token expiry)
- User sessions
- Rate limiting counters

### Reporting Service
- Dashboard metrics (5 minutes TTL)
- Aggregated reports cache

## Cache Patterns

### Cache-Aside (Lazy Loading)
```go
// Check cache
value, err := redisClient.Get(ctx, key).Result()
if err == redis.Nil {
    // Cache miss - fetch from DB
    value = fetchFromDB()
    // Set in cache
    redisClient.Set(ctx, key, value, ttl)
}
return value
```

### Write-Through
```go
// Update DB
db.Update(data)
// Update cache
redisClient.Set(ctx, key, data, ttl)
```

### Write-Behind (Async)
```go
// Update cache immediately
redisClient.Set(ctx, key, data, ttl)
// Queue DB update
kafkaProducer.Send("db.update", data)
```

## Access

Redis CLI:
```bash
kubectl exec -it redis-0 -n redis -- redis-cli
```

Redis Commander (GUI):
```bash
kubectl port-forward svc/redis-commander 8081:8081 -n redis
```

## Configuration

```yaml
# redis.conf
maxmemory 2gb
maxmemory-policy allkeys-lru
appendonly yes
appendfsync everysec
```

## Monitoring

Metrics:
- Hit rate
- Memory usage
- Connected clients
- Commands per second
- Evictions

## Common Commands

Set key:
```bash
SET product:123 "{\"name\":\"Product\"}" EX 3600
```

Get key:
```bash
GET product:123
```

Check TTL:
```bash
TTL product:123
```

Flush cache:
```bash
FLUSHDB
```
