# wvu-purchases

**NOTE**: This is currently a work in progress.

This project provides a SQLite(3) database of purchase records from [West Virginia University](https://wvu.edu), and a simple command-line tool to help manage record ingestion from the various formats (in Excel) used by the [Office of Procurement Contracting & Payment Services](https://procurement.wvu.edu/).

### Data

All information has been sourced from WVU over the course of many years via WVFOIA (W.Va. Code § 29B-1-1), and the intent is to keep the database updated on a rolling fiscal year basis as the University releases information.

Current data spans from **Jan 26, 2017** to **Dec 2, 2021**.

| procurement ||
|---|---|
| | |
| date | count(records) |
|---|---|
| 2017-Jan-26 to 2017-Dec-31 | 3,898 |
| 2018-Jan-1 to 2018-Dec-31 | 1,204 |
| 2019-Jan-1 to 2019-Dec-31 | 1,313 |
| 2020-Jan-1 to 2020-Dec-31 | 1,728 |
| 2021-Jan-1 to 2021-Dec-2 | 482 |
| | **total**: 8,625 |

| pcard ||
|---|---|
| _still in progress_ | |


### License

See [LICENSE.txt](./LICENSE.txt).