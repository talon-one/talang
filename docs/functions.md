# Embedded Functions

### !(String, Any...)Any
Resolve a template
```lisp
(! Template1)                                                    ; executes the Template1
(! Template2 "Hello World")                                      ; executes Template2 with "Hello World" as parameter
```

### !=(Atom, Atom, Atom...)Boolean
Tests if the arguments are not the same
```lisp
(!= 1 1)                                                         ; compares decimals, returns false
(!= "Hello World" "Hello World")                                 ; compares strings, returns false
(!= true true)                                                   ; compares booleans, returns false
(!= 2006-01-02T15:04:05Z 2006-01-02T15:04:05Z)                   ; compares time, returns false
(!= 1 "1")                                                       ; returns false
(!= "Hello" "Bye")                                               ; returns true
(!= "Hello" "Hello" "Bye")                                       ; returns false
```

### *(Decimal, Decimal, Decimal...)Decimal
Multiplies the arguments
```lisp
(* 1 2)                                                          ; returns 2
(* 1 2 3)                                                        ; returns 6
```

### +(String, String, String...)String
Concat strings
```lisp
(+ "Hello" " " "World")                                          ; returns "Hello World"
(+ "Hello" " " (toString (+ 1 2)))                               ; returns "Hello 3"
```

### +(Decimal, Decimal, Decimal...)Decimal
Adds the arguments
```lisp
(+ 1 1)                                                          ; returns 2
(+ 1 2 3)                                                        ; returns 6
```

### -(Decimal, Decimal, Decimal...)Decimal
Subtracts the arguments
```lisp
(- 1 1)                                                          ; returns 0
(- 1 2 3)                                                        ; returns -4
```

### .(Atom, Atom...)Any
Access a variable in the binding
```lisp
(. Key1)                                                         ; returns the data assigned to Key1
(. Key2 SubKey1)                                                 ; returns the data assigned to SubKey1 in the Map Key2
```

### .|(Any, Token)Any
Safe read a binding
```lisp
.| boo (. List)                                                  ; returns "XJK_992" (assuming $List = "XJK_992")
.| boo (. Meh)                                                   ; returns "boo" (assuming $Meh is not set)
```

### /(Decimal, Decimal, Decimal...)Decimal
Divides the arguments
```lisp
(/ 1 2)                                                          ; returns 0.5
(/ 1 2 3)                                                        ; returns 0.166666
```

### <(Time, Time, Time...)Boolean
Tests if the first argument is less then the following
```lisp
(< 2006-01-02T15:04:05Z 2007-01-02T15:04:05Z)                    ; returns true
(< 2007-01-02T15:04:05Z 2007-01-02T15:04:05Z)                    ; returns false
(< 2008-01-02T15:04:05Z 2007-01-02T15:04:05Z)                    ; returns false
```

### <(Decimal, Decimal, Decimal...)Boolean
Tests if the first argument is less then the following
```lisp
(< 0 1)                                                          ; returns true
(< 1 1)                                                          ; returns false
(< 2 1)                                                          ; returns false
```

### <=(Decimal, Decimal, Decimal...)Boolean
Tests if the first argument is less or equal then the following
```lisp
(<= 0 1)                                                         ; returns true
(<= 1 1)                                                         ; returns true
(<= 2 1)                                                         ; returns false
```

### <=(Time, Time, Time...)Boolean
Tests if the first argument is less or equal then the following
```lisp
(<= 2006-01-02T15:04:05Z 2007-01-02T15:04:05Z)                   ; returns true
(<= 2007-01-02T15:04:05Z 2007-01-02T15:04:05Z)                   ; returns true
(<= 2008-01-02T15:04:05Z 2007-01-02T15:04:05Z)                   ; returns false
```

