// A simple non-concurrent file cache
package cache

import "testing"

const testPath = "/Users/salah/Dropbox/code/go/src/github.com/drgo/core/cache/tmpl"

func TestFileCache_PopulateList(t *testing.T) {
	c := NewFileCache(testPath)
	if c.Len() != 3 {
		t.Errorf("Len() failed expected %d got %d", 3, c.Len())
	}
	contents, ok := c.Get(testPath + "/hello.gohtml")
	if !ok {
		t.Errorf("Get(/hello.gohtml) failed: /hello.gohtml not found")
	}
	if contents != "hello" {
		t.Errorf("Get(/hello.gohtml) failed: expected %s got %s", "hello", contents)
	}
	contents, ok = c.Get(testPath + "/empty.gohtml")
	if contents != "" {
		t.Errorf("Get(/empty.gohtml) failed: expected empty file got %s", contents)
	}
	contents, ok = c.Get(testPath + "/nested/nested.gohtml")
	if !ok {
		t.Errorf("Get(/nested/nested.gohtml) failed: /nested/nested.gohtml not found")
	}
	if contents != "nested" {
		t.Errorf("Get(/nested/nested.gohtml) failed: expected %s got %s", "nested", contents)
	}
	contents, ok = c.Get(testPath + "/doesnotExist")
	if ok {
		t.Errorf("Get(/doesnotExist) failed: it returned ok for non-existing file")
	}

}
