# Embedded Functions

### !(String, Any...)Any
Resolve a template
```
```

### !=(Atom, Atom, Atom...)Bool
Tests if the arguments are not the same
```(!= 1 1)                                                         // compares decimals, returns false
(!= "Hello World" "Hello World")                                 // compares strings, returns false
(!= true true)                                                   // compares booleans, returns false
(!= "2006-01-02T15:04:05Z" "2006-01-02T15:04:05Z")               // compares time, returns false
(!= 1 "1")                                                       // returns false
(!= "Hello" "Bye")                                               // returns true
(!= "Hello" "Hello" "Bye")                                       // returns false
```

### *(Decimal, Decimal, Decimal...)Decimal
Multiplies the arguments
```
```

### +(String, String, String...)String
Concat strings
```
```

### +(Decimal, Decimal, Decimal...)Decimal
Adds the arguments
```
```

### -(Decimal, Decimal, Decimal...)Decimal
Subtracts the arguments
```
```

### .(Atom, Atom...)Any
Access a variable in the binding
```
```

### /(Decimal, Decimal, Decimal...)Decimal
Divides the arguments
```
```

### <(Decimal, Decimal, Decimal...)Bool
Tests if the first argument is less then the following
```(< 0 1)                                                         // returns true
(< 1 1)                                                         // returns false
(< 2 1)                                                         // returns false
```

### <(Time, Time, Time...)Bool
Tests if the first argument is less then the following
```(< "2006-01-02T15:04:05Z" "2007-01-02T15:04:05Z")               // returns true
(< "2007-01-02T15:04:05Z" "2007-01-02T15:04:05Z")               // returns false
(< "2008-01-02T15:04:05Z" "2007-01-02T15:04:05Z")               // returns false
```

### <=(Decimal, Decimal, Decimal...)Bool
Tests if the first argument is less or equal then the following
```(<= 0 1)                                                        // returns true
(<= 1 1)                                                        // returns true
(<= 2 1)                                                        // returns false
```

### <=(Time, Time, Time...)Bool
Tests if the first argument is less or equal then the following
```(<= "2006-01-02T15:04:05Z" "2007-01-02T15:04:05Z")              // returns true
(<= "2007-01-02T15:04:05Z" "2007-01-02T15:04:05Z")              // returns true
(<= "2008-01-02T15:04:05Z" "2007-01-02T15:04:05Z")              // returns false
```

### =(Atom, Atom, Atom...)Bool
Tests if the arguments are the same
```(= 1 1)                                                         // compares decimals, returns true
(= "Hello World" "Hello World")                                 // compares strings, returns true
(= true true)                                                   // compares booleans, returns true
(= "2006-01-02T15:04:05Z" "2006-01-02T15:04:05Z") // compares time, returns true
(= 1 "1")                                                       // returns true
(= "Hello" "Bye")                                               // returns false
(= "Hello" "Hello" "Bye")                                       // returns false
```

### >(Decimal, Decimal, Decimal...)Bool
Tests if the first argument is greather then the following
```(> 0 1)                                                         // returns false
(> 1 1)                                                         // returns false
(> 2 1)                                                         // returns true
```

### >(Time, Time, Time...)Bool
Tests if the first argument is greather then the following
```(> "2006-01-02T15:04:05Z" "2007-01-02T15:04:05Z")               // returns false
(> "2007-01-02T15:04:05Z" "2007-01-02T15:04:05Z")               // returns false
(> "2008-01-02T15:04:05Z" "2007-01-02T15:04:05Z")               // returns true
```

### >=(Time, Time, Time...)Bool
Tests if the first argument is greather or equal then the following
```(>= "2006-01-02T15:04:05Z" "2007-01-02T15:04:05Z")              // returns false
(>= "2007-01-02T15:04:05Z" "2007-01-02T15:04:05Z")              // returns true
(>= "2008-01-02T15:04:05Z" "2007-01-02T15:04:05Z")              // returns true
```

