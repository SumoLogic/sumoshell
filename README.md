# sumoshell
[![Join the chat at https://gitter.im/SumoLogic/sumoshell](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/SumoLogic/sumoshell?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

`Sumoshell` is collection of utilities to improve analyzing log files written in Go. `grep` can't tell that some log lines span multiple individual lines. Parsing out fields is cumbersome. Aggregating is basically impossible, and there is no good way to view the results. In Sumoshell, each individual command acts as a phase in a pipeline to get the answer you want. Sumoshell brings a lot of the functionality of [Sumo Logic](www.sumologic.com) to the command line.

Commands should start with
`sumo search [filter]` which will transform logs into the json format `sumoshell` uses. Commands should end with `render` `render-basic` or `graph` which render the output to the terminal. Each operator is a stand-alone binary allowing them to be easily composed.

## Installation
[OSX and Linux binaries are provided for sumoshell](https://github.com/SumoLogic/sumoshell/releases). Simply extract the archive and place the binaries on your path. 

If you run a different OS or would prefer to install from source, it's easy to build from source. Given a working [go](https://golang.org/doc/install) installation, run:
```
go get github.com/SumoLogic/sumoshell
cd $GOPATH/src/github.com/SumoLogic/sumoshell # Will warn about `no buildable go source`
go get ./...
go install ./...
```

## Usage
Like [SumoLogic](https://www.sumologic.com), sumoshell enables you pass log data through a series of transformations to get your final result. Pipelines start with a source (`tail`, `cat`, etc.) followed by the `sumo` operator. An example pipeline might be:

```tail -f logfile | sumo search "ERROR" | sumo parse "thread=*]" | sumo count thread | render-basic```

This would produce a count of log messages matching `ERROR` by thead. In the basic renderer, the output would look like:
```
_Id   _count   thread   
0     4        C        
1     4        A        
2     1        B      
```
### The `sumo search` operator
`sumo search` takes an optional filter parameter to allow for basic searching. The sumo operator performs 3 steps: 

1. Break a text file into logical log messages. This merges things like stack traces into a single message for easy searching.
2. Allow basic searching.
3. Transforms the log message into the sumoshell internal json format.

### Displaying results

After using the `sumo` operator, the output will be in JSON. To re-render the output in a human-readable form, `|` the results of your query into one of the three `render` operators.

1. `render-basic`: Capable of rendering aggregate and non-aggregate data. Mimics curses style CLIs by calculating the terminal height and printing new lines to the end to keep your text aligned. Add `nowraw` to drop the raw data when an aggregate isn't present.
2. `render`: Curses based renderer for rendering tabular data.
3. `graph`: Curses based renderer for rendering tabular data as a bar chart.


### Parsing Data

sumoshell supports a basic parse operator similar to the `parse` operator in `SumoLogic`. Queries take the form:
```
... | sumo parse "[pattern=*] pattern2:'*' morePatterns=(*)" as pattern, pattern2, more | ...
```

### Filtering Data

sumoshell supports a filter operator similar to the `where` operator in `SumoLogic`. Queries take the form:
```
... | sumo parse "[host=*]" as host | sumo filter host = server1 
```

This will drop any log lines that don't have `server1` as the host.

### Aggregating Data

sumoshell currently supports 3 aggregate operators:

1. `count`   Example queries:
  ```
  ... | sumo count           # number of rows
  ... | sumo count key       # number of rows per key
  ... | sumo count key value # number of rows per the cartesian product of (key, value)

  ```

2. `sum`     Example queries:
  ```
  ... | sumo sum k      # sum of all k's
  ... | sumo sum v by k # sum of all v's by k

  ```

3. `average` Example queries:
  ```
  ... | sumo average k      # average of all k's
  ... | sumo average v by k # average of all v's by k

  ```
