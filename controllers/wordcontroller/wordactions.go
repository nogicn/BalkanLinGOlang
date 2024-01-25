package wordcontroller

import (
	"BalkanLinGO/models/userworddb"
	"BalkanLinGO/models/worddb"
	"math/rand"
)

func fillWordList(words []userworddb.Word, finalWords []worddb.Word, n int) []worddb.Word {
	// get 4 random words
	for i := 0; i < n; i++ {
		////fmt.Println("generating ")
		// generate random number between 1 and number of words in dictionary
		random := rand.Intn(len(words))
		duplicate := false
		for j := 0; j < len(finalWords); j++ {
			if finalWords[j].ID == words[random].ID {
				duplicate = true
				break
			}
		}
		if duplicate {
			i--
			continue
		}
		finalWords = append(finalWords, worddb.Word(words[random]))
	}
	return finalWords
}
