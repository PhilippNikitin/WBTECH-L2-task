package command

/*
	Реализовать паттерн «комманда».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Command_pattern


Паттерн "Команда" (Command) является поведенческим паттерном проектирования, который инкапсулирует запрос в виде объекта,
позволяя тем самым параметризовать клиентские объекты с различными запросами, ставить запросы в очередь или протоколировать их,
а также поддерживать отмену операций.

Применимость паттерна "Команда" (Command)
1. Управление запросами. Паттерн "Команда" полезен, когда нужно обрабатывать запросы как объекты, позволяя параметризовать операции,
ставить их в очередь, и управлять выполнением.
2. Отмена и повторение операций. Использование паттерна позволяет легко реализовать отмену (undo) и повторение (redo) операций,
так как каждая команда может знать, как отменить своё действие.
3. Отделение отправителя от получателя: Команда убирает прямую связь между объектами, которые инициируют запрос (отправители)
и теми, которые его исполняют (получатели). Это способствует уменьшению связанности и увеличению гибкости системы.
4. Очереди и логирование. Паттерн позволяет создавать очереди команд для последовательного выполнения или логировать выполненные
команды для последующего анализа или восстановления состояния.

Плюсы паттерна "Команда":
1. Разделение ответственности: каждая команда отвечает только за выполнение одного действия, что способствует улучшению структуры
кода.
2. Отмена и повторение операций: легко добавлять функции отмены и повтора, так как каждая команда знает, как выполнять и отменять
своё действие.
3. Гибкость системы: позволяет динамически изменять набор операций, добавлять новые команды без изменения существующего кода
отправителя.

Минусы паттерна "Команда":
1. Увеличение количества классов: внедрение команд может привести к созданию большого количества классов,
особенно если операций много и они детализированы.
2. Сложность отладки: усложняется отладка системы из-за большего количества классов и взаимодействий между ними.
3. Возможное замедление: использование очередей команд может привести к задержкам при выполнении операций,
особенно если операции требуют мгновенного выполнения.
*/

/*
Пример реализации паттерна "Команда"
Система для обработки заказов, где каждый заказ может быть обработан различными способами (например, создание, обновление,
удаление заказа). Можно использовать паттерн "Команда" для инкапсуляции каждой операции над заказом в отдельный объект команды.
*/

import "fmt"

// OrderService - сервис для работы с заказами
// содержит поле orders, в котором ключ - ID заказа, а значение - указатель на соответствующий заказ
type OrderService struct {
	orders map[int]*Order
}

// Order - структура заказа
type Order struct {
	ID     int
	Status string
}

// Command - интерфейс команды
type Command interface {
	Execute() error
	Undo() error
	GetOrderID() int
}

// CreateOrderCommand - команда для создания нового заказа
type CreateOrderCommand struct {
	orderService *OrderService
	order        *Order
}

func (cmd *CreateOrderCommand) Execute() error {
	if _, exists := cmd.orderService.orders[cmd.order.ID]; exists {
		return fmt.Errorf("Заказ с ID %d уже существует", cmd.order.ID)
	}
	cmd.orderService.orders[cmd.order.ID] = cmd.order
	fmt.Printf("Создан заказ с ID %d\n", cmd.order.ID)
	return nil
}

func (cmd *CreateOrderCommand) Undo() error {
	delete(cmd.orderService.orders, cmd.order.ID)
	fmt.Printf("Удалён заказ с ID %d\n", cmd.order.ID)
	return nil
}

func (cmd *CreateOrderCommand) GetOrderID() int {
	return cmd.order.ID
}

// UpdateOrderStatusCommand - команда для обновления статуса заказа
type UpdateOrderStatusCommand struct {
	orderService   *OrderService
	orderID        int
	newStatus      string
	previousStatus string
}

func (cmd *UpdateOrderStatusCommand) Execute() error {
	orderToUpdate, exists := cmd.orderService.orders[cmd.orderID]
	if !exists {
		return fmt.Errorf("Заказ с ID %d не найден", cmd.orderID)
	}

	// Сохраним предыдущий статус
	cmd.previousStatus = orderToUpdate.Status

	// Обновим статус
	orderToUpdate.Status = cmd.newStatus
	fmt.Printf("Обновлен статус заказа с ID %d на %s\n", cmd.orderID, cmd.newStatus)
	return nil
}

func (cmd *UpdateOrderStatusCommand) Undo() error {
	orderToUpdate, exists := cmd.orderService.orders[cmd.orderID]
	if !exists {
		return fmt.Errorf("Заказ с ID %d не найден", cmd.orderID)
	}

	// Возвращаем предыдущий статус
	orderToUpdate.Status = cmd.previousStatus
	fmt.Printf("Статус заказа с ID %d возвращён к %s\n", cmd.orderID, cmd.previousStatus)
	return nil
}

func (cmd *UpdateOrderStatusCommand) GetOrderID() int {
	return cmd.orderID
}

// CommandInvoker - инвокер команд
type CommandInvoker struct {
	history map[int][]Command
}

func (invoker *CommandInvoker) StoreAndExecute(cmd Command) error {
	err := cmd.Execute()
	if err == nil {
		orderID := cmd.GetOrderID()
		invoker.history[orderID] = append(invoker.history[orderID], cmd)
	}
	return err
}

func (invoker *CommandInvoker) UndoLastCommand(orderID int) error {
	commands, exists := invoker.history[orderID]
	if !exists || len(commands) == 0 {
		return fmt.Errorf("Нет команд для заказа с ID %d", orderID)
	}

	lastCommand := commands[len(commands)-1]
	err := lastCommand.Undo()
	if err == nil {
		invoker.history[orderID] = commands[:len(commands)-1]
	}
	return err
}

func main() {
	orderService := &OrderService{orders: make(map[int]*Order)}
	invoker := &CommandInvoker{history: make(map[int][]Command)}

	// Создаем новый заказ
	newOrder := &Order{ID: 1, Status: "новый"}
	createOrderCmd := &CreateOrderCommand{orderService: orderService, order: newOrder}
	invoker.StoreAndExecute(createOrderCmd)

	// Обновляем статус заказа
	updateOrderCmd := &UpdateOrderStatusCommand{orderService: orderService, orderID: 1, newStatus: "в обработке"}
	invoker.StoreAndExecute(updateOrderCmd)

	// Печатаем историю выполненных команд
	fmt.Println("\nИстория выполненных команд для заказа 1:")
	for _, cmd := range invoker.history[1] {
		fmt.Printf("- %T\n", cmd)
	}

	// Пример отмены последней команды (откат на предыдущий статус)
	fmt.Println("\nОтмена последней команды для заказа 1:")
	invoker.UndoLastCommand(1)

	// Печатаем обновленный статус заказа
	fmt.Println("\nОбновленный статус заказа после отмены:")
	if order, exists := orderService.orders[1]; exists {
		fmt.Printf("- Заказ ID %d имеет статус: %s\n", order.ID, order.Status)
	}
}
