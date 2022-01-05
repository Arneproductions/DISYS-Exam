# Distributed Systems - Exam
Name:       Andreas Nicolaj Tietgen

Date: 05/01/22

Course code: BSDISYS1KU

I hereby declare that this submission was created in its entirety by me and only me.

## 1. Questions

1. Question 
   * Answer: 1

2. Question 
   * Answer: 2

3. Question 
   * Answer: 3

4. Question 
   * Answer: 1

5. Question 
   * Answer: 2

6. Question 
   * Answer: 1

7. Question 
   * Answer: 2

8. Question 
   * Answer: 3

9. Question 
   * Answer: 2

10. Question 
    * Answer: 1



## 2. Implementation
I have in my implementation taken some inspiration from my [Mock exam](https://github.com/Arneproductions/DISYS-Passive-Replication) attemt which i have created together with ABSO. I have also reused the parts about RPC from our [mandatory exercise 1](https://github.com/Arneproductions/DISYS-Exercise-1). At last but not least i have reused the same idea of creating an auto client when running it through Docker from our [mandatory exercise 3](https://github.com/Arneproductions/DISYS-Mini-Project-03/blob/master/cmd/client/main.go)

## 3. Implementation discussion

### 1. Question
The system model i am assuming is a model with crash failure and an asynchronous network. I assume that when a replica has crashed then it wont restart or recover. There is no upper bound delay since the network is asynchronous.

### 2. Question
Since the system has to tolerate a crash failure then we have to use replication. For that reason i have choosen to implement passive replication with leader election. Passive replication contains a leader and n amount of replica that the leader updates for each `Put()` request. If a replica does not get the `Put()` request from the leader, it then forwards the request to the leader and then replies to the client when the leader is done. 

The reliability property is quaranteed with the system having two nodes. However to show that the election works, i have then choosen to create 3 nodes. Then we can still be able to write and keep a backup if one of the nodes dies.

### 3. Question
How do i handle crash failures? Since i use passive replication then we have to tolerate if either a replica or a leader dies:
- **Replica dies:** The leader noticies that the replica does not respond to the heartbeat. It then removes the node from its list and notifies the other replicas that the replica is dead by providing a new list of nodes through the heartbeat.
- **Leader dies:** The replica noticies that it haven't recieved a heartbeat for a while and now wants an election. Therefor it sends a request for an election and recieves the others processId and compares it to itself. The node with the highest process id is the winner. That is being declared by the ElectedRequest. The replica which is being declared the leader gets it values replicated to the other replicas who is alive.

Thus the system can recover from a leader or a replica crash.

### 4. Question
The passive replication implementation is using a synchronous style of replication. This ensures that when a client wants to call `Put()` then the leader ensures that it has recieved an acknowledgement from the replicas before replying that the put was succesful. 

When the system wants to update the value with the `Put()` operation, and the call returns successful, then a get request on the same key can be done in two ways:
- **Get request to leader:** If the leader gets the `Get()` request then it reads the value after the Read lock has been taken. The leader is the one responsible for updating the value. Therefor it has the update value for that particular key.
- **Get request to replica:** If the replica gets the `Get()` operation then it calls `Get()` on the leader. By the explanation of what happens when the leader recieves the `Get()` request then it can quarantee that the update value is the same as the newest `Put()`.

### 5. Question
The initialization of the hashmap is done by GoLang´s own implementation called `Map`. By the documentiaton of the `Map` in GoLang, *"If key is not in the map, then a elem is the zero value for the map´s element type"*([maps - A tour of go](https://go.dev/tour/moretypes/22)). The zero value of an int is zero. 

Therefor at the initialization of the map all values that could possibly be called by the `Get()` operation is set to zero by default.

### 6. Question
Since the passive replication is using a synchronous style of replication then all replicas has to acknowledge to the leader if the value is updated. Only after the replicas has acknowledged then the leader replies with `Success: true`.

If a replica is crashing while the leader is updating then the leader wont get a response and declare it dead. If the leader recieves a false on the update then the leader will declare that replica dead. 
Only replicas that is up-to-date is being kept alive. 

**How do we get the value?**

If the leader recieves the `Get()` request then it will send its value back. The leader is the node that is being kept up-to-date since it handles all the ´Put()´ requests.
If the replica recieves the `Get()` request then it will ask the leader for the value, thus the replica will also send back the same value as the leader.

Thus every node alive node returns the same value.


### 7. Question
Since we use Leader election then we have the following scenarios:
- **Everything is fine:** then every call should work as expected. A `Put()` request returns true when the value is updated and the `Get()` request returns the value that the client wants to recieve.
- **An election has started:** If an election has started then false will be returned by all the nodes that does not have a status of either *Replica* or *Leader*. The `Get()` operation will return an error while it is not in a state of either *Replica* or *Leader*. But when the election is over then the operations will function as normal.

**How does the election work?**

The election work by a node issuing an election. It sends out an election request to every node where every node replies with its `processId` and `value`. The node with the highest `processId` wins. When the election is over, an `Elected` request is being sent to all replicas with the `value` and `ip` from the winning leader.

If there is only one node left then it will elect itself as a leader and only update itself.

So for that reason the liveness property is fullfilled even if there is an election. The system will respond with either and error or a proper response no matter what.

### 8. Question
The system uses passive replication with leader election thus by at least having to nodes(one leader and one replica) in the system then the system quarantees the Reliability property. 

How the system handles crash failures is the same explanation as the answer in the **3. Question**.


