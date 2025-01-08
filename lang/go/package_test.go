package pgsgo

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	pgs "github.com/joakim-brannstrom/protoc-gen-star"
)

func TestPackageName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		dir      string
		expected pgs.Name
	}{
		{"keyword", "_package"},              // go keywords are prefixed with _
		{"package", "my_package"},            // use the go_package option
		{"import", "bar"},                    // uses the basename if go_package contains a /
		{"override", "baz"},                  // if go_package contains ;, use everything to the right
		{"mapped", "unaffected"},             // M mapped params are ignored for build targets
		{"import_path_mapped", "go_package"}, // mixed import_path and M parameters should lose to go_package
		{"transitive_package", "foobar"},     // go_option gets picked up from other files if present
		{"path_dash", "path_dash"},           // if basename of go_package contains invalid characters, replace with _
	}

	for _, test := range tests {
		tc := test
		t.Run(tc.dir, func(t *testing.T) {
			t.Parallel()

			ast := buildGraph(t, "names", tc.dir)
			ctx := loadContext(t, "names", tc.dir)

			for _, target := range ast.Targets() {
				assert.Equal(t, tc.expected, ctx.PackageName(target))
			}
		})
	}
}

func TestImportPath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		dir string

		fully       pgs.FilePath
		unqualified pgs.FilePath
		none        pgs.FilePath
	}{
		{ // no params changing the behavior of the import paths
			"no_options",
			"example.com/packages/targets/fully_qualified",
			"targets/unqualified",
			"targets/none",
		},
		{ // M params provided for each imported package
			"mapped",
			"example.com/foo/bar",
			"example.com/fizz/buzz",
			"example.com/quux",
		},
	}

	for _, test := range tests {
		tc := test
		t.Run(tc.dir, func(t *testing.T) {
			t.Parallel()

			ast := buildGraph(t, "packages", tc.dir)
			ctx := loadContext(t, "packages", tc.dir)

			pkgs := map[string]pgs.FilePath{
				"packages.targets.fully_qualified": tc.fully,
				"packages.targets.unqualified":     tc.unqualified,
				"packages.targets.none":            tc.none,
			}

			for pkg, expected := range pkgs {
				t.Run(pkg, func(t *testing.T) {
					p, ok := ast.Packages()[pkg]
					require.True(t, ok, "package not found")
					f := p.Files()[0]
					assert.Equal(t, expected, ctx.ImportPath(f))
				})
			}
		})
	}
}

func TestOutputPath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		dir, file string
		expected  pgs.FilePath
	}{
		{"none", "none.proto", "none.pb.go"},
		{"none_srcrel", "none.proto", "none.pb.go"},
		{"unqualified", "unqualified.proto", "unqualified.pb.go"},
		{"unqualified_srcrel", "unqualified.proto", "unqualified.pb.go"},
		{"qualified", "qualified.proto", "example.com/qualified/qualified.pb.go"},
		{"qualified_srcrel", "qualified.proto", "qualified.pb.go"},
		{"mapped", "mapped.proto", "mapped.pb.go"},
		{"mapped_srcrel", "mapped.proto", "mapped.pb.go"},
	}

	for _, test := range tests {
		tc := test
		t.Run(tc.dir, func(t *testing.T) {
			t.Parallel()

			ast := buildGraph(t, "outputs", tc.dir)
			ctx := loadContext(t, "outputs", tc.dir)
			f, ok := ast.Lookup(tc.file)
			require.True(t, ok, "file not found")
			assert.Equal(t, tc.expected, ctx.OutputPath(f))
		})
	}

}
