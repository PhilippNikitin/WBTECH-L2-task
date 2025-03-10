package facade

/*
	Реализовать паттерн «фасад».
Объяснить применимость паттерна, его плюсы и минусы,а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Facade_pattern
*/

/*
Паттерн «Фасад» (Facade) предоставляет простой интерфейс к сложной системе классов, библиотеке или фреймворку.
Он упрощает взаимодействие с этой системой, скрывая её сложность и предоставляя более удобный интерфейс.

Применимость паттерна "Фасад".
Паттерн Фасад применяется в следующих ситуациях:
1. Для упрощения использования сложной системы: когда система состоит из множества взаимосвязанных классов,
паттерн Фасад предоставляет упрощённый интерфейс для их использования.
2. Для сокрытия деталей реализации: паттерн Фасад позволяет скрыть сложность внутренней системы и предоставлять
только необходимые методы для работы с ней.
3. Для снижения зависимости между подсистемами: Фасад уменьшает количество прямых зависимостей между клиентами
и сложной системой, предоставляя единый интерфейс.
4. Упрощение тестирования: при наличии сложных взаимодействий между компонентами фасад может обеспечить
упрощенный интерфейс для тестирования.

Плюсы паттерна "Фасад":
1. Упрощение интерфейса: Фасад предоставляет простой и удобный интерфейс для взаимодействия с подсистемой,
скрывая её сложность.
2. Инкапсуляция деталей реализации: скрывает детали реализации подсистемы от клиента,
предоставляя ему только необходимые методы.
3. Снижение зависимости: клиенты зависят только от фасада, а не от всех компонентов подсистемы,
что упрощает изменение и поддержку кода.
4. Улучшение читаемости кода: Улучшает читаемость и поддерживаемость кода за счёт предоставления более
понятного интерфейса.

Минусы паттерна "Фасад":
1. Ограниченная функциональность: Фасад может ограничивать доступ ко всей функциональности подсистемы,
предоставляя только часть возможностей.
2. Создание дополнительного уровня абстракции: Это может привести к дублированию кода или усложнению системы,
если не используется правильно.
3. Риск нарушения принципа единой ответственности: Фасад может стать слишком большим и сложным, если он управляет
слишком многими компонентами, нарушая принцип единственной ответственности.

Применение паттерна Фасад на практике:
1. Фасад может применяться в различных библиотеках и фреймворках для упрощения работы с ними.
2. Фасад может применяться в API интеграциях путем предоставления унифицированного интерфейса для взаимодействия с ними,
скрывая различия в API.
*/

/*
Пример применения паттерна Фасад:

Сервис, работающий с различными сторонними API для получения данных о погоде, новостях и курсах валют.
Вместо того чтобы взаимодействовать с каждым API отдельно в клиентском коде, можно создать фасад для
всех этих сервисов.
*/

// Подсистемы
// Подсистема для получения данных о погоде

import "fmt"

type WeatherService struct{}

func (ws *WeatherService) GetWeather(location string) string {
	// Логика получения данных о погоде
	return "Sunny"
}

// Подсистема для получения новостей
type NewsService struct{}

func (ns *NewsService) GetLatestNews() string {
	// Логика получения новостей
	return "Breaking News: Facade pattern in use!"
}

// Подсистема для получения данных о курсах валют
type CurrencyService struct{}

func (cs *CurrencyService) GetExchangeRate(from, to string) float64 {
	// Логика получения курса валют
	return 1.13
}

// Фасад - унифицированный интерфейс для взаимодействия с подсистемами
type ServiceFacade struct {
	weatherService  *WeatherService
	newsService     *NewsService
	currencyService *CurrencyService
}

func NewServiceFacade() *ServiceFacade {
	return &ServiceFacade{
		weatherService:  &WeatherService{},
		newsService:     &NewsService{},
		currencyService: &CurrencyService{},
	}
}

func (sf *ServiceFacade) GetWeather(location string) string {
	return sf.weatherService.GetWeather(location)
}

func (sf *ServiceFacade) GetLatestNews() string {
	return sf.newsService.GetLatestNews()
}

func (sf *ServiceFacade) GetExchangeRate(from, to string) float64 {
	return sf.currencyService.GetExchangeRate(from, to)
}

// Клиентский код - взаимодействует только с фасадом,
// не зная о внутренней структуре и реализации подсистем.

func main() {
	facade := NewServiceFacade()

	fmt.Println("Weather:", facade.GetWeather("New York"))
	fmt.Println("News:", facade.GetLatestNews())
	fmt.Println("Exchange Rate (USD to EUR):", facade.GetExchangeRate("USD", "EUR"))
}
