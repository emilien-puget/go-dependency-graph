package package_list

import (
	"path/filepath"
	"testing"

	"github.com/emilien-puget/go-dependency-graph/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestGetPackagesToParse(t *testing.T) {
	tests := map[string]struct {
		pathDir string
		want    []string
		wantErr bool
	}{
		"inter": {
			pathDir: "../testdata/inter",
			want:    []string{"pa", "inter"},
			wantErr: false,
		},
		"external_dep": {
			pathDir: "../testdata/ext_dep",
			want:    []string{"ext_dep"},
			wantErr: false,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			abs, err := filepath.Abs(tt.pathDir)
			require.NoError(t, err)
			got, err := GetPackagesToParse(abs, []string{config.VendorDir})
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPackagesToParse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			gotPackagesName := make([]string, 0, len(got))
			for i := range got {
				gotPackagesName = append(gotPackagesName, got[i].Name)
			}
			require.Equal(t, tt.want, gotPackagesName)
		})
	}
}

func BenchmarkGetPackagesToParse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = GetPackagesToParse("../testdata/inter", nil)
	}
}
