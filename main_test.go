package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

func TestAddAndListTodosAcrossProcessInvocations(t *testing.T) {
	binary := buildTodo(t)
	database := filepath.Join(t.TempDir(), "todos.json")

	stdout, stderr, err := runTodo(binary, database, "add", "  Buy milk  ")
	if err != nil {
		t.Fatalf("first add failed: %v\nstderr: %s", err, stderr)
	}
	if stdout != "added 1\n" {
		t.Fatalf("first add stdout = %q, want %q", stdout, "added 1\n")
	}
	if stderr != "" {
		t.Fatalf("first add stderr = %q, want empty", stderr)
	}

	stdout, stderr, err = runTodo(binary, database, "add", "Walk dog")
	if err != nil {
		t.Fatalf("second add failed: %v\nstderr: %s", err, stderr)
	}
	if stdout != "added 2\n" {
		t.Fatalf("second add stdout = %q, want %q", stdout, "added 2\n")
	}
	if stderr != "" {
		t.Fatalf("second add stderr = %q, want empty", stderr)
	}

	stdout, stderr, err = runTodo(binary, database, "list")
	if err != nil {
		t.Fatalf("list failed: %v\nstderr: %s", err, stderr)
	}
	if stdout != "1\tactive\tBuy milk\n2\tactive\tWalk dog\n" {
		t.Fatalf("list stdout = %q", stdout)
	}
	if stderr != "" {
		t.Fatalf("list stderr = %q, want empty", stderr)
	}
}

func TestEmptyTitleFailsWithoutModifyingDatabase(t *testing.T) {
	binary := buildTodo(t)
	database := filepath.Join(t.TempDir(), "todos.json")

	_, stderr, err := runTodo(binary, database, "add", "Keep me")
	if err != nil {
		t.Fatalf("seed add failed: %v\nstderr: %s", err, stderr)
	}
	before, err := os.ReadFile(database)
	if err != nil {
		t.Fatalf("read database before invalid add: %v", err)
	}

	stdout, stderr, err := runTodo(binary, database, "add", " \t\n ")
	if err == nil {
		t.Fatal("whitespace-only add succeeded, want non-zero exit")
	}
	if stdout != "" {
		t.Fatalf("invalid add stdout = %q, want empty", stdout)
	}
	if stderr != "title must not be empty\n" {
		t.Fatalf("invalid add stderr = %q, want %q", stderr, "title must not be empty\n")
	}
	after, err := os.ReadFile(database)
	if err != nil {
		t.Fatalf("read database after invalid add: %v", err)
	}
	if !bytes.Equal(after, before) {
		t.Fatalf("invalid add modified database\nbefore: %s\nafter: %s", before, after)
	}
}

func buildTodo(t *testing.T) string {
	t.Helper()

	binary := filepath.Join(t.TempDir(), "todo")
	if runtime.GOOS == "windows" {
		binary += ".exe"
	}
	command := exec.Command("go", "build", "-o", binary, ".")
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("build todo: %v\n%s", err, output)
	}
	return binary
}

func runTodo(binary, database string, args ...string) (string, string, error) {
	command := exec.Command(binary, args...)
	command.Env = append(os.Environ(), "TODO_DB="+database)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr
	err := command.Run()
	return stdout.String(), stderr.String(), err
}
