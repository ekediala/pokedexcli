package cache

import (
	"testing"
	"time"
)
var store *Cache
func TestMain(m *testing.M) {
	store = New(time.Minute * 5)
	m.Run()
}