### =(Atom, Atom, Atom...)Boolean
Tests if the arguments are the same
```lisp
(= 1 1)                                                          ; compares decimals, returns true
(= "Hello World" "Hello World")                                  ; compares strings, returns true
(= true true)                                                    ; compares booleans, returns true
(= 2006-01-02T15:04:05Z 2006-01-02T15:04:05Z)                    ; compares time, returns true
(= 1 "1")                                                        ; returns true
(= "Hello" "Bye")                                                ; returns false
(= "Hello" "Hello" "Bye")                                        ; returns false
```

### >(Decimal, Decimal, Decimal...)Boolean
Tests if the first argument is greather then the following
```lisp
(> 0 1)                                                          ; returns false
(> 1 1)                                                          ; returns false
(> 2 1)                                                          ; returns true
```

### >(Time, Time, Time...)Boolean
Tests if the first argument is greather then the following
```lisp
(> 2006-01-02T15:04:05Z 2007-01-02T15:04:05Z)                    ; returns false
(> 2007-01-02T15:04:05Z 2007-01-02T15:04:05Z)                    ; returns false
(> 2008-01-02T15:04:05Z 2007-01-02T15:04:05Z)                    ; returns true
```

### >=(Time, Time, Time...)Boolean
Tests if the first argument is greather or equal then the following
```lisp
(>= 2006-01-02T15:04:05Z 2007-01-02T15:04:05Z)                   ; returns false
(>= 2007-01-02T15:04:05Z 2007-01-02T15:04:05Z)                   ; returns true
(>= 2008-01-02T15:04:05Z 2007-01-02T15:04:05Z)                   ; returns true
```

### >=(Decimal, Decimal, Decimal...)Boolean
Tests if the first argument is greather or equal then the following
```lisp
(>= 0 1)                                                         ; returns false
(>= 1 1)                                                         ; returns true
(>= 2 1)                                                         ; returns true
```

### addDuration(Time, Decimal, String)Time
Extract days from now from time
```lisp
(addDuration 2018-03-18T00:04:05Z 3 minutes)                     ; returns "2018-03-18T00:07:05Z"
(addDuration 2018-03-18T00:04:05Z 2 hours)                       ; returns "2018-03-18T02:04:05Z"
(addDuration 2018-03-18T00:04:05Z 18 days)                       ; returns "2018-04-05T00:04:05Z"
```

### after(Time, Time)Boolean
Checks whether time A is after B
```lisp
(after 2006-01-02T19:04:05Z 2006-01-02T15:04:05Z)               ; returns "true"
(after 2006-01-01T19:04:05Z 2006-01-02T15:04:05Z)               ; returns "false"
```

### and(Collection|Atom...)Boolean
Evaluates whether a series of predicates are all true
```lisp
(and false (> 2 1))                                              ; returns true
(and false false)                                                ; returns false
```

### append(List, Collection|Atom, Collection|Atom...)List
Adds an item to the list and returns the list
```lisp
(append (list "Hello World" "Hello Universe") "Hello Human")     ; returns a list containing "Hello World", "Hello Universe" and "Hello Human"
(append (list 1 2) 3 4)                                          ; returns a list containing 1, 2, 3 and 4
```

### before(Time, Time)Boolean
Checks whether time A is before B
```lisp
(before 2006-01-02T19:04:05Z 2006-01-02T15:04:05Z)              ; returns "false"
(before 2006-01-01T19:04:05Z 2006-01-02T15:04:05Z)              ; returns "true"
```

### between(Decimal, Decimal, Decimal, Decimal...)Boolean
Tests if the arguments are between the second last and the last argument
```lisp
(between 1 0 3)                                                  ; returns true, (1 is between 0 and 3)
(between 1 2 0 3)                                                ; returns true, (1 and 2 are between 0 and 3)
(between 0 0 2)                                                  ; returns false
(between 2 0 2)                                                  ; returns false
(between 1 4 0 3)                                                ; returns false, (1 is between 0 and 3, 4 is not)
```

