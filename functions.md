# Embedded Functions

### !(String, Any...)Any
    Resolve a template

### !=(Any...)Bool
    Tests if the arguments are not the same

### *(Decimal...)Decimal
    Multiplies the arguments

### +(String...)String
    Concat strings

### +(Decimal, Decimal...)Decimal
    Adds the arguments

### -(Decimal...)Decimal
    Subtracts the arguments

### .(Atom...)Any
    Access a variable in the binding

### /(Decimal...)Decimal
    Divides the arguments

### <(Decimal...)Bool
    Tests if the first argument is less then the following

### <(Time...)Bool
    Tests if the first argument is less then the following

### <=(Time...)Bool
    Tests if the first argument is less or equal then the following

### <=(Decimal...)Bool
    Tests if the first argument is less or equal then the following

### =(Any...)Bool
    Tests if the arguments are the same

### >(Decimal...)Bool
    Tests if the first argument is greather then the following

### >(Time...)Bool
    Tests if the first argument is greather then the following

### >=(Decimal...)Bool
    Tests if the first argument is greather or equal then the following

### >=(Time...)Bool
    Tests if the first argument is greather or equal then the following

### between(Decimal...)Bool
    Tests if the arguments are between the second last and the last argument

### between(Time...)Bool
    Tests if the arguments are between the second last and the last argument

### ceil(Decimal)Decimal
    Ceil the decimal argument

### concat(String...)String
    Concat strings

### contains(String, String...)Bool
    Returns wether the first argument exists in the following arguments

### drop(List)List
    Create a list containing all but the last item in the input list

### endswith(String, String...)Bool
    Returns wether the first argument is the suffix of the following arguments

### floor(Decimal)Decimal
    Floor the decimal argument

### head(List)Any
    Returns the first item in the list

### item(List, Decimal)Any
    Returns a specific item from a list

### kv(Block...)Map
    Create a map with any key value pairs passed as arguments.

### list(Atom...)List
    Create a list out of the children

### mod(Decimal...)Decimal
    Modulo the arguments

### noop()Any
    No operation

### notcontains(String, String...)Bool
    Returns wether the first argument does not exist in the following arguments

### startswith(String, String...)Bool
    Returns wether the first argument is the prefix of the following arguments

### tail(List)List
    Returns list without the first item

### tostring(Kind(15))String
    Converts the parameter to a string

### ~(String, String...)Bool
    Returns wether the first argument (regex) matches all of the following arguments

