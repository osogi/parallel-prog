# Task-1: Стек Трайбера
## Цель эксперимента

Исследовать "эффективность" реализации стека с элиминацией представленной в [статье](https://people.csail.mit.edu/shanir/publications/Lock_Free.pdf) по сравнению с классической реализацией стека Трайбера.

## Условия эксперимента
Эксперимент проводился на машине с ОС `Ubuntu 20.04.5 LTS x86_64`, на процессоре `AMD Ryzen 7 4700U`, программа собиралась и оценивалась `go version go1.22.1`, запуск тестов для проверки эффективности проводился следующей командой.

```bash
go test -v -bench=. -cpu=2,4,8 -count=20 -timeout=1h
```

Стеки проверялись на 5 разных видах тестов со следующими параметрами:
- goroutines=10000,ops=1000,pushChance=0,500000,sleep=0s
- goroutines=4000,ops=1000,pushChance=0,500000,sleep=5μs
- goroutines=2000,ops=1000,pushChance=0,500000,sleep=50μs
- goroutines=2000,ops=1000,pushChance=0,500000,sleep=500μs
- goroutines=10000,ops=5000,pushChance=0,250000,sleep=0s

Где:
- **goroutines**: кол-во горутин которые будут запущены
- **ops**: кол-во операций стека для каждой горутины
- **pushChance**: вероятность что как операция стека выберется push
-  **sleep**: задержка между каждой операцией стека (для симуляции работы остальной системы между операциями стека)

Как результат выполнения теста берётся время (в наносекундах) затраченное с того момента как инициализировались все горутины, до того как последняя из них завершила работу со стеком.

## Теоретические ожидания о результатах эксперимента**

До проведения эксперимента ожидалось, что стек с элиминацией будет работать быстрее на тестах, где вероятность push и pop равны. И судя из [статьи](https://people.csail.mit.edu/shanir/publications/Lock_Free.pdf) при отношении вероятностей push и pop 25/75, стек с элиминацией все равно должен показывать более быстрые результаты по сравнению со стеком Трайбера. 

Численная оценка ожидаемого ускорения не рассчитывалась.
  
  
## Результаты эксперимента
Результаты запуска тестов производительности приведены в файле [res.txt](./res.txt).

После устранения выбросов для дальнейших расчётов были взяты средние арифметические значения оставшихся экспериментов. Итоговые результаты приведены в [таблице](https://docs.google.com/spreadsheets/d/1ZRjbH-eyYRZ4Adr47Zl38RsnJQRIa3hqHO1Ug0BRVz0/edit?usp=sharing)

<details>
<summary>Изображение таблицы</summary>

![изображение-таблицы](https://github.com/osogi/it-math/assets/66139162/61f58c73-3184-4559-9c4f-b06f777a276e)

</details>
  

## Оценка полученного результата
Полученные результаты эксперимента не совпадают с теоретическим ожиданием. На всех тестах кроме одного с 50% вероятностью push, стек с элиминацией оказался медленнее обычного стека Трайбера, даже если учитывать 95% доверительный интервал. На тестах с 25% вероятностью push замедление оказалось ещё больше, что было ожидаемо.

**Замечания:**
- на тестах с `sleep` > 0 стек с элиминацией имеет меньшее замедление по сравнению с тестами у которых `sleep` = 0. Следовательно, стек с элиминацией меньше подвержен замедлению из-за сторонней работы;
- если отсортировать количество выделенных ядер процессора по возрастанию замедления, то:
	- на тестах со `sleep` > 0 это будет: 2, 4, 8;
	- на тестах со `sleep` = 0 это будет: 4, 2, 8.
    
    Возможно подобные результаты связаны с тем, что стек с элиминацией при выполнении операций иногда уходит в ожидание, тем самым уменьшая нагрузку на процессор и не используя все данные ему мощности по сравнению со стеком Трайбера.


**Гипотезы почему результаты не совпали с ожиданием:**
- го и горутины положительно влияют на скорость выполнения стека трайбера
- подобраны плохие стартовые и граничные значения для переменных внутри реализации стека с элиминацией
- реализация представленная в [статье](https://people.csail.mit.edu/shanir/publications/Lock_Free.pdf), плохо подходит для её имплементации на языке го (мне пришлось слегка усложнить её из-за специфики языка). Так как судя из [экспериментов коллег](https://github.com/toadharvard/parallel-course/tree/main/stack) другая реализация как минимум показывает ускорение
- была допущена неявная ошибка при реализации стека