# Embedded Functions

### !(String, Any...)Any
Resolve a template
```
(! Template1)                                                  // executes the Template1
(! Template2 "Hello World")                                    // executes Template2 with "Hello World" as parameter
```

### !=(Atom, Atom, Atom...)Bool
Tests if the arguments are not the same
```
(!= 1 1)                                                         // compares decimals, returns false
(!= "Hello World" "Hello World")                                 // compares strings, returns false
(!= true true)                                                   // compares booleans, returns false
(!= 2006-01-02T15:04:05Z 2006-01-02T15:04:05Z)               // compares time, returns false
(!= 1 "1")                                                       // returns false
(!= "Hello" "Bye")                                               // returns true
(!= "Hello" "Hello" "Bye")                                       // returns false
```

### *(Decimal, Decimal, Decimal...)Decimal
Multiplies the arguments
```
(* 1 2)                                                         // returns 2
(* 1 2 3)                                                       // returns 6
```

### +(Decimal, Decimal, Decimal...)Decimal
Adds the arguments
```
(+ 1 1)                                                         // returns 2
(+ 1 2 3)                                                       // returns 6
```

### +(String, String, String...)String
Concat strings
```
(+ "Hello" " " "World")                                           // returns "Hello World"
(+ "Hello" " " (toString (+ 1 2)))                                // returns "Hello 3"
```

### -(Decimal, Decimal, Decimal...)Decimal
Subtracts the arguments
```
(- 1 1)                                                         // returns 0
(- 1 2 3)                                                       // returns -4
```

### .(Atom, Atom...)Any
Access a variable in the binding
```
(. Key1)                                                       // returns the data assigned to Key1
(. Key2 SubKey1)                                               // returns the data assigned to SubKey1 in the Map Key2
```

### /(Decimal, Decimal, Decimal...)Decimal
Divides the arguments
```
(/ 1 2)                                                         // returns 0.5
(/ 1 2 3)                                                       // returns 0.166666
```

### <(Decimal, Decimal, Decimal...)Bool
Tests if the first argument is less then the following
```
(< 0 1)                                                         // returns true
(< 1 1)                                                         // returns false
(< 2 1)                                                         // returns false
```

### <(Time, Time, Time...)Bool
Tests if the first argument is less then the following
```
(< 2006-01-02T15:04:05Z 2007-01-02T15:04:05Z)               // returns true
(< 2007-01-02T15:04:05Z 2007-01-02T15:04:05Z)               // returns false
(< 2008-01-02T15:04:05Z 2007-01-02T15:04:05Z)               // returns false
```

### <=(Decimal, Decimal, Decimal...)Bool
Tests if the first argument is less or equal then the following
```
(<= 0 1)                                                        // returns true
(<= 1 1)                                                        // returns true
(<= 2 1)                                                        // returns false
```

### <=(Time, Time, Time...)Bool
Tests if the first argument is less or equal then the following
```
(<= 2006-01-02T15:04:05Z 2007-01-02T15:04:05Z)              // returns true
(<= 2007-01-02T15:04:05Z 2007-01-02T15:04:05Z)              // returns true
(<= 2008-01-02T15:04:05Z 2007-01-02T15:04:05Z)              // returns false
```

### =(Atom, Atom, Atom...)Bool
Tests if the arguments are the same
```
(= 1 1)                                                         // compares decimals, returns true
(= "Hello World" "Hello World")                                 // compares strings, returns true
(= true true)                                                   // compares booleans, returns true
(= 2006-01-02T15:04:05Z 2006-01-02T15:04:05Z)               // compares time, returns true
(= 1 "1")                                                       // returns true
(= "Hello" "Bye")                                               // returns false
(= "Hello" "Hello" "Bye")                                       // returns false
```

### >(Time, Time, Time...)Bool
Tests if the first argument is greather then the following
```
(> 2006-01-02T15:04:05Z 2007-01-02T15:04:05Z)               // returns false
(> 2007-01-02T15:04:05Z 2007-01-02T15:04:05Z)               // returns false
(> 2008-01-02T15:04:05Z 2007-01-02T15:04:05Z)               // returns true
```

### >(Decimal, Decimal, Decimal...)Bool
Tests if the first argument is greather then the following
```
(> 0 1)                                                         // returns false
(> 1 1)                                                         // returns false
(> 2 1)                                                         // returns true
```

### >=(Time, Time, Time...)Bool
Tests if the first argument is greather or equal then the following
```
(>= 2006-01-02T15:04:05Z 2007-01-02T15:04:05Z)              // returns false
(>= 2007-01-02T15:04:05Z 2007-01-02T15:04:05Z)              // returns true
(>= 2008-01-02T15:04:05Z 2007-01-02T15:04:05Z)              // returns true
```

### >=(Decimal, Decimal, Decimal...)Bool
Tests if the first argument is greather or equal then the following
```
(>= 0 1)                                                        // returns false
(>= 1 1)                                                        // returns true
(>= 2 1)                                                        // returns true
```

