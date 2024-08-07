Что выведет программа? Объяснить вывод программы.

```go
package main

type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}

func test() *customError {
	{
		// do something
	}
	return nil
}

func main() {
	var err error
	err = test()
	if err != nil {
		println("error")
		return
	}
	println("ok")
}
```

Ответ:
```
Результат вывода программы:
error

Объяснение:
1. В функции main мы создаем переменную err интерфейсного типа error.
2. Далее мы записываем в данную переменную результат работы функции test, в результате чего
динамический тип переменной err становится равным *customError, а значение данного
динамического типа устанавливается равным nil.
3. Далее мы сравниваем значение err с nil. Err не равно nil, т.к. у данной переменной 
установлен динамический тип (*customError). Интерфейс равен nil, только если
и динамический тип, и значение динамического типа равны nil.
4. Т.о., программа будет заходить в блок `if err != nil` и выводить "error".

```
