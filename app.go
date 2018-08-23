package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/jaredwarren/curl/swagger"
	"github.com/spf13/cobra"
)

func main() {
	// var echoTimes int

	// 	var cmdPrint = &cobra.Command{
	// 		Use:   "print [string to print]",
	// 		Short: "Print anything to the screen",
	// 		Long: `print is for printing anything back to the screen.
	// For many years people have printed back to the screen.`,
	// 		Args: cobra.MinimumNArgs(1),
	// 		Run: func(cmd *cobra.Command, args []string) {
	// 			fmt.Println("Print: " + strings.Join(args, " "))
	// 		},
	// 	}

	// 	var cmdEcho = &cobra.Command{
	// 		Use:   "echo [string to echo]",
	// 		Short: "Echo anything to the screen",
	// 		Long: `echo is for echoing anything back.
	// Echo works a lot like print, except it has a child command.`,
	// 		Args: cobra.MinimumNArgs(1),
	// 		Run: func(cmd *cobra.Command, args []string) {
	// 			fmt.Println("Print: " + strings.Join(args, " "))
	// 		},
	// 	}

	// 	var cmdTimes = &cobra.Command{
	// 		Use:   "times [# times] [string to echo]",
	// 		Short: "Echo anything to the screen more times",
	// 		Long: `echo things multiple times back to the user by providing
	// a count and a string.`,
	// 		Args: cobra.MinimumNArgs(1),
	// 		Run: func(cmd *cobra.Command, args []string) {
	// 			for i := 0; i < echoTimes; i++ {
	// 				fmt.Println("Echo: " + strings.Join(args, " "))
	// 			}
	// 		},
	// 	}

	// 	cmdTimes.Flags().IntVarP(&echoTimes, "times", "t", 1, "times to echo the input")

	var rootCmd = &cobra.Command{
		Use: "file [path] [method]",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			if len(args) < 1 {
				return errors.New("File Missing")
			}
			fmt.Printf("root: %+v\n", args)
			path := ""
			if len(args) > 1 {
				path = args[1]
			}
			method := ""
			if len(args) > 2 {
				method = args[2]
			}

			// Open our jsonFile
			jsonFile, err := os.Open(args[0])
			// if we os.Open returns an error then handle it
			if err != nil {
				return
			}
			// defer the closing of our jsonFile so that we can parse it later on
			defer jsonFile.Close()

			// read our opened xmlFile as a byte array.
			byteValue, err := ioutil.ReadAll(jsonFile)
			if err != nil {
				return
			}

			var sw swagger.Swagger

			// we unmarshal our byteArray which contains our
			// jsonFile's content into 'users' which we defined above
			err = json.Unmarshal(byteValue, &sw)

			err = swagger.SwagToCurl(sw, path, method)
			return
		},
	}
	// rootCmd.AddCommand(cmdPrint, cmdEcho)
	// cmdEcho.AddCommand(cmdTimes)
	rootCmd.Execute()
}
