package main

import (
	"os"
	"path"

	"github.com/mxcd/avatar-generator/pkg/font"
	"github.com/mxcd/avatar-generator/pkg/initials"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		panic("No command provided: expected 'single' or 'all'")
	}

	if args[0] == "single" {
		if len(args) != 3 {
			panic("Expected 2 arguments: 'single', initials and the output file")
		}
		initials := args[1]
		outputFile := args[2]
		generateSingle(initials, outputFile)

	} else if args[0] == "all" {
		if len(args) != 2 {
			panic("Expected 2 arguments: 'all' and the output directory")
		}
		outputDirectory := args[1]
		generateAll(outputDirectory)
	} else {
		panic("Invalid command: expected 'single' or 'all'")
	}
}

func generateAll(outputDirectory string) {
	generator, err := initials.NewInitialsAvatarGenerator(&initials.InitialsAvatarOptions{
		AvaratWidth:  256,
		AvatarHeight: 256,
		FontSize:     160,
		FontPath:     "fonts/theboldfont.ttf",
		FontFS:       &font.Fonts,
	})

	if err != nil {
		panic(err)
	}

	err = os.MkdirAll(outputDirectory, 0755)
	if err != nil {
		panic(err)
	}

	initialsList := initials.GetAllInitials()
	for _, initialsString := range initialsList {
		data, err := generator.GenerateInitialsAvatar(initialsString)
		if err != nil {
			panic(err)
		}
		outputPath := path.Join(outputDirectory, initialsString+".png")
		err = os.WriteFile(outputPath, data, 0644)
		if err != nil {
			panic(err)
		}
	}
}

func generateSingle(initialsString string, outputFile string) {
	generator, err := initials.NewInitialsAvatarGenerator(&initials.InitialsAvatarOptions{
		AvaratWidth:  256,
		AvatarHeight: 256,
		FontSize:     160,
		FontPath:     "fonts/theboldfont.ttf",
		FontFS:       &font.Fonts,
	})

	if err != nil {
		panic(err)
	}

	data, err := generator.GenerateInitialsAvatar(initialsString)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(outputFile, data, 0644)
	if err != nil {
		panic(err)
	}
}
