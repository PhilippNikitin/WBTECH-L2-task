package chain_of_resp

/*
	Реализовать паттерн «цепочка вызовов».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Chain-of-responsibility_pattern

Паттерн "Цепочка вызовов" (Chain of Responsibility) используется для передачи запроса по цепочке потенциальных
обработчиков, пока один из них не обработает запрос.

Данный паттерн может применяться в следующих ситуациях:
1. Разделение обязанностей: когда требуется разделить обработку запроса между различными объектами, каждый из которых отвечает
за определённый тип обработки.
2. Динамическая настройка обработки: когда необходимо динамически определять, какой обработчик будет обработать запрос,
основываясь на условиях, которые могут изменяться во время выполнения программы.
3. Избежание жесткой связанности: когда необходимо ослабить связанность между отправителем запроса и его получателями,
так чтобы отправитель не знал, какой именно объект обработает запрос.
4. Расширяемость: если необходимо добавить новые обработчики без изменения существующего кода отправителя и других обработчиков.

Плюсы паттерна "Цепочка вызовов":
1. Гибкость обработки запросов: запрос может быть обработан различными объектами в зависимости от их состояния и логики,
что делает систему более гибкой.
2. Ослабление связанности: отправитель запроса не знает, какой обработчик его обработает, что уменьшает связанность между
компонентами.
3. Легкость добавления новых обработчиков: легко добавлять новые обработчики, не изменяя существующий код, что способствует
расширяемости системы.
4. Разделение ответственности: обработка запроса может быть распределена между различными объектами, каждый из которых отвечает
только за свою часть обработки.

Минусы паттерна "Цепочка вызовов":
1. Неопределенность обработчика: запрос может пройти по всей цепочке и остаться необработанным, если ни один из обработчиков
не подходит для его обработки.
2. Трудности отладки: отладка цепочки вызовов может быть сложной, так как трудно отследить, какой обработчик в итоге обработал
запрос или почему он не был обработан.
3. Потенциальная деградация производительности: Если цепочка слишком длинная, то передача запроса через множество обработчиков
может снизить производительность системы.
*/

/*
Реализация паттерна "Цепочка вызовов" на примере обработки запросов на основе уровня логирования (INFO, DEBUG, ERROR).
*/

import "fmt"

// LogLevel определяет уровень логирования
type LogLevel int

const (
	INFO LogLevel = iota
	DEBUG
	ERROR
)

// Logger - интерфейс обработчика запросов
type Logger interface {
	SetNext(Logger)              // установка следующего обработчика
	LogMessage(LogLevel, string) // обработка сообщения
}

// BaseLogger - базовый логгер, реализующий установку следующего обработчика
type BaseLogger struct {
	next Logger
}

func (b *BaseLogger) SetNext(next Logger) {
	b.next = next
}

func (b *BaseLogger) LogMessage(level LogLevel, message string) {
	if b.next != nil {
		b.next.LogMessage(level, message)
	}
}

// InfoLogger - обработчик логов уровня INFO
type InfoLogger struct {
	BaseLogger
}

func (l *InfoLogger) LogMessage(level LogLevel, message string) {
	if level == INFO {
		fmt.Printf("INFO: %s\n", message)
	} else {
		l.BaseLogger.LogMessage(level, message)
	}
}

// DebugLogger - обработчик логов уровня DEBUG
type DebugLogger struct {
	BaseLogger
}

func (l *DebugLogger) LogMessage(level LogLevel, message string) {
	if level == DEBUG {
		fmt.Printf("DEBUG: %s\n", message)
	} else {
		l.BaseLogger.LogMessage(level, message)
	}
}

// ErrorLogger - обработчик логов уровня ERROR
type ErrorLogger struct {
	BaseLogger
}

func (l *ErrorLogger) LogMessage(level LogLevel, message string) {
	if level == ERROR {
		fmt.Printf("ERROR: %s\n", message)
	} else {
		l.BaseLogger.LogMessage(level, message)
	}
}

func main() {
	infoLogger := &InfoLogger{}
	debugLogger := &DebugLogger{}
	errorLogger := &ErrorLogger{}

	// Строим цепочку: INFO -> DEBUG -> ERROR
	infoLogger.SetNext(debugLogger)
	debugLogger.SetNext(errorLogger)

	// Передаем запросы разным обработчикам
	infoLogger.LogMessage(INFO, "This is an info message.")
	infoLogger.LogMessage(DEBUG, "This is a debug message.")
	infoLogger.LogMessage(ERROR, "This is an error message.")
}
