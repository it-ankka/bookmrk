# Bookmrk - The self-hostable bookmark manager

Bookmrk is entirely built around Pocketbase and uses standard Go html templating capabilities for it's user interface.

### Developing - Getting started

1. Make sure all needed dependencies have been installed by running `go mod tidy` and `npm run install`
2. If you are gonna be developing the UI, start the tailwind file watcher by running `npm run tailwind:watch`
2. Then you can run the project by running `go run main.go serve` or if you want hot reloading capabilities you can install [Air](https://github.com/air-verse/air) and run the project using the command `air`
