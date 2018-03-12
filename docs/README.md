# Talang reference guide [![talang](https://img.shields.io/badge/talang-reference-5272B4.svg)](#)

[Introduction](#introduction)  
[Embedded Functions](#embedded-functions)  
[Types](#types)  

----

## Introduction

## Embedded Functions

You can see a list of all embedded functions in [functions.md](functions.md)


## Types
| Name        | Description                                                                    | Example                               |
|-------------|--------------------------------------------------------------------------------|---------------------------------------|
| Decimal     |                                                                                | `1.2`                                 |
| String      |                                                                                | `Hello`                               |
|             |                                                                                | `"Hello World"`                       |
| Bool        | `true` or `false`                                                              | `true`                                |
|             |                                                                                | `false`                               |
| Time        |                                                                                | `Mon Jan 2 15:04:05 MST 2006`         |
| Null        |                                                                                |                                       |
| List        |                                                                                | `list 1 true "Hello World"`           |
| Map         |                                                                                | `kv (Key1 true) (Key2 "Hello World")` |
| Block       |                                                                                |                                       |
| Atom        | Reserved Type that can be one of `Decimal`, `String`, `Bool`, `Time` or `Null` |                                       |
| Collection  | Reserved Type that can be one of `List` or `Map`                               |                                       |
| Any         | Reserved Type that can be one of `Atom`, `Block` or `Collection`               |                                       |