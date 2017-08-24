package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"sort"
	// "fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	listenAddr := flag.String("http.addr", ":8090", "http listen address")
	flag.Parse()

	http.HandleFunc("/primes", handler([]int{2, 3, 5, 7, 11, 13}))
	http.HandleFunc("/fibo", handler([]int{1, 1, 2, 3, 5, 8, 13, 21}))
	http.HandleFunc("/odd", handler([]int{1, 3, 5, 7, 9, 11, 13, 15, 17, 19, 21, 23}))
	http.HandleFunc("/rand", handler([]int{5, 17, 3, 19, 76, 24, 1, 5, 10, 34, 8, 27, 7}))
	http.HandleFunc("/numbers", numberHandler())

	log.Fatal(http.ListenAndServe(*listenAddr, nil))
}

func removeDuplicates(elements []float64) []float64 {
	// Use map to record duplicates as we find them.
	encountered := map[float64]bool{}
	result := []float64{}

	for v := range elements {
		if encountered[elements[v]] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[elements[v]] = true
			// Append to result slice.
			result = append(result, elements[v])
		}
	}
	// Return the new slice.
	return result
}

//This is the handler for /numbers endpoint
func numberHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var data map[string][]interface{}
		var mergedNumbers []float64
		queryParam := r.URL.Query()["u"]

		for _, value := range queryParam {
			response, err := http.Get(value)

			if err != nil {
				log.Printf("Error: Query incorrect")
			} else {
				defer response.Body.Close()
				body, err2 := ioutil.ReadAll(response.Body)
				if err2 != nil {
					log.Fatal(err2)
				}
				err := json.Unmarshal(body, &data)
				if err != nil {
					log.Printf("error decoding response: %v", err)
					if e, ok := err.(*json.SyntaxError); ok {
						log.Printf("syntax error at byte offset %d", e.Offset)
					}
					log.Printf("response: %q", body)
				}

				temp := data["Numbers"]
				numbers := make([]float64, len(temp))

				for i, value := range temp {
					switch typedValue := value.(type) {
					case float64:
						numbers[i] = typedValue
						break
					default:
						fmt.Println("Not an int: ", value)
					}
				}

				mergedNumbers = append(mergedNumbers, numbers...)
			}
		}

		tempNum := removeDuplicates(mergedNumbers)
		uniqueNumbers := make([]int, len(tempNum))
		for i := range tempNum {
			uniqueNumbers[i] = int(tempNum[i])
		}
		sort.Ints(uniqueNumbers)
		log.Printf("Unique Numbers list: %v", uniqueNumbers)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"Numbers": uniqueNumbers})

	}
}

func handler(numbers []int) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		waitPeriod := rand.Intn(550)

		log.Printf("%s: waiting %dms.", r.URL.Path, waitPeriod)

		time.Sleep(time.Duration(waitPeriod) * time.Millisecond)

		x := rand.Intn(100)
		if x < 10 {
			http.Error(w, "service unavailable", http.StatusServiceUnavailable)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(map[string]interface{}{"Numbers": numbers})
	}
}
