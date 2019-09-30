package cloudstorage

import "time"

type StorageEvent struct {
	Id             string            `json:"id"`
	Bucket         string            `json:"bucket"`
	Name           string            `json:"name"`
	Metageneration string            `json:"metageneration"`
	ResourceState  string            `json:"resourceState"`
	TimeCreated    time.Time         `json:"timeCreated"`
	Updated        time.Time         `json:"updated"`
	Metadata       map[string]string `json:"metadata"`
}
