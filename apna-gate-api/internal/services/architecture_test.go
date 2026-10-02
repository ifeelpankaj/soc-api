package service

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// Paths are checked independently of aliases used in source code.
func TestPersistenceDependencyBoundary(t *testing.T) {
	for _, root := range []string{".", "../repositories/contracts"} {
		err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if entry.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
				return nil
			}
			file, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
			if err != nil {
				return err
			}
			for _, imp := range file.Imports {
				name, err := strconv.Unquote(imp.Path.Value)
				if err != nil {
					return err
				}
				if name == "go-server/internal/db" || strings.HasPrefix(name, "go-server/internal/db/") || name == "go-server/pkg/database" || strings.HasPrefix(name, "github.com/jackc/") || name == "database/sql" || name == "go-server/internal/repositories" {
					t.Errorf("%s imports persistence implementation %q", path, name)
				}
			}
			ast.Inspect(file, func(n ast.Node) bool {
				if sel, ok := n.(*ast.SelectorExpr); ok {
					switch sel.Sel.Name {
					case "GetQueries", "OutboxQueries", "FinancialQueries", "WithQueries":
						t.Errorf("%s exposes query access through %s", path, sel.Sel.Name)
					}
				}
				return true
			})
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	}
}
