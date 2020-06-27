# lunchserv
Server to easily distribute daily Doordash order links

```
> go run main.go 8080 &
[1] 2468

> curl localhost:8080
Order hasn't been created yet

> curl localhost:8080 -XPOST -d 'https://drd.sh/cart/bJ0TmX/'

> curl localhost:8080
<a href="https://drd.sh/cart/bJ0TmX/">Temporary Redirect</a>.

>
```
