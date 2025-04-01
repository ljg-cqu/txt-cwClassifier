// Description: This program categorizes Chinese text into various linguistic categories
// Features:
// - Extracts Chinese characters, nouns, verbs, adjectives, adverbs, idioms, and slang
// - Categorizes text into noun phrases and verb phrases
// - Counts frequency of each category
// - Outputs results to separate text files
// Workflow:
// 1. Select an input text file containing Chinese text.
// 2. Select an output directory for the categorized files.
// 3. The program reads the input file, categorizes the text, and writes the results to output files.
// 4. Each category is saved in a separate text file, sorted by frequency of occurrence.
// 5. The program handles errors gracefully and provides user feedback.

package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode"

	"github.com/jdkato/prose/v2"
	"github.com/sqweek/dialog"
)

// Checks if a given string contains only Chinese characters
func isChineseText(text string) bool {
	for _, r := range text {
		if !unicode.Is(unicode.Han, r) && r != ' ' && r != '-' { // Allow spaces and hyphens
			return false
		}
	}
	return true
}

// Extracts and returns individual Chinese characters from a string
func extractChineseCharacters(text string) []string {
	var characters []string
	for _, r := range text {
		if unicode.Is(unicode.Han, r) {
			characters = append(characters, string(r))
		}
	}
	return characters
}

// Capitalizes the first character of each word or phrase
func capitalizePhrase(phrase string) string {
	runes := []rune(phrase)
	if len(runes) > 0 {
		runes[0] = unicode.ToUpper(runes[0])
	}
	return string(runes)
}

// Counts appearances of items and stores them in a frequency map
func countFrequencies(content []string) map[string]int {
	counts := make(map[string]int)
	for _, item := range content {
		capitalizedItem := capitalizePhrase(item)
		counts[capitalizedItem]++
	}
	return counts
}

// Converts frequency map to sorted slice (only items, sorted by frequency)
func sortByFrequency(counts map[string]int) []string {
	type itemFrequency struct {
		Item      string
		Frequency int
	}
	var items []itemFrequency
	for item, freq := range counts {
		items = append(items, itemFrequency{Item: item, Frequency: freq})
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].Frequency > items[j].Frequency
	})
	var sortedItems []string
	for _, entry := range items {
		sortedItems = append(sortedItems, entry.Item)
	}
	return sortedItems
}

// Extracts noun phrases using Chinese POS rules
func extractNounPhrases(tokens []prose.Token) []string {
	var nounPhrases []string
	var currentPhrase []string

	for _, tok := range tokens {
		if isChineseText(tok.Text) {
			switch tok.Tag {
			case "DT", "NN", "JJ": // Determiners, Nouns, Adjectives
				currentPhrase = append(currentPhrase, tok.Text)
			default:
				if len(currentPhrase) > 0 {
					nounPhrases = append(nounPhrases, strings.Join(currentPhrase, " "))
					currentPhrase = nil
				}
			}
		}
	}

	if len(currentPhrase) > 0 {
		nounPhrases = append(nounPhrases, strings.Join(currentPhrase, " "))
	}

	return nounPhrases
}

// Extracts verb phrases using Chinese POS rules
func extractVerbPhrases(tokens []prose.Token) []string {
	var verbPhrases []string
	var currentPhrase []string

	for _, tok := range tokens {
		if isChineseText(tok.Text) {
			switch tok.Tag {
			case "VB", "RB", "MD": // Verbs, Adverbs, Modals
				currentPhrase = append(currentPhrase, tok.Text)
			default:
				if len(currentPhrase) > 0 {
					verbPhrases = append(verbPhrases, strings.Join(currentPhrase, " "))
					currentPhrase = nil
				}
			}
		}
	}

	if len(currentPhrase) > 0 {
		verbPhrases = append(verbPhrases, strings.Join(currentPhrase, " "))
	}

	return verbPhrases
}

