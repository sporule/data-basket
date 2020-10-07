# Data Basket

![Go](https://github.com/sporule/data-basket/workflows/Go/badge.svg?branch=master)

> A test data generator through Regex

## Latest Update

### v0.2

- Added the option to control one to many and one to one relationship between columns
- Added the option to set sample size rather than generate randomly

### v0.1

- Initial Release

## Features

- Easy to use, generate data through json config file
- Parallelism (Default to CPU Cores * 4)

## Known Bugs

- It may generate duplicate records in windows/

## Quick Start

1. Download the latest binary from the [release page](https://github.com/sporule/data-basket/releases)
2. Create your config file. The config file is in json format with below options.

| Option   | Description                                                                                                                                                                                                                                                                     | Value Type |
| -------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------- |
| fileName | The name of the output file                                                                                                                                                                                                                                                     | string     |
| rows     | Number of rows you want to generate                                                                                                                                                                                                                                             | int        |
| columns  | Columns to generate                                                                                                                                                                                                                                                             | obj        |
| Pattern  | The regex pattern                                                                                                                                                                                                                                                               | string     |
| Size     | How many unique values the generator will produce                                                                                                                                                                                                                               | int        |
| Group    | Use for relationship. If two columns are in the same group then Size option will be use to determind relationship. The example below shows one to one relationship between product code and price, it also has one to many relationship between product group and product code. | int        |

```python
{
    "fileName": "output.csv",
    "rows": 1000000,
    "columns": {
        "ProductCode": {
            "Pattern":"[A-Z]{5}\\d{5}",
            "Size":100,
            "Group":1
        },
        "ProductGroup": {
            "Pattern":"(Fruit)|(Meat)|(Vegetable)",
            "Size":3,
            "Group":1
        },
        "Price": {
            "Pattern":"[1-9][0-9]{5}",
            "Size":100,
            "Group":1
        },
        "Notes": {
            "Pattern":"Notes:\\w{3}"
        },
        "Picker": {
            "Pattern":"Picker:\\w{3}",
            "Size":10
        }

    }
}
```

3. Length is required when you defined the regex, see above examples. Otherwise each columns will return the length of Max Int8. 
4. Run the binary in terminal and pass the config file as argument, it will use config.json as default config if you leave the argument empty:

```bash
./data-basket config.json
```

5. You will see output from the terminal, and the data will be stored in the output file you specified:

```bash
Set up Go Routine Size: 16 
Could not find config file, loading default config.json file 
Writing 1000000 rows to output.csv with column names:
 [Notes Picker Price ProductCode ProductGroup] 
Preparing Relationships Model... 
50.00 % : Written 500000 rows to file output.csv 
100.00 % : Written 1000000 rows to file output.csv 
Total Time used: 16.577820146s
```


## Package Used

[Reggen](https://github.com/lucasjones/reggen)
