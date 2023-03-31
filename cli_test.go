package poker_test

import (
	"github.com/jaymorelli96/game-tracker"
	"strings"
	"testing"
)

func TestCLI(t *testing.T) {
	t.Run("assert win for Jay", func(t *testing.T) {
		in := strings.NewReader("Jay wins\n")
		playerStore := &poker.StubPlayerStore{}

		cli := poker.NewCLI(playerStore, in)
		cli.PlayPoker()

		poker.AssertPlayerWin(t, playerStore, "Jay")
	})

	t.Run("assert win for Joe", func(t *testing.T) {
		in := strings.NewReader("Joe wins\n")
		playerStore := &poker.StubPlayerStore{}

		cli := poker.NewCLI(playerStore, in)
		cli.PlayPoker()

		poker.AssertPlayerWin(t, playerStore, "Joe")
	})
}
