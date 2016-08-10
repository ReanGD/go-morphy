go-morphy
===

Порт морфологического анализатора [pymorphy2](https://github.com/kmike/pymorphy2) ([v0.8](https://github.com/kmike/pymorphy2/releases/tag/0.8)) и его составной части [DAWG-Python](https://github.com/pytries/DAWG-Python) ([v0.7.2](https://github.com/pytries/DAWG-Python/releases/tag/0.7.2)) на Golang.

Тесты и бенчмарки
---
Для прохождения требуются тестовые данные из оригинального pymorphy2, поэтому нужно клонировать проект вместе с submodule:
```
git clone --recursive https://github.com/ReanGD/go-morphy.git
```

- Запуск тестов:
```
cd go-morphy
go test -v ./...
```

- Запуск бенчмарков:
```
cd go-morphy/benchmarks
go test -bench=BenchmarkDAWG
```
