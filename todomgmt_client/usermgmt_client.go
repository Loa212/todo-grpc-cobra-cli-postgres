package main

import "github.com/Loa212/todo-grpc/todomgmt_client/cmd"

// const (
// 	address = "localhost:50051"
// )

func main() {
	cmd.Execute()

// 	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
// 	if err != nil {
// 		log.Fatalf("did not connect: %v", err)
// 	}
// 	defer conn.Close()

// 	c := pb.NewTodoServiceClient(conn)

// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
// 	defer cancel()

// 	var new_todos = make(map[string]string)
// 	new_todos["Todo 1"] = "Giulia puzza"
// 	new_todos["Todo 2"] = "Miao miao miao"

// 	for title, description := range new_todos {
// 		r, err := c.CreateTodo(ctx, &pb.NewTodo{Title: title, Description: description})
// 		if err != nil {
// 			log.Fatalf("could not create todo: %v", err)
// 		}
// 		log.Printf(`Created Todo: [
// Title: %s
// Description: %v
// ]`, r.GetTitle(), r.GetDescription())
// 	}

}