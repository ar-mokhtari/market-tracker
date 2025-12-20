# Market Tracker - راهنمای استفاده

پروژه دریافت و مدیریت قیمت‌های طلا، ارز و ارزهای دیجیتال

## 📋 پیش‌نیازها

- Go 1.21 یا بالاتر
- Docker & Docker Compose (فقط برای MySQL)
- API Key از [BrsApi.ir](https://brsapi.ir)

## 🚀 راه‌اندازی سریع

### 1. کلون و نصب
```bash
git clone <repository-url>
cd market-tracker
go mod download
```

### 2. تنظیم محیط
```bash
cp .env.example .env
```

فایل `.env` را ویرایش کنید و `AUTH_KEY` خود را قرار دهید:
```env
AUTH_KEY=YourApiKeyFromBrsApiMinimum20Chars
```

### 3. راه‌اندازی MySQL
```bash
# راه‌اندازی MySQL با Docker
docker compose up -d

# یا با Makefile
make db-up
```

### 4. اجرای برنامه
```bash
go run main.go

# یا با Makefile
make run
```

### 5. یک دستوری! 🎯
```bash
make start
# این دستور MySQL را راه‌اندازی کرده و برنامه را اجرا می‌کند
```

## 🔧 دستورات Makefile

### دیتابیس
```bash
make db-up          # راه‌اندازی MySQL
make db-down        # متوقف کردن
make db-restart     # Restart
make db-logs        # مشاهده لاگ‌ها
make db-shell       # اتصال به MySQL
make db-clean       # پاک کردن همه چیز
```

### برنامه
```bash
make run            # اجرای برنامه
make build          # Build برنامه
make dev            # راه‌اندازی DB + اجرا
make test           # تست‌ها
make start          # شروع سریع (DB + App)
```

### API
```bash
make health         # بررسی سلامت
make fetch          # دریافت داده‌ها
make prices         # نمایش قیمت‌ها
```

## 📁 ساختار پروژه

```
github.com/ar-mokhtari/market-tracker
├── adapter/
│   └── storage/
│       └── mysql/          # لایه دیتابیس
├── config/
│   └── env/                # تنظیمات
├── delivery/
│   └── http/
│       └── v1/             # HTTP Handlers
├── dto/                    # Data Transfer Objects
├── entity/                 # Domain Models
├── usecase/                # Business Logic
├── validation/             # Validations
├── docker-compose.yml      # MySQL با Docker
├── Makefile               # دستورات کمکی
├── .env                   # تنظیمات محیطی
└── main.go                # نقطه شروع
```

## 🌐 API Endpoints

### Base URL
```
http://localhost:8080/api/v1
```

### دریافت قیمت‌ها
```bash
# همه قیمت‌ها
curl http://localhost:8080/api/v1/prices

# فیلتر براساس نوع
curl http://localhost:8080/api/v1/prices?type=gold
curl http://localhost:8080/api/v1/prices?type=currency
curl http://localhost:8080/api/v1/prices?type=cryptocurrency

# یک قیمت خاص
curl http://localhost:8080/api/v1/prices/1
```

### عملیات CRUD
```bash
# ایجاد
curl -X POST http://localhost:8080/api/v1/prices \
  -H "Content-Type: application/json" \
  -d '{
    "symbol": "TEST",
    "name_en": "Test Price",
    "name_fa": "قیمت تست",
    "price": "100",
    "unit": "تومان",
    "type": "currency"
  }'

# بروزرسانی
curl -X PUT http://localhost:8080/api/v1/prices/1 \
  -H "Content-Type: application/json" \
  -d '{
    "price": "150",
    "change_percent": 0.5
  }'

# حذف
curl -X DELETE http://localhost:8080/api/v1/prices/1
```

### دریافت داده‌های جدید
```bash
curl -X POST http://localhost:8080/api/v1/prices/fetch
```

## 🗄️ دیتابیس

### اتصال به MySQL
```bash
# با Makefile
make db-shell

# یا مستقیم
mysql -h127.0.0.1 -P3306 -umarket_user -pmarket_secure_password_2024 market_tracker
```

### مشخصات اتصال
- **Host**: localhost
- **Port**: 3306
- **Database**: market_tracker
- **User**: market_user
- **Password**: (در فایل .env)

### جدول prices
```sql
-- ساختار جدول
CREATE TABLE prices (
    id INT AUTO_INCREMENT PRIMARY KEY,
    date VARCHAR(20),
    time VARCHAR(20),
    time_unix BIGINT,
    symbol VARCHAR(50) NOT NULL,
    name_en VARCHAR(100),
    name_fa VARCHAR(100),
    price VARCHAR(50),
    change_value VARCHAR(50),
    change_percent DECIMAL(10, 2),
    unit VARCHAR(20),
    type VARCHAR(20) NOT NULL,
    market_cap BIGINT,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY unique_symbol_type (symbol, type)
);
```

## 🔄 دریافت خودکار داده‌ها

برنامه به صورت خودکار هر **5 دقیقه** داده‌های جدید را دریافت می‌کند.

برای تغییر فاصله زمانی، فایل `main.go` را ویرایش کنید:
```go
// main.go - خط 45
go startPeriodicFetch(priceUseCase, 5*time.Minute) // تغییر به دلخواه
```

## 🛠️ توسعه

### نصب ابزارها
```bash
# golangci-lint برای linting
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# air برای hot reload
go install github.com/cosmtrek/air@latest
```

### Hot Reload با Air
```bash
# ایجاد فایل .air.toml
air init

# اجرا با hot reload
air
```

### فرمت کردن کد
```bash
make fmt
```

### بررسی کد
```bash
make lint
```

### تست
```bash
# اجرای تست‌ها
make test

# با coverage
make test-coverage
```

## 🐛 عیب‌یابی

### MySQL متصل نمی‌شود
```bash
# بررسی وضعیت
docker compose ps

# بررسی لاگ‌ها
make db-logs

# Restart
make db-restart
```

### پورت اشغال است
در فایل `.env` پورت را تغییر دهید:
```env
SERVICE_PORT=8081
DB_PORT=3307
```

### برنامه داده نمی‌گیرد
1. بررسی API Key در `.env`
2. بررسی اتصال اینترنت
3. مشاهده لاگ‌های برنامه

### پاک کردن و شروع دوباره
```bash
make db-clean
make db-up
go run main.go
```

## 📊 مثال‌های استفاده

### دریافت قیمت دلار
```bash
curl -s http://localhost:8080/api/v1/prices?type=currency | \
  jq '.data[] | select(.symbol=="USD")'
```

### دریافت قیمت بیت‌کوین
```bash
curl -s http://localhost:8080/api/v1/prices?type=cryptocurrency | \
  jq '.data[] | select(.symbol=="BTC")'
```

### نمایش 10 قیمت اخیر
```bash
curl -s http://localhost:8080/api/v1/prices | \
  jq '.data[:10]'
```

## 📦 Build برای Production

### Build باینری
```bash
make build
# خروجی: bin/market-tracker
```

### اجرای باینری
```bash
./bin/market-tracker
```

### Build با flags بیشتر
```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
  -ldflags="-w -s" \
  -o market-tracker \
  main.go
```

## 🔒 امنیت

- ✅ API Key را در `.env` نگه دارید
- ✅ `.env` را به git اضافه نکنید
- ✅ پسورد MySQL را قوی انتخاب کنید
- ✅ در production از HTTPS استفاده کنید
- ✅ Rate limiting را فعال کنید

## 📝 نکات مهم

1. **اولین اجرا**: MySQL باید آماده باشد (30 ثانیه صبر کنید)
2. **API Key**: حتماً از BrsApi.ir دریافت کنید
3. **Go Version**: حداقل Go 1.21
4. **Memory**: برنامه حدود 20-30MB RAM مصرف می‌کند

## 🤝 مشارکت

1. Fork کنید
2. Branch بسازید: `git checkout -b feature/amazing`
3. Commit کنید: `git commit -m 'Add amazing feature'`
4. Push کنید: `git push origin feature/amazing`
5. Pull Request بزنید

## 📄 لایسنس

MIT License

---

**ساخته شده با ❤️ برای جامعه Go ایران**
