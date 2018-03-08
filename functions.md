# Embedded Functions

### !(String, Any...)Any
    Resolve a template

### !=(Atom, Atom, Atom...)Bool
    Tests if the arguments are not the same

### *(Decimal, Decimal, Decimal...)Decimal
    Multiplies the arguments

### +(Decimal, Decimal, Decimal...)Decimal
    Adds the arguments

### +(String, String, String...)String
    Concat strings

### -(Decimal, Decimal, Decimal...)Decimal
    Subtracts the arguments

### .(Atom, Atom...)Any
    Access a variable in the binding

### /(Decimal, Decimal, Decimal...)Decimal
    Divides the arguments

### <(Time, Time, Time...)Bool
    Tests if the first argument is less then the following

### <(Decimal, Decimal, Decimal...)Bool
    Tests if the first argument is less then the following

### <=(Time, Time, Time...)Bool
    Tests if the first argument is less or equal then the following

### <=(Decimal, Decimal, Decimal...)Bool
    Tests if the first argument is less or equal then the following

### =(Atom, Atom, Atom...)Bool
    Tests if the arguments are the same

### >(Time, Time, Time...)Bool
    Tests if the first argument is greather then the following

### >(Decimal, Decimal, Decimal...)Bool
    Tests if the first argument is greather then the following

### >=(Decimal, Decimal, Decimal...)Bool
    Tests if the first argument is greather or equal then the following

### >=(Time, Time, Time...)Bool
    Tests if the first argument is greather or equal then the following

### between(Time, Time, Time, Time...)Bool
    Tests if the arguments are between the second last and the last argument

### between(Decimal, Decimal, Decimal, Decimal...)Bool
    Tests if the arguments are between the second last and the last argument

### ceil(Decimal)Decimal
    Ceil the decimal argument

### concat(String, String, String...)String
    Concat strings

### contains(String, String, String...)Bool
    Returns wether the first argument exists in the following arguments

### drop(List)List
    Create a list containing all but the last item in the input list

### endswith(String, String, String...)Bool
    Returns wether the first argument is the suffix of the following arguments

### floor(Decimal)Decimal
    Floor the decimal argument

### head(List)Any
    Returns the first item in the list

### item(List, Decimal)Any
    Returns a specific item from a list

### kv(Block...)Map
    Create a map with any key value pairs passed as arguments.

### list(Atom, Atom...)List
    Create a list out of the children

### map(List, String, Block)List
    Create a new list by evaluating the given block for each item in the input list

### mod(Decimal, Decimal, Decimal...)Decimal
    Modulo the arguments

### noop()Any
    No operation

### notcontains(String, String, String...)Bool
    Returns wether the first argument does not exist in the following arguments

### push(List, Kind(127), Kind(127)...)List
    Adds an item to the list and returns the list

### set(String, Kind(127), Kind(127)...)Null
    Set a variable in the binding

### startswith(String, String, String...)Bool
    Returns wether the first argument is the prefix of the following arguments

### tail(List)List
    Returns list without the first item

### tostring(Kind(15))String
    Converts the parameter to a string

### ~(String, String, String...)Bool
    Returns wether the first argument (regex) matches all of the following arguments

