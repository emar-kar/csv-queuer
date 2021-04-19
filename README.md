# Welcome to a CSV-QUEUER!

A simple tool which allows you to perform some sort of the SQL queues requests on a csv file.

## Set up
To get started you can use the stable build artifact for **macOS** and **linux** systems.

First you need to configure some values for the proper use.

* Create ./configs/config.yml file;
* Add 3 variables:

    * *csv_separator* - this value defaults to "," but if you have another file separator, then you can customize it. (**Do not use dots as separators, float numbers might defined incorrectly in this case!** )
    * *request_timeout* - in seconds. This value defaults to 5. Defines each request timeout. In case of timeout deadline parsed data will be printed.
    * *log_folder* - this value defaults to "./logs", but can be set-up manually.

## Request language
CSV-queuer parses given request string and gets specified fields for you.
Each request should end with **" ; "**. This is built-in separator.

Your request have to have those fields:
* **SELECT** - what fields you would like to extract.
* **FROM** - path to a csv file.
* **WHERE** - options of the request.

**This options + NOT, AND, OR (see below) should be always in capital letters!**

## SELECT:
This field cannot be omitted! You always need to specify it.

This field supports **" * "**. It means you would like to get all fields from the CSV file.

Results of the request will be printed in order of the elements defined in this field.

## FROM
This field cannot be omitted! You always need to specify it.

## WHERE
This field can be omitted. It means you would like to get all specified fields from the CSV file without any requirements.

Conditions have several key words which help you to clearly represent your idea of what you would like to receive. They are:

* *AND* - defines that next condition is strictly wanted.
* *OR* - defines that next condition can be omitted.

Options for variables:

* *NOT* - not equal.
* *=* - equal.
* *>* - greater.
* *<* - less.
* *>=* - greater or equal.
* *<=* - less or equal.

Currently app reads the string from left to right and does not support brackets.

That is why you need to think a little bit more about your request, for example:

```
SELECT location, new_cases, date 
FROM path/to/your/file.csv
WHERE location = Russia OR location = Ukraine AND date >= 2020-04-20 AND date <= 2020-04-30 AND new_cases > 500 AND new_cases < 5500;
```

and

```
SELECT location, new_cases, date 
FROM path/to/your/file.csv
WHERE location = Russia OR location = Ukraine AND date >= 2020-04-20 AND date <= 2020-04-30 AND new_cases > 500 OR new_cases < 5500;
```

are two different requests which leads to two different outputs.

The first one is correct, so the result will be:

```
=====================================
| location | new_cases |    date    |
=====================================
|  Russia  |  4268.0   | 2020-04-20 |
|  Russia  |  5236.0   | 2020-04-22 |
|  Russia  |  4774.0   | 2020-04-23 |
| Ukraine  |   578.0   | 2020-04-23 |
| Ukraine  |   540.0   | 2020-04-30 |
=====================================
```

But the result of the second:

```
=====================================
| location | new_cases |    date    |
=====================================
|  Russia  |  4268.0   | 2020-04-20 |
|  Russia  |  5642.0   | 2020-04-21 |
|  Russia  |  5236.0   | 2020-04-22 |
|  Russia  |  4774.0   | 2020-04-23 |
|  Russia  |  5849.0   | 2020-04-24 |
|  Russia  |  5966.0   | 2020-04-25 |
|  Russia  |  6361.0   | 2020-04-26 |
|  Russia  |  6198.0   | 2020-04-27 |
|  Russia  |  6411.0   | 2020-04-28 |
|  Russia  |  5841.0   | 2020-04-29 |
|  Russia  |  7099.0   | 2020-04-30 |
| Ukraine  |   578.0   | 2020-04-23 |
| Ukraine  |   540.0   | 2020-04-30 |
=====================================
```

The missunderstanding came from the option `new_cases`.

Since we have `AND new_cases > 500` this means that this criterion is strict and should always be true. 
But next we can write `OR new_cases < 5500` and it will take no effect. 
Because this option is not strict and can be omitted. On the other hand if you put `AND new_cases < 5500` 
this will means that this condition should be true as well, so the result will be more specific.

Right now the app can understand several datatypes:

* integer number e.g. 4
* float number e.g. 4.5 (float numbers have to be devided by dot!)
* string values
* date values e.g. 2020-11-18 (app can understand dates in **YYYY-MM-DD** format)
