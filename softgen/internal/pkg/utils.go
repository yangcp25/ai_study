package pkg

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

/* ----------------- 辅助函数 ----------------- */

func BuildPrompt(promptTmpl, example, softwareName string) string {
	promptTmpl = strings.ReplaceAll(promptTmpl, "[doc_name]", softwareName)
	return fmt.Sprintf(promptTmpl, example)
}
func BuildCodePrompt(promptTmpl, example string) string {
	return fmt.Sprintf(promptTmpl, example)
}

func SaveMarkdown(typ FileType, content, outputDir string) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	// , sanitizeFileName(softwareName)
	file := fmt.Sprintf("%s.md", typ)
	path := filepath.Join(outputDir, file)

	return os.WriteFile(path, []byte(content), 0644)
}
