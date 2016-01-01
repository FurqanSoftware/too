// Copyright 2014 The Too Authors. All rights reserved.

package too_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/hjr265/too"
)

func ExampleEngine() {
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
}

func ExampleBatch() {
	te, err := too.New("redis://localhost", "movies")
	if err != nil {
		log.Fatal(err)
	}

	err = te.Likes.Batch([]too.BatchRaterOp{
		{
			User: "Sonic",
			Items: []too.Item{
				"The Shawshank Redemption",
				"The Godfather",
				"The Dark Knight",
				"Pulp Fiction",
			},
		},
		{
			User: "Mario",
			Items: []too.Item{
				"The Godfather",
				"The Dark Knight",
				"The Shawshank Redemption",
				"The Prestige",
				"The Matrix",
			},
		},
		{
			User: "Peach",
			Items: []too.Item{
				"The Godfather",
				"Inception",
				"Fight Club",
				"WALL·E",
				"Princess Mononoke",
			},
		}, {
			User: "Luigi",
			Items: []too.Item{
				"The Prestige",
				"The Dark Knight",
			},
		},
	}, true)

	if err != nil {
		log.Fatal(err)
	}

	items, _ := te.Suggestions.For("Luigi", 2)
	for _, item := range items {
		fmt.Println(item)
	}

	// Output:
	// The Shawshank Redemption
	// The Matrix
}

func BenchmarkNoBatch(b *testing.B) {
	te, err := too.New("redis://localhost", "movies")
	if err != nil {
		log.Fatal(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
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
	}
}

func BenchmarkUsingBatchWithoutAutoUpdate(b *testing.B) {
	te, err := too.New("redis://localhost", "movies")
	if err != nil {
		log.Fatal(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		te.Likes.Batch([]too.BatchRaterOp{
			{
				User: "Sonic",
				Items: []too.Item{
					"The Shawshank Redemption",
					"The Godfather",
					"The Dark Knight",
					"Pulp Fiction",
				},
			},
			{
				User: "Mario",
				Items: []too.Item{
					"The Godfather",
					"The Dark Knight",
					"The Shawshank Redemption",
					"The Prestige",
					"The Matrix",
				},
			},
			{
				User: "Peach",
				Items: []too.Item{
					"The Godfather",
					"Inception",
					"Fight Club",
					"WALL·E",
					"Princess Mononoke",
				},
			},
		}, false)
	}
}

func BenchmarkUsingBatchWithAutoUpdate(b *testing.B) {
	te, err := too.New("redis://localhost", "movies")
	if err != nil {
		log.Fatal(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		te.Likes.Batch([]too.BatchRaterOp{
			{
				User: "Sonic",
				Items: []too.Item{
					"The Shawshank Redemption",
					"The Godfather",
					"The Dark Knight",
					"Pulp Fiction",
				},
			},
			{
				User: "Mario",
				Items: []too.Item{
					"The Godfather",
					"The Dark Knight",
					"The Shawshank Redemption",
					"The Prestige",
					"The Matrix",
				},
			},
			{
				User: "Peach",
				Items: []too.Item{
					"The Godfather",
					"Inception",
					"Fight Club",
					"WALL·E",
					"Princess Mononoke",
				},
			},
		}, true)
	}
}