### between(Time, Time, Time, Time...)Boolean
Tests if the arguments are between the second last and the last argument
```lisp
(between 2007-01-02T00:00:00Z 2006-01-02T00:00:00Z 2009-01-02T00:00:00Z)                        ; returns true, (2007-01-02T00:00:00Z is between 2006-01-02T00:00:00Z and 3)
(between 2007-01-02T00:00:00Z 2008-01-02T00:00:00Z 2006-01-02T00:00:00Z 2009-01-02T00:00:00Z)   ; returns true, (2007-01-02T00:00:00Z and 2008-01-02T00:00:00Z are between 2006-01-02T00:00:00Z and 2009-01-02T00:00:00Z)
(between 2006-01-02T00:00:00Z 2006-01-02T00:00:00Z 2008-01-02T00:00:00Z)                        ; returns false
(between 2008-01-02T00:00:00Z 2006-01-02T00:00:00Z 2008-01-02T00:00:00Z)                        ; returns false
(between 2007-01-02T00:00:00Z 2010-01-02T00:00:00Z 2006-01-02T00:00:00Z 2009-01-02T00:00:00Z)   ; returns false, (2007-01-02T00:00:00Z is between 2006-01-02T00:00:00Z and 2009-01-02T00:00:00Z, 2010-01-02T00:00:00Z is not)
```

### betweenTimes(Time, Time, Time)Boolean
Evaluates whether a timestamp is between minTime and maxTime
```lisp
(betweenTimes 2006-01-02T19:04:05Z 2006-01-01T15:04:05Z 2006-01-03T19:04:05Z)                                ; returns "false"
(betweenTimes 2006-01-01T19:04:05Z 2006-01-02T15:04:05Z 2006-01-03T19:04:05Z)                                ; returns "true"
```

### catch(Any, Any)Any
Evaluate & return the second argument. If any errors occur, return the first argument instead
```lisp
catch "Edward" (. Profile Name)                                  ; returns "Edward"
catch 22 (. Profile Age)                                         ; returns 46
catch 22 2                                                       ; returns 22
```

### ceil(Decimal)Decimal
Ceil the decimal argument
```lisp
(ceil 2)                                                         ; returns 2
(ceil 2.4)                                                       ; returns 3
(ceil 2.5)                                                       ; returns 3
(ceil 2.9)                                                       ; returns 3
(ceil -2.7)                                                      ; returns -2
(ceil -2)                                                        ; returns -2
```

### concat(String, String, String...)String
Concat strings
```lisp
(+ "Hello" " " "World")                                          ; returns "Hello World"
(+ "Hello" " " (toString (+ 1 2)))                               ; returns "Hello 3"
```

### contains(String, String, String...)Boolean
Returns wether the first argument exists in the following arguments
```lisp
(contains "Hello" "Hello World")                                 ; returns true
(contains "Hello" "World")                                       ; returns false
(contains "Hello" "Hello World" "Hello Universe")                ; returns true
(contains "World" "Hello World" "Hello Universe")                ; returns false
```

### count(List)Decimal
Return the number of items in the input list
```lisp
(count (list 1 2 3 4))                                           ; returns "4"
(count (list 1))                                                 ; returns "1"
```

### date(Time)String
Extract the date in YYYY-MM-DD format from a time.
```lisp
(betweenTimes 2006-01-02T19:04:05Z 2006-01-01T15:04:05Z 2006-01-03T19:04:05Z)                                ; returns "false"
```

### days(Time)Decimal
Extract days from now from time
```lisp
(days 2018-03-18T00:04:05Z)                                      ; returns "3.423892107645601701193527333089150488376617431640625" results vary as the function is relative to the current date.
```

### daysBetween(Time, Time)Decimal
Calculates the difference in days between 2 dates
```lisp
daysBetween 2006-01-02T19:04:05Z 2006-01-02T22:19:05Z            ; returns "0.13541666666666666"
```

### do(Collection|Atom, Token)Any
Apply a block to a value
```lisp
do (list 1 2 3) ((Item) (. Item)))                               ; returns 1 2 3
```

### do(Collection|Atom, String, Token)Any
Apply a block to a value
```lisp
do (list 1 2 3) Item (. Item))                                   ; returns 1 2 3
```

