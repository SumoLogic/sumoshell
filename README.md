# sumoshell
Sumo shell is a partial implementation of the Sumo Logic stream pipeline written in Go. Commands should start with
`sumo` which will transform logs into the json format sumoshell uses. Commands should end with `render` `render-basic` or `graph` which render the output to the terminal. 

## Installation
Currently no binaries are provided for Sumoshell, however it's easy to build from source. Given a working [go](https://golang.org/doc/install) installation, run:
```
go get github.com/SumoLogic/sumoshell
cd $GOPATH/src/github.com/SumoLogic/sumoshell # Will warn about `no buildable go source`
go get ./...
go install ./...
```

## Usage
Like [SumoLogic](www.sumologic.com), sumoshell enables you pass log data through a series of transformations to get your final result. Pipelines start with a source (`tail`, `cat`, etc.) followed by the `sumo` operator. An example pipeline might be:

```tail -f logfile | sumo "ERROR" | parse "thread=*]" | count thread | render-basic```

This would produce a count of log messages matching error by thead. In the basic renderer, the output would look like:
```
_Id   _count   thread   
0     4        C        
1     4        A        
2     1        B      
```
### The `sumo` operator
The sumo operator performs 3 steps: 

1. Break a text file into logical log messages. This merges things like stack traces into a single message for easy searching.
2. Allow basic searching.
3. Transforms the log message into the sumoshell internal json format.

### Displaying results

After using the `sumo` operator, the output will be in JSON. To re-render the output in a human-readable form, `|` the results of your query into one of the three `render` operators.

1. `render-basic`: Capable of rendering aggregate and non-aggregate data. Mimics curses style CLIs by calculating the terminal height and printing new lines to the end to keep your text aligned.
2. `render`: Curses based renderer for rendering tabular data
3. `graph`: Curses based renderer for rendering tabular data as a bar chart.
