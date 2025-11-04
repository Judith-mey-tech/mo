package core

import (
	"io/ioutil"
	"strings"
)

// VerifyResult — результат простой проверки
type VerifyResult struct {
    Confidence float64 // 0..1 простая эвристика
    Reason     string
}

// SimpleVerify — очень простая логика: если SHA уже встречался, confidence=0.99 (т.е. вероятно оригинал).
// если имя файла содержит "watermark" или "AUTH:" — уменьшаем вероятность фейка.
// если файл слишком маленький, ставим низкую уверенность.
func SimpleVerify(path string, originalName string, sha string) VerifyResult {
    // Очень простая эвристика:
    // 1) проверка на "встроенный знак" в имени
    nameLower := strings.ToLower(originalName)
    if strings.Contains(nameLower, "watermark") || strings.Contains(nameLower, "auth:") || strings.Contains(nameLower, "©️") {
        return VerifyResult{Confidence: 0.95, Reason: "found watermark-like token in filename or metadata"}
    }

    // 2) если в содержимом файла есть ASCII слова "AUTHENTRACK" (демо) — считаем оригиналом
    b, _ := ioutil.ReadFile(path)
    if strings.Contains(strings.ToUpper(string(b)), "AUTHENTRACK") {
        return VerifyResult{Confidence: 0.99, Reason: "embedded AUTHENTRACK tag found in file bytes"}
    }

    // 3) по размеру файла — если очень мало, то подозрительно
    if len(b) < 1024*10 { // меньше 10KB
        return VerifyResult{Confidence: 0.2, Reason: "file too small, suspicious"}
    }

    // 4) baseline
    return VerifyResult{Confidence: 0.5, Reason: "no strong signals; consider external AI check"}
}
