# Optional
A GO's "generics" version of `Optional`s, pretty much borrowed from [Java's JDK11](https://github.com/AdoptOpenJDK/openjdk-jdk11/blob/master/src/java.base/share/classes/java/util/Optional.java), 
but probably more "idiomatic", to avoid panics at runtime.

## Usage
```go
import (
    "github.com/nebiros/optional"
)

func main() {
    var (
        v   *string
        tmp = "something"
    )
    
    v = &tmp
    
    ov := optional.OfNullable(v)

    err := doSomething(ov)
    if err != nil {
        panic(err)
    }
}

func doSomething(s Optional[string]) error {
    if !s.IsPresent() {
        return fmt.Errorf("v not present")
    }
    
    sv, err := s.Get()
    if err != nil {
        return err
    }
    
    fmt.Println("sv: " + sv)
}
```

