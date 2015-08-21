# sumoshell
Sumo shell is a partial implementation of the Sumo Logic stream pipeline written in Go. Commands should start with
`sumo` which will transform logs into the json format sumoshell uses. Commands should end with `render` `render-basic` or `graph` which render the output to the terminal. 

Examples:

```
cat log | sumo "search term" | parse "responseTime=*]" as resp | avg resp | render
cat log | sumo "search term" | parse "responseTime=*][host=*]" as resp, host | avg resp by host | render
```