### drop(List)List
Create a list containing all but the last item in the input list
```lisp
(drop (list "Hello World" "Hello Universe"))                     ; returns a list containing "Hello World"
(drop (list 1 true Hello))                                       ; returns a list containing 1 and true
```

### endsWith(String, String, String...)Boolean
Returns wether the first argument is the suffix of the following arguments
```lisp
(endsWith "World" "Hello World")                                 ; returns true
(endsWith "World" "Hello Universe")                              ; returns false
(endsWith "World" "Hello World" "Hello Universe")                ; returns false
(endsWith "World" "Hello World" "By World")                      ; returns true
```

### every(List, String, Token)Boolean
Test if every item in a list matches a predicate
```lisp
every (. Items) ((x) (= 1 (. x Price)))                          ; returns 1 with the right binding in the scope
```

### every(List, Token)Boolean
Test if every item in a list matches a predicate
```lisp
every (. Items) ((x) (= 1 (. x Price)))                          ; returns 1 with the right binding in the scope
```

### exists(List, Token)Boolean
Test if any item in a list matches a predicate
```lisp
exists (list hello world) ((Item) (= (. Item) "hello"))          ; returns true
exists (list hello world) ((Item) (= (. Item) "hey!!"))          ; returns false
```

### exists(List, String, Token)Boolean
Test if any item in a list matches a predicate
```lisp
exists (list hello world) Item (= (. Item) "hello")              ; returns true
exists (list hello world) Item (= (. Item) "hey!!")              ; returns false
```

### filter(List, Token)List
Create a new list containing items from the input list for which the block evaluates to true
```lisp
filter (list 1 4 7 12 24 48) ((x) (> (. x) 10))                                                                                    ; returns "[12 24 48]"
filter (list "Sasquatch" "Front squats" "Caramel" "Cart items") ((x) (contains (. x) "squat"))                                     ; returns "["Sasquatch" "Front squats"]"
```

### firstName(String)String
Extract all but the last word (space-separated) from a string
```lisp
(firstName "Alex Unger")                                         ; returns "Alex"
(firstName "Mr Foo Bar")                                         ; returns "Mr"
```

### floor(Decimal)Decimal
Floor the decimal argument
```lisp
(floor 2)                                                        ; returns 2
(floor 2.4)                                                      ; returns 2
(floor 2.5)                                                      ; returns 2
(floor 2.9)                                                      ; returns 2
(floor -2.7)                                                     ; returns -3
(floor -2)                                                       ; returns -2
```

### formatTime(Time)String
Create an RFC3339 timestamp, the inverse of parseTime
```lisp
(formatTime 2018-01-02T19:04:05Z)                                ; returns "2018"
```

### head(List)Any
Returns the first item in the list
```lisp
(head (list "Hello World" "Hello Universe"))                     ; returns "Hello World"
(head (list 1 true Hello))                                       ; returns 1
```

### hour(Time)String
Extract the hour (00-23) from a time
```lisp
(hour 2018-01-14T19:04:05Z)                                      ; returns "19"
```

### isEmpty(List)Boolean
Check if a list is empty
```lisp
isEmpty (list hello world)                                       ; returns "false"
isEmpty (list)                                                   ; returns "true"
```

### item(List, Decimal)Any
Returns a specific item from a list
```lisp
(item (list "Hello World" "Hello Universe") 0)                   ; returns "Hello World"
(item (list 1 true Hello) 1)                                     ; returns true
(item (list 1 true Hello) 3)                                     ; fails
```

### join(List, String)String
Create a string by joining together a list of strings with `glue`
```lisp
(join (list hello world) "-")                                    ; returns "hello-world"
(join (list hello world) ",")                                    ; returns "hello,world"
```

### kv(Token...)Map
Create a map with any key value pairs passed as arguments.
```lisp
(kv (Key1 "Hello World") (Key2 true) (Key3 123))                 ; returns a Map with the keys key1, key2, key3
```

