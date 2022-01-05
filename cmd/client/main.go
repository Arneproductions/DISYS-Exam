package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	grpcUtil "exam/internal/grpc"
	proto "exam/proto"

	"google.golang.org/grpc"
)

var (
	serverAddr = flag.String("serverAddr", "localhost:5001", "Server to connect to")
	random     = flag.Bool("random", false, "Randomly send data")
	nodeIndex  = 0
	nodes      []string
)

func main() {
	flag.Parse()

	nodes = strings.Split(*serverAddr, ",")

	if *random {
		autoInteract()
	} else {
		interact()
	}
}

func autoInteract() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var key, value int64
	key = r.Int63n(10)
	value = r.Int63n(10)

	if valueResult, ok := get(key); ok {
		log.Printf("Value at key %d is: %d", key, valueResult)
	} else {
		log.Printf("Could not get value at this time... try again")
	}

	for {
		if ok := put(key, value); !ok {
			continue
		}

		if valueResult, ok := get(key); ok {
			log.Printf("Value at key %d is: %d", key, valueResult)
		} else {
			log.Printf("Could not get value at this time... try again")
		}
		time.Sleep(time.Duration(r.Intn(10)) * time.Second) // Sleep between 0 and 10 seconds

		key = r.Int63n(10)
		value = r.Int63n(10)
	}
}

func interact() {
	for {
		var choice int32
		var key, value int64
		log.Printf("Do you want to get value by key (1) or put (2): ")
		if _, err := fmt.Scanf("%d\n", &choice); err != nil {
			log.Printf("Invalid input, dying: %v\n", err)
		}

		switch choice {
		case 1:
			log.Printf("What key do you want to get value from?\n> ")

			if _, err := fmt.Scanf("%d\n", &key); err != nil {
				log.Printf("Invalid input, dying: %v\n", err)
			}

			if valueResult, ok := get(key); ok {
				log.Printf("Value at key %d is: %d", key, valueResult)
			} else {
				log.Printf("Could not get value at this time... try again")
			}

			break
		case 2:
			log.Printf("What key do you want to put?\n> ")

			if _, err := fmt.Scanf("%d\n", &key); err != nil {
				log.Printf("Invalid input, dying: %v\n", err)
			}

			log.Printf("What value to you want to put at %d?\n> ", key)

			if _, err := fmt.Scanf("%d\n", &value); err != nil {
				log.Printf("Invalid input, dying: %v\n", err)
			}

			result := put(key, value)
			log.Printf("Put was successful: %v", result)

			break
		default:
			log.Printf("Invalid choice\n")
			return
		}

	}
}

func put(key int64, value int64) bool {

	log.Printf("Put value %d at key %d\n", value, key)
	log.Printf("Dialing %s\n", *serverAddr)

	conn, err := grpc.Dial(nodes[nodeIndex], grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Printf("Could not connect: %v\n", err)
		nextNode()
	}
	defer conn.Close()

	c := proto.NewReplicationClient(conn)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	reply, err := c.Put(ctx, &proto.PutMessage{
		Key:   key,
		Value: value,
	})

	if err != nil {
		log.Printf("Failed to put: %v\n", err)
		nextNode()
		return false
	}

	log.Printf("Successfully put %d at key %d\n", value, key)

	return reply.GetSuccess()
}

func get(key int64) (int64, bool) {

	log.Printf("Dialing %s\n", *serverAddr)

	conn, err := grpc.Dial(nodes[nodeIndex], grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Printf("Could not connect: %v\n", err)
		nextNode()
	}
	defer conn.Close()

	c := proto.NewReplicationClient(conn)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	reply, err := c.Get(ctx, &proto.GetMessage{Key: key})
	if err != nil {
		if !grpcUtil.IsAnElectionError(err) {

			nextNode()
		}

		return 0, false
	}

	return reply.GetValue(), true
}

func nextNode() {
	nodeIndex++
	if nodeIndex == len(nodes) {
		log.Fatalln("There are no more nodes left...")
	}
}
