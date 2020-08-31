package main

import (
	"time"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

const daysInLastSixMonths = 183

/**
stats generates a graph of git contributions.
*/
func Stats(e string, r []string) {
	if len(r) == 0 {
		print("no commits to explore")
		return
	}
	commits := getCommits(e, r)
	print(commits)
}

func getCommits(e string, r []string) map[int]int {
	c := make(map[int]int, daysInLastSixMonths)
	for i := daysInLastSixMonths; i > 0; i-- {
		c[i] =0
	}
	for _, path := range r {
		c = fillCommits(e, path, c)
	}

	return c
}

/**
For a given repo path, get the commits related to an email
and update the map of commits with the quantity,
 */
func fillCommits(em string, rp string, cm map[int]int) map[int]int {
	repo, err := git.PlainOpen(rp)
	if err != nil {
		panic(err)
	}

	ref, err := repo.Head()
	if err != nil {
		panic(err)
	}

	iterator, err := repo.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		panic(err)
	}

	offset := calcOffset()
	err = iterator.ForEach(func(c *object.Commit) error {
		daysAgo := daysSinceDate(c.Author.When) + offset
		if c.Author.Email != em {
			return nil
		}
		if daysAgo < daysInLastSixMonths {
			cm[daysAgo]++
		}
		return nil
	})

	if err != nil {
		panic(err)
	}

	return cm
}

func daysSinceDate(when time.Time) int {
	return int(time.Since(when).Hours() / 24)
}

/**
Place the commit in its correct place.
 */
func calcOffset() int {
	var offset int
	weekday := time.Now().Weekday()

	switch weekday {
	case time.Sunday:
		offset = 7
	case time.Monday:
		offset = 6
	case time.Tuesday:
		offset = 5
	case time.Wednesday:
		offset = 4
	case time.Thursday:
		offset = 3
	case time.Friday:
		offset = 2
	case time.Saturday:
		offset = 1
	}

	return offset
}