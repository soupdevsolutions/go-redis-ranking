# go-redis-ranking

This repository contains a basic setup using Go, Gin and Redis to build a ranking/leaderboard system.  
You can find a full tutorial on how to build this on the [soup.dev blog](https://www.soup.dev/post/building-a-ranking-system-with-go-and-redis).

## Setup 

You should have Go and either Docker or Redis installed on your system.  

To run the project, use the following commands:  

```shell
git clone https://github.com/humbertodias/go-redis-ranking.git
cd src
go get .
docker run -d -p 6379:6379 redis:latest
go run .
```
or

```shell
cd `mktemp -d` && git clone https://github.com/humbertodias/go-redis-ranking.git .
docker compose up -d
```

## Interacting with the system

You can add new entries into the system by sending a POST request:
```shell
curl -X POST http://localhost:8080/register -d '{"name":"test"}'
```
```
{
   "id": "e309ab7e-6e02-4b60-a374-52cd5c2a41dd",
   "name": "user",
   "rank": 6,
   "score": 4
} 
```

You can then use the id from the response to query the data for that entry:
```shell
curl -X GET “http://localhost:8080/rank?id=e309ab7e-6e02-4b60-a374-52cd5c2a41dd”
```
```
{
       "id": "e309ab7e-6e02-4b60-a374-52cd5c2a41dd",
       "name": "test",
       "rank": 6,
       "score": 4
}
```

Finally, you can get slices of the leaderboard by sending:
```shell
curl -X GET http://localhost:8080/ranks?loffset=1&imit=3
```
```
[
    {
           "id": "0a363580-98f9-4d4c-b57b-e30259807871",
           "name": "AUserWithQuiteALongName",
           "rank": 1,
           "score": 23
    },
    {
           "id": "629e846e-7d8c-4a4a-b4dc-72c98b3038cc",
           "name": "AUserWithAShorterName",
           "rank": 2,
           "score": 21
    },
    {
           "id": "06894035-f00b-439f-9837-8096af59de51",
           "name": "AnEvenShorterName",
           "rank": 3,
           "score": 17
    }
]
```
