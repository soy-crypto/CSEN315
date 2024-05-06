# Overview

The goal of this assignment is to create a `recursive caching DNS resolver`. There is 
a lot of freely available code online that you can _take inspiration from_. Just 
remember that _you_ must be able to explain all your code in person to the professor. 
I will ask not just the what, but also the why your code is structured as is. 

# Restrictions

- go standard library except system calls which invoke an OS'es stub resolver like `gethostbyname`
- do not change the _logic_ of the main function although you may choose to move out the `RecordType` to a new package

# Grading

The grading is based on the completion of the following criteria _and_ your ability
to explain your code. I suggest you leave many comments that explain the what and why of the code
so you're prepared for when I ask you about it. 

Remember that you are _not_ required to complete each criteria in this assignment.

| Points | ID          | Test Criteria                                                                                 |
| -----: | ----------- | --------------------------------------------------------------------------------------------- |
|     20 | A_RECORD    | Print the IPv4 address of each site if it exists                                              |
|      5 | AAAA_RECORD | Print the IPv6 address of each site if it exists                                              |
|      5 | NS_RECORD   | Print the domain name and IP address of an NS of each site                                    |
|      5 | CNAME       | If the given name is a CNAME, resolves the IPv4 of the canonical name. If not, prints nothing |
|      5 | TXT         | When a TXT record is requested, return a random number in the given range                     |
|      5 | CACHING     | Cache all response records for as long as the TTL specifies and serve back cached records     |

# Submission

- Include all relevant code, any go.mod file, and a copy of the `grading.md` file into a zip on Camino.
- Mark the rows of functionality you'd like me to test

# Expected Functionality

For each output, if multiple records are served, you only need to send _one of_ 

By default, an A record is given 

``` bash
> go run main.go x.com facebook.com google.com
x.com,104.244.42.129
facebook.com,157.240.22.35
google.com,142.251.214.142

> go run main.go -t AAAA x.com facebook.com google.com
x.com,   # empty if x.com does not have an AAAA record
facebook.com,2a03:2880:f131:83:face:b00c:0:25de
google.com,2607:f8b0:4005:80d::200e # sometimes IPv6 records don't display sections that are all 0's

> go run main.go -t NS x.com facebook.com google.com
x.com,c.r10.twtrdns.net,205.251.194.151
facebook.com,b.ns.facebook.com,185.89.219.12
google.com,ns1.google.com,216.239.32.10

> go run main.go -t CNAME facebook.com 
facebook.com, # empty because facebook.com is not a CNAME

> go run main.go -t CNAME www.facebook.com
www.facebook.com,star-mini.c10r.facebook.com,157.240.22.35
```

# Help

The best resource for this is [Implement DNS in a Weekend](https://implement-dns.wizardzines.com/) from 
Julia Evans. She goes through the implementation of a DNS resolver in one weekend in Python. There are
some hairy corners relating to Python. 

Another good resource in Go (if you're struggling) is [this library](https://pkg.go.dev/github.com/miekg/dns)
which has implemented much more than you're required to. 

For socket programming, I suggest using [Network Programming with Go](https://encore.scu.edu/iii/encore/record/C__Rb3797465__S%22Network%20Programming%20with%20Go%22__Orightresult__U__X7?lang=eng&suite=def) available through the university.
If the link doesn't work for you, you can directly search the ISBN (9781098128890) in the [library 
catalog](scu.edu/library) and access through O'Reilly.