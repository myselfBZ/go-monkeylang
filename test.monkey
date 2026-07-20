let adder = fn(y) {
    let c = fn(x) {
        return x + y
    }
    return c
}
let addOne = adder(1);
addOne(6);