### lastName(String)String
Extract the last word (space-separated) from a string
```lisp
(lastName "Alex Unger")                                          ; returns "Unger"
(lastName "Mr Foo Bar")                                          ; returns "Bar"
```

### list(Atom...)List
Create a list out of the children
```lisp
(list "Hello World" "Hello Universe")                            ; returns a list with string items
(list 1 true Hello)                                              ; returns a list with an int, bool and string
```

### map(List, Token)List
Create a new list by evaluating the given block for each item in the input list
```lisp
(map (list "World" "Universe") ((x) (+ "Hello " (. x))))         ; returns a list containing "Hello World" and "Hello Universe"
```

### map(List, String, Token)List
Create a new list by evaluating the given block for each item in the input list
```lisp
(map (list "World" "Universe") x (+ "Hello " (. x)))             ; returns a list containing "Hello World" and "Hello Universe"
```

### matchTime(Time, Time, String)Boolean
Checks if two times match for a given layout
```lisp
matchTime 2018-03-11T00:04:05Z 2018-03-11T00:04:05Z YYYY-MM-DD   ; returns "true"
```

### max(List)Decimal
Find the largest number in the list
```lisp
(max  (list 3 4 1 3 7 1 17 15 2))                                ; returns 17
(max  (list 4 2 9 2 27 1 2 422))                                 ; returns 422
```

### min(List)Decimal
Find the lowest number in the list
```lisp
(min  (list 3 4 1 3 7 1 17 15 2))                                ; returns 1
(min  (list 3 4 -1 3 7 1 17 0 2))                                ; returns -1
```

### minute(Time)String
Extract the hour (00-23) from a time
```lisp
(minute 2018-01-14T19:04:05Z)                                    ; returns "04"
```

### mod(Decimal, Decimal, Decimal...)Decimal
Modulo the arguments
```lisp
(mod 1 2)                                                        ; returns 1
(mod 3 8 2)                                                      ; returns 1
```

### month(Time)String
Extract the month (1-12) from a time
```lisp
(month 2018-01-02T19:04:05Z)                                     ; returns "1"
```

### monthDay(Time)String
Extract the day (1-31) from a time
```lisp
(monthDay 2018-01-14T19:04:05Z)                                  ; returns "14"
```

### noop()Any
No operation
```lisp
(noop)
```

### not(Boolean)Boolean
Inverts the argument
```lisp
(not false)                                                      ; returns "true"
(not (not false))                                                ; returns "false"
```

### notContains(String, String, String...)Boolean
Returns wether the first argument does not exist in the following arguments
```lisp
(notContains "Hello" "Hello World")                              ; returns false
(notContains "Hello" "World")                                    ; returns true
(notContains "Hello" "Hello World" "Hello Universe")             ; returns false
(notContains "World" "Hello World" "Hello Universe")             ; returns false
```

### or(Collection|Atom...)Boolean
Evaluates whether at least one predicate is true
```lisp
(or false false false true false)                                ; returns true
(or false false)                                                 ; returns false
```

### parseTime(String, String...)Time
Evaluates whether a timestamp is between minTime and maxTime
```lisp
(parseTime "2018-01-02T19:04:05Z")                               ; returns "2018-01-02 19:04:05 +0000 UTC"
(parseTime "20:04:05Z" "HH:mm:ss")                               ; returns "2018-01-02 20:04:05 +0000 UTC"
```

### push(List, Collection|Atom, Collection|Atom...)List
Adds an item to the list and returns the list
```lisp
(push (list "Hello World" "Hello Universe") "Hello Human")       ; returns a list containing "Hello World", "Hello Universe" and "Hello Human"
(push (list 1 2) 3 4)                                            ; returns a list containing 1, 2, 3 and 4
```

### reverse(List)List
Reverses the order of items in a given list
```lisp
(reverse (list 1 2 3 4))                                         ; returns "4 3 2 1"
(reverse (list 1))                                               ; returns "1"
```

