package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/google/go-github/v43/github"
)

func main() {
	fmt.Println(selective("Which?", []string{"hogehoge", "fugafuga", "piyopiyo"}))
}

func listenSource() {
	choices := []string{"GitHub"}
	i := selective("Select source", choices)
	// GitHub
	if i == 0 {
		repo_name := question("Enter repository name (e.g. flat35hd99/ig)")
		// repository service -> ListReleases
		client := github.NewClient(nil)
		release, res, err := github.RepositoriesServiceGetRelease(context.Background(), "flat35hd99", "ig", 0)
		if err != nil {
			printf("hogehoge")
		}
	}
}

// Show question and return answer
func question(q string) string {
	fmt.Println("> ", q)

	var result string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		result = scanner.Text()
		break
	}
	return result
}

/*
Show question and choices and return index choiced
*/
func selective(q string, choices []string) int {
	fmt.Println("> ", q)
	for i, choice := range choices {
		fmt.Printf("%d: %s\n", i, choice)
	}

	var result_index int
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()

		// Check string inputs
		is_in, index := contain(choices, input)
		if is_in {
			result_index = index
			break
		}

		// Check index inputs
		int_input, err := strconv.Atoi(input)
		if err == nil {
			result_index = int_input
			break
		} else {
			fmt.Println("cannot recognize")
		}
	}
	return result_index
}

// Return bool and index
// If it does not contain, return -1 as index
func contain(list []string, subject string) (bool, int) {
	for i, v := range list {
		if v == subject {
			return true, i
		}
	}
	return false, -1
}
