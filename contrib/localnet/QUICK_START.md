# Quick Start Guide - Localnet with Dashboards

## 🚀 Get Started in 5 Minutes

### **Prerequisites:**
- Docker Desktop installed and running
- Git repository cloned

### **Step 1: Start Services**
```bash
cd /Users/peterlee/go/src/github.com/zeta-chain/node/contrib/localnet
docker compose --profile monitoring up -d
```

### **Step 2: Wait for Services**
```bash
# Wait for services to be healthy (about 2-3 minutes)
docker compose logs -f --tail=10
```

### **Step 3: Access Grafana**
- Open http://localhost:3000
- Login: admin/admin (change password on first login)
- Navigate to "Dashboards" section

### **Step 4: Explore Dashboards**
Look for dashboards with "Localnet -" prefix:
- **Localnet - CCTXs analytics**
- **Localnet - Consensus metrics**
- **Localnet - Cosmos Dashboard**
- **Localnet - E2E testing**
- **Localnet - Gas Stats & Metrics**
- **Localnet - Log Browser**
- **Localnet - Rate limiter**
- **Localnet - Transaction Data**
- **Localnet - TSS metrics**
- **Localnet - ZetaClient orchestrator**
- **Localnet - ZRC20 data**

## 🔍 Quick Verification

### **Check Service Status:**
```bash
# All services should be running
docker ps | grep -E "(grafana|prometheus|zetachain-exporter|redis)"
```

### **Test API Endpoints:**
```bash
# Grafana
curl -s http://localhost:3000 | grep -q "Grafana" && echo "✅ Grafana OK"

# Prometheus  
curl -s http://localhost:9091 | grep -q "Prometheus" && echo "✅ Prometheus OK"

# ZetaChain Exporter
curl -s http://localhost:9015/metrics | grep -q "zetachain_exporter" && echo "✅ Exporter OK"

curl -s http://localhost:9015/missed-inbounds

# Redis
docker exec zetachain-exporter-redis redis-cli ping && echo "✅ Redis OK"
```

### **Check Dashboard Data:**
```bash
# Check if metrics are being collected
curl -s http://localhost:9091/api/v1/targets | jq '.data.activeTargets[] | select(.health == "up")'
```

## 🎯 What You'll See

### **Grafana Dashboards:**
- Real-time metrics from ZetaChain localnet
- Chain-specific data (Ethereum: 1337, ZetaChain: 101, etc.)
- Network variable set to "localnet"
- Time range selectors and filters

### **Prometheus Metrics:**
- ZetaChain consensus metrics
- ZetaClient metrics  
- ZetaChain Exporter metrics
- System metrics

### **Redis Data:**
- Missed inbounds storage