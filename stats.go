package main

import (
	"fmt"
	"sort"
	"time"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

const (
	daysInLastSixMonths = 183
	weeksInLastSixMonths = 26
)

type column []int

/**
stats generates a graph of git contributions.
*/
func Stats(e string, r []string) {
	if len(r) == 0 {
		println("no commits to explore")
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

func print(c map[int]int){
	keys := getKeys(c)
	sort.Ints(keys)
	cols := buildColumns(keys, c)
	printCells(cols)
}

func getKeys(m map[int]int) []int {
	var keys []int
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func buildColumns(keys []int, commits map[int]int) map[int]column {
	cols := make(map[int]column)
	col := column{}

	for _, k := range keys {
		week := k / 7
		dayinweek := k % 7

		if dayinweek == 0 { //reset
			col = column{}
		}

		col = append(col, commits[k])

		if dayinweek == 6 {
			cols[week] = col
		}
	}

	return cols
}


func printCells(cols map[int]column) {
	printMonths()
	for j := 6; j >= 0; j-- {
		for i := weeksInLastSixMonths + 1; i >= 0; i-- {
			if i == weeksInLastSixMonths+1 {
				printDayCol(j)
			}
			if col, ok := cols[i]; ok {
				//special case today
				if i == 0 && j == calcOffset()-1 {
					printCell(col[j], true)
					continue
				} else {
					if len(col) > j {
						printCell(col[j], false)
						continue
					}
				}
			}
			printCell(0, false)
		}
		fmt.Printf("\n")
	}
}

func printMonths() {
	week := getBeginningOfDay(time.Now()).Add(-(daysInLastSixMonths * time.Hour * 24))
	month := week.Month()
	fmt.Printf("         ")
	for {
		if week.Month() != month {
			fmt.Printf("%s ", week.Month().String()[:3])
			month = week.Month()
		} else {
			fmt.Printf("    ")
		}

		week = week.Add(7 * time.Hour * 24)
		if week.After(time.Now()) {
			break
		}
	}
	fmt.Printf("\n")
}

func getBeginningOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	startOfDay := time.Date(year, month, day, 0, 0, 0, 0, t.Location())
	return startOfDay
}


func printDayCol(day int) {
	out := "     "
	switch day {
	case 1:
		out = " Mon "
	case 3:
		out = " Wed "
	case 5:
		out = " Fri "
	}

	fmt.Printf(out)
}

func printCell(val int, today bool) {
	escape := "\033[0;37;30m"
	switch {
	case val > 0 && val < 5:
		escape = "\033[1;30;42m"
	case val >= 5 && val < 10:
		escape = "\033[1;30;43m"
	case val >= 10:
		escape = "\033[1;30;42m"
	}

	if today {
		escape = "\033[1;37;45m"
	}

	if val == 0 {
		fmt.Printf(escape + "  - " + "\033[0m")
		return
	}

	str := "  %d "
	switch {
	case val >= 10:
		str = " %d "
	case val >= 100:
		str = "%d "
	}

	fmt.Printf(escape+str+"\033[0m", val)
}

