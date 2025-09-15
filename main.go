package main

import (
	"bufio"
	"fmt"
	"os"
)

var taskLists = make(map[string][]string)

func main() {
	fmt.Println("Welcome to Golang To-Do List!")

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("\nChoose an option:")
		fmt.Println("1. Create a new task list")
		fmt.Println("2. Add task to a list")
		fmt.Println("3. View task lists")
		fmt.Println("4. Exit")

		fmt.Print("Enter choice: ")
		scanner.Scan()
		choice := scanner.Text()

		switch choice {
		case "1":
			createTaskList(scanner)
		case "2":
			addTasks(scanner)
		case "3":
			printTasks()
		case "4":
			fmt.Println("Exited")
			return
		default:
			fmt.Println("Invalid choice, try again.")
		}
	}
}

func createTaskList(scanner *bufio.Scanner) {
	fmt.Print("Enter name of new task list: ")
	scanner.Scan()
	listName := scanner.Text()

	if _, exists := taskLists[listName]; exists {
		fmt.Println("Task list already exists.")
	} else {
		taskLists[listName] = []string{}
		fmt.Println("Task list created:", listName)
	}
}

func addTasks(scanner *bufio.Scanner) {
	if len(taskLists) == 0 {
		fmt.Println("No task lists available. Create one first!")
		return
	}

	fmt.Println("Available task lists:")
	for name := range taskLists {
		fmt.Println("-", name)
	}

	fmt.Print("Enter task list to add task to: ")
	scanner.Scan()
	listName := scanner.Text()

	tasks, exists := taskLists[listName]
	if !exists {
		fmt.Println("Task list does not exist.")
		return
	}

	fmt.Print("Enter task: ")
	scanner.Scan()
	task := scanner.Text()

	taskLists[listName] = append(tasks, task)
	fmt.Println("Task added to", listName)
}

func printTasks() {
	if len(taskLists) == 0 {
		fmt.Println("No task lists found.")
		return
	}

	fmt.Println("\nAll Task Lists:")
	for name, tasks := range taskLists {
		fmt.Println("List:", name)
		if len(tasks) == 0 {
			fmt.Println("  (No tasks)")
		} else {
			for i, task := range tasks {
				fmt.Printf("  %d. %s\n", i+1, task)
			}
		}
	}
}
