/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/Masamerc/prompt/data"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "create a new note entry",
	Long:  `create a new note entry`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("new called")
		createNewNote()
	},
}

type promptContent struct {
	errorMsg string
	label    string
}

func init() {
	noteCmd.AddCommand(newCmd)
}

func promptGetInput(p promptContent) string {
	validate := func(input string) error {
		if len(input) <= 0 {
			return errors.New(p.errorMsg)
		}
		return nil
	}

	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }}",
		Valid:   "{{ . | green }}",
		Invalid: "{{ . | red }}",
		Success: "{{ . | bold }}",
	}

	prompt := promptui.Prompt{
		Label:     p.label,
		Templates: templates,
		Validate:  validate,
	}

	result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed! %v\n", err)
		os.Exit(1)
	}

	return result
}

func promptSelectCategory(p promptContent) string {
	options := []string{"animal", "human", "food", "object"}
	index := -1

	var result string
	var err error

	for index < 0 {
		prompt := promptui.SelectWithAdd{
			Label:    p.label,
			Items:    options,
			AddLabel: "Other",
		}

		index, result, err = prompt.Run()

		if index == -1 {
			options = append(options, result)
		}

	}

	if err != nil {
		log.Fatal("Prompt failed.")
		os.Exit(1)
	}

	fmt.Println(result)
	return result

}

func createNewNote() {
	wordPromptContent := promptContent{
		errorMsg: "please provide a valid word",
		label:    "what word would you like to make a note of: ",
	}

	word := promptGetInput(wordPromptContent)

	definitionPromptContent := promptContent{
		errorMsg: "please provide a valid definition",
		label:    fmt.Sprintf("what is the definition of %s: ", word),
	}

	def := promptGetInput(definitionPromptContent)

	categorySelectPromptContent := promptContent{
		errorMsg: "please select a category",
		label:    "select a category: ",
	}

	category := promptSelectCategory(categorySelectPromptContent)

	data.InsertNote(word, def, category)
	fmt.Printf("Inserted word: %s, %s - category: %s\n", word, def, category)

}
