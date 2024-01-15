# Usa una imagen base de Golang
FROM golang:latest

# Establece el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copia el código fuente al contenedor
COPY . .

# Ejecuta el comando go run para ejecutar la aplicación
CMD ["go", "run", "main.go"]