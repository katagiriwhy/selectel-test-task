package analyzer

import (
	"go/ast"
	"go/token"
	"strconv"

	"golang.org/x/tools/go/analysis"
)

func isLoggerCall(call *ast.CallExpr) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	switch sel.Sel.Name {
	case "Info", "Error", "Warn", "Debug":
		return true
	default:
		return false
	}
}

func extractMessage(call *ast.CallExpr) (string, bool) {
	if len(call.Args) == 0 {
		return "", false
	}

	lit, ok := call.Args[0].(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return "", false
	}

	msg, err := strconv.Unquote(lit.Value)
	if err != nil {
		return "", false
	}

	return msg, true
}

func checkCall(pass *analysis.Pass, call *ast.CallExpr) {
	if !isLoggerCall(call) {
		return
	}

	msg, ok := extractMessage(call)
	if !ok {
		return
	}

	if !checkLowercase(msg) {
		pass.Reportf(call.Pos(), "log message must start with lowercase letter")
	}

	if !checkEnglish(msg) {
		pass.Reportf(call.Pos(), "log message must be in English")
	}

	if !checkNoSpecialChars(msg) {
		pass.Reportf(call.Pos(), "log message contains special characters or emoji")
	}

	if !checkSensitive(msg) {
		pass.Reportf(call.Pos(), "log message may contain sensitive data")
	}
}
