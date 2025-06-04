package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func SelectDC(recommended string, uniqueDCs map[string]struct{}) string {
	reader := bufio.NewReader(os.Stdin)
	caser := cases.Title(language.English)

	for {
		fmt.Println("\033[1;33mSelect DC_LOCATION:\033[0m")
		i := 1
		options := make(map[int]string)
		defaultChoice := 0
		for dc := range uniqueDCs {
			parts := strings.Split(dc, ":")
			label := caser.String(parts[1])
			if dc == recommended {
				label = fmt.Sprintf("\033[0;32m%s (recommended, matches hostname: %s)\033[0m", label, parts[0])
				defaultChoice = i
			}
			fmt.Printf("  %d) %s\n", i, label)
			options[i] = dc
			i++
		}

		prompt := "Enter the number of your choice"
		if defaultChoice > 0 {
			prompt += fmt.Sprintf(" [\033[0;32m%d\033[0m]", defaultChoice)
		}
		prompt += ": "

		fmt.Print(prompt)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "" && defaultChoice > 0 {
			input = fmt.Sprintf("%d", defaultChoice)
		}

		var choice int
		_, err := fmt.Sscanf(input, "%d", &choice)
		if err != nil || choice < 1 || choice >= i {
			fmt.Println("\033[0;31mInvalid selection. Please enter a valid number.\033[0m")
			continue
		}

		selected := options[choice]
		key := strings.Split(selected, ":")[0]
		fmt.Printf("\033[1;33mTo confirm, type the location label -- \033[0;31m%s\033[1;33m: \033[0m", key)
		confirm, _ := reader.ReadString('\n')
		confirm = strings.TrimSpace(confirm)

		if strings.EqualFold(confirm, key) {
			return key
		}
		fmt.Println("\033[0;31mInput does not match. Please try again.\033[0m")
	}
}
