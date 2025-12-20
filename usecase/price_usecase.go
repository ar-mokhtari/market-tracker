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

var customClient = &http.Client{
	Timeout: 30 * time.Second,
	Transport: &http.Transport{
		TLSNextProto: make(map[string]func(authority string, c *tls.Conn) http.RoundTripper),
	},
}

type Repo interface {
	Upsert(p entity.Price) error
	List(pType string) ([]entity.Price, error)
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

func (uc *PriceUseCase) FetchFromExternal() error {
	url := fmt.Sprintf("%s?key=%s", uc.baseURL, uc.apiKey)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36")

	resp, err := customClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var data dto.BrsResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}

	// تابع داخلی برای ذخیره دسته‌بندی‌ها
	process := func(items []dto.BrsItem, pType string) {
		for _, item := range items {
			// تبدیل قیمت (چه عدد چه رشته) به رشته برای ذخیره در دیتابیس
			priceStr := fmt.Sprintf("%v", item.Price)

			uc.repo.Upsert(entity.Price{
				Symbol: item.Symbol,
				NameFa: item.Name, // اصلاح شد: استفاده از Name به جای NameFa
				Price:  priceStr,
				Unit:   item.Unit,
				Type:   pType,
				Date:   item.Date,
				Time:   item.Time,
			})
		}
	}

	process(data.Gold, "gold")
	process(data.Currency, "currency")
	process(data.Cryptocurrency, "cryptocurrency")

	fmt.Printf("✅ داده‌ها در تاریخ %s بروزرسانی شدند.\n", time.Now().Format("15:04:05"))
	return nil
}

func (uc *PriceUseCase) GetPrices(pType string) ([]entity.Price, error) {
	return uc.repo.List(pType)
}