### set(String, Collection|Atom, Collection|Atom...)Null
Set a variable in the binding
```lisp
(set Key1 "Hello World")                                         ; sets Key1 to "Hello World"
(set Key2 SubKey1 true)                                          ; sets SubKey1 in map Key2 to true
```

### setTemplate(String, Token)Any
Set a template
```lisp
(setTemplate "plus(Decimal, Decimal)Decimal" (+ (# 0) (# 1)))    ; creates a template with the signature plus(Decimal, Decimal)Decimal
```

### sort(List, Boolean...)List
Sort a list ascending, set the second argument to true for descending order
```lisp
(sort  (list "World" "Universe"))                                ; returns a list containing "Universe" and "World"
(sort  (list "World" "Universe") true)                           ; returns a list containing "World" and "Universe"
```

### sortByNumber(List, Token, Boolean)List
Sort a list numerically by value
```lisp
sortByNumber (list 2 4 3 1) ((Item) (. Item)) true               ; returns [4, 3, 2, 1]
sortByNumber (list 2 4 3 1) ((Item) (. Item)) false              ; returns [1, 2, 3, 4]
```

### sortByString(List, Token, Boolean)List
Sort a list alphabetically
```lisp
sortByString (list "b" "a" "z" "t") ((Item) (. Item)) true       ; returns [a, b, t, z]
sortByString (list "b" "a" "z" "t") ((Item) (. Item)) false      ; returns [z, t, b, a]
```

### split(String, String)List
Create a list of strings by splitting the given string at each occurrence of `sep`
```lisp
(split "1,2,3,a" ",")                                            ; returns "1 2 3 a"
(split "1-2-3-a" "-")                                            ; returns "1 2 3 a"
```

### startsWith(String, String, String...)Boolean
Returns wether the first argument is the prefix of the following arguments
```lisp
(startsWith "Hello" "Hello World")                               ; returns true
(startsWith "Hello" "World")                                     ; returns false
(startsWith "Hello" "Hello World" "Hello Universe")              ; returns true
(startsWith "Hello" "Hello World" "Hell Universe")               ; returns false
```

### subDuration(Time, Decimal, String)Time
Extract days from now from time
```lisp
(subDuration 2018-03-18T00:04:05Z 12 minutes)                    ; returns "2018-03-17T23:52:05Z"
(subDuration 2018-03-18T00:04:05Z 17 hours)                      ; returns "2018-03-17T07:04:05Z"
(subDuration 2018-03-18T00:04:05Z 22 days)                       ; returns "2018-02-24T00:04:05Z"
```

### sum(List, String, Token)Decimal
Test if any item in a list matches a predicate
```lisp
sum (. List) Item (. Item Price)                                 ; returns 4 With the binding "$Items" containing prices: [2, 2]
```

### tail(List)List
Returns list without the first item
```lisp
(tail (list "Hello World" "Hello Universe"))                     ; returns a list containing "Hello Universe"
(tail (list 1 true Hello))                                       ; returns a list containing true and Hello
```

### toString(Decimal|String|Boolean|Time)String
Converts the parameter to a string
```lisp
(toString 1)                                                     ; returns "1"
(toString true)                                                  ; returns "true"
```

### weekday(Time)String
Extract the week day (0-6) from a time
```lisp
(weekDay 2018-01-14T19:04:05Z)                                   ; returns "3"
```

### year(Time)String
Extract the year from a time
```lisp
(year 2018-01-02T19:04:05Z)                                      ; returns "2018"
```

### ~(String, String, String...)Boolean
Returns wether the first argument (regex) matches all of the following arguments
```lisp
(~ "[a-z\s]*" "Hello World")                                     ; returns true
(~ "[a-z\s]*" "Hello W0rld")                                     ; returns false
(~ "[a-z\s]*" "Hello World" "Hello Universe")                    ; returns true
(~ "[a-z\s]*" "Hello W0rld" "Hello Universe")                    ; returns false
```

