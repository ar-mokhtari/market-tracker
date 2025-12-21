package usecase

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/ar-mokhtari/market-tracker/dto"
	"github.com/ar-mokhtari/market-tracker/entity"
)

var httpClient = &http.Client{
	Timeout: 30 * time.Second,
	Transport: &http.Transport{
		ForceAttemptHTTP2: false, // اجبار به HTTP/1.1
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			},
		},
	},
}

type PriceUseCase struct {
	repo    Repo
	apiKey  string
	baseURL string
}

func NewPriceUseCase(repo Repo, apiKey string) *PriceUseCase {
	baseUrl := os.Getenv("API_BASE_URL")
	if baseUrl == "" {
		baseUrl = "https://BrsApi.ir/Api/Market/Gold_Currency.php"
	}
	return &PriceUseCase{repo: repo, apiKey: apiKey, baseURL: baseUrl}
}

// func (uc *PriceUseCase) FetchFromExternal() error {
// 	url := fmt.Sprintf("%s?key=%s", uc.baseURL, uc.apiKey)

// 	req, _ := http.NewRequest("GET", url, nil)

// 	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
// 	req.Header.Set("Accept", "text/html,application/json,application/xhtml+xml,application/xml q=0.9,image/avif,image/webp,*")
// 	req.Header.Set("Accept-Language", "en-US,en;q=0.9,fa;q=0.8")
// 	req.Header.Set("Connection", "keep-alive")

// 	resp, err := httpClient.Do(req)
// 	if err != nil {
// 		return err
// 	}
// 	defer resp.Body.Close()

// 	var data dto.BrsResponse
// 	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
// 		return err
// 	}

// 	// تابع داخلی برای ذخیره دسته‌بندی‌ها
// 	process := func(items []dto.BrsItem, pType string) {
// 		for _, item := range items {
// 			// تبدیل قیمت (چه عدد چه رشته) به رشته برای ذخیره در دیتابیس
// 			priceStr := fmt.Sprintf("%v", item.Price)

// 			uc.repo.Upsert(entity.Price{
// 				Symbol: item.Symbol,
// 				NameFa: item.Name, // اصلاح شد: استفاده از Name به جای NameFa
// 				Price:  priceStr,
// 				Unit:   item.Unit,
// 				Type:   pType,
// 				Date:   item.Date,
// 				Time:   item.Time,
// 			})
// 		}
// 	}

// 	process(data.Gold, "gold")
// 	process(data.Currency, "currency")
// 	process(data.Cryptocurrency, "cryptocurrency")

// 	fmt.Printf("✅ داده‌ها در تاریخ %s بروزرسانی شدند.\n", time.Now().Format("15:04:05"))
// 	return nil
// }

func (uc *PriceUseCase) FetchFromExternal() error {
	// خواندن فایل به جای درخواست شبکه
	fileData, err := os.ReadFile("data.json")
	if err != nil {
		return fmt.Errorf("خطا در خواندن فایل محلی: %v", err)
	}

	var response dto.BrsResponse
	if err := json.Unmarshal(fileData, &response); err != nil {
		return fmt.Errorf("خطا در پارس کردن JSON: %v", err)
	}

	processAndSave := func(items []dto.BrsItem, category string) {
		for _, item := range items {
			priceStr := ""
			switch v := item.Price.(type) {
			case float64:
				priceStr = fmt.Sprintf("%.0f", v) // حذف نماد علمی
			default:
				priceStr = fmt.Sprintf("%v", v)
			}
			uc.repo.Upsert(entity.Price{
				Symbol: item.Symbol,
				NameFa: item.Name,
				Price:  priceStr,
				Unit:   item.Unit,
				Type:   category,
				Date:   item.Date,
				Time:   item.Time,
			})
		}
	}

	processAndSave(response.Gold, "gold")
	processAndSave(response.Currency, "currency")
	processAndSave(response.Cryptocurrency, "cryptocurrency")

	fmt.Println("✅ داده‌ها از فایل محلی با موفقیت به دیتابیس تزریق شدند.")
	return nil
}

func (uc *PriceUseCase) GetPrices(pType string) ([]entity.Price, error) {
	return uc.repo.List(pType)
}

func (uc *PriceUseCase) GetSymbolTimeline(symbol string) ([]entity.Price, error) {
	const defaultLimit = 24 // Last 24 records for hourly timeline
	return uc.repo.GetHistory(symbol, defaultLimit)
}
