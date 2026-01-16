package main

import (
"fmt"
)

func main() {
gen := NewSpecificationGenerator("test")

patterns := []string{
".{0,128}",
".{1,}",
".{32}",
"[0-9A-Fa-f]{32}",
"[0-9A-Fa-f]{512}",
"^#(?:[0-9a-fA-F]{3}){1,2}$",
"wan|wan2|lan",
"^[a-z]+$",
}

for _, pattern := range patterns {
validators := gen.buildStringValidators(pattern)
fmt.Printf("Pattern: %-35s → ", pattern)
if len(validators) == 0 {
fmt.Println("(no validators)")
} else {
for i, v := range validators {
if i > 0 {
fmt.Print(", ")
}
if v.Custom != nil {
fmt.Print(v.Custom.SchemaDefinition)
}
}
fmt.Println()
}
}
}
