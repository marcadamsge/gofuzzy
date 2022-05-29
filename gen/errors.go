package gen

type RandIntGenerator interface {
	// Intn returns a random int in the interval [0, n)
	Intn(n int) int
}

type errorType int

const (
	insert  errorType = 0
	remove            = 1
	replace           = 2
	swap              = 3
)

func RandomFuzzyErrors(word string, randGen RandIntGenerator, distance int, alphabet []rune) string {
	out := []rune(word)

	for i := 0; i < distance; i++ {
		var randomError errorType

		if len(out) == 0 {
			// only insert error is possible in that case
			randomError = insert
		} else if len(out) == 1 {
			// insert, remove or replace is possible
			randomError = errorType(randGen.Intn(3))
		} else {
			// any error is possible
			randomError = errorType(randGen.Intn(4))
		}

		switch randomError {
		case insert:
			out = InsertCharacterError(out, randGen, alphabet)
		case remove:
			out = RemoveCharacterError(out, randGen)
		case replace:
			out = ReplaceCharacterError(out, randGen, alphabet)
		case swap:
			out = SwapAdjacentCharacterError(out, randGen)
		}
	}

	return string(out)
}

func InsertCharacterError(word []rune, randGen RandIntGenerator, alphabet []rune) []rune {
	wordLength := len(word)
	randomCharacter := alphabet[randGen.Intn(len(alphabet))]
	position := randGen.Intn(wordLength + 1)

	if position == wordLength {
		return append(word, randomCharacter)
	}

	word = append(word, word[wordLength-1])
	copy(word[position+1:wordLength], word[position:wordLength-1])
	word[position] = randomCharacter
	return word
}

func RemoveCharacterError(word []rune, randGen RandIntGenerator) []rune {
	if len(word) >= 1 {
		wordLength := len(word)
		characterToRemove := randGen.Intn(wordLength)
		copy(word[characterToRemove:wordLength-1], word[characterToRemove+1:])
		return word[:wordLength-1]
	}

	return word
}

func ReplaceCharacterError(word []rune, randGen RandIntGenerator, alphabet []rune) []rune {
	if len(word) >= 1 {
		randomCharacter := alphabet[randGen.Intn(len(alphabet))]
		position := randGen.Intn(len(word))
		word[position] = randomCharacter
	}

	return word
}

func SwapAdjacentCharacterError(word []rune, randGen RandIntGenerator) []rune {
	if len(word) > 1 {
		position := randGen.Intn(len(word) - 1)
		word[position], word[position+1] = word[position+1], word[position]
	}

	return word
}
