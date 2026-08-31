central static-ip rendezvous hub for sf.


Basically interacts with both the client and the worker-server


Kinda like this:

CLIENT -----client_conn----------- CENTRAL_SERVER ----- w_conn --------- WORKER


Since, this is a middleman, it houses both the protocols for interacting with
client-side as well as worker-side, tbh they are just plain tcp transfers




Core of central:
At its core, it just has the following structure



   | <<<<<<<<<<< JOB QUEUE >>>>>>>>>>>>>>> |
   |                                       | 
Client Pushed                         Worker is sent
   |									the latest job
   |								       |
Client Blocked till						   |
   <---------send back the response<---After done
Recv result
		