### after(Time, Time)Bool
Checks whether time A is after B
```
(after 2006-01-02T19:04:05Z 2006-01-02T15:04:05Z)                                // returns "true"
(after 2006-01-01T19:04:05Z 2006-01-02T15:04:05Z)                                // returns "false"
```

### append(List, Kind(127), Kind(127)...)List
Adds an item to the list and returns the list
```
(push (list "Hello World" "Hello Universe") "Hello Human")        // returns a list containing "Hello World", "Hello Universe" and "Hello Human"
(push (list 1 2) 3 4)                                             // returns a list containing 1, 2, 3 and 4
```

### before(Time, Time)Bool
Checks whether time A is before B
```
(before 2006-01-02T19:04:05Z 2006-01-02T15:04:05Z)                                // returns "false"
(before 2006-01-01T19:04:05Z 2006-01-02T15:04:05Z)                                // returns "true"
```

### between(Time, Time, Time, Time...)Bool
Tests if the arguments are between the second last and the last argument
```
(between 2007-01-02T00:00:00Z 2006-01-02T00:00:00Z 2009-01-02T00:00:00Z)                        // returns true, (2007-01-02T00:00:00Z is between 2006-01-02T00:00:00Z and 3)
(between 2007-01-02T00:00:00Z 2008-01-02T00:00:00Z 2006-01-02T00:00:00Z 2009-01-02T00:00:00Z) // returns true, (2007-01-02T00:00:00Z and 2008-01-02T00:00:00Z are between 2006-01-02T00:00:00Z and 2009-01-02T00:00:00Z)
(between 2006-01-02T00:00:00Z 2006-01-02T00:00:00Z 2008-01-02T00:00:00Z)                        // returns false
(between 2008-01-02T00:00:00Z 2006-01-02T00:00:00Z 2008-01-02T00:00:00Z)                        // returns false
(between 2007-01-02T00:00:00Z 2010-01-02T00:00:00Z 2006-01-02T00:00:00Z 2009-01-02T00:00:00Z) // returns false, (2007-01-02T00:00:00Z is between 2006-01-02T00:00:00Z and 2009-01-02T00:00:00Z, 2010-01-02T00:00:00Z is not)
```

### between(Decimal, Decimal, Decimal, Decimal...)Bool
Tests if the arguments are between the second last and the last argument
```
(between 1 0 3)                                                 // returns true, (1 is between 0 and 3)
(between 1 2 0 3)                                               // returns true, (1 and 2 are between 0 and 3)
(between 0 0 2)                                                 // returns false
(between 2 0 2)                                                 // returns false
(between 1 4 0 3)                                               // returns false, (1 is between 0 and 3, 4 is not)
```

### betweentimes(Time, Time, Time)Bool
Evaluates whether a timestamp is between minTime and maxTime
```
(betweenTimes 2006-01-02T19:04:05Z 2006-01-01T15:04:05Z 2006-01-03T19:04:05Z)                                // returns "false"
(betweenTimes 2006-01-01T19:04:05Z 2006-01-02T15:04:05Z 2006-01-03T19:04:05Z)                                // returns "true"
```

### ceil(Decimal)Decimal
Ceil the decimal argument
```
(ceil 2)                                                          // returns 2
(ceil 2.4)                                                        // returns 3
(ceil 2.5)                                                        // returns 3
(ceil 2.9)                                                        // returns 3
(ceil -2.7)                                                       // returns -2
(ceil -2)                                                         // returns -2
```

### concat(String, String, String...)String
Concat strings
```
(+ "Hello" " " "World")                                           // returns "Hello World"
(+ "Hello" " " (toString (+ 1 2)))                                // returns "Hello 3"
```

### contains(String, String, String...)Bool
Returns wether the first argument exists in the following arguments
```
(contains "Hello" "Hello World")                                  // returns true
(contains "Hello" "World")                                        // returns false
(contains "Hello" "Hello World" "Hello Universe")                 // returns true
(contains "World" "Hello World" "Hello Universe")                 // returns false
```

### count(List)Decimal
Return the number of items in the input list
```
(count (list 1 2 3 4))											// returns "4"
(count (list 1))												// returns "1"
```

### drop(List)List
Create a list containing all but the last item in the input list
```
(drop (list "Hello World" "Hello Universe"))                    // returns a list containing "Hello World"
(drop (list 1 true Hello))                                      // returns a list containing 1 and true
```

### endswith(String, String, String...)Bool
Returns wether the first argument is the suffix of the following arguments
```
(endsWith "World" "Hello World")                                   // returns true
(endsWith "World" "Hello Universe")                                // returns false
(endsWith "World" "Hello World" "Hello Universe")                  // returns false
(endsWith "World" "Hello World" "By World")                        // returns true
```

### floor(Decimal)Decimal
Floor the decimal argument
```
(floor 2)                                                         // returns 2
(floor 2.4)                                                       // returns 2
(floor 2.5)                                                       // returns 2
(floor 2.9)                                                       // returns 2
(floor -2.7)                                                      // returns -3
(floor -2)                                                        // returns -2
```

