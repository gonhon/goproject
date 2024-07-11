add
curl -X POST -H "Content-Type:application/json" -d '{"id": "978-7-111-55842-2", "name": "The Go Programming Language", "authors":["Alan A.A.Donovan", "Brian W. Kergnighan"],"press": "Pearson Education"}' localhost:8090/book

get
curl -X GET -H "Content-Type:application/json" localhost:8090/book/978-7-111-55842-2