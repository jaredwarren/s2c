package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/jaredwarren/curl/swagger"
)

func main() {
	// Open our jsonFile
	jsonFile, err := os.Open("openapi.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// read our opened xmlFile as a byte array.
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
	}

	var sw swagger.Swagger

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	err = json.Unmarshal(byteValue, &sw)
	if err != nil {
		fmt.Println("  ERROR::: ", err)
	}

	// fmt.Println(reflect.TypeOf(sw))
	// fmt.Println(reflect.TypeOf(sw.Paths))
	// fmt.Println(reflect.TypeOf(sw.Paths["/color"]))
	// fmt.Println(reflect.TypeOf(sw.Paths["/color"]["get"]))

	methods := *sw.Paths["/color"].Methods
	method := methods["get"]

	// m := sw.Paths["/color"].Methods
	// fmt.Printf("%+v\n", sw.Paths)
	c, _ := method.ToCurl()
	fmt.Println("")
	fmt.Println("")
	fmt.Println(c)

	// fmt.Printf("%+v\n", sw.Paths["/color"])

	// TODO: run through Paths and generate curl

}
