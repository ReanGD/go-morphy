go-morphy
===
[![Build Status](https://travis-ci.org/ReanGD/go-morphy.svg?branch=master)](https://travis-ci.org/ReanGD/go-morphy) [![codecov](https://codecov.io/gh/ReanGD/go-morphy/branch/master/graph/badge.svg)](https://codecov.io/gh/ReanGD/go-morphy)
[![Go Report Card](https://goreportcard.com/badge/github.com/ReanGD/go-morphy)](https://goreportcard.com/report/github.com/ReanGD/go-morphy)

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

[Результаты бенчмарков](https://github.com/ReanGD/go-morphy/wiki/%D0%A0%D0%B5%D0%B7%D1%83%D0%BB%D1%8C%D1%82%D0%B0%D1%82%D1%8B-%D0%B1%D0%B5%D0%BD%D1%87%D0%BC%D0%B0%D1%80%D0%BA%D0%BE%D0%B2-%D0%B4%D0%BB%D1%8F-%D1%81%D1%82%D1%80%D1%83%D0%BA%D1%82%D1%83%D1%80%D1%8B-DAWG)
