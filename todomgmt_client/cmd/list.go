/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"log"
	"time"

	pb "github.com/Loa212/todo-grpc/todomgmt"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "this lists the todos",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {

		conn, err := grpc.Dial(Address, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()

		c := pb.NewTodoServiceClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

	
		r, err := c.GetTodos(ctx, &pb.GetTodosParams{})

		if err != nil {
			log.Fatalf("could not get todos: %v", err)
		}

		for _, todo := range r.GetTodos() {
			log.Printf(`Todo: {
Id: %d,
Title: %s,
Description: %v,
Done: %v
}`, todo.GetId(), todo.GetTitle(), todo.GetDescription(), todo.GetDone())
		}
		

	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
