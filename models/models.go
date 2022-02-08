package models

type Config struct {
	DB_USERNAME string `yaml:"DB_USERNAME"`
	DB_PASSWORD string `yaml:"DB_PASSWORD"`
	DB_IP       string `yaml:"DB_IP"`
	DB_PORT     string `yaml:"DB_PORT"`
	DB_NAME     string `yaml:"DB_NAME"`
}


type Api_body struct {


	Apikey   string  `json:"api_key"`
	Senderid string  `json:"sender_id"`
	Contacts   string  `json:"contacts"`
	MessageType   string  `json:"type"`
	Message   string  `json:"msg"`

}
type User struct {
	ID     int
	Name   string
	Email  string
	Status int
}

type Queue struct {
	ID     int
	Name   string
	Email  string
	Status int
}
