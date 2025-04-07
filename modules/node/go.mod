module github.com/DALDA-IITJ/libr/modules/node

go 1.24.1

require (
	github.com/gorilla/mux v1.8.1
	github.com/joho/godotenv v1.5.1
	github.com/lib/pq v1.10.9
	gopkg.in/yaml.v2 v2.4.0
)

require github.com/DALDA-IITJ/libr/modules/core v0.0.0

replace github.com/DALDA-IITJ/libr/modules/core => ../core
