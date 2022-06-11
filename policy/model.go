package policy

import (
	mapset "github.com/deckarep/golang-set"
)

type ReplyType int

const (
	Failed ReplyType = iota
	Text
	Image
)

type Reply struct {
	Type ReplyType
	Body []byte
}

type ContentItem struct {
	Img      string   `json:"img"`
	Emoji    []string `json:"emoji"`
	Emoticon []string `json:"emoticon"`
	Keyword  []string `json:"keyword"`
	Role     []string `json:"role"`
}
type ImgDef struct {
	Version string        `json:"version"`
	Content []ContentItem `json:"content"`
}

func (d *ImgDef) GetMaps() SymbolMaps {
	content := d.Content
	// var imgNameSet mapset.Set = mapset.NewSet()
	var emojiMap map[string]mapset.Set = make(map[string]mapset.Set)
	var emotiMap map[string]mapset.Set = make(map[string]mapset.Set)
	var kwordMap map[string]mapset.Set = make(map[string]mapset.Set)
	var roleMap map[string]mapset.Set = make(map[string]mapset.Set)
	for _, item := range content {
		keysValueToMap(item.Emoji, item.Img, emojiMap)
		keysValueToMap(item.Emoticon, item.Img, emotiMap)
		keysValueToMap(item.Keyword, item.Img, kwordMap)
		keysValueToMap(item.Role, item.Img, roleMap)
	}
	return SymbolMaps{emojiMap, emotiMap, kwordMap, roleMap}
}

type SymbolMaps struct {
	EmojiMap map[string]mapset.Set
	EmotiMap map[string]mapset.Set
	KwordMap map[string]mapset.Set
	RoleMap  map[string]mapset.Set
}

func (s *SymbolMaps) ContainsEmoji(key string) bool {
	return mapContainsKey(s.EmojiMap, key)
}

func (s *SymbolMaps) ContainsEmoti(key string) bool {
	return mapContainsKey(s.EmotiMap, key)
}

func (s *SymbolMaps) ContainsKword(key string) bool {
	return mapContainsKey(s.KwordMap, key)
}

func (s *SymbolMaps) ContainsRole(key string) bool {
	return mapContainsKey(s.RoleMap, key)
}

func (s *SymbolMaps) ContainsAny(key string) bool {
	return s.ContainsEmoji(key) || s.ContainsEmoti(key) || s.ContainsKword(key) || s.ContainsRole(key)
}

/* func (s *SymbolMaps) GetAng(key string) mapset.Set {
	emojiSet := s.EmojiMap[key]
	emotiSet := s.EmotiMap[key]
	kwordSet := s.KwordMap[key]
	roleSet := s.RoleMap[key]
	emojiSet.Union(emotiSet).Union(kwordSet).Union(roleSet)
} */
