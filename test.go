package main

import (
    "fmt"
    "encoding/json"
)

var temp = `{"Name":"rjpj","Body":"hello","Time":1241241254124124}`

var b = []byte(`
    {
        "Name":"Wednesday",
        "Age":6,
        "Parents":["Gomez","Morticia"]
    }`)

type Person struct{
    Name string
    Age int64
    Parents []string
}

func main(){
    //p := Person{"rjpj", 19,[]string{"ram","pab"}}
    
    var p interface{}
    json.Unmarshal(b,&p)

    m := p.(map[string]interface{})

    for k,v := range m{
        fmt.Println(k,v)
    }

    fmt.Println(p)
}