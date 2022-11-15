/*
Copyright © 2022 Ömer AKBAŞ
*/
package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "bbb-downloader",
	Short: "BigBlueButton downloader.",
	Long:  `It is a tool used to download content from Bigbluebutton.`,
	Run: func(cmd *cobra.Command, args []string) {
		createNew()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func promptGetInput(pc promptContent) string {
	validate := func(input string) error {
		if len(input) <= 0 {
			return errors.New(pc.errorMsg)
		}
		return nil
	}

	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }} ",
		Valid:   "{{ . | green }} ",
		Invalid: "{{ . | red }} ",
		Success: "{{ . | bold }} ",
	}

	prompt := promptui.Prompt{
		Label:     pc.label,
		Templates: templates,
		Validate:  validate,
	}

	result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}

	return result
}

func createNew() {
	folderPromptContent := promptContent{
		"Please specify the folder path.",
		"Your folder road?",
	}
	folder := promptGetInput(folderPromptContent)

	linkPromptContent := promptContent{
		"Please specify the url address of the course video.",
		"Url address?",
	}
	link := promptGetInput(linkPromptContent)

	err := bbbContentProcess(folder, link)
	if err != nil {
		fmt.Println("start err :", err.Error())
		os.Exit(1)
	}
	fmt.Println("Downloads completed.")
}
