# mercari mackerel plugin chocon

this is the mackerel plugin for [chocon](https://github.com/kazeburo/chocon)

this plugin is for chocon >= 0.9.0

## Usage

```bash
$ ./mackerel-plugin-chocon -h
```

## Example

```bash
$ ./mackerel-plugin-chocon
chocon.http.requests.count	13131.428571	1548920564
chocon.http.requests.count_200	13122.857143	1548920564
chocon.http.requests.count_403	0	1548920564
chocon.http.requests.count_4xx	17.142857	1548920564
chocon.http.requests.count_5xx	0	1548920564
chocon.http.latency.avg_time	0.056403	1548920564
chocon.http.latency.time_percentile_90	0.167347	1548920564
chocon.http.latency.time_percentile_95	0.234521	1548920564
chocon.http.latency.time_percentile_99	0.623367	1548920564
```
