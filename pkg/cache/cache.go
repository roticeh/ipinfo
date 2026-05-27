package cache

import (
	"sync"
	"time"
)

// Item: Önbellekte tutulan verinin kendisi ve son kullanma tarihi
type Item[T any] struct {
	Value      T
	Expiration int64
}

// Config: Önbellek ayarları
type Config struct {
	TTL           time.Duration // Veri ne kadar süre kalacak?
	SweepInterval time.Duration // Temizlikçi ne sıklıkla çalışacak?
	MaxEntries    int           // RAM koruması: Maksimum kaç kayıt tutulabilir?
}

// Store: Jenerik Önbellek Yapısı
type Store[T any] struct {
	items      map[string]Item[T]
	mu         sync.RWMutex
	maxEntries int
	defaultTTL time.Duration
}

// New: Yeni bir önbellek havuzu oluşturur
func New[T any](cfg Config) *Store[T] {
	s := &Store[T]{
		items:      make(map[string]Item[T]),
		maxEntries: cfg.MaxEntries,
		defaultTTL: cfg.TTL,
	}

	go s.janitor(cfg.SweepInterval)
	return s
}

func (s *Store[T]) Set(key string, value T) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Kapasite Kontrolü ve Rastgele Yer Açma (Random Eviction)
	if len(s.items) >= s.maxEntries {
		// Anahtar zaten varsa kapasiteyi artırmaz, sadece günceller. O yüzden önce kontrol et.
		if _, exists := s.items[key]; !exists {
			// Kapasite gerçekten dolmuş. Rastgele birini feda et ki yeni gelene yer açılsın.
			for k := range s.items {
				delete(s.items, k)
				break
			}
		}
	}

	// Veriyi Yaz
	s.items[key] = Item[T]{
		Value:      value,
		Expiration: time.Now().Add(s.defaultTTL).UnixNano(),
	}
}

func (s *Store[T]) Get(key string) (T, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item, found := s.items[key]
	if !found {
		var zero T
		return zero, false
	}

	// Süresi dolmuş mu?
	if time.Now().UnixNano() > item.Expiration {
		var zero T
		return zero, false
	}

	return item.Value, true
}

// Delete: İşlemi biten veriyi anında silmek için eklendi (Örn: OTP başarılı olunca)
func (s *Store[T]) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.items, key)
}

func (s *Store[T]) janitor(interval time.Duration) {
	ticker := time.NewTicker(interval)
	// Program çalıştığı sürece arka planda temizlik yapacak
	for range ticker.C {
		s.mu.Lock()
		now := time.Now().UnixNano()
		for key, item := range s.items {
			if now > item.Expiration {
				delete(s.items, key)
			}
		}
		s.mu.Unlock()
	}
}
