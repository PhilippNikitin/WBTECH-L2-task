package main

/*
=== Поиск анаграмм по словарю ===

Напишите функцию поиска всех множеств анаграмм по словарю.
Например:
'пятак', 'пятка' и 'тяпка' - принадлежат одному множеству,
'листок', 'слиток' и 'столик' - другому.

Входные данные для функции: ссылка на массив - каждый элемент которого - слово на русском языке в кодировке utf8.
Выходные данные: Ссылка на мапу множеств анаграмм.
Ключ - первое встретившееся в словаре слово из множества
Значение - ссылка на массив, каждый элемент которого, слово из множества. Массив должен быть отсортирован по возрастанию.
Множества из одного элемента не должны попасть в результат.
Все слова должны быть приведены к нижнему регистру.
В результате каждое слово должно встречаться только один раз.

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

import (
	"fmt"
	"sort"
	"strings"
)

func findAnagramSets(words []string) map[string][]string {
	anagramSets := make(map[string][]string)

	// Создаем мапу для хранения сортированных слов и их оригинальных версий
	sortedWords := make(map[string]string)

	// Проходим по всем словам во входном массиве
	for _, word := range words {
		// Приводим слово к нижнему регистру
		word = strings.ToLower(word)

		// Создаем сортированную версию слова для использования в качестве ключа мапы
		sortedWord := sortString(word)

		// Если такой ключ уже существует в мапе, добавляем оригинальное слово в соответствующий массив
		if _, ok := sortedWords[sortedWord]; ok {
			anagramSets[sortedWords[sortedWord]] = append(anagramSets[sortedWords[sortedWord]], word)
		} else {
			// Если ключа еще нет, создаем новую запись в мапе и в мапе для анаграмм
			sortedWords[sortedWord] = word
			anagramSets[word] = []string{word}
		}
	}

	// Удаляем все множества из одного элемента из результата
	for key, value := range anagramSets {
		if len(value) <= 1 {
			delete(anagramSets, key)
		} else {
			// Сортируем массив анаграмм по возрастанию перед добавлением в результат
			sort.Strings(value)
			anagramSets[key] = value
		}
	}

	return anagramSets
}

// Вспомогательная функция для сортировки символов в строке
func sortString(s string) string {
	chars := strings.Split(s, "")
	sort.Strings(chars)
	return strings.Join(chars, "")
}

func main() {
	// Пример использования функции
	words := []string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик", "кот", "ток", "окт"}

	anagramSets := findAnagramSets(words)

	// Вывод результата
	for key, value := range anagramSets {
		fmt.Printf("%s: %v\n", key, value)
	}
}
