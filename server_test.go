package poker

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestGETPlayer(t *testing.T) {
	playerStore := StubPlayerStore{
		map[string]int{
			"John":  20,
			"Maria": 30,
		},
		[]string{},
		League{},
	}

	server, err := NewPlayerServer(&playerStore, &GameSpy{})
	if err != nil {
		t.Fatalf("didn't expected error when creating a server but got one, %v", err)
	}

	t.Run("returns John's score", func(t *testing.T) {
		request := newGetScoreRequest("John")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response, http.StatusOK)
		assertResponseBody(t, response.Body.String(), "20")
	})

	t.Run("returns Maria's score", func(t *testing.T) {
		request := newGetScoreRequest("Maria")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response, http.StatusOK)
		assertResponseBody(t, response.Body.String(), "30")
	})

	t.Run("returns 404 for missing players", func(t *testing.T) {
		request := newGetScoreRequest("Missing Player")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response, http.StatusNotFound)
	})
}

func TestStoreWins(t *testing.T) {
	playerStore := &StubPlayerStore{
		map[string]int{
			"John":  20,
			"Maria": 30,
		},
		[]string{},
		League{},
	}

	server, err := NewPlayerServer(playerStore, &GameSpy{})
	if err != nil {
		t.Fatalf("didn't expected error when creating a server but got one, %v", err)
	}

	t.Run("record win for POST", func(t *testing.T) {
		player := "John"
		request := newPostWinRequest(player)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response, http.StatusAccepted)

		AssertPlayerWin(t, playerStore, player)
	})
}

func TestLeague(t *testing.T) {

	t.Run("it returns league table as JSON", func(t *testing.T) {
		wantedLeague := League{
			{"Leo", 32},
			{"John", 40},
			{"Maria", 12},
		}

		store := StubPlayerStore{nil, nil, wantedLeague}
		server, err := NewPlayerServer(&store, &GameSpy{})
		if err != nil {
			t.Fatalf("didn't expected error when creating a server but got one, %v", err)
		}

		request, _ := http.NewRequest(http.MethodGet, "/league", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := getLeagueFromResponse(t, response.Body)
		assertStatus(t, response, http.StatusOK)
		assertContentType(t, response, jsonContentType)
		assertLeague(t, got, wantedLeague)

	})
}

func TestGame(t *testing.T) {
	t.Run("GET /game returns 200", func(t *testing.T) {
		server, err := NewPlayerServer(&StubPlayerStore{}, &GameSpy{})
		if err != nil {
			t.Fatalf("didn't expected error when creating a server but got one, %v", err)
		}

		request := newGameRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response, http.StatusOK)
	})

	t.Run("when we get a message over a websocket it is a winner of a game", func(t *testing.T) {
		game := &GameSpy{}
		store := &StubPlayerStore{}
		winner := "Joe"
		playerServer, err := NewPlayerServer(store, game)
		if err != nil {
			t.Fatalf("didn't expected error when creating a server but got one, %v", err)
		}

		server := httptest.NewServer(playerServer)
		defer server.Close()

		wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"

		ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		if err != nil {
			t.Fatalf("could not open a ws connection on %s %v", wsURL, err)
		}
		defer ws.Close()
		n := "5"
		if err := ws.WriteMessage(websocket.TextMessage, []byte(n)); err != nil {
			t.Fatalf("could not send message over ws connection %v", err)
		}
		if err := ws.WriteMessage(websocket.TextMessage, []byte(winner)); err != nil {
			t.Fatalf("could not send message over ws connection %v", err)
		}

		time.Sleep(10 * time.Millisecond)
		assertGameStartedWith(t, game, 5)
		assertFinishCalledWith(t, game, winner)
	})
}

func assertGameStartedWith(t testing.TB, game *GameSpy, numberOfPlayersWanted int) {
	t.Helper()
	if game.StartCalledWith != numberOfPlayersWanted {
		t.Errorf("wanted Start called with %d but got %d", numberOfPlayersWanted, game.StartCalledWith)
	}
}
func assertFinishCalledWith(t testing.TB, game *GameSpy, winner string) {
	t.Helper()
	if game.FinishCalledWith != winner {
		t.Errorf("expected finish called with %q but got %q", winner, game.FinishCalledWith)
	}
}

func newGameRequest() *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/game", nil)
	return request
}

func assertContentType(t testing.TB, response *httptest.ResponseRecorder, want string) {
	t.Helper()
	if response.Result().Header.Get("content-type") != want {
		t.Errorf("response did not have content-type of %s, got %v", want, response.Result().Header)
	}
}

func assertLeague(t testing.TB, got, want []Player) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func assertStatus(t testing.TB, got *httptest.ResponseRecorder, want int) {
	t.Helper()
	if got.Result().StatusCode != want {
		t.Errorf("did not get correct status, got %d, want %d", got.Result().StatusCode, want)
	}
}

func getLeagueFromResponse(t testing.TB, body io.Reader) (league League) {
	t.Helper()
	err := json.NewDecoder(body).Decode(&league)

	if err != nil {
		t.Fatalf("Unable to parse response from server %q into slice of Player, '%v'", body, err)
	}

	return
}

func assertResponseBody(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response body is wrong, got %q want %q", got, want)
	}
}

func newLeagueRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/league", nil)
	return req
}

func newGetScoreRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/players/"+name, nil)
	return req
}

func newPostWinRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/players/%s", name), nil)
	return req
}
