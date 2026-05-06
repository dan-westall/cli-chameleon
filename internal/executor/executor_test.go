package executor

import "testing"

func TestRunSingleCommand(t *testing.T) {
	r := Run([]string{"echo hello"}, false)
	if r.Err != nil {
		t.Fatal(r.Err)
	}
	if r.Output != "hello\n" {
		t.Fatalf("expected 'hello\\n', got %q", r.Output)
	}
}

func TestRunSequential(t *testing.T) {
	r := Run([]string{"echo one", "echo two"}, false)
	if r.Err != nil {
		t.Fatal(r.Err)
	}
	if r.Output != "one\ntwo\n" {
		t.Fatalf("unexpected output: %q", r.Output)
	}
}

func TestRunParallel(t *testing.T) {
	r := Run([]string{"echo a", "echo b"}, true)
	if r.Err != nil {
		t.Fatal(r.Err)
	}
	// Both should have run (order not guaranteed)
	if len(r.Output) == 0 {
		t.Fatal("expected output from parallel commands")
	}
}

func TestRunFailure(t *testing.T) {
	r := Run([]string{"false"}, false)
	if r.Err == nil {
		t.Fatal("expected error from 'false' command")
	}
}

func TestRunSequentialStopsOnError(t *testing.T) {
	r := Run([]string{"false", "echo should-not-run"}, false)
	if r.Err == nil {
		t.Fatal("expected error")
	}
	if r.Output != "" {
		t.Fatalf("expected no output after failure, got %q", r.Output)
	}
}