### >=(Decimal, Decimal, Decimal...)Bool
Tests if the first argument is greather or equal then the following
```(>= 0 1)                                                        // returns false
(>= 1 1)                                                        // returns true
(>= 2 1)                                                        // returns true
```

### between(Time, Time, Time, Time...)Bool
Tests if the arguments are between the second last and the last argument
```(between "2007-01-02T00:00:00Z" "2006-01-02T00:00:00Z" "2009-01-02T00:00:00Z")                        // returns true, (2007-01-02T00:00:00Z is between 2006-01-02T00:00:00Z and 3)
(between "2007-01-02T00:00:00Z" "2008-01-02T00:00:00Z" "2006-01-02T00:00:00Z" "2009-01-02T00:00:00Z") // returns true, (2007-01-02T00:00:00Z and 2008-01-02T00:00:00Z are between 2006-01-02T00:00:00Z and 2009-01-02T00:00:00Z)
(between "2006-01-02T00:00:00Z" "2006-01-02T00:00:00Z" "2008-01-02T00:00:00Z")                        // returns false
(between "2008-01-02T00:00:00Z" "2006-01-02T00:00:00Z" "2008-01-02T00:00:00Z")                        // returns false
(between "2007-01-02T00:00:00Z" "2010-01-02T00:00:00Z" "2006-01-02T00:00:00Z" "2009-01-02T00:00:00Z") // returns false, (2007-01-02T00:00:00Z is between 2006-01-02T00:00:00Z and 2009-01-02T00:00:00Z, 2010-01-02T00:00:00Z is not)
```

### between(Decimal, Decimal, Decimal, Decimal...)Bool
Tests if the arguments are between the second last and the last argument
```(between 1 0 3)                                                 // returns true, (1 is between 0 and 3)
(between 1 2 0 3)                                               // returns true, (1 and 2 are between 0 and 3)
(between 0 0 2)                                                 // returns false
(between 2 0 2)                                                 // returns false
(between 1 4 0 3)                                               // returns false, (1 is between 0 and 3, 4 is not)
```

### ceil(Decimal)Decimal
Ceil the decimal argument
```
```

### concat(String, String, String...)String
Concat strings
```
```

### contains(String, String, String...)Bool
Returns wether the first argument exists in the following arguments
```
```

### drop(List)List
Create a list containing all but the last item in the input list
```
```

### endswith(String, String, String...)Bool
Returns wether the first argument is the suffix of the following arguments
```
```

### floor(Decimal)Decimal
Floor the decimal argument
```
```

### head(List)Any
Returns the first item in the list
```
```

### item(List, Decimal)Any
Returns a specific item from a list
```
```

### kv(Block...)Map
Create a map with any key value pairs passed as arguments.
```
```

### list(Atom, Atom...)List
Create a list out of the children
```
```

### map(List, String, Block)List
Create a new list by evaluating the given block for each item in the input list
```
```

### max(List)Decimal
Find the largest number in the list
```
```

### min(List)Decimal
Find the lowest number in the list
```
```

### mod(Decimal, Decimal, Decimal...)Decimal
Modulo the arguments
```
```

### noop()Any
No operation
```
```

### notcontains(String, String, String...)Bool
Returns wether the first argument does not exist in the following arguments
```
```

### push(List, Kind(127), Kind(127)...)List
Adds an item to the list and returns the list
```
```

### set(String, Kind(127), Kind(127)...)Null
Set a variable in the binding
```
```

### sort(List, Bool...)List
Sort a list ascending, set the second argument to true for descending order
```
```

### startswith(String, String, String...)Bool
Returns wether the first argument is the prefix of the following arguments
```
```

### tail(List)List
Returns list without the first item
```
```

### tostring(Kind(15))String
Converts the parameter to a string
```
```

### ~(String, String, String...)Bool
Returns wether the first argument (regex) matches all of the following arguments
```
```

