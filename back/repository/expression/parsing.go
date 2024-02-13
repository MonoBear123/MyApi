package expression

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/MonoBear123/MyApi/back/model"
)

// ParseExpressionToTree парсит математическое выражение в постфиксную форму и строит дерево
func ParseExpressionToTree(s string) (*model.Node, error) {
	var (
		outputQueue   []string
		operatorStack []string
	)

	// Функция, определяющая, является ли символ оператором
	isOperator := func(char string) bool {
		return strings.ContainsAny(char, "+-*/")
	}

	// Функция для определения приоритета оператора
	precedence := map[string]int{
		"+": 1,
		"-": 1,
		"*": 2,
		"/": 2,
	}
	scribe := 0
	// Проход по каждому символу в выражении
	for _, char := range s {
		token := string(char)

		// Пропускаем пробелы
		if unicode.IsSpace(char) {
			continue
		}

		// Проверяем, является ли символ допустимым
		if !(unicode.IsDigit(char) || isOperator(token) || token == "(" || token == ")") {
			return nil, errors.New("недопустимый символ в выражении")
		}

		// Если символ - число, добавляем его в очередь вывода
		if unicode.IsDigit(char) {
			outputQueue = append(outputQueue, token)
		} else if isOperator(token) {
			// Если символ - оператор, обрабатываем стек операторов
			for len(operatorStack) > 0 && precedence[operatorStack[len(operatorStack)-1]] >= precedence[token] {
				outputQueue = append(outputQueue, operatorStack[len(operatorStack)-1])
				operatorStack = operatorStack[:len(operatorStack)-1]
			}
			operatorStack = append(operatorStack, token)
		} else if token == "(" {
			scribe++
			// Есл	и символ - открывающая скобка, добавляем её в стек операторов
			operatorStack = append(operatorStack, token)
		} else if token == ")" {
			if scribe < 0 {
				return &model.Node{}, fmt.Errorf("error")
			}
			scribe--
			// Если символ - закрывающая скобка, выталкиваем операторы из стека в очередь до открывающей скобки
			for len(operatorStack) > 0 && operatorStack[len(operatorStack)-1] != "(" {
				outputQueue = append(outputQueue, operatorStack[len(operatorStack)-1])
				operatorStack = operatorStack[:len(operatorStack)-1]
			}
			// Убираем открывающую скобку из стека
			if len(operatorStack) > 0 && operatorStack[len(operatorStack)-1] == "(" {
				operatorStack = operatorStack[:len(operatorStack)-1]
			} else {
				return nil, errors.New("неверное выражение, лишняя закрывающая скобка")
			}
		}
	}

	// Оставшиеся операторы добавляем в очередь вывода
	for len(operatorStack) > 0 {
		outputQueue = append(outputQueue, operatorStack[len(operatorStack)-1])
		operatorStack = operatorStack[:len(operatorStack)-1]
	}

	// Строим дерево из постфиксной формы
	var nodeStack []*model.Node
	for _, token := range outputQueue {
		if unicode.IsDigit(rune(token[0])) {
			// Если токен - число, добавляем его в стек нод
			value, err := strconv.ParseFloat(token, 64)
			if err != nil {
				return nil, err
			}
			nodeStack = append(nodeStack, &model.Node{Value: value})
		} else if isOperator(token) {
			// Если токен - оператор, выталкиваем две верхние ноды из стека и создаем новую с оператором
			if len(nodeStack) < 2 {
				return nil, errors.New("неверное количество операндов для оператора")
			}
			right := nodeStack[len(nodeStack)-1]
			nodeStack = nodeStack[:len(nodeStack)-1]
			left := nodeStack[len(nodeStack)-1]
			nodeStack = nodeStack[:len(nodeStack)-1]
			node := &model.Node{Left: left, Right: right, Operator: token}
			nodeStack = append(nodeStack, node)
		}
	}

	// В конечном итоге в стеке должна остаться только одна нода - корень дерева
	if len(nodeStack) != 1 || scribe != 0 {
		return nil, errors.New("неверное количество нод в дереве")
	}

	return nodeStack[0], nil
}
