# Embedded Functions

### !=(Any...)Bool
    Tests if the arguments are not the same

### *(Decimal...)Decimal
    Multiplies the arguments

### +(String...)String
    Concat strings

### +(Decimal...)Decimal
    Adds the arguments

### -(Decimal...)Decimal
    Substracts the arguments

### .(Atom...)Block
    Access a variable in the binding

### /(Decimal...)Decimal
    Divides the arguments

### <(Time...)Bool
    Tests if the first argument is less then the following

### <(Decimal...)Bool
    Tests if the first argument is less then the following

### <=(Decimal...)Bool
    Tests if the first argument is less or equal then the following

### <=(Time...)Bool
    Tests if the first argument is less or equal then the following

### =(Any...)Bool
    Tests if the arguments are the same

### >(Decimal...)Bool
    Tests if the first argument is greather then the following

### >(Time...)Bool
    Tests if the first argument is greather then the following

### >=(Time...)Bool
    Tests if the first argument is greather or equal then the following

### >=(Decimal...)Bool
    Tests if the first argument is greather or equal then the following

### between(Time...)Bool
    Tests if the arguments are between the second last and the last argument

### between(Decimal...)Bool
    Tests if the arguments are between the second last and the last argument

### ceil(Decimal)Decimal
    Ceil the decimal argument

### concat(String...)String
    Concat strings

### contains(String...)Bool
    Returns wether the first argument exists in the following arguments

### drop(Block)Block
    Create a list containing all but the last item in the input list

### endsWith(String...)Bool
    Returns wether the first argument is the suffix of the following arguments

### floor(Decimal)Decimal
    Floor the decimal argument

### head(Block)Block
    Returns the first item in the list

### item(Block, Decimal)Block
    Returns a specific item from a list

### list(Atom...)Block
    Create a list out of the children

### misc3(Block)Kind(0)
    

### mod(Decimal...)Decimal
    Modulo the arguments

### noop()Any
    No operation

### notContains(String...)Bool
    Returns wether the first argument does not exist in the following arguments

### startsWith(String...)Bool
    Returns wether the first argument is the prefix of the following arguments

### tail(Block)Block
    Returns list without the first item

### toString(Any)String
    Converts the parameter to a string

### ~(String, String)Bool
    Returns wether the first argument matches the regular expression in the second argument

