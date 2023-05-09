package bookquery

import (
	"bytes"
	"crypto/sha1"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"sync"
	"time"
)

type cacheItem struct {
	Data      []byte
	ExpiresAt time.Time
}

type Cache struct {
	data          map[string]*cacheItem
	cacheToMemory bool
	cacheFile     string
	mu            sync.Mutex
}

func (c *Cache) Set(key string, data []byte, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = &cacheItem{
		Data:      data,
		ExpiresAt: time.Now().Add(ttl),
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, exists := c.data[key]
	if !exists {
		return nil, false
	}

	if time.Now().After(item.ExpiresAt) {
		delete(c.data, key)
		return nil, false
	}

	return item.Data, true
}

func (c *Cache) loadFromFile() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	file, err := os.Open(c.cacheFile)
	if err != nil {
		return err
	}
	defer file.Close()

	dec := gob.NewDecoder(file)
	err = dec.Decode(&c.data)
	if err != nil {
		return err
	}

	// Clean up expired cache entries
	now := time.Now()
	for k, v := range c.data {
		if now.After(v.ExpiresAt) {
			delete(c.data, k)
		}
	}

	return nil
}

func (c *Cache) saveToFile() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	file, err := os.Create(c.cacheFile)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := gob.NewEncoder(file)
	err = enc.Encode(c.data)
	if err != nil {
		return err
	}

	return nil
}

type CachedHTTPClient struct {
	Client  *http.Client
	Cache   *Cache
	cookies []*http.Cookie
	mu      sync.Mutex
}

func NewCachedHTTPClient(cacheToMemory bool, cacheFile string) (*CachedHTTPClient, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Transport: &http.Transport{},
		Jar:       jar,
	}

	cache := &Cache{
		data:          make(map[string]*cacheItem),
		cacheToMemory: cacheToMemory,
		cacheFile:     cacheFile,
		mu:            sync.Mutex{},
	}

	if !cacheToMemory {
		err = cache.loadFromFile()
		if err != nil {
			return nil, err
		}
	}

	return &CachedHTTPClient{
		Client:  client,
		Cache:   cache,
		cookies: nil,
		mu:      sync.Mutex{},
	}, nil
}

func (chc *CachedHTTPClient) setCookies(cookies []*http.Cookie) {
	chc.mu.Lock()
	defer chc.mu.Unlock()
	chc.cookies = cookies
}

func (chc *CachedHTTPClient) getCookies() []*http.Cookie {
	chc.mu.Lock()
	defer chc.mu.Unlock()
	return chc.cookies
}

func (chc *CachedHTTPClient) generateCacheKey(req *http.Request) (string, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	err := enc.Encode(req.Header)
	if err != nil {
		return "", err
	}

	hash := sha1.New()
	hash.Write(buf.Bytes())
	hash.Write([]byte(req.URL.String()))

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func (chc *CachedHTTPClient) CachedHttpGet(urlStr string, headers http.Header, ttl time.Duration) ([]byte, error) {
	_, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, err
	}

	req.Header = headers
	for _, cookie := range chc.getCookies() {
		req.AddCookie(cookie)
	}

	cacheKey, err := chc.generateCacheKey(req)
	if err != nil {
		return nil, err
	}

	cacheData, isCached := chc.Cache.Get(cacheKey)
	if isCached {
		return cacheData, nil
	}

	resp, err := chc.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	chc.Cache.Set(cacheKey, body, ttl)
	if !chc.Cache.cacheToMemory {
		err = chc.Cache.saveToFile()
		if err != nil {
			return nil, err
		}
	}

	// Set cookies from server response
	chc.setCookies(resp.Cookies())

	return body, nil
}

func (chc *CachedHTTPClient) CachedHttpPost(url string, headers http.Header, body []byte, ttl time.Duration) ([]byte, error) {
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header[k] = v
	}

	cookiesToSend := chc.getCookies()
	for _, cookie := range cookiesToSend {
		req.AddCookie(cookie)
	}

	cacheKey, err := chc.generateCacheKey(req)
	if err != nil {
		return nil, err
	}

	if data, ok := chc.Cache.Get(cacheKey); ok {
		return data, nil
	}

	res, err := chc.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	chc.mu.Lock()
	defer chc.mu.Unlock()
	chc.Cache.Set(cacheKey, data, ttl)

	return data, nil
}
