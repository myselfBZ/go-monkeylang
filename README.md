# An interpreter for the monkey language

Based off of the Writing an Interpreter in Go book, by Thorsten Ball.


Examples:


Variable declaration
```monkey
let x = 1;
let b = true;
let b = false;
```
Boolean expressions
```monkey
let x = 1 == 1; // true
let x = 1 >= 2; // true
let x = 1 <= 2; // true
let x = 2 != 1; // false
```
Arithmetic expressions
```monkey
let x = 1 + 1; // 2
let x = (-12) + 2; // -10
let x = 4 * 2 + 1; // 9
let x = 16 / 2 - 1; // 7 
let x = (2 * 2) + 3; // 7
```

Conditional `if` expressions
```monkey
let x = 12;
if x > 10 {
    x = x + 1;
} else {
    x = x - 1;
}

let x = if true {
        1
    } else {
        0
    }
```

Functions
```monkey
let adder = fn (x) {
    return fn (y) {
        return x + y
    }
}

let addOne = adder(1);
addOne(6);
```
