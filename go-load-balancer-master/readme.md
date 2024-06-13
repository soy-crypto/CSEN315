# Overview

The goal of this assignment is to evaluate different load balancing strategies. You must simulate and display the results of the simulation only using go (and associated libraries).

# Restrictions

- anything in the go standard library is okay
- you may use a charting library such as https://github.com/go-echarts/go-echarts or https://github.com/gonum/plot

# Grading

The grading is based on the completion of the following criteria _and_ your ability
to explain your code. I suggest you leave many comments that explain the what and why of the code
so you're prepared for when I ask you about it. 

You are required to complete each criteria in this assignment. 

| Points | ID     | Test Criteria                                                                        |
| -----: | ------ | ------------------------------------------------------------------------------------ |
|      5 | CHARTS | Your code generates charts (a long with instructions on how to regenerate them)      |
|      5 | PO2    | Create an experiment that pits the "Power of 2 Random Choices" against "Round Robin" |

You can include other load balancing algorithms, as long as you explain the algorithm.  Each additional algorithm and explanation will be 5 extra points. Although the explanation is about distributed load balancing, your implementation can assume a single load balancer.

# Test Criteria

- Charts - a command should be provided in `grading.md` which can be used to generate charts
  - the command will likely be `go run main.go` 
  - show the difference in load between the two algorithms
  - Power of 2 Random Choices should have a much more even distribution than Round Robin
  - For reference, I've included my own demonstration in this file, you can open the `.html` file in your browser
    - Remember your code should be able to generate this _on demand_
