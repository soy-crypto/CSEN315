# RESTful URL Shortener 

- Due Date: Thursday / Week 4 
- You should not need to import any additional libraries
  - to get started, you may need to read a little on [go modules](https://go.dev/ref/mod)

The task for this assignment is to finish the shortlink portion of the RESTful URL
Shortener we started in class. For full credit, you'll need to: 

1. Complete datastore method and add route to retrieve all links for a user 
2. Complete datastore method and add routes to deleta a link for a user
3. Add relationship links to the new and existing entries such that a client can use those links to replace hardcoded links. 

## Homework Criteria (Total Points Available: 15)

For full credit, appropriate status codes must be sent in the HTTP response and handle
errors. 

| Points | ID          | Test Criteria                                               |
| -----: | ----------- | ----------------------------------------------------------- |
|      5 | USERS_LINKS | Retrieve all links for user, links for non existant user    |
|      5 | DELETE_LINK | Delete a link only if the owner matches the supplied user   |
|      5 | HATEOAS     | Add relationship links to the Create and Get link functions |

## What you don't need to do

- You should only need to changethe existing code to complete the HATEOAS portion of the assignment. But you should be able to understand the existing code sicne it'll help in completing the other tasks. 
- You should not need to update the front end portion, but if you run the server and connect to the front end, the CreateLink should already be working. So if you'd like to create multiple links from there, feel free. 