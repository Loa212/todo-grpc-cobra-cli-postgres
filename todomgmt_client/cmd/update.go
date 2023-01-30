/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"log"
	"strconv"
	"time"

	pb "github.com/Loa212/todo-grpc/todomgmt"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update [id] [title?] [description?] [done?]",
	Short: "update a todo",
	Long: ``,
	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.ParseInt(args[0], 10, 32) 
		if err != nil {
			log.Fatalf("could not convert id to int: %v", err)
		}

		conn, err := grpc.Dial(Address, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()

		c := pb.NewTodoServiceClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()




	
		r, err := c.UpdateTodo(ctx, &pb.Todo{Id: int32(id), Title: args[1], Description: args[2], Done: (args[3] == "true")})

		if err != nil {
			log.Fatalf("could not update todo: %v", err)
		}

log.Printf(`Updated Todo: {
Id: %d,
Title: %s,
Description: %v,
Done: %v
}`, id, r.GetTitle(), r.GetDescription(), r.GetDone())
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
