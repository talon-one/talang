[![Travis](https://img.shields.io/travis/talon-one/talang.svg)](https://travis-ci.org/talon-one/talang) [![Codecov](https://img.shields.io/codecov/c/github/talon-one/talang.svg)](https://codecov.io/gh/talon-one/talang) [![GoDoc](https://godoc.org/github.com/talon-one/talang?status.svg)](https://godoc.org/github.com/talon-one/talang) [![go-report](https://goreportcard.com/badge/github.com/talon-one/talang)](https://goreportcard.com/report/github.com/talon-one/talang)

Talang
==================

  - [Usage](#usage)
  - [Playing Around](#playing-around)
  - [Contributing](#contributing)


Talang is a custom programming language, specifically a [lisp](https://en.wikipedia.org/wiki/Lisp_(programming_language))-dialect, implemented in Go, that we developed and use internally at [Talon.One](https://talon.one). 


### *Motivation*

The motivation behind developing our own custom language is that the very nature of what our customers are doing in the core of our product is a form of programming. We needed a language to represent these "programs" with certain characteristics:

- Easy to parse, interpret and manipulate
- Easy to represent in JSON
- Safe runtime properties (Talang has no recursion, no infinite loops)
- Type-safe expressions

No existing language would fit these requirements and was also easy to integrate.

### *How do we use Talang at Talon.One?*

Within our product, we let our customers write their own sets of what we call "*Rules*". These are concluded of two main types of expressions:

 - *Conditions*: representing predicates which evaluate to a boolean reault
 - *Effects*: representing expressions which return side effects handled by the promotions
 
We needed a way to give our customers the flexibility and control when composing such "Rules".


At the core of our product we developed a processing engine which takes these predefined Rules, gives them an execution context, and evaluates them. The return result(s) from this engine then returned to our customers to process / work with.

### *Why would you need it?*

Talang is a very simple, fast and flexible language which is very easy to represents in JSON objects.
When you desire a language, customizable, expandable, with no compiling time and fast evaluation time - Talang can be a very good choice for your requirements.

## Usage

Get the package:

    $ go get github.com/talon-one/talang

Then import it and use the Interpreter:

```go
interp := talang.MustNewInterpreter()
result, err := interp.LexAndEvaluate(`(+ 1 2)`)
if err != nil {
	panic(err)
}
fmt.Println(result.Stringify()) // 3
```


You can refer to the [examples](https://github.com/talon-one/talang/tree/master/examples) folder for more examples and usages.

[Here](https://talon-one.github.io/talang/docs/functions) you can see a list of the embedded function in the language.

## Playing Around

You can get a feeling of how is it to write some Talang using the integrated [CLI](https://github.com/talon-one/talang/tree/master/cmd/talang-cli) tool.

## Contributing

We have collected notes on how to contribute to this project in [CONTRIBUTING.md](https://github.com/talon-one/talang/tree/master/CONTRIBUTING.md).


## License

Talang is released under the [MIT License](https://github.com/talon-one/talang/tree/master/LICENSE).
