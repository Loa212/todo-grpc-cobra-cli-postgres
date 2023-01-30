package main

import (
	"context"
	"log"
	"net"
	"os"

	pb "github.com/Loa212/todo-grpc/todomgmt"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

type TodoManagementServer struct {
	conn *pgx.Conn
	pb.UnimplementedTodoServiceServer
}

func (s *TodoManagementServer) CreateTodo(ctx context.Context, in *pb.NewTodo) (*pb.Todo, error) {

	log.Printf("Received: %v", in.GetTitle())

	//create the table if the table does not exist
	//it is not a good practice to create the table in the code
	//and would be better to use the command line
	//but for the sake of simplicity, I will do it here

	createSql := `CREATE TABLE IF NOT EXISTS todos (
		id SERIAL PRIMARY KEY,
		title TEXT NOT NULL,
		description TEXT NOT NULL,
		done BOOLEAN NOT NULL
		);`
		
	_, err := s.conn.Exec(context.Background(), createSql)

	if err != nil {
		log.Fatalf("failed to create table: %v", err)
		os.Exit(1)
	}

	created_todo := &pb.Todo{
		Title: in.GetTitle(),
		Description: in.GetDescription(),
		Done: in.GetDone(),
	}

	tx, err := s.conn.Begin(context.Background())

	if err != nil {
		log.Fatalf("failed to begin transaction: %v", err)
	}

	_, err = tx.Exec(context.Background(), "INSERT INTO todos (title, description, done) VALUES ($1, $2, $3)", created_todo.GetTitle(), created_todo.GetDescription(), created_todo.GetDone())

	if err != nil {
		log.Fatalf("failed to insert todo: %v", err)
	}

	err = tx.Commit(context.Background())

	if err != nil {
		log.Fatalf("failed to commit transaction: %v", err)
	}

	return created_todo, nil
}


func (s *TodoManagementServer) GetTodos(ctx context.Context, in *pb.GetTodosParams) (*pb.TodosList, error) {

	var todos *pb.TodosList = &pb.TodosList{}

	rows, err := s.conn.Query(context.Background(), "SELECT * FROM todos")

	if err != nil {
		log.Fatalf("failed to query todos: %v", err)
	}

	for rows.Next() {
		var todo pb.Todo

		err = rows.Scan(&todo.Id, &todo.Title, &todo.Description, &todo.Done)

		if err != nil {
			return nil, err
		}

		todos.Todos = append(todos.Todos, &todo)
	}

	return todos, nil
}

// update todo
func (s *TodoManagementServer) UpdateTodo(ctx context.Context, in *pb.Todo) (*pb.Todo, error) {
	var todo = pb.Todo {
		Id: in.GetId(),
		Title: in.GetTitle(),
		Description: in.GetDescription(),
		Done: in.GetDone(),
	}

	tx, err := s.conn.Begin(context.Background())

	if err != nil {
		log.Fatalf("failed to begin transaction: %v", err)
	}

	_, err = tx.Exec(context.Background(), `
	UPDATE todos
	SET title = $1, description = $2, done = $3
	WHERE id = $4;`, in.GetTitle(), in.GetDescription(), in.GetDone(), in.GetId())

	if err != nil {
		return nil, err
	}

	return &todo, nil
}


// delete todo
func (s *TodoManagementServer) DeleteTodo(ctx context.Context, in *pb.GetTodoParams) (*pb.GetTodoParams, error) {
	tx, err := s.conn.Begin(context.Background())

	if err != nil {
		log.Fatalf("failed to begin transaction: %v", err)
	}

	_, err = tx.Exec(context.Background(), `DELETE FROM todos WHERE id = $1;`,  in.GetId())

	if err != nil {
		return nil, err
	}

	return in, nil
}


func (server *TodoManagementServer) Run() error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterTodoServiceServer(s, server)
	log.Printf("server listening at %v", lis.Addr())

	return s.Serve(lis)
}

func main() {
	connString := "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"

	conn, err := pgx.Connect(context.Background(), connString)

	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	defer conn.Close(context.Background())

	server := &TodoManagementServer{conn: conn}

	if err := server.Run(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}

}