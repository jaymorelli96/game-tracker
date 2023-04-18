# Poker Game Tracker ![in progress badge](https://badgen.net/badge/status/in%20progress/blue)

Poker Game Tracker is a very simple application that for tracking a poker game. It also calculate and display the blind values updating automatically everytime it is increased.

This application can either run on the web or in the cli. 

# Table of Contents
* [Motivation](#motivation)
* [Key Takeways](#key-takeaways)
* [Setup](#setup)
* [Usage](#usage)
* [Credits](#credits)
* [Road Map](#road-map)
* [License](#license)

# Motivation
It is not easy to find resources on TDD. Usually it covers simple problems, but Chris James does an incredible job with the book [Learn Go With Tests](https://quii.gitbook.io/learn-go-with-tests/)*.


*_I highly recommend you to check his book. James not only shows you how to think and approach the problem using TDD but also show you the power of refactoring._

# Key Takeaways
* TDD
* Mocks, Stubs & Spys
* HTML Templating
* JSON Decoding
* CLI
* Web Sockets


# Setup
### Requirements
* Go 1.18

### Web Server
From root directory:
```
$ go run ./cmd/webserver
```
The server will run on port :5000
### CLI
From root directory:
```
$ go run ./cmd/cli
```
# Usage
The application is currently using a file system store for the persistent layer. The file is `game.db.json`.

When you start your application (cli or web) you will be prompt to enter the number of players. This initial value will be used to calculate the schedule of blind values. 
### Blind Values
* Base amount of time is 5 minutes
* For every player, 1 minute is added

Values: 100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000.

Example: 
* 7 players
* Starts with 100 chips
* At 12 minutes (5 + 7) is 200
* At 24 minutes is 300
* At 36 minutes is 400
* ...
 

## API
<table>
<thead>
<tr>
<th>Method</th>
<th>Pattern</th>
<th>Action</th>
</tr>
</thead>

<tbody>
<tr>
<td>GET</td>
<td>/league</td>
<td>Display the home page</td>
</tr>

<tr>
<td>GET</td>
<td><span>/players/:name</span></td>
<td>Show player score</td>
</tr>

<tr>
<td>POST</td>
<td><span>/players/:name</span></td>
<td>Store player win</td>
</tr>

<tr>
<td>GET</td>
<td>/game</td>
<td>Home for the Web App</td>
</tr>

<tr>
<td>GET</td>
<td>/ws</td>
<td>Create a websocket connection</td>
</tr>

</tbody>
</table>

# Credits
Game Tracker is part of the book [Learn Go With Tests](https://quii.gitbook.io/learn-go-with-tests/) by Chris James.

# Road Map
- [ ] Web page styling
- [ ] Blind scheduler as CLI argument

# License
MIT License

Copyright (c) 2023 Jean Morelli

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.