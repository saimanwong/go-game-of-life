# Game of Life

An implementation of Conway's Game of Life using Go and Ebiten.

- [Conway's Game of Life wiki](https://en.wikipedia.org/wiki/Conway%27s_Game_of_Life)
- [Ebiten - A dead simple 2D game library for Go](https://ebiten.org/)

## Key Bindings

| Key(s)         | Description                    |
| -------------- | -------------------------------|
| P              | Randomly populates the world   |
| R              | Resets the world               |
| UP / DOWN      | Increase or decrease TPS       |
| H + Left Click | Inserts a [Heavyweight Spaceship](https://en.wikipedia.org/wiki/Conway%27s_Game_of_Life#/media/File:Animated_Hwss.gif) |

## Getting Started

```
$ go run main.go world.go
```

### Using Docker and Browser

Open http://localhost in browser.

#### Remote from Docker Hub

```
$ PORT=8080; docker run --pull always -it -p 80:$PORT -e PORT=$PORT saimanwong/go-game-of-life:latest
```

#### Local Build and Run `Dockerfile.dockerhub`

```
$ docker build -t gol-local -f Dockerfile.dockerhub .
$ PORT=8080; docker run -it -p 80:$PORT -e PORT=$PORT --expose $PORT gol-local
```

## Screenshots

![](screenshots/1.png)
