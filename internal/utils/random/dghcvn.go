package random

import "fmt"

var (
	cities    = []string{"SG", "HN", "DDN", "BD", "DNA", "HP"}
	districts = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
	wards     = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
)

func RandomCity() string {
	return cities[r.Intn(len(cities))]
}

func RandomDistrict() string {
	return districts[r.Intn(len(districts))]
}

func RandomWard() string {
	return wards[r.Intn(len(wards))]
}

func RandomAddress() string {
	return fmt.Sprintf("Số %d, %s, %s, %s", r.Intn(200), RandomWard(), RandomDistrict(), RandomCity())
}
