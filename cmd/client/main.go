package main

import(
	"flag"
	"log"
	"context"
	"google.golang.org/grpc"
	proto "exam/proto"
)

var (
	serverAddr = flag.String("serverAddr", "localhost:5001", "Server to connect to")
)

func main(){
	flag.Parse()

	conn, err := grpc.Dial(*serverAddr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Printf("Could not connect: %v\n", err)
	}
	defer conn.Close()

	c := proto.NewReplicationClient(conn)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	put(c, ctx)
	
}

func put(c proto.ReplicationClient, ctx context.Context) (bool,error) {
	reply, err := c.Put(ctx, &proto.PutMessage{Key: 1, Value: 10})

	if err != nil {
		log.Printf("Failed to increment: %v\n", err)
		return false,err
	}

	return reply.Success, nil
}