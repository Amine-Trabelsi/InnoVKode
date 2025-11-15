package config

import (
	"time"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	MaxBotToken       string        `env:"MAX_BOT_TOKEN,required"`
	BackendBaseURL    string        `env:"BACKEND_BASE_URL" envDefault:"http://localhost:8001"`
	HTTPTimeout       time.Duration `env:"HTTP_TIMEOUT" envDefault:"10s"`
	OTPExpiry         time.Duration `env:"OTP_EXPIRY" envDefault:"5m"`
	Environment       string        `env:"ENVIRONMENT" envDefault:"local"`
	LogLevel          string        `env:"LOG_LEVEL" envDefault:"info"`
	AdmissionsEmail   string        `env:"ADMISSIONS_EMAIL" envDefault:"admissions@univ.ru"`
	AdmissionsPhone   string        `env:"ADMISSIONS_PHONE" envDefault:"+7 (812) 555-0101"`
	AdmissionsOffice  string        `env:"ADMISSIONS_OFFICE" envDefault:"Main Campus, Office 204"`
	DormPaymentURL    string        `env:"DORM_PAYMENT_URL" envDefault:"https://pay.univ.ru/dorm"`
	TuitionPaymentURL string        `env:"TUITION_PAYMENT_URL" envDefault:"https://pay.univ.ru/tuition"`
	ELibraryURL       string        `env:"E_LIBRARY_URL" envDefault:"https://library.univ.ru/ebooks"`
	SupportEmail      string        `env:"SUPPORT_EMAIL" envDefault:"support@univ.ru"`
}

func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
