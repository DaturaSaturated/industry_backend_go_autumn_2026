# Промышленная backend-разработка на Go — autumn 2026

11 заданий: `task_00`–`task_10`. Go 1.27.1, только стандартная библиотека.

## Задания

| Задание | Статус | Название |
|---|---|---|
| [task_00](tasks/task_00/README.md) | ![task_00](badges/tasks/task_00.svg) | Hello, world |
| [task_01](tasks/task_01/README.md) | ![task_01](badges/tasks/task_01.svg) | Приветствие с нормализацией |
| [task_02](tasks/task_02/README.md) | ![task_02](badges/tasks/task_02.svg) | Циклический сдвиг UTF-8 |
| [task_03](tasks/task_03/README.md) | ![task_03](badges/tasks/task_03.svg) | FizzBuzz для знаковых чисел |
| [task_04](tasks/task_04/README.md) | ![task_04](badges/tasks/task_04.svg) | Агрегация изменений |
| [task_05](tasks/task_05/README.md) | ![task_05](badges/tasks/task_05.svg) | Generic-кэш с отказом при заполнении |
| [task_06](tasks/task_06/README.md) | ![task_06](badges/tasks/task_06.svg) | LRU: чтение обновляет давность |
| [task_07](tasks/task_07/README.md) | ![task_07](badges/tasks/task_07.svg) | Потокобезопасный LRU с теми же правилами |
| [task_08](tasks/task_08/README.md) | ![task_08](badges/tasks/task_08.svg) | Token bucket: атомарное списание |
| [task_09](tasks/task_09/README.md) | ![task_09](badges/tasks/task_09.svg) | Worker pool: независимые результаты |
| [task_10](tasks/task_10/README.md) | ![task_10](badges/tasks/task_10.svg) | In-memory репозиторий и HTTP API задач |

## Начало работы

1. Сделайте форк `ippaveln/industry_backend_go_autumn_2026`. Основная ветка — `master`.
2. Включите Actions в своём форке: вкладка **Actions** → **I understand my workflows, go ahead and enable them**. После этого каждый push в `master` форка прогоняет тесты и обновляет бейджи в таблице заданий.
3. Выполните в клоне форка `git config pull.rebase true`. CI форка коммитит бейджи в `master`, поэтому после push локальная ветка отстаёт от форка. С этой настройкой `git pull` пересаживает ваши новые коммиты поверх коммита бота; без неё git остановится с ошибкой `Need to specify how to reconcile divergent branches`.
4. Установите Go 1.27.1 (https://go.dev/dl/) и проверьте `go version`. Для `-race` нужен cgo: на Linux и Windows должен быть установлен gcc.
5. Реализуйте функции в `tasks/task_XX/solution.go`. Заготовки компилируются, но вызывают `panic("TODO: ...")`, поэтому до реализации тесты падают.
6. Проверяйте решение локально:

   | Команда | Что делает |
   |---|---|
   | `make task_XX` | тесты задания с флагами официальной проверки: `go test -race -count=4 -timeout=60s ./tasks/task_XX` |
   | `make check` | то же для всех заданий |
   | `make run-task_XX` | `go run ./tasks/task_XX` — пример из условия |
   | `make grade` | проверяющая программа; правило файлов локально не проверяется |
   | `make help` | список команд |

7. Откройте PR из форка в `master` этого репозитория. Заголовок — `ИСУ <номер>`, например `ИСУ 123456`. Результат по каждому заданию — отдельный check `task_XX` в PR. Как устроена проверка — [GRADING.md](GRADING.md).

## Правила

README задания — полная спецификация: тесты проверяют только описанное в нём. Все тесты опубликованы в репозитории, эталонных решений в нём нет.

Менять можно только `tasks/task_00/solution.go` … `tasks/task_10/solution.go`. Изменения `.gitignore` и `badges/tasks/*.svg` проверка игнорирует. Любое другое изменение — новый, удалённый или изменённый файл, символическая ссылка — проваливает все задания.

В `solution.go` запрещены:

- пакеты вне стандартной библиотеки и cgo;
- импорты `os`, `syscall`, `runtime`, `unsafe`, `testing`, `flag`, `plugin`, `io/ioutil`, `log/syslog` и их подпакетов;
- директивы `//go:` и build tags;
- использование функций и переменных из `*_test.go`: решение должно собираться без тестов.

AI-инструментами пользоваться можно, но на экзамене решение нужно уметь объяснить.

## Срок

Сдаётся один PR со всеми заданиями. Срок — **2 ноября 2026, 23:59 МСК (UTC+3)**.

## Обновления

Если в этом репозитории обновятся условия или тесты, подтяните изменения в форк.

Через GitHub: в форке нажмите **Sync fork** → **Update branch**, затем локально выполните `git pull`.

Или из командной строки:

```sh
git remote add upstream https://github.com/ippaveln/industry_backend_go_autumn_2026.git  # один раз
git pull                                        # коммиты бота из форка
git pull --no-rebase --no-edit upstream master  # изменения курса
git push
```

`--no-rebase` обязателен: с `pull.rebase true` git пересадит уже запушенные коммиты, и `git push` будет отклонён как non-fast-forward. Если git сообщил о конфликте, выполните `git merge --abort` и напишите преподавателю.
