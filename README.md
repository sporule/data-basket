# Data Basket

> A test data generator through Regex

## Latest Update

### 2020-06-06

- Initial Release

## Features

- Easy to use, generate data through json config file
- Parallelism (Default to CPU Cores * 4)

## Credits To

[Reggen](https://github.com/lucasjones/reggen)

## Quick Start - Set up

### Fork This Repo to your GitHub Account

![Fork](https://i.imgur.com/VSqrEHf.png)

### Update the Configuration Files

#### General Configuration File

_config.js is the configuration file for your site.

![Config](https://i.imgur.com/9Rl3J3B.png)

**Available Configs:**

| Name              | Value                                                                                               | Example                           | Type    |
| ----------------- | --------------------------------------------------------------------------------------------------- | --------------------------------- | ------- |
| Site              | The name of the site                                                                                | "Sporule"                         | string  |
| url               | The link to your site                                                                               | "https://www.sporule.com"         | string  |
| description       | short description about your site                                                                   | "Sporule is a micro blog system"  | string  |
| keywords          | keywords for SEO purpose                                                                            | "blog,post"                       | string  |
| logo              | The logo                                                                                            | "https://i.imgur.com/MrRLOC4.png" | string  |
| disqusShortname   | Disqus is a comment system, you will get a shortname after having an account with them              | NA                                | string  |
| postPerPage       | How many posts do you want to show per page?                                                        | 8                                 | int     |
| googleAnaltics    | Google Analytics Tag if you are using Google analytics                                              | UA-11402457-1                     | string  |
| alwaysRefreshPost | The system will always get the latest content of the post rather than using the cache if it is true | false                             | boolean |
| gh_custom_domain  | Put this to true if you are using github pages with custom domain                                   | false                             | boolean |

#### Theme Configurations

These configuration will be used for that specific theme only, they are located under templates/_templateConfig.js. Look at the theme documentation to understand what is available. You will not be able to change the template configuration unless you fork the template repo. Please see change templates section for more details.

![TemplateConfig](https://i.imgur.com/mVoIG2w.png)

## Quick Start - Your Content

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
3. Length is required when you defined the regex, see above examples. Otherwise each columns will return the length of Max Int64. 
4. Run the binary in terminal and pass the config file as argument:

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

