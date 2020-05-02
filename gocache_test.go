package gocache

import (
	"testing"
	"time"
)

func TestCache_GetAndSet(t *testing.T) {
	testCache := NewCache(5*time.Second, 0)

	var found bool
	var item interface{}

	_, found = testCache.Get("key")
	if found {
		t.Error("Key Should not be exist")
	}

	_, found = testCache.Get("key1")
	if found {
		t.Error("Key Should not be exist")
	}

	testCache.Set("key", "value", testCache.DefaultLife)
	testCache.Set("key1", 1, testCache.DefaultLife)

	item, found = testCache.Get("key")
	if !found {
		t.Error("Key Should be exist")
	}

	if item == nil {
		t.Error("item should not be nil")
	}

	item, found = testCache.Get("key1")
	if !found {
		t.Error("Key Should be exist")
	}

	if item == nil {
		t.Error("item should not be nil")
	}
}

func TestCache_GC(t *testing.T) {

	var found bool
	var item interface{}

	testCache := NewCache(30*time.Millisecond, 1*time.Millisecond)

	testCache.Set("a", "valueA", testCache.DefaultLife)
	testCache.Set("b", "valueB", Eternal)
	testCache.Set("c", "valueC", 20*time.Millisecond)
	testCache.Set("d", "valueD", 40*time.Millisecond)

	<-time.After(30 * time.Millisecond)

	_, found = testCache.Get("c")
	if found {
		t.Error("it should have been deleted")
	}

	<-time.After(5 * time.Millisecond)
	_, found = testCache.Get("a")
	if found {
		t.Error("it should have been deleted")
	}

	<-time.After(20 * time.Millisecond)
	_, found = testCache.Get("d")
	if found {
		t.Error("it should have been deleted")
	}

	<-time.After(20 * time.Millisecond)
	item, found = testCache.Get("b")
	if !found {
		t.Error("it should be exist")
	}

	if item == nil {
		t.Error("it should have value")
	}

	if v, ok := item.(string); ok {
		if v != "valueB" {
			t.Error("it should have same value")
		}
	}
}

func TestCache_Add(t *testing.T) {
	testCache := NewCache(3*time.Second, 0)

	added, _ := testCache.Add("foo", "bar", testCache.DefaultLife)
	if !added {
		t.Error("It should have been added")
	}

	added, _ = testCache.Add("foo", "bar", testCache.DefaultLife)
	if added {
		t.Error("it should have returned an error")
	}
}

func TestCache_ClearCache(t *testing.T) {
	testCache := NewCache(3*time.Second, 0)

	added, _ := testCache.Add("foo", "bar2", testCache.DefaultLife)
	if !added {
		t.Error("It should have been added")
	}

	added, _ = testCache.Add("foo2", "bar2", testCache.DefaultLife)
	if !added {
		t.Error("It should have been added")
	}

	testCache.ClearCache()

	_, found := testCache.Get("c")
	if found {
		t.Error("it should have been deleted")
	}
}
