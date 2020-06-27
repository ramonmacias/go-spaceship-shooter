# go-spaceship-shooter

One of my motivations when I started my studies on computer engineering were the videogames, understand the logic that moves each game under the roof. With this project I'm trying to learn some game basics as well practice some parts of the Go language that I didn't use to use every day. I took as an inspiration this blog https://mortenson.coffee/blog/making-multiplayer-game-go-and-grpc/ and try to do a similar game but on my own way.

## Description

Spaceship shooter, as the name says, is a shooter with spaceships, you are managing your spaceship with a laser gun and the main goal of the game is to kill all the bots. These bots are spaceships as well with their own guns. You have a life on this games as well the bots have, so the main idea is to kill everyone before they kill you.

We have different strategies for the bots.

* **No movement** this is doing nothing, just be on the same place all the time
* **Only movement** this is going to move the bot in a random directions
* **Only shooting** this means that the bot is going to shoot the laser in random directions
* **Shoot and move** this is a mix of only movement and only shooting strategies, means the bot will shoot and move in random directions

By now we are pre define the map, how many players and bots will have the engine to run, but the game engine is ready to receive a combinations for all these fields. Take a quick look on **/cmd/spaceshipShooter/main.go** for see a sample of how we manage this configuration.

We used https://github.com/rivo/tview for manage all the stuff related with the view, in our case we execute the view directly on the terminal.

## Controls

- <kbd>←</kbd> <kbd>→</kbd> <kbd>↑</kbd> <kbd>↓</kbd> movement
- <kbd>w</kbd> <kbd>a</kbd> <kbd>s</kbd> <kbd>d</kbd> shoot laser on specific direction
- <kbd>Ctrl</kbd>+<kbd>C</kbd> exit game
- <kbd>p</kbd> show score
- <kbd>Esc</kbd> close score modal

## How to run

There a simple make file, which has two commands.

```
// Will do a go run, this will build the go binary and will run the game
$ make run

// This command will execute all the tests
$ make test
```
