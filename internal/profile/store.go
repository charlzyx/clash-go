package profile

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// CatalogItem is one subscription in profiles.json. It is not ImportHTTP's
// Item: that type is the fetch result, this one is the saved index.
type CatalogItem struct {
	UID         string        `json:"uid"`
	Name        string        `json:"name"`
	Description string        `json:"desc,omitempty"`
	URL         string        `json:"url,omitempty"`
	Type        string        `json:"type"`
	Home        string        `json:"home,omitempty"`
	Updated     time.Time     `json:"updated,omitempty"`
	Used        int64         `json:"used,omitempty"`
	Total       int64         `json:"total,omitempty"`
	Expire      time.Time     `json:"expire,omitempty"`
	File        string        `json:"file,omitempty"`
	Option      CatalogOption `json:"option"`
}

// CatalogOption is the remote form and the extension slots stored in the index.
type CatalogOption struct {
	UpdateInterval  int    `json:"update_interval,omitempty"`
	UseSystemProxy  bool   `json:"use_system_proxy,omitempty"`
	UseClashProxy   bool   `json:"use_clash_proxy,omitempty"`
	AcceptInvalid   bool   `json:"allow_insecure,omitempty"`
	AllowAutoUpdate bool   `json:"allow_auto_update,omitempty"`
	Rules           string `json:"rules,omitempty"`
	Proxies         string `json:"proxies,omitempty"`
	Groups          string `json:"groups,omitempty"`
	Merge           string `json:"merge,omitempty"`
	Script          string `json:"script,omitempty"`
	UserAgent       string `json:"user_agent,omitempty"`
	TimeoutSeconds  int    `json:"timeout,omitempty"`
}

// File is profiles.json.
type File struct {
	Items   []CatalogItem `json:"items"`
	Current string        `json:"current"`
}

// Store loads and saves the profile list inside a config directory.
type Store struct {
	dir string
}

// NewStore returns a store for dir.
func NewStore(dir string) *Store { return &Store{dir: dir} }

func (s *Store) path() string { return filepath.Join(s.dir, "profiles.json") }

// Load returns an empty file when profiles.json is missing.
func (s *Store) Load() (File, error) {
	if s == nil || s.dir == "" {
		return File{}, errors.New("profile store has no directory")
	}
	b, err := os.ReadFile(s.path())
	if errors.Is(err, os.ErrNotExist) {
		return File{}, nil
	}
	if err != nil {
		return File{}, err
	}
	if len(strings.TrimSpace(string(b))) == 0 {
		return File{}, nil
	}
	var f File
	if err := json.Unmarshal(b, &f); err != nil {
		return File{}, err
	}
	return f, nil
}

// Save writes profiles.json.
func (s *Store) Save(f File) error {
	if s == nil || s.dir == "" {
		return errors.New("profile store has no directory")
	}
	if err := os.MkdirAll(s.dir, 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		return err
	}
	b = append(b, '\n')
	return os.WriteFile(s.path(), b, 0o644)
}

// ErrInvalidURL is returned when a subscription link is not http(s).
var ErrInvalidURL = errors.New("无效的订阅链接，请输入以 http:// 或 https:// 开头的地址")

// ValidURL reports whether raw is an http or https URL with a host.
func ValidURL(raw string) bool {
	raw = strings.TrimSpace(raw)
	if !strings.HasPrefix(raw, "http://") && !strings.HasPrefix(raw, "https://") {
		return false
	}
	rest := strings.TrimPrefix(strings.TrimPrefix(raw, "https://"), "http://")
	host, _, _ := strings.Cut(rest, "/")
	return host != "" && !strings.Contains(host, " ")
}

// ImportRemote fetches url, stores the body, and appends a remote item.
// The new item becomes current when the file had none.
func (s *Store) ImportRemote(raw string, fetch func(string) ([]byte, error)) (CatalogItem, error) {
	raw = strings.TrimSpace(raw)
	if !ValidURL(raw) {
		return CatalogItem{}, ErrInvalidURL
	}
	if fetch == nil {
		return CatalogItem{}, errors.New("no fetcher")
	}
	body, err := fetch(raw)
	if err != nil {
		return CatalogItem{}, err
	}
	f, err := s.Load()
	if err != nil {
		return CatalogItem{}, err
	}
	uid := uniqueID(f.Items)
	name := hostOf(raw)
	dir := filepath.Join(s.dir, "profiles")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return CatalogItem{}, err
	}
	rel := filepath.Join("profiles", uid+".yaml")
	if err := os.WriteFile(filepath.Join(s.dir, rel), body, 0o644); err != nil {
		return CatalogItem{}, err
	}
	item := CatalogItem{
		UID:     uid,
		Name:    name,
		URL:     raw,
		Type:    "remote",
		Updated: time.Now().UTC().Truncate(time.Second),
		File:    rel,
		Option:  CatalogOption{UpdateInterval: 1440, AllowAutoUpdate: true, TimeoutSeconds: 60},
	}
	f.Items = append(f.Items, item)
	if f.Current == "" {
		f.Current = uid
	}
	if err := s.Save(f); err != nil {
		return CatalogItem{}, err
	}
	return item, nil
}

// ImportLocal stores a YAML file dropped on the window.
func (s *Store) ImportLocal(filename string, body []byte) (CatalogItem, error) {
	if !strings.HasSuffix(strings.ToLower(filename), ".yaml") && !strings.HasSuffix(strings.ToLower(filename), ".yml") {
		return CatalogItem{}, errors.New("仅支持 YAML 文件")
	}
	f, err := s.Load()
	if err != nil {
		return CatalogItem{}, err
	}
	uid := uniqueID(f.Items)
	name := filename
	if i := strings.LastIndexAny(name, `/\`); i >= 0 {
		name = name[i+1:]
	}
	if name == "" {
		name = "New Profile"
	}
	dir := filepath.Join(s.dir, "profiles")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return CatalogItem{}, err
	}
	rel := filepath.Join("profiles", uid+".yaml")
	if err := os.WriteFile(filepath.Join(s.dir, rel), body, 0o644); err != nil {
		return CatalogItem{}, err
	}
	item := CatalogItem{UID: uid, Name: name, Type: "local", File: rel, Updated: time.Now().UTC().Truncate(time.Second)}
	f.Items = append(f.Items, item)
	if f.Current == "" {
		f.Current = uid
	}
	if err := s.Save(f); err != nil {
		return CatalogItem{}, err
	}
	return item, nil
}

func uniqueID(items []CatalogItem) string {
	n := len(items) + 1
	for {
		id := "p" + itoa(n)
		taken := false
		for _, it := range items {
			if it.UID == id {
				taken = true
				break
			}
		}
		if !taken {
			return id
		}
		n++
	}
}

func hostOf(raw string) string {
	rest := strings.TrimPrefix(strings.TrimPrefix(raw, "https://"), "http://")
	host, _, _ := strings.Cut(rest, "/")
	if host == "" {
		return "订阅"
	}
	return host
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var b [16]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	return string(b[i:])
}
