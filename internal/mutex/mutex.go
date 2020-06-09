package mutex

// TODO mutex
/*
The algorithm is based upon logical clocks – Lamport clocks to be precise. We will use the Lamport clock only for the algorithm and it will have no further clock (in the sense of periodically increments) and no other interaction besides the mutex algorithm shall touch the Lamport clock.
A client has to know three states:
• released: he does not hold a lock and does not want a lock
• wanting: he does not hold a lock but wants it
• held: he does hold a lock and may enter the critical section
There are two messages in the algorithm:
• request: the sending client tells the others his desire to enter the critical section
• reply-ok: message to tell the other client that he is okay with the previous request
The request is send to ALL clients including the requesting client himself. The client may only enter the critical section when he received a reply-ok from ALL clients (including himself). The reply-ok is only send by a client when he is inactive (not wanting the mutex at all) or does want the mutex, but the received request has a lower logical clock than his own request. On a draw the clientid (and in our case the username) tips the tide and the lower
BAI5-VSP
Praktikum Verteilte Systeme – Aufgabenblatt 3
AZI/KSS
SoSe 20
Step 3: le grand final!
2/2
clientid wins. If the client already holds a lock or wants a lock and his request has a lower logical clock than the receiving one, he stores the request and continues with his work until he finishes his critical section and sending the reply-ok to the stored requests afterwards.
And it goes like this with 3 clients (A,B,C):
• Client A wants to enter the critical section
• A sends request with his clock to A,B,C
• B is currently in the critical section, does store the request
• C is idle and sends reply-ok
• A sends himself an reply-ok
• C wants to enter the critical section & sends request to A, B, C
• A waits for the mutex and his request has a lower clock, therefore stores the request
• B is in the critical section, therefore stores the request
• B finishes his critical section
• B sends reply-ok to the stored requests of A and C
• A got all required reply-ok and may now enter the critical section
• C still waits.
• A has finished his critical section and sends reply-ok to the stored request of C
• C got all required reply-ok and may now enter the critical section
Remember that you have to increase your clock every time you send AND receive an message.
Again we don’t care about message overhead and just use json for convenience.
{
"msg":"<the message: request or reply-ok>",
"time": <int, the lamport clock>,
"reply":"<url to the endpoint where responses shall be send>", "user":"<url to user sending this message>"
}
We shield our self against faulty clients and introduce an extra entpoint “mutexstate” where we can ask for the current state of the client. The answer must be:
{
"state":"<current state: released, wanting, held>", "time": <int, the current lamport clock>
}
If a client does not sent his reply-ok in a reasonable time, one may check back with his state. If he is wanting or held, you just have to wait some more. If he does either not respond at all or does respond with a state of released, one may assume he forgot to send his reply-ok. This is to avoid complete lockups.
The algorithm is to be implemented and shall work automatically when required. Critical sections are marked with critical_section:True within the previous response.
Please test your algorithm with unittest while in development to not lock other clients. On requests, you have to consider all adventurers that claim the capability “mutex”.
 */