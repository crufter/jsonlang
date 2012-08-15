jsonlang
========
JSONLang is the most retarded scripting language ever intended for actual use.

Some fun facts:
- Pure Go.
- No loops, ifs, function calls, only ASM-like jumps.
- Easy nested JSON access syntax.
- Can read and modify outside variables.
- Extendable with custom Go functions (they need a 5-line wrapper each).
- The syntax is a crazy disguised version of JSON.*

*Lets take a look at the example source code:
```
isOdd(&is_odd, 3);
set(&x, { "val": is_odd });
set(&f, {"lol":{"trolol":22}});
println(f.lol.trolol);
println(x);
```

Now, if we replace "(" with ", " and ")" with nothing, we get:
```
isOdd, &is_odd, 3;
set, &x, { "val": is_odd };
set, &f, {"lol":{"trolol":22}};
println, f.lol.trolol;
println, x;
```

If we quote the unquoted strings and mark them with a special char:
```
".isOdd", ".&is_odd", 3;
".set", ".&x", { "val": ".is_odd" };
".set", ".&f", {"lol":{"trolol":22}};
".println", ".f.lol.trolol";
".println", ".x";
```

We can distingish between string constants and variable/function names.

We only have to replace ";" with "],\n [":
```
".isOdd", ".&is_odd", 3],
[".set", ".&x", { "val": ".is_odd" }],
[".set", ".&f", {"lol":{"trolol":22}}],
[".println", ".f.lol.trolol"],
[".println", ".x"
```

We add "[[" to the beginning and "]]" to the end, run it trough jsonlint.org:
```
[
    [
        ".isOdd",
        ".&is_odd",
        3
    ],
    [
        ".set",
        ".&x",
        {
            "val": ".is_odd"
        }
    ],
    [
        ".set",
        ".&f",
        {
            "lol": {
                "trolol": 22
            }
        }
    ],
    [
        ".println",
        ".f.lol.trolol"
    ],
    [
        ".println",
        ".x"
    ]
]
```

And we can use json.Unmarshal to decode our primitive language into a primitive AST.
Thank you for watching Useless Projects TV and don't forget to watch all the commercials, because taking a piss while they run is like robbing the cable company.