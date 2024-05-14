# Используем официальный образ Golang в качестве базового образа
FROM golang:latest

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем go.mod и go.sum и загружаем зависимости
COPY go.mod go.sum ./
RUN go mod download
RUN go mod tidy

# Копируем все файлы из текущего каталога хоста в рабочую директорию внутри контейнера
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main

# Указываем порт, который будет слушать приложение
EXPOSE 8080

# Команда для запуска вашего приложения при старте контейнера
CMD ["/cmd/server/app.go"]
