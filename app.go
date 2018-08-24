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
	var rootCmd = &cobra.Command{
		Use: "file [path] [method]",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			if len(args) < 1 {
				return errors.New("File Missing")
			}
			path := ""
			if len(args) > 1 {
				path = args[1]
			}
			method := ""
			if len(args) > 2 {
				method = args[2]
			}

			jsonFile, err := os.Open(args[0])
			if err != nil {
				return err
			}
			defer jsonFile.Close()

			byteValue, err := ioutil.ReadAll(jsonFile)
			if err != nil {
				return err
			}

			var sw swagger.Swagger

			err = json.Unmarshal(byteValue, &sw)
			if err != nil {
				return err
			}

			if path != "" {
				p := sw.FindPath(path)
				if p == nil {
					fmt.Printf("path \"%s\" not found\n", path)
				} else {
					if method != "" {
						m := p.FindMethod(method)
						if m == nil {
							fmt.Printf("method \"%s\" not found\n", method)
						} else {
							fmt.Println(m.ToCurl(sw.Host))
						}
					} else {
						fmt.Println(p.ToCurl(sw.Host))
					}
				}
			} else {
				fmt.Println(sw.ToCurl())
			}

			return
		},
	}

	rootCmd.Execute()
}
