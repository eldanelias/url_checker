package urls_checker

import (
	"github.com/patrickmn/go-cache"
	"time"
	"net/http"
)

type UrlsChecker struct {
	urlsCache *cache.Cache
	stringLocks StringLocks
}

func NewUrlsChecker() UrlsChecker {
	return UrlsChecker {
		urlsCache: cache.New(10*time.Minute, 10*time.Minute),
		stringLocks: NewStringLocks(),
	}
}

func (urlsChecker UrlsChecker) Exists(url string) bool {
	exists, isInCache := urlsChecker.urlsCache.Get(url)
	if isInCache { return exists.(bool) }

	lock := urlsChecker.stringLocks.GetOrAdd(url)
	lock.Lock()
	defer lock.Unlock()

	exists, isInCache = urlsChecker.urlsCache.Get(url)
	if isInCache { return exists.(bool) }

	exists = fetchExists(url)
	urlsChecker.urlsCache.Set(url, exists, cache.DefaultExpiration)

	return exists.(bool)
}

func fetchExists(url string) (exists bool) {
	defer catchBadUrl(&exists)

	res, err := http.Get(url)
	if err != nil { panic(err) }

	return res.StatusCode != http.StatusNotFound
}

func catchBadUrl(exists *bool) {
	if r := recover(); r != nil {
		*exists = false
	}
}
