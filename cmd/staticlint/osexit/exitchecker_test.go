package osexit_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"aprokhorov-praktikum/cmd/staticlint/osexit"
)

func TestOSExitAnalyzer(t *testing.T) {
	// функция analysistest.Run применяет тестируемый анализатор ErrCheckAnalyzer
	// к пакетам из папки testdata и проверяет ожидания
	analysistest.Run(t, analysistest.TestData(), osexit.OSExitAnalyzer, "./...")
}
