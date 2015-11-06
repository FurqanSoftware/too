# Too

[![Build Status](https://drone.io/github.com/hjr265/too/status.png)](https://drone.io/github.com/hjr265/too/latest)

Too is a simple recommendation engine built on top of Redis in Go.

## Installation

Install Too using the go get command:

    $ go get github.com/hjr265/too

The only dependencies are Go distribution and [Redis](http://redis.io).

## Usage

```go
te, err := too.New("redis://localhost", "movies")
if err != nil {
	log.Fatal(err)
}

te.Likes.Add("Sonic", "The Shawshank Redemption")
te.Likes.Add("Sonic", "The Godfather")
te.Likes.Add("Sonic", "The Dark Knight")
te.Likes.Add("Sonic", "Pulp Fiction")

te.Likes.Add("Mario", "The Godfather")
te.Likes.Add("Mario", "The Dark Knight")
te.Likes.Add("Mario", "The Shawshank Redemption")
te.Likes.Add("Mario", "The Prestige")
te.Likes.Add("Mario", "The Matrix")

te.Likes.Add("Peach", "The Godfather")
te.Likes.Add("Peach", "Inception")
te.Likes.Add("Peach", "Fight Club")
te.Likes.Add("Peach", "WALL·E")
te.Likes.Add("Peach", "Princess Mononoke")

te.Likes.Add("Luigi", "The Prestige")
te.Likes.Add("Luigi", "The Dark Knight")

items, _ := te.Suggestions.For("Luigi", 2)
for _, item := range items {
	fmt.Println(item)
}

// Output:
// The Shawshank Redemption
// The Matrix
```

## Documentation

- [Reference](http://godoc.org/github.com/hjr265/too)

## Contributing

Contributions are welcome.

## License

Too is available under the [BSD (3-Clause) License](http://opensource.org/licenses/BSD-3-Clause).

## Inspiration

This project is inspired by the very existence of the awesome project [Recommendable](http://davidcel.is/recommendable/).
