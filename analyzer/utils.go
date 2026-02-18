package analyzer

import (
	"go/ast"
	"go/token"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/enescakir/emoji"
	"golang.org/x/tools/go/analysis"

	"selectel/config"
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

func extractMessage(call *ast.CallExpr) (string, *ast.BasicLit, bool) {
	if len(call.Args) == 0 {
		return "", nil, false
	}

	lit, ok := call.Args[0].(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return "", nil, false
	}

	msg, err := strconv.Unquote(lit.Value)
	if err != nil {
		return "", nil, false
	}

	return msg, lit, true
}

func lowercaseFix(msg string, lit *ast.BasicLit) *analysis.SuggestedFix {
	if msg == "" {
		return nil
	}

	r, size := utf8.DecodeRuneInString(msg)
	if !unicode.IsLetter(r) || unicode.IsLower(r) {
		return nil
	}

	lowered := string(unicode.ToLower(r)) + msg[size:]
	quoted := strconv.Quote(lowered)

	return &analysis.SuggestedFix{
		Message: "make log message start with lowercase letter",
		TextEdits: []analysis.TextEdit{
			{
				Pos:     lit.Pos(),
				End:     lit.End(),
				NewText: []byte(quoted),
			},
		},
	}
}

func specialCharsFix(msg string, lit *ast.BasicLit, cfg *config.Config) *analysis.SuggestedFix {
	if cfg == nil {
		cfg = config.Load("")
	}

	changed := false
	out := make([]rune, 0, len(msg))
	for _, r := range msg {
		if emoji.Exist(string(r)) {
			changed = true
			continue
		}
		if strings.ContainsRune(cfg.ForbiddenSymbols, r) {
			changed = true
			continue
		}
		out = append(out, r)
	}

	if !changed {
		return nil
	}

	newMsg := string(out)
	quoted := strconv.Quote(newMsg)

	return &analysis.SuggestedFix{
		Message: "remove special characters or emoji from log message",
		TextEdits: []analysis.TextEdit{
			{
				Pos:     lit.Pos(),
				End:     lit.End(),
				NewText: []byte(quoted),
			},
		},
	}
}

func checkCall(pass *analysis.Pass, call *ast.CallExpr) {
	if !isLoggerCall(call) {
		return
	}

	msg, lit, ok := extractMessage(call)
	if !ok {
		return
	}

	cfg := config.Load("")

	// lowercase
	if cfg.Rules.Lowercase && !checkLowercase(msg) {
		diag := analysis.Diagnostic{
			Pos:     call.Pos(),
			Message: "log message must start with lowercase letter",
		}
		if fix := lowercaseFix(msg, lit); fix != nil {
			diag.SuggestedFixes = append(diag.SuggestedFixes, *fix)
		}
		pass.Report(diag)
	}

	// english only
	if cfg.Rules.EnglishOnly && !checkEnglish(msg) {
		pass.Reportf(call.Pos(), "log message must be in English")
	}

	// special chars / emoji
	if cfg.Rules.NoSpecialSymbols && !checkNoSpecialChars(msg, cfg) {
		diag := analysis.Diagnostic{
			Pos:     call.Pos(),
			Message: "log message contains special characters or emoji",
		}
		if fix := specialCharsFix(msg, lit, cfg); fix != nil {
			diag.SuggestedFixes = append(diag.SuggestedFixes, *fix)
		}
		pass.Report(diag)
	}

	// sensitive data
	if cfg.Rules.SensitiveData && !checkSensitive(msg, cfg) {
		pass.Reportf(call.Pos(), "log message may contain sensitive data")
	}
}
