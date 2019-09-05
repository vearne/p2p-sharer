## randomchoice
randomly select M elements from a slice(which has N elements)


### Install:
```
go get -u github.com/vearne/randomchoice
```

### Import:
```
import "github.com/vearne/randomchoice"
```

### Quick Start
```
package main

import (
    rc "github.com/vearne/randomchoice"
    "fmt"
)

type Car struct{
    color string
    name  string
}

func (c *Car) String() string{
    return c.color + "-" + c.name
}

func main(){
    // example 1
    var children []string
    children = []string{"lily", "rose","lisa"}
    // randomly select 2 kids from children
    idxSlice := rc.RandomChoice(len(children), 2)
    result := make([]string, 0, 2)
    for _, v := range idxSlice{
        result = append(result, children[v])
    }
    // result: selected kids
    fmt.Println(result)

    // example 2
    var carSlice []*Car= make([]*Car, 0, 3)
    bmw := Car{color:"black", name:"bmw"}
    carSlice = append(carSlice, &bmw)
    buick := Car{color:"silvery", name:"buick"}
    carSlice = append(carSlice, &buick )
    skoda := Car{color:"white", name:"skoda"}
    carSlice = append(carSlice, &skoda )
    // random select 2 kinds from carSlice
    idxSlice = rc.RandomChoice(len(carSlice), 2)
    cars := make([]*Car, 0, 2)
    for _, v := range idxSlice{
        cars = append(cars, carSlice[v])
    }
    fmt.Println(cars)
}
```

### Performance
`CPU Model Name`: 2.3 GHz Intel Core i5    
`CPU Processors`: 4    
`Memory`: 8GB    

### Test Results
|N|M|ns/op|
|:---|:---|:---|
|10|3|135|
|100|3|267|
|1000|3|1502|
|10|5|197|
|100|5|336|
|1000|5|1582|
