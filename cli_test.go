package poker_test

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/jaymorelli96/game-tracker"
)

var dummyBlindAlerter = &poker.SpyBlindAlerter{}
var dummyPlayerStore = &poker.StubPlayerStore{}
var dummyStdOut = &bytes.Buffer{}

func TestCLI(t *testing.T) {
	t.Run("prompts the user to enter the number of players and starts the game", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		in := strings.NewReader("10\n")
		game := &poker.GameSpy{}

		cli := poker.NewCLI(in, stdout, game)
		cli.PlayPoker()

		gotPrompt := stdout.String()
		wantPrompt := poker.PlayerPrompt

		if gotPrompt != wantPrompt {
			t.Errorf("got %q, want %q", gotPrompt, wantPrompt)
		}

		if game.StartCalledWith != 10 {
			t.Errorf("wanted Start called with 10 but got %d", game.StartCalledWith)
		}
	})

	t.Run("finish game with 'Jay' as winner", func(t *testing.T) {
		in := strings.NewReader("1\nJay wins\n")
		game := &poker.GameSpy{}
		cli := poker.NewCLI(in, dummyStdOut, game)

		cli.PlayPoker()

		if game.FinishCalledWith != "Jay" {
			t.Errorf("expected finish called with 'Jay' but got %q", game.FinishCalledWith)
		}
	})

	t.Run("record 'John' win from user input", func(t *testing.T) {
		in := strings.NewReader("1\nJohn wins\n")
		game := &poker.GameSpy{}
		cli := poker.NewCLI(in, dummyStdOut, game)

		cli.PlayPoker()

		if game.FinishCalledWith != "John" {
			t.Errorf("expected finish called with 'John' but got %q", game.FinishCalledWith)
		}
	})

	t.Run("do not read beyond the first newline", func(t *testing.T) {
		in := failOnEndReader{
			t,
			strings.NewReader("1\nJay wins\n hello there"),
		}

		game := poker.NewTexasHoldem(dummyBlindAlerter, dummyPlayerStore)

		cli := poker.NewCLI(in, dummyStdOut, game)
		cli.PlayPoker()
	})
}

type failOnEndReader struct {
	t   *testing.T
	rdr io.Reader
}

func (m failOnEndReader) Read(p []byte) (n int, err error) {

	n, err = m.rdr.Read(p)

	if n == 0 || err == io.EOF {
		m.t.Fatal("Read to the end when you shouldn't have")
	}

	return n, err
}

func assertScheduledAlert(t testing.TB, got, want poker.ScheduledAlert) {
	amountGot := got.Amount
	if amountGot != want.Amount {
		t.Errorf("got amount %d, want %d", amountGot, want.Amount)
	}

	gotScheduledTime := got.At
	if gotScheduledTime != want.At {
		t.Errorf("got scheduled time of %v, want %v", gotScheduledTime, want.At)
	}

}
