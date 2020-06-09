# Data Basket

![Go](https://github.com/sporule/data-basket/workflows/Go/badge.svg?branch=master)

> A test data generator through Regex

## Latest Update

### v0.1

- Initial Release

## Features

- Easy to use, generate data through json config file
- Parallelism (Default to CPU Cores * 4)

## Credits To

[Reggen](https://github.com/lucasjones/reggen)

## Quick Start

1. Download the latest binary from the [release page](https://github.com/sporule/data-basket/releases)
2. Create your config file. The config file is in json format with below options.

| Option   | Description                         | Value Type                      |
| -------- | ----------------------------------- | ------------------------------- |
| fileName | The name of the output file         | string                          |
| rows     | Number of rows you want to generate | int                             |
| columns  | Columns to generate                 | "columnName":"Regex Definition" |

```python
{
"fileName": "data.csv",
    "rows": 1000000,
    "columns": {
        "Column1": "Column1\\d{1,3}[A-Z]{1,3}",
        "Column2": "Column2\\d{1,3}[A-Z]{1,3}"
    }
}
```
3. Length is required when you defined the regex, see above examples. Otherwise each columns will return the length of Max Int8. 
4. Run the binary in terminal and pass the config file as argument, it will use config.json as default config if you leave the argument empty:

```bash
./data-basket config.json
```
1. You will see output from the terminal, and the data will be stored in the output file you specified:

```bash
Writing 5000000 rows to output.csv with column names: [Column1 Column2] 
Go Routine Size: 16 
10.00 % : Written 500000 rows to file output.csv 
20.00 % : Written 1000000 rows to file output.csv 
30.00 % : Written 1500000 rows to file output.csv 
40.00 % : Written 2000000 rows to file output.csv 
50.00 % : Written 2500000 rows to file output.csv 
60.00 % : Written 3000000 rows to file output.csv 
70.00 % : Written 3500000 rows to file output.csv 
80.00 % : Written 4000000 rows to file output.csv 
90.00 % : Written 4500000 rows to file output.csv 
100.00 % : Written 5000000 rows to file output.csv 
Total Time used: 10.769763323s
```