// Categorizes text into linguistic categories, focusing exclusively on Chinese content
func categorizeChineseText(inputFile string, outputDir string) error {
	file, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("failed to open input file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var content string
	for scanner.Scan() {
		content += scanner.Text() + " "
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading input file: %v", err)
	}

	doc, err := prose.NewDocument(content)
	if err != nil {
		return fmt.Errorf("error creating Prose document: %v", err)
	}

	categoryFiles := map[string]string{
		"ChineseCharacters":       "ChineseCharacters.txt",
		"ChineseAdjectives":       "ChineseAdjectives.txt",
		"ChineseAdverbs":          "ChineseAdverbs.txt",
		"ChineseCommonPhrases":    "ChineseCommonPhrases.txt",
		"ChineseIdioms":           "ChineseIdioms.txt",
		"ChineseNouns":            "ChineseNouns.txt",
		"ChineseNounPhrases":      "ChineseNounPhrases.txt",
		"ChineseSlang":            "ChineseSlang.txt",
		"ChineseVerbPhrases":      "ChineseVerbPhrases.txt",
		"ChineseVerbs":            "ChineseVerbs.txt",
		"ChineseOtherExpressions": "ChineseOtherExpressions.txt",
	}

	idioms := []string{"井底之蛙", "守株待兔", "画蛇添足", "纸上谈兵"}
	slang := []string{"吃土", "学霸", "宅男", "高富帅"}

	results := make(map[string][]string)

	// Extracting and categorizing tokens
	for _, tok := range doc.Tokens() {
		text := tok.Text
		if isChineseText(text) {
			// Extract individual characters
			results["ChineseCharacters"] = append(results["ChineseCharacters"], extractChineseCharacters(text)...)

			switch tok.Tag {
			case "NN":
				results["ChineseNouns"] = append(results["ChineseNouns"], text)
			case "VB":
				results["ChineseVerbs"] = append(results["ChineseVerbs"], text)
			case "JJ":
				results["ChineseAdjectives"] = append(results["ChineseAdjectives"], text)
			case "RB":
				results["ChineseAdverbs"] = append(results["ChineseAdverbs"], text)
			default:
				results["ChineseOtherExpressions"] = append(results["ChineseOtherExpressions"], text)
			}
			if matchesPhraseList(text, idioms) {
				results["ChineseIdioms"] = append(results["ChineseIdioms"], text)
			}
			if matchesPhraseList(text, slang) {
				results["ChineseSlang"] = append(results["ChineseSlang"], text)
			}
		}
	}

	// Extract phrases
	results["ChineseNounPhrases"] = extractNounPhrases(doc.Tokens())
	results["ChineseVerbPhrases"] = extractVerbPhrases(doc.Tokens())

	// Output results
	for category, filename := range categoryFiles {
		filePath := filepath.Join(outputDir, filename)
		file, err := os.Create(filePath)
		if err != nil {
			return fmt.Errorf("failed to create output file for %s: %v", category, err)
		}
		defer file.Close()

		writer := bufio.NewWriter(file)
		countedContent := countFrequencies(results[category])
		sortedContent := sortByFrequency(countedContent)
		for _, item := range sortedContent {
			writer.WriteString(item + "\n")
		}
		writer.Flush()
	}

	return nil
}

func matchesPhraseList(phrase string, list []string) bool {
	for _, item := range list {
		if strings.EqualFold(item, phrase) {
			return true
		}
	}
	return false
}

func main() {
	fmt.Println("Select the input text file:")
	inputFile, err := dialog.File().Title("Select Input File").Filter("Text Files (*.txt)", "txt").Load()
	if err != nil || inputFile == "" {
		fmt.Println("No file selected or error occurred:", err)
		return
	}

	fmt.Println("Select the output directory:")
	outputDir, err := dialog.Directory().Title("Select Output Directory").Browse()
	if err != nil || outputDir == "" {
		fmt.Println("No directory selected or error occurred:", err)
		return
	}

	err = categorizeChineseText(inputFile, outputDir)
	if err != nil {
		fmt.Println("Error during categorization:", err)
		return
	}

	fmt.Println("Chinese content has been categorized and written to output files.")
}
