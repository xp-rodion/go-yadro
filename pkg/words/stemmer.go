package words

import (
	"flag"
	"fmt"
	"github.com/kljensen/snowball"
	"regexp"
	"strings"
)

// IsSeparateRune Функция для проверки разделительного символа
func IsSeparateRune(r rune) bool {
	switch r {
	// Ввел самые частые разделительные знаки
	case ' ', ':', '!', ',':
		return true
	}
	return false
}

// ReadCLIArgs Чтение CLI ввода с флагом -s
func ReadCLIArgs() []string {
	var sentence string
	flag.StringVar(&sentence, "s", "", "Введите ваше предложение: ")
	flag.Parse()
	if sentence == "" {
		fmt.Println("Слова не найдены")
		return nil
	}
	return strings.FieldsFunc(sentence, IsSeparateRune)
}

// Split Разбиение строки на слайс строк
func Split(words string) []string {
	return strings.FieldsFunc(words, IsSeparateRune)
}

// IsStopWord Проверка на стоп слово
func IsStopWord(word string) bool {
	switch word {
	case "a", "about", "above", "after", "again", "against", "all", "am", "an",
		"and", "any", "are", "as", "at", "be", "because", "been", "before",
		"being", "below", "between", "both", "but", "by", "can", "did", "do",
		"does", "doing", "don", "down", "during", "each", "few", "for", "from",
		"further", "had", "has", "have", "having", "he", "her", "here", "hers",
		"herself", "him", "himself", "his", "how", "i", "if", "in", "into", "is",
		"it", "its", "itself", "just", "me", "more", "most", "my", "myself",
		"no", "nor", "not", "now", "of", "off", "on", "once", "only", "or",
		"other", "our", "ours", "ourselves", "out", "over", "own", "s", "same",
		"she", "should", "so", "some", "such", "t", "than", "that", "the", "their",
		"theirs", "them", "themselves", "then", "there", "these", "they",
		"this", "those", "through", "to", "too", "under", "until", "up",
		"very", "was", "we", "were", "what", "when", "where", "which", "while",
		"who", "whom", "why", "will", "with", "you", "your", "yours", "yourself",
		"yourselves", "i'll":
		return true
	}
	return false
}

// NormalizeWord Нормализатор слова, удаление небуквенных символов
func NormalizeWord(word string) (normalizeWord string) {
	// pattern для regex
	re := regexp.MustCompile(`[^a-zA-Z']+`)
	normalizeWord = strings.ToLower(strings.TrimSpace(
		re.ReplaceAllString(word, ""),
	))
	return
}

// ValidateStopWords Валидатор, который возвращает слайс без стоп слов
func ValidateStopWords(words []string) (validateWords []string) {
	for _, word := range words {
		word = NormalizeWord(word)
		if !IsStopWord(word) {
			validateWords = append(validateWords, word)
		}
	}
	return
}

// StemWord Стемм введеного слова
func StemWord(word string) (string, error) {
	return snowball.Stem(word, "english", true)
}

// StemWords Обработка нормализированных слов, выдача результата;
func StemWords(words []string) (stemmingWords []string) {
	uniqueWords := make(map[string]bool)
	for _, word := range words {
		word, err := StemWord(word)

		// явно указываю проверки перед основным функционалом + исключаю вложенность условий

		if err != nil {
			continue
		}
		if uniqueWords[word] {
			continue
		}

		uniqueWords[word] = true
		stemmingWords = append(stemmingWords, word)
	}
	return
}

// Stemmer Скомпанованная функция
func Stemmer(sentence string) []string {
	words := Split(sentence)
	validateWords := ValidateStopWords(words)
	return StemWords(validateWords)
}
