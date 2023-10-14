package files

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func testCloseAndName(t *testing.T) {
	tests := []struct {
		f       *os.File
		dir     string
		base    string
		ext     string
		want    string
		wantErr bool
	}{
		{
			dir:  "./test",
			base: "test",
			ext:  ".txt",
		},
		{
			dir:  "",
			base: "test",
			ext:  ".txt",
		},
	}
	for _, tt := range tests {
		t.Run(tt.base, func(t *testing.T) {
			tt.f, _ = ioutil.TempFile("", "")
			got, err := CloseAndName(tt.f, filepath.Join(tt.dir, tt.base+tt.ext))

			if (err != nil) != tt.wantErr {
				t.Errorf("CloseAndName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CloseAndName() = %v, want %v", got, "/Users/salah/local/git/core/files/test/test1.txt")
			}
		})
	}
}
