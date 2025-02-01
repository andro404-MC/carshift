# Default url: http://localhost:7331
live/templ:
	templ generate --watch --proxy="http://localhost:8000" --open-browser=false

live/server:
	go run github.com/air-verse/air@latest \
	--build.cmd "go build -o tmp/bin/main" --build.bin "tmp/bin/main" --build.delay "100" \
	--build.exclude_dir "node_modules" \
	--build.include_ext "go" \
	--build.stop_on_error "false" \
	--misc.clean_on_exit true

live/tailwind:
	npx @tailwindcss/cli -i ./static/css/input.css -o ./static/css/output.css --watch -m

live/sync_assets:
	go run github.com/air-verse/air@latest \
	--build.cmd "templ generate --notify-proxy" \
	--build.bin "true" \
	--build.delay "100" \
	--build.exclude_dir "" \
	--build.include_dir "assets" \
	--build.include_ext "js,css"

live: 
	make -j5 live/templ live/server

build/templ:
	templ generate .

build/tailwindcss:
	npm install
	npx @tailwindcss/cli -i ./static/css/input.css -o ./build/static/css/output.css -m

build/go:
	go mod tidy
	go build -ldflags "-s -w" -o ./build/carshift main.go

build/mini:
	upx ./build/carshift

build:
	mkdir build
	cp static build -r
	make build/templ build/tailwindcss build/go build/mini
