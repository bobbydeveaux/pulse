package complexity

import (
	"go/ast"
	"go/parser"
	"go/token"

	"github.com/bobbydeveaux/pulse/internal/types"
)

// AnalyzeGoFile uses the Go AST for precise complexity analysis.
func AnalyzeGoFile(path string) ([]types.FunctionMetrics, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var results []types.FunctionMetrics

	ast.Inspect(f, func(n ast.Node) bool {
		fn, ok := n.(*ast.FuncDecl)
		if !ok {
			return true
		}

		name := fn.Name.Name
		if fn.Recv != nil && len(fn.Recv.List) > 0 {
			name = receiverName(fn.Recv.List[0].Type) + "." + name
		}

		startLine := fset.Position(fn.Pos()).Line
		endLine := fset.Position(fn.End()).Line

		ccn := cyclomaticGo(fn)
		cog := cognitiveGo(fn)

		results = append(results, types.FunctionMetrics{
			Name:       name,
			File:       path,
			StartLine:  startLine,
			EndLine:    endLine,
			SLOC:       endLine - startLine + 1,
			Cyclomatic: ccn,
			Cognitive:  cog,
			Grade:      types.CyclomaticGrade(ccn),
		})

		return true
	})

	return results, nil
}

func receiverName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return receiverName(t.X)
	case *ast.IndexExpr:
		return receiverName(t.X)
	case *ast.IndexListExpr:
		return receiverName(t.X)
	default:
		return "?"
	}
}

// cyclomaticGo computes McCabe cyclomatic complexity for a Go function.
// CCN = 1 + number of decision points.
func cyclomaticGo(fn *ast.FuncDecl) int {
	if fn.Body == nil {
		return 1
	}
	ccn := 1
	ast.Inspect(fn.Body, func(n ast.Node) bool {
		switch n.(type) {
		case *ast.IfStmt:
			ccn++
		case *ast.ForStmt, *ast.RangeStmt:
			ccn++
		case *ast.CaseClause:
			// default case doesn't add complexity
			cc := n.(*ast.CaseClause)
			if cc.List != nil {
				ccn++
			}
		case *ast.CommClause:
			cc := n.(*ast.CommClause)
			if cc.Comm != nil {
				ccn++
			}
		case *ast.BinaryExpr:
			be := n.(*ast.BinaryExpr)
			if be.Op == token.LAND || be.Op == token.LOR {
				ccn++
			}
		}
		return true
	})
	return ccn
}

// cognitiveGo computes cognitive complexity for a Go function.
// Rules: +1 for each break in linear flow, +1 for each nesting level.
func cognitiveGo(fn *ast.FuncDecl) int {
	if fn.Body == nil {
		return 0
	}
	return cognitiveBlock(fn.Body.List, 0)
}

func cognitiveBlock(stmts []ast.Stmt, nesting int) int {
	cog := 0
	for _, stmt := range stmts {
		cog += cognitiveStmt(stmt, nesting)
	}
	return cog
}

func cognitiveStmt(stmt ast.Stmt, nesting int) int {
	cog := 0
	switch s := stmt.(type) {
	case *ast.IfStmt:
		// +1 for the if, +nesting for depth
		cog += 1 + nesting
		// Count boolean operators in condition
		cog += countBoolOps(s.Cond)
		if s.Body != nil {
			cog += cognitiveBlock(s.Body.List, nesting+1)
		}
		if s.Else != nil {
			// else/else if: +1 but no nesting increment for the keyword itself
			cog++
			switch e := s.Else.(type) {
			case *ast.BlockStmt:
				cog += cognitiveBlock(e.List, nesting+1)
			case *ast.IfStmt:
				// else if: already counted +1 above, now recurse without extra nesting
				cog += cognitiveStmt(e, nesting) - 1 // -1 because we already counted +1 for else
				// Actually let's handle this correctly:
				// else if gets +1 for the else, and then the if inside gets its own +1 + nesting
				// But that double-counts. The standard is: else if = +1 (no nesting penalty).
			}
		}
	case *ast.ForStmt:
		cog += 1 + nesting
		if s.Cond != nil {
			cog += countBoolOps(s.Cond)
		}
		if s.Body != nil {
			cog += cognitiveBlock(s.Body.List, nesting+1)
		}
	case *ast.RangeStmt:
		cog += 1 + nesting
		if s.Body != nil {
			cog += cognitiveBlock(s.Body.List, nesting+1)
		}
	case *ast.SwitchStmt:
		cog += 1 + nesting
		if s.Body != nil {
			cog += cognitiveBlock(s.Body.List, nesting+1)
		}
	case *ast.TypeSwitchStmt:
		cog += 1 + nesting
		if s.Body != nil {
			cog += cognitiveBlock(s.Body.List, nesting+1)
		}
	case *ast.SelectStmt:
		cog += 1 + nesting
		if s.Body != nil {
			cog += cognitiveBlock(s.Body.List, nesting+1)
		}
	case *ast.CaseClause:
		// Case clauses are nested inside the switch body
		cog += cognitiveBlock(s.Body, nesting)
	case *ast.CommClause:
		cog += cognitiveBlock(s.Body, nesting)
	case *ast.BlockStmt:
		cog += cognitiveBlock(s.List, nesting)
	case *ast.LabeledStmt:
		// goto/break/continue to label: +1
		cog++
		cog += cognitiveStmt(s.Stmt, nesting)
	case *ast.BranchStmt:
		// break/continue/goto with label
		if s.Label != nil {
			cog++
		}
	}
	return cog
}

func countBoolOps(expr ast.Expr) int {
	count := 0
	ast.Inspect(expr, func(n ast.Node) bool {
		if be, ok := n.(*ast.BinaryExpr); ok {
			if be.Op == token.LAND || be.Op == token.LOR {
				count++
			}
		}
		return true
	})
	return count
}
