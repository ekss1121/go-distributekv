package config

type Partition struct {
	Name  string
	Index int
	Host  string
}

type Partitions struct {
	Partitions []Partition
}
