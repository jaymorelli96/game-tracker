package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type StubPlayerStore struct {
	scores   map[string]int
	winCalls []string
	league   League
}

func (ps *StubPlayerStore) GetPlayerScore(name string) int {
	score, ok := ps.scores[name]
	if !ok {
		return 0
	}

	return score
}

func (ps *StubPlayerStore) RecordWin(name string) {
	ps.winCalls = append(ps.winCalls, name)
}

func (ps *StubPlayerStore) GetLeague() League {
	return ps.league
}

func TestGETPlayer(t *testing.T) {
	playerStore := StubPlayerStore{
		map[string]int{
			"John":  20,
			"Maria": 30,
		},
		[]string{},
		League{},
	}

	server := NewPlayerServer(&playerStore)

	t.Run("returns John's score", func(t *testing.T) {
		request := newGetScoreRequest("John")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), "20")
	})

	t.Run("returns Maria's score", func(t *testing.T) {
		request := newGetScoreRequest("Maria")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), "30")
	})

	t.Run("returns 404 for missing players", func(t *testing.T) {
		request := newGetScoreRequest("Missing Player")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusNotFound)
	})
}

func TestStoreWins(t *testing.T) {
	playerStore := StubPlayerStore{
		map[string]int{
			"John":  20,
			"Maria": 30,
		},
		[]string{},
		League{},
	}

	server := NewPlayerServer(&playerStore)

	t.Run("record win for POST", func(t *testing.T) {
		player := "John"
		request := newPostWinRequest(player)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusAccepted)

		if len(playerStore.winCalls) != 1 {
			t.Errorf("got %d calls to RecordWin want %d", len(playerStore.winCalls), 1)
		}

		if playerStore.winCalls[0] != player {
			t.Errorf("did not store correct winner got %q want %q", playerStore.winCalls[0], player)
		}
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
		server := NewPlayerServer(&store)

		request, _ := http.NewRequest(http.MethodGet, "/league", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := getLeagueFromResponse(t, response.Body)
		assertStatus(t, response.Code, http.StatusOK)
		assertContentType(t, response, jsonContentType)
		assertLeague(t, got, wantedLeague)

	})
}

func assertContentType(t testing.TB, response *httptest.ResponseRecorder, want string) {
	t.Helper()
	if response.Result().Header.Get("content-type") != want {
		t.Errorf("response did not have content-type of %s, got %v", want, response.Result().Header)
	}
}

func assertLeague(t testing.TB, got, want League) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
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

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get the correct status code, got %d want %d", got, want)
	}
}

func assertResponseBody(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response body is wrong, got %q want %q", got, want)
	}
}
