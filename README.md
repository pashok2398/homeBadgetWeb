docker run -p 8081:8080 -v ./data.csv:/data.csv -it go
docker run -p 8080:8080 -v ./data.csv:/usr/src/app/data.csv -it test