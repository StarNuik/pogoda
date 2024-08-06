# PoGoDa
> A straight-forward Go Postgres driver implementation.

The primary goal of the project is to conform to the database/sql "interface" with as little programmer effort as possible. It MAYBE is a nice showcase of the postgres wire protocol, as there are no optimizations that could ruin the understanding of the flow.

## Pogoda is
* Simple
* Easy to understand (hopefully)

## Pogoda isn't
* Efficient / fast

## Possible improvements
* Buffer scratch array (pq)
* Flyweight (pgx)
* Sql caching (pgx, mb pq)