# wvu-purchases 💸

![Creative Commons Zero v1.0](https://licensebuttons.net/p/zero/1.0/88x15.png)

**NOTE**: This is currently a work in progress.

This project provides and a SQLite(3) database of purchase records from [West Virginia University](https://wvu.edu), a web-accessible way to browse and query the database using the [Datsette](https://datasette.io) project, and a simple command-line tool to help project admins manage database updates from the various formats (in Excel) used by the [WVU Office of Procurement Contracting & Payment Services](https://procurement.wvu.edu/).

### Data

All information has been sourced from WVU over the course of many years via WVFOIA (W.Va. Code § 29B-1-1), and the intent is to keep the database updated on a rolling fiscal year basis as the University releases information.

Current data spans from **Oct 3, 2010** to **Dec 2, 2021**⁺, with over $2 billion dollars worth of transactions.

<table>
<tr><th>procurement_records</th></tr>
<tr><td>

|year |date_range                |count|total_amount   |
|-----|--------------------------|-----|---------------|
|2013 |2013-01-02 thru 2013-12-30|16,073|$127,484,829.84|
|2014 |2014-01-02 thru 2014-12-31|24,047|$218,340,699.54|
|2015 |2015-01-05 thru 2015-12-30|21,952|$247,358,331.70|
|2016 |2016-01-04 thru 2016-09-27|18,337|$496,400,373.52|
|2017 |2017-01-26 thru 2017-12-28|3,898|$87,505,408.53 |
|2018 |2018-01-02 thru 2018-12-21|1,204|$79,222,514.55 |
|2019 |2019-01-02 thru 2019-12-20|1,313|$57,588,736.12 |
|2020 |2020-01-07 thru 2020-12-29|1,728|$128,618,446.67|
|2021 |2021-01-04 thru 2021-11-30|482  |$73,434,112.34 |
|Total|                          |**89,034**|**$1,515,953,452.81**|
<details> 
  <summary>View SQL Query</summary>

   ```sql
SELECT
    strftime('%Y', approved_date) AS year,
    printf('%s thru %s', MIN(approved_date), MAX(approved_date)) AS date_range,
    printf('%,d', count(*)) as count,
    printf('$%,.2f', SUM(amount)) AS total_amount
FROM
    procurement_records
GROUP BY
    strftime('%Y', approved_date)
UNION ALL
SELECT
    'Total' AS year,
    NULL as date_range,
    printf('%,d', count(*)) AS count,
    printf('$%,.2f', SUM(amount)) AS total_amount
FROM
    procurement_records
ORDER BY
    year;
   ```
</details>
</td></tr>
<th>pcard_records</th>
<tr><td>

**** 
|year |date_range                |count|total_amount   |
|-----|--------------------------|-----|---------------|
|2010 |2010-10-03 thru 2010-12-31|492  |$125,628.12    |
|2011 |2011-01-01 thru 2011-12-31|334,977|$100,737,264.66|
|2012 |2012-01-01 thru 2012-12-31|353,595|$105,301,445.82|
|2013 |2013-01-01 thru 2013-12-31|329,023|$87,282,066.05 |
|2014 |2014-01-01 thru 2014-12-31|348,491|$90,847,153.18 |
|2015 |2015-01-01 thru 2015-12-31|356,010|$76,880,023.01 |
|2016 |2016-01-01 thru 2016-12-31|170,050|$31,360,041.57 |
|2017 |2017-01-01 thru 2017-12-31|229,248|$53,004,052.84 |
|2018 |2018-01-01 thru 2018-03-13|46,021|$10,373,883.79 |
|**Total**|                          |**2,167,907**|**$555,911,559.04**|
<details> 
  <summary>View SQL Query</summary>

   ```sql
SELECT
    strftime('%Y', trans_date) AS year,
    printf('%s thru %s', MIN(trans_date), MAX(trans_date)) AS date_range,
    printf('%,d', count(*)) as count,
    printf('$%,.2f', SUM(trans_amount)) AS total_amount
FROM
    pcard_records
GROUP BY
    strftime('%Y', trans_date)
UNION ALL
SELECT
    'Total' AS year,
    NULL as date_range,
    printf('%,d', count(*)) AS count,
    printf('$%,.2f', SUM(trans_amount)) AS total_amount
FROM
    pcard_records
ORDER BY
    year;
   ```
</details>
</td></tr> </table>

### License

This project is released into the public domain via Creative Commons Zero. See the [LICENSE](./LICENSE) file.