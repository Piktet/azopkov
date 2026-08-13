/*
Package staticlint предоставляет multichecker для статического анализа Go-кода.

Multichecker включает:
  - Стандартные анализаторы из golang.org/x/tools/go/analysis/passes
  - Все анализаторы класса SA из staticcheck.io
  - Анализаторы других классов из staticcheck.io
  - Кастомный анализатор exitchecker, запрещающий использование  в функции main пакета main

Запуск:

	go run ./cmd/staticlint ./...
*/
package main

import (
	"github.com/Piktet/azopkov.git/cmd/staticlint/noexit"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"honnef.co/go/tools/simple"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"
)

func main() {
	var analyzers []*analysis.Analyzer

	analyzers = append(analyzers,
		copylock.Analyzer,
		loopclosure.Analyzer,
		lostcancel.Analyzer,
		printf.Analyzer,
		structtag.Analyzer,
		unreachable.Analyzer,
		shadow.Analyzer,
	)

	for _, v := range staticcheck.Analyzers {
		if len(v.Analyzer.Name) >= 2 && v.Analyzer.Name[:2] == "SA" {
			analyzers = append(analyzers, v.Analyzer)
		}
	}

	for _, v := range stylecheck.Analyzers {
		analyzers = append(analyzers, v.Analyzer)
	}

	count := 0
	for _, v := range simple.Analyzers {
		if count >= 2 {
			break
		}
		analyzers = append(analyzers, v.Analyzer)
		count++
	}

	analyzers = append(analyzers, noexit.Analyzer)

	multichecker.Main(analyzers...)
}
