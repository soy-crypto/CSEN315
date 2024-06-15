# Execution Command

The way to execute my code is:  

``` bash
go run main.go
```

# Extra Algorithms

You may implement extra algorithms such as: 
- Least Recently Contacted
- Weighted Round Robin

For each extra algorithm you're implementing, please fill out the following template: 

/******** Round Robin  *******/
- Algorithm Name: Round Robin
- Explanation:
  - Round Robin works by doing fairness with randomness and less reliance on centralized information 
- Comparison with Power of 2 in a distributed environment:
  - Round Robin is worse than Power of 2 in a distrbuted environment because of two reasons.
    Firstly, P2R selects two random servers and chooses the less loaded one. This randomness helps avoid all load balancers converging on the same server, leading to fairer distribution across servers even with varying capacities.
    Secondly, P2R doesn't require a central load balancer to maintain server load information. Each load balancer can make independent decisions based on the two chosen servers, reducing reliance on centralized data.
    

/******** Weighted Round Robin *******/
- Algorithm Name: Weighted Round Robin
- Explanation:
- Weighted Round Robin works by assigning weights to servers in a pool
- Comparison with Power of 2 in a distributed environment:
  Weighted Round Robin might be better or worse than Power of 2. It depends on different scenarios. The details are as follows.
  If you have servers with different processing power or memory limitations and want to distribute workloads accordingly,
  WRR is a good choice.
  If simplicity and adaptability are priorities, and you have a dynamic environment with fluctuating server loads or unknown server capabilities, P2R is a better option.
  If you have a mix of these factors, consider the trade-offs between control, efficiency, and complexity based on your specific needs.


/******** Least Recently Contacted ******/
- Algorithm Name: Least Recently Contacted
- Explanation:
    Least Recently Contacted works by maintaining a limited-size cache and following these principles.
    Firstly, when an item is accessed (e.g., a server in a load balancing context), it's promoted to the most recently used position in the cache.
    Secondly, if the cache reaches its capacity and a new item needs to be added, the least recently contacted item is removed to make space.

- Comparison with Power of 2 in a distributed environment:
    Least Recently Contacted performs worse than Power of 2 in a distrbuted environment because of two reasons.The details are as        follows.
    Firstly, ach load balancer maintains its own LRU cache. A server that was recently used by one load balancer might not be in the cache of another, leading to uneven load distribution.
    Seconldy, f multiple requests arrive at the same time across different load balancers, they might all choose the same recently used server (from their individual caches), causing a bottleneck.



# Grading

The grading is based on the completion of the following criteria _and_ your ability
to explain your code. I suggest you leave many comments that explain the what and why of the code
so you're prepared for when I ask you about it. 

You are required to complete each criteria in this assignment. 

| Points | ID     | Test Criteria                                                                        |
| -----: | ------ | ------------------------------------------------------------------------------------ |
|      5 | CHARTS | Your code generates charts (a long with instructions on how to regenerate them)      |
|      5 | PO2    | Create an experiment that pits the "Power of 2 Random Choices" against "Round Robin" |

You can include other load balancing algorithms, as long as you explain the algorithm. Each additional algorithm and explanation will be 5 extra points. Although the explanation is about distributed load balancing, your implementation can assume a single load balancer.