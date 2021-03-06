package gqlschema

import (
	"bytes"
	"html/template"
	"log"
	"strings"

	"github.com/ahmedalhulaibi/substance/substancegen"
)

var graphqlDataTypes map[string]string

func init() {
	graphqlDataTypes = make(map[string]string)
	graphqlDataTypes["int"] = "Int"
	graphqlDataTypes["int8"] = "Int"
	graphqlDataTypes["int16"] = "Int"
	graphqlDataTypes["int32"] = "Int"
	graphqlDataTypes["int64"] = "Int"
	graphqlDataTypes["uint"] = "Int"
	graphqlDataTypes["uint8"] = "Int"
	graphqlDataTypes["uint16"] = "Int"
	graphqlDataTypes["uint32"] = "Int"
	graphqlDataTypes["uint64"] = "Int"
	graphqlDataTypes["byte"] = "Int"
	graphqlDataTypes["rune"] = "Int"
	graphqlDataTypes["bool"] = "Boolean"
	graphqlDataTypes["string"] = "String"
	graphqlDataTypes["float32"] = "Float"
	graphqlDataTypes["float64"] = "Float"
}

/*OutputGraphqlSchema Returns a buffer containing a GraphQL schema in the standard GraphQL schema syntax*/
func OutputGraphqlSchema(gqlObjectTypes map[string]substancegen.GenObjectType) bytes.Buffer {
	var buff bytes.Buffer
	GenerateGraphqlSchemaTypes(gqlObjectTypes, &buff)

	return buff
}

/*GenerateGraphqlSchemaTypes generates graphql types in graphql sstandard syntax*/
func GenerateGraphqlSchemaTypes(gqlObjectTypes map[string]substancegen.GenObjectType, buff *bytes.Buffer) {
	for _, object := range gqlObjectTypes {
		for _, propVal := range object.Properties {
			if propVal.IsObjectType {
				propVal.AltScalarType["gqlschema"] = propVal.ScalarType
			} else {
				propVal.AltScalarType["gqlschema"] = graphqlDataTypes[propVal.ScalarType]
			}
		}
	}

	tmpl := template.New("graphqlSchemaTypes")
	tmpl, err := tmpl.Parse(graphqlSchemaTypesTemplate)
	if err != nil {
		log.Fatal("Parse: ", err)
		return
	}
	//print schema
	err1 := tmpl.Execute(buff, gqlObjectTypes)
	if err1 != nil {
		log.Fatal("Execute: ", err1)
		return
	}
}

/*GenerateGraphqlSchemaInputTypes generates graphql input types in graphql sstandard syntax*/
func GenerateGraphqlSchemaInputTypes(gqlObjectTypes map[string]substancegen.GenObjectType, buff *bytes.Buffer) {
	for _, object := range gqlObjectTypes {
		for _, propVal := range object.Properties {
			if propVal.IsObjectType {
				propVal.AltScalarType["gqlschema"] = propVal.ScalarType
			} else {
				propVal.AltScalarType["gqlschema"] = graphqlDataTypes[propVal.ScalarType]
			}
		}
	}

	tmpl := template.New("graphqlSchemaInputTypes")
	tmpl, err := tmpl.Parse(graphqlSchemaInputTypesTemplate)
	if err != nil {
		log.Fatal("Parse: ", err)
		return
	}
	//print schema
	err1 := tmpl.Execute(buff, gqlObjectTypes)
	if err1 != nil {
		log.Fatal("Execute: ", err1)
		return
	}
}

/*GenerateGraphqlQueries generates graphql queries and mutations in graphql standard syntax*/
func GenerateGraphqlQueries(gqlObjectTypes map[string]substancegen.GenObjectType, buff *bytes.Buffer) {
	buff.WriteString("type Query {")
	GenerateGraphqlGetQueries(gqlObjectTypes, buff)
	buff.WriteString("}")
}

/*GenerateGraphqlGetQueries generates graphql GET queries in graphql standard syntax*/
func GenerateGraphqlGetQueries(gqlObjectTypes map[string]substancegen.GenObjectType, buff *bytes.Buffer) {
	for _, object := range gqlObjectTypes {
		for _, propVal := range object.Properties {
			if propVal.IsObjectType {
				propVal.AltScalarType["gqlschema"] = propVal.ScalarType
			} else {
				propVal.AltScalarType["gqlschema"] = graphqlDataTypes[propVal.ScalarType]
			}
		}
	}

	tmpl := template.New("graphqlSchemaGet")
	tmpl, err := tmpl.Parse(graphqlSchemaGetQueriesTemplate)
	if err != nil {
		log.Fatal("Parse: ", err)
		return
	}
	//print schema
	err1 := tmpl.Execute(buff, gqlObjectTypes)
	var tmpbuff bytes.Buffer
	tmpbuff.WriteString(strings.Replace(buff.String(), ", )", ")", 1))
	buff.Reset()
	*buff = tmpbuff
	//fmt.Println(fixedstr)
	if err1 != nil {
		log.Fatal("Execute: ", err1)
		return
	}
}
