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

const (
	Address = "localhost:50051"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add [title] [description] [done?]",
	Short: "Adds a todo. If done is not specified, it will be set to false.",
	Long: ``,
	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

		if(len(args) < 3) {
			args = append(args, "false")
		}

		newTodo := &pb.NewTodo{Title: args[0], Description: args[1], Done:(args[2] == "true")}

		conn, err := grpc.Dial(Address, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()

		c := pb.NewTodoServiceClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

	
		r, err := c.CreateTodo(ctx, newTodo)
		if err != nil {
			log.Fatalf("could not create todo: %v", err)
		}
		log.Printf(`Created Todo: {
Title: %s,
Description: %v,
Done: %v
}`, r.GetTitle(), r.GetDescription(), r.GetDone())},

}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