### head(List)Any
Returns the first item in the list
```
(head (list "Hello World" "Hello Universe"))                    // returns "Hello World"
(head (list 1 true Hello))                                      // returns 1
```

### item(List, Decimal)Any
Returns a specific item from a list
```
(item (list "Hello World" "Hello Universe") 0)                    // returns "Hello World"
(item (list 1 true Hello) 1)                                      // returns true
(item (list 1 true Hello) 3)                                      // fails
```

### join(List, String)String
Create a string by joining together a list of strings with `glue`
```
(join (list hello world) "-")										// returns "hello-world"
(join (list hello world) ",")										// returns "hello,world"
```

### kv(Block...)Map
Create a map with any key value pairs passed as arguments.
```
(kv (Key1 "Hello World") (Key2 true) (Key3 123))               // returns a Map with the keys key1, key2, key3
```

### list(Atom, Atom...)List
Create a list out of the children
```
(list "Hello World" "Hello Universe")                           // returns a list with string items
(list 1 true Hello)                                             // returns a list with an int, bool and string
```

### map(List, String, Block)List
Create a new list by evaluating the given block for each item in the input list
```
(map  (list "World" "Universe") x (+ "Hello " (. x)))             // returns a list containing "Hello World" and "Hello Universe"
```

### max(List)Decimal
Find the largest number in the list
```
(max  (list 3 4 1 3 7 1 17 15 2))                              // returns 17
```

### min(List)Decimal
Find the lowest number in the list
```
(min  (list 3 4 1 3 7 1 17 15 2))                              // returns 1
```

### mod(Decimal, Decimal, Decimal...)Decimal
Modulo the arguments
```
(mod 1 2)                                                         // returns 1
(mod 3 8 2)                                                       // returns 1
```

### noop()Any
No operation
```
(noop)
```

### not(Bool)Bool
Inverts the argument
```
(not false)                                                      // returns "true"
(not (not false))                                                // returns "false"
```

### notcontains(String, String, String...)Bool
Returns wether the first argument does not exist in the following arguments
```
(notContains "Hello" "Hello World")                                  // returns false
(notContains "Hello" "World")                                        // returns true
(notContains "Hello" "Hello World" "Hello Universe")                 // returns false
(notContains "World" "Hello World" "Hello Universe")                 // returns false
```

### parsetime(String, String...)Time
Evaluates whether a timestamp is between minTime and maxTime
```
(parseTime "2018-01-02T19:04:05Z")                              // returns "2018-01-02 19:04:05 +0000 UTC"
(parseTime "20:04:05Z" "HH:mm:ss")                              // returns "2018-01-02 20:04:05 +0000 UTC"
```

### push(List, Kind(127), Kind(127)...)List
Adds an item to the list and returns the list
```
(push (list "Hello World" "Hello Universe") "Hello Human")        // returns a list containing "Hello World", "Hello Universe" and "Hello Human"
(push (list 1 2) 3 4)                                             // returns a list containing 1, 2, 3 and 4
```

### reverse(List)List
Reverses the order of items in a given list
```
(reverse (list 1 2 3 4))										// returns "(4 3 2 1)"
(reverse (list 1))												// returns "(1)"
```

### set(String, Kind(127), Kind(127)...)Null
Set a variable in the binding
```
(set Key1 "Hello World")                                       // sets Key1 to "Hello World"
(set Key2 SubKey1 true)                                        // sets SubKey1 in map Key2 to true
```

### sort(List, Bool...)List
Sort a list ascending, set the second argument to true for descending order
```
(sort  (list "World" "Universe"))                                 // returns a list containing "Universe" and "World"
(sort  (list "World" "Universe") true)                            // returns a list containing "World" and "Universe"
```

### startswith(String, String, String...)Bool
Returns wether the first argument is the prefix of the following arguments
```
(startsWith "Hello" "Hello World")                                   // returns true
(startsWith "Hello" "World")                                         // returns false
(startsWith "Hello" "Hello World" "Hello Universe")                  // returns true
(startsWith "Hello" "Hello World" "Hell Universe")                   // returns false
```

### tail(List)List
Returns list without the first item
```
(tail (list "Hello World" "Hello Universe"))                    // returns a list containing "Hello Universe"
(tail (list 1 true Hello))                                      // returns a list containing true and Hello
```

### tostring(Kind(15))String
Converts the parameter to a string
```
(toString 1)                                                      // returns "1"
(toString true)                                                   // returns "true"
```

### ~(String, String, String...)Bool
Returns wether the first argument (regex) matches all of the following arguments
```
(~ "[a-z\s]*" "Hello World")                                       // returns true
(~ "[a-z\s]*" "Hello W0rld")                                       // returns false
(~ "[a-z\s]*" "Hello World" "Hello Universe")                      // returns true
(~ "[a-z\s]*" "Hello W0rld" "Hello Universe")                      // returns false
```

