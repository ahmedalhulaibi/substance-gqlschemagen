package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/ahmedalhulaibi/substance"
	"github.com/ahmedalhulaibi/substance-gqlschemagen/gqlschema"
	"github.com/ahmedalhulaibi/substance/substancegen"
)

func main() {
	helpText := `Usage: 
		gengqlschema -db="dbtype" -cnstr="connection:string@locahost:9999 -file="schema.graphql"
			or
		gengqlschema -jsonsrc="path-to-substance-objects.json -file="schema.graphql""
			or
		gengqlschema -file="schema.graphql"	<-- Defaults to -jsonsrc=substance-objects.json
			or
		gengqlschema							<-- Default -jsonsrc=substance-objects.json and file=schema.graphql
`
	dbtype := flag.String("db", "", "Database driver name.\nSupported databases types:\n\t- mysql\n\t- postgres \n\t- sqlite3\n")
	connString := flag.String("cnstr", "", "Connection string to connect to database.")
	jsonSourceFilePath := flag.String("jsonsrc", "substance-objects.json", "JSON substance-objects.json file describing the database objects. This can be used as an alternative to providing connection info.")
	outputSrcFilePath := flag.String("file", "schema.graphql", "File to output graphql schema file. Defaults to schema.graphql . If file=\"\" outputs to stdout.")
	flag.Parse()

	var objects map[string]substancegen.GenObjectType

	if jsonSourceFilePath != nil {
		jsonFile, err := os.Open(*jsonSourceFilePath)
		if err != nil {
			fmt.Printf(helpText)
			log.Panicf(err.Error())
		}
		log.Printf("Opened %s successfully", *jsonSourceFilePath)
		byteVal, _ := ioutil.ReadAll(jsonFile)
		log.Printf("Read %s successfully", *jsonSourceFilePath)
		json.Unmarshal(byteVal, &objects)
		log.Printf("Unmarshalled %s successfully", *jsonSourceFilePath)
	} else if dbtype != nil && connString != nil {
		results, err := substance.DescribeDatabase(*dbtype, *connString)
		if err != nil {
			fmt.Printf(helpText)
			log.Panicf(err.Error())
		}
		if len(results) > 0 {
			log.Println("Database: ", results[0].DatabaseName)
		}
		var tables []string
		for _, result := range results {
			log.Printf("Table: %s\n", result.TableName)
			tables = append(tables, result.TableName)
		}
		log.Println("=====================")

		objects = substancegen.GetObjectTypesFunc(*dbtype, *connString, tables)
	}
	if objects != nil {
		log.Println("printing objects")
		log.Println(objects)

		var outputBuff1 bytes.Buffer
		gqlschema.GenerateGraphqlSchemaTypes(objects, &outputBuff1)

		var outputBuff2 bytes.Buffer
		gqlschema.GenerateGraphqlSchemaInputTypes(objects, &outputBuff2)

		var outputBuff3 bytes.Buffer
		gqlschema.GenerateGraphqlQueries(objects, &outputBuff3)

		outputBuff := append(outputBuff1.Bytes(), append(outputBuff2.Bytes(), outputBuff3.Bytes()...)...)
		err := ioutil.WriteFile(*outputSrcFilePath, outputBuff, 0664)
		if err != nil {
			fmt.Printf(helpText)
			fmt.Printf(outputBuff1.String())
		}
	}
}
