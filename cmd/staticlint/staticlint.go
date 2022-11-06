// Пакет staticlint предназначен для проверки кода линтерами.
package main

import (
	"strings"

	"github.com/kisielk/errcheck/errcheck"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/deepequalerrors"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/stdmethods"
	"golang.org/x/tools/go/analysis/passes/stringintconv"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"golang.org/x/tools/go/analysis/passes/usesgenerics"
	"honnef.co/go/tools/staticcheck"

	"aprokhorov-praktikum/cmd/staticlint/osexit"
)

// Собирает необходимые Анализаторы
// 		assign - detects useless assignments.
//		atomic - checks for common mistakes using the sync/atomic package.
//		bools - detects common mistakes involving boolean operators.
//		deepequalerrors - checks for the use of reflect.DeepEqual with error values.
//		httpresponse - checks for mistakes using HTTP responses.
//		loopclosure - checks for references to enclosing loop variables from within nested functions.
//		nilfunc - checks for useless comparisons against nil.
//		printf - checks consistency of Printf format strings and arguments.
//		shadow - checks for shadowed variables.
//		stdmethods - checks for misspellings in the signatures of methods similar to well-known interfaces.
//		stringintconv - flags type conversions from integers to strings.
//		structtag - checks struct field tags are well formed.
//		unmarshal - checks for passing non-pointer or non-interface types to unmarshal and decode functions.
//		unreachable - checks for unreachable code.
//		unusedresult - checks for unused results of calls to certain pure functions.
//		usesgenerics - checks for usage of generic features added in Go 1.18.
//		osexit.OSExitAnalyzer - checks for usage of os.Exit() in func main() of package main
func main() {
	totalchecks := []*analysis.Analyzer{
		assign.Analyzer,
		atomic.Analyzer,
		bools.Analyzer,
		deepequalerrors.Analyzer,
		httpresponse.Analyzer,
		loopclosure.Analyzer,
		nilfunc.Analyzer,
		printf.Analyzer,
		shadow.Analyzer,
		stdmethods.Analyzer,
		stringintconv.Analyzer,
		structtag.Analyzer,
		unmarshal.Analyzer,
		unreachable.Analyzer,
		unusedresult.Analyzer,
		usesgenerics.Analyzer,
	}

	// Ошибки при работе, пришлось пока выключить SA-проверки

	for _, v := range staticcheck.Analyzers {
		if strings.HasPrefix(v.Name, "ST") { // || strings.HasPrefix(v.Name, "SA") {
			totalchecks = append(totalchecks, v)
		}
	}

	totalchecks = append(totalchecks, errcheck.Analyzer)

	totalchecks = append(totalchecks, osexit.OSExitAnalyzer)

	multichecker.Main(
		totalchecks...,
	)
}
