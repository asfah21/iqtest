//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	questions := []string{
		"q_mtx_001", "q_mtx_002", "q_mtx_003", "q_mtx_004", "q_mtx_005", "q_mtx_006",
		"q_seq_001", "q_seq_002", "q_seq_003", "q_seq_004", "q_seq_005",
		"q_spa_001", "q_spa_002", "q_spa_003", "q_spa_004", "q_spa_005",
		"q_anl_001", "q_anl_002", "q_anl_003", "q_anl_004",
	}

	options := []string{"opt_a", "opt_b", "opt_c", "opt_d", "opt_a2", "opt_b2", "opt_c2", "opt_d2"}

	outDir := filepath.Join("assets", "images")
	os.MkdirAll(outDir, 0755)

	// Question SVGs — circle + pattern
	qSvg := `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 150 150" width="150" height="150">
<rect width="150" height="150" fill="#f8fafc" rx="8" stroke="#e2e8f0" stroke-width="1"/>
<circle cx="75" cy="75" r="40" fill="#eef2ff" stroke="#6366f1" stroke-width="2"/>
<circle cx="75" cy="75" r="15" fill="#6366f1" opacity="0.5"/>
<text x="75" y="135" text-anchor="middle" font-family="Inter,sans-serif" font-size="9" fill="#94a3b8">%s</text>
</svg>`

	for _, name := range questions {
		svg := fmt.Sprintf(qSvg, name)
		os.WriteFile(filepath.Join(outDir, name+".svg"), []byte(svg), 0644)
	}

	// Option SVGs — rectangle + pattern
	oSvg := `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 120 120" width="120" height="120">
<rect width="120" height="120" fill="#ffffff" rx="6" stroke="#e2e8f0" stroke-width="1"/>
<rect x="30" y="30" width="60" height="60" rx="4" fill="#eef2ff" stroke="#6366f1" stroke-width="1.5"/>
<text x="60" y="110" text-anchor="middle" font-family="Inter,sans-serif" font-size="8" fill="#94a3b8">%s</text>
</svg>`

	for _, name := range options {
		svg := fmt.Sprintf(oSvg, name)
		os.WriteFile(filepath.Join(outDir, name+".svg"), []byte(svg), 0644)
	}

	fmt.Println("Generated", len(questions)+len(options), "SVG files in", outDir)
}
