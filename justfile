css:
    tailwindcss -i ./templates/static/style.css -o ./static/style.css 

generate:
    templ generate

go-build: generate css
    go build -o ./tmp/web ./cmd/web

air:
    air -c .air.toml