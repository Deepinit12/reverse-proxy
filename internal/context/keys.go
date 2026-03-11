package context

// contextKey — приватный тип, чтобы избежать коллизий ключей с другими пакетами
type contextKey string

// Экспортируемые ключи (используй их везде через эти константы)
var (
	RequestIDKey = contextKey("request_id")
	RealIPKey    = contextKey("real_ip")
	// Дальше можно добавить:
	// FingerprintKey   = contextKey("fingerprint")
	// AnomalyScoreKey  = contextKey("anomaly_score")
	// StartTimeKey     = contextKey("start_time")
)
