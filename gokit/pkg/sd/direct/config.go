package direct

type Config struct {
    Urls []string `json:"urls" yaml:"urls" db:"urls"`
    Name string   `json:"name" yaml:"name" db:"name"`
}
