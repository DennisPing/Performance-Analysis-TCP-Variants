package pkg

import (
	"encoding/csv"
	"os"
	"strconv"
	"strings"
)

// Record the results of a single experiment trial and save it as a CSV file
// Parameters: x slice, y slice, x name, y name, filename
func Record(x, y []float64, xname, yname string, fname string) {
	if !strings.HasSuffix(fname, ".csv") {
		panic("filename must end with .csv")
	}
	file, err := os.Create(fname)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write the header row
	writer.Write([]string{xname, yname})
	// Write the data rows
	for i := 0; i < len(x); i++ {
		writer.Write([]string{strconv.FormatFloat(x[i], 'f', 10, 64), strconv.FormatFloat(y[i], 'f', 10, 64)})
	}
}
