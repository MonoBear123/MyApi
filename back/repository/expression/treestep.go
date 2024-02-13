package expression

import "github.com/MonoBear123/MyApi/back/model"

// EvaluateExpression вычисляет значение выражения, представленного в виде дерева
func EvaluateExpression(node *model.Node) float64 {
	if node == nil {
		return 0
	}

	if node.Operator == "" {
		// Если узел не является оператором, возвращаем значение узла
		return node.Value
	}

	// Рекурсивно вычисляем значения для левого и правого поддеревьев
	leftValue := EvaluateExpression(node.Left)
	rightValue := EvaluateExpression(node.Right)

	// Применяем оператор к значениям
	switch node.Operator {
	case "+":
		return leftValue + rightValue
	case "-":
		return leftValue - rightValue
	case "*":
		return leftValue * rightValue
	case "/":
		if rightValue != 0 {
			return leftValue / rightValue
		}
		// Обработка деления на ноль
		return 0
	default:
		// Неизвестный оператор
		return 0
	}
}
