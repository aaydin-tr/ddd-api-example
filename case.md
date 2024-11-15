Backend Case Study

We ask you to implement a ticket allocating and purchasing service using a REST Api. The details are explained in the following chapter. Try to fulfill every requirement but a partial solution is still very useful.

Problem definition
The following three routes need to be built to enable allocating tickets to multiple purchases.

The solution needs to ensure that the allocation does not drop below 0, and the purchased amounts are not greater than the allocation given.

Expect multiple requests to be made against this API concurrently.

Taking payment is out of scope for this problem.

Create Ticket
Create an event with an allocation of tickets available to purchase:

POST /ticketsuser.

Request Body:
{
  "name": "example",
  "desc": "sample description",
  "allocation": 100
}

Response Body:
{
  "id": 1,
  "name": "example",
  "desc": "sample description",
  "allocation": 100
}







Get Ticket

Get ticket by id:

GET /tickets/{id}

(No request body)
Response Body:
{
  "id": 1,
  "name": "example",
  "desc": "sample description",
  "allocation": 100
}



Purchase Ticket
Purchase a quantity of tickets from the allocation of the given tickets:

POST /tickets/{id}/purchases


Request body:
{
  "quantity": 2,
  "user_id": "406c1d05-bbb2-4e94-b183-7d208c2692e1"
}


(No Response body)

A 2xx status code must be returned on success.

A 4xx status code must be returned on any request that attempts to purchase more tickets than are available. In this case, no tickets should be purchased for that request.
P.S.: The user_id is a randomly generated field. For this case study, it is not necessary to create a real user entity and link it to the ticket table. This field is merely a dummy field added for the scenario. 


Languages and frameworks
You need to use golang for writing this service. You can use any framework you want
Database
Postgresql
Additional Requirements
You need to dockerize your app
You need to write Unit Tests to your app
Also open api documentation required
Timeline
You have 4 days, Good luck.


