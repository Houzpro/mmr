# Реформирование системы ограничений на турнирах: от лимита сборников к балансу MMR

## Введение

В современных реалиях во внутривузком берспортивном сообществе вопрос сбалансированного формирования команд на турнирах остаётся открытым. Традиционный подход с ограничением количества "сборников" (игроков высокого ранга) в команде имеет ряд существенных недостатков. В данной статье предлагается альтернативный метод — внедрение ограничения на суммарный MMR команды, который способен обеспечить более справедливый баланс сил и разнообразие составов.

## Проблемы существующей системы

Текущая система ограничений по количеству сборников сталкивается с несколькими принципиальными проблемами:

1. **Бинарная классификация игроков**: деление на "сборников" и "не-сборников" создаёт искусственный порог, не отражающий плавную градацию навыков. Игрок с рейтингом чуть выше рейтига "сборника" не приравнивается к ним а считается обычным игроком.

2. **Неравномерность составов**: команды с одинаковым количеством сборников могут иметь радикально разный уровень силы. Например, команда с двумя игроками на нижней границе категории "сборник" значительно слабее команды с двумя топ-игроками.

## Преимущества системы ограничений по суммарному MMR

Переход к системе с ограничением суммарного MMR команды предлагает следующие улучшения:

1. **Гибкость формирования составов**: команды смогут включать как высокоранговых, так и игроков среднего уровня в различных комбинациях, не выходя за рамки общего ограничения.

2. **Справедливый баланс**: команды с одинаковой суммой MMR теоретически должны иметь сопоставимый уровень навыков.

3. **Снижение мотивации к манипуляциям**: поскольку учитывается каждый пункт MMR, преимущество от незначительного снижения рейтинга минимально.

## Механизм реализации

Реализация системы ограничений по суммарному MMR может выглядеть следующим образом:

1. **Определение базового MMR**: использование официального рейтинга игры или специальной формулы для расчёта индивидуального MMR каждого игрока.

2. **Установление лимита для турнира**: организаторы определяют максимально допустимую сумму MMR для команды исходя из уровня турнира.

3. **Верификация составов**: перед началом турнира проверяется соответствие суммарного MMR команды установленному лимиту.

## Примеры применения

**Пример 1**: Турнир с ограничением в 25,000 MMR на команду из 5 человек.
- Команда А: 6000 + 5500 + 5000 + 4500 + 4000 = 25,000 MMR (сбалансированный состав)
- Команда Б: 7000 + 6000 + 5000 + 4000 + 3000 = 25,000 MMR (два сильных игрока + поддержка)
- Команда В: 8000 + 7000 + 4000 + 3000 + 3000 = 25,000 MMR (два топ-игрока + игроки поддержки)

Все три команды теоретически сбалансированы по общей силе, но используют разные  подходы к формированию состава.

## Решение потенциальных проблем

1. **Определение точного MMR**: возможно использование комбинированной метрики из официального рейтинга и дополнительных показателей.

2. **Смурф-аккаунты**: ужесточение правил верификации игроков и регулярная проверка истории аккаунтов.

3. **Специализация ролей**: внедрение коэффициентов для различных позиций/ролей в игре, если это необходимо для баланса.

