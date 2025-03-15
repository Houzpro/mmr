package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/xuri/excelize/v2"
)

// Максимальное количество комбинаций для расчета, чтобы избежать проблем с памятью и производительностью
const maxCombinationsLimit = 1000000

func main() {
	// Проверка, указан ли путь к файлу
	if len(os.Args) < 2 {
		log.Fatal("Использование: go run main.go <путь_к_excel_файлу> [имя_столбца] [макс_комбинаций]")
	}

	filePath := os.Args[1]
	columnName := "A" // Столбец по умолчанию
	if len(os.Args) > 2 {
		columnName = os.Args[2]
	}

	// Установка максимального количества комбинаций из командной строки или использование значения по умолчанию
	maxCombinations := maxCombinationsLimit
	if len(os.Args) > 3 {
		var err error
		maxCombinations, err = strconv.Atoi(os.Args[3])
		if err != nil {
			maxCombinations = maxCombinationsLimit
		}
	}

	// Открытие файла Excel
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		log.Fatalf("Не удается открыть файл: %v", err)
	}
	defer f.Close()

	// Получение всех имен листов
	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		log.Fatal("В файле Excel не найдены листы")
	}

	// Использование первого листа
	sheetName := sheets[0]

	// Чтение всех строк в указанном столбце
	rows, err := f.GetRows(sheetName)
	if err != nil {
		log.Fatalf("Не удается прочитать строки: %v", err)
	}

	var sum float64
	var count int
	var mmrValues []float64

	// Поиск индекса столбца по имени
	colIndex := columnNameToIndex(columnName)

	// Извлечение значений MMR
	for _, row := range rows {
		if len(row) > colIndex {
			// Попытка преобразовать значение ячейки в число
			mmrStr := row[colIndex]
			mmr, err := strconv.ParseFloat(mmrStr, 64)
			if err == nil {
				mmrValues = append(mmrValues, mmr)
				sum += mmr
				count++
			}
		}
	}

	// Расчет начального среднего значения
	if count == 0 {
		fmt.Println("В указанном столбце не найдены допустимые значения MMR")
		return
	}

	initialMean := sum / float64(count)
	roundedInitialMean := math.Round(initialMean)
	fmt.Printf("Исходные значения MMR: %v\n", mmrValues)
	fmt.Printf("Количество исходных значений MMR: %d\n", count)
	fmt.Printf("Исходное среднее MMR: %.2f\n", initialMean)
	fmt.Printf("Округленное среднее MMR: %.0f\n", roundedInitialMean)

	// Создание нового скорректированного набора данных
	var adjustedMmrValues []float64
	var belowMeanCount int

	// Сначала добавляем все значения, которые больше или равны среднему
	for _, mmr := range mmrValues {
		if mmr >= roundedInitialMean {
			adjustedMmrValues = append(adjustedMmrValues, mmr)
		} else {
			belowMeanCount++
		}
	}

	// Добавление одного экземпляра среднего значения, чтобы представить все значения ниже среднего
	if belowMeanCount > 0 {
		adjustedMmrValues = append(adjustedMmrValues, roundedInitialMean)
	}

	// Расчет скорректированной суммы и среднего значения
	var adjustedSum float64
	for _, mmr := range adjustedMmrValues {
		adjustedSum += mmr
	}

	adjustedCount := len(adjustedMmrValues)
	adjustedMean := adjustedSum / float64(adjustedCount)

	fmt.Printf("\nСкорректированный набор данных (все %d значений ниже %.0f заменены на один экземпляр): %v\n",
		belowMeanCount, roundedInitialMean, adjustedMmrValues)
	fmt.Printf("Количество элементов в скорректированном наборе данных: %d\n", adjustedCount)
	fmt.Printf("Скорректированное среднее MMR: %.2f\n", adjustedMean)

	// Расчет комбинаций из 5 игроков
	fmt.Println("\n--- Анализ команды (с использованием скорректированных значений MMR) ---")
	if len(adjustedMmrValues) < 5 {
		fmt.Println("Для формирования команды необходимо минимум 5 игроков")
	} else {
		// Расчет общего количества возможных комбинаций
		totalPossibleCombinations := calculateBinomialCoefficient(len(adjustedMmrValues), 5)
		fmt.Printf("Всего возможных команд из 5 игроков: %d\n", totalPossibleCombinations)
		fmt.Printf("Использование максимум %d комбинаций для анализа\n", maxCombinations)

		// Генерация ограниченного числа комбинаций с использованием скорректированных значений MMR
		teamInfos := generateLimitedTeamInfos(adjustedMmrValues, 5, maxCombinations)

		// Добавление гарантированно максимальной по MMR команды
		maxTeamInfo := findMaximumTeam(adjustedMmrValues, 5)
		teamInfos = append(teamInfos, maxTeamInfo)

		// Сортировка информации о командах по сумме MMR
		sort.Slice(teamInfos, func(i, j int) bool {
			return teamInfos[i].mmrSum < teamInfos[j].mmrSum
		})

		// Расчет статистики
		var totalSum float64
		for _, info := range teamInfos {
			totalSum += info.mmrSum
		}
		meanTeamMmr := totalSum / float64(len(teamInfos))
		minTeamMmr := teamInfos[0].mmrSum
		maxTeamMmr := teamInfos[len(teamInfos)-1].mmrSum

		// Расчет дисперсии
		var sumSquaredDiff float64
		for _, info := range teamInfos {
			diff := info.mmrSum - meanTeamMmr
			sumSquaredDiff += diff * diff
		}
		variance := sumSquaredDiff / float64(len(teamInfos))
		stdDev := math.Sqrt(variance)

		fmt.Printf("Статистика MMR команд (на основе %d образцов):\n", len(teamInfos))
		fmt.Printf("  Среднее MMR команды: %.2f\n", meanTeamMmr)
		fmt.Printf("  Минимальное MMR команды: %.2f\n", minTeamMmr)
		fmt.Printf("  Максимальное MMR команды: %.2f\n", maxTeamMmr)
		fmt.Printf("  Разница между максимальным и минимальным: %.2f\n", maxTeamMmr-minTeamMmr)
		fmt.Printf("  Стандартное отклонение MMR: %.2f\n", stdDev)

		// Показать примеры комбинаций команд (первые 5 и последние 5)
		if len(teamInfos) > 0 {
			fmt.Println("\nПримеры комбинаций команд:")
			fmt.Println("Команды с низким MMR:")
			for i := 0; i < min(1, len(teamInfos)); i++ {
				fmt.Printf("  Команда с MMR %.2f: %v\n", teamInfos[i].mmrSum, teamInfos[i].team)
			}

			fmt.Println("Команды с высоким MMR:")
			// Начинаем с индекса, который даст максимум 5 последних элементов
			startIndex := max(0, len(teamInfos)-5)
			for i := startIndex; i < len(teamInfos); i++ {
				fmt.Printf("  Команда с MMR %.2f: %v\n", teamInfos[i].mmrSum, teamInfos[i].team)
			}

			// Явно показать команду с максимальным MMR (это должна быть последняя в отсортированном списке)
			maxTeam := teamInfos[len(teamInfos)-1]
			fmt.Printf("\nКоманда с максимальным MMR %.2f: %v\n", maxTeam.mmrSum, maxTeam.team)
		}
	}
}

// TeamInfo хранит состав команды и сумму MMR
type TeamInfo struct {
	team   []float64
	mmrSum float64
}

// Преобразование имени столбца типа "A", "B", "AA" в индекс с нулевой базой
func columnNameToIndex(name string) int {
	result := 0
	for i := 0; i < len(name); i++ {
		result = result*26 + int(name[i]-'A'+1)
	}
	return result - 1
}

// Вычисление биномиального коэффициента (n выбрать k)
func calculateBinomialCoefficient(n, k int) int {
	if k > n-k {
		k = n - k
	}

	result := 1
	for i := 0; i < k; i++ {
		result = result * (n - i) / (i + 1)
	}

	return result
}

// Генерация ограниченного количества комбинаций и возврат информации о командах
func generateLimitedTeamInfos(values []float64, k int, maxCombinations int) []TeamInfo {
	totalCombinations := calculateBinomialCoefficient(len(values), k)

	// Если возможных комбинаций меньше, чем наш лимит, генерируем все
	if totalCombinations <= maxCombinations {
		var result []TeamInfo
		combinationHelperWithSum(values, k, 0, []float64{}, &result)
		return result
	} else {
		// Используем случайную выборку, если комбинаций слишком много
		return generateRandomTeamInfos(values, k, maxCombinations)
	}
}

// Генерация случайных комбинаций для ограничения вычислений
func generateRandomTeamInfos(values []float64, k int, count int) []TeamInfo {
	rand.Seed(time.Now().UnixNano())
	n := len(values)
	result := make([]TeamInfo, count)

	for i := 0; i < count; i++ {
		// Генерация случайной комбинации
		indices := rand.Perm(n)[:k]
		team := make([]float64, k)
		teamSum := 0.0

		for j, idx := range indices {
			team[j] = values[idx]
			teamSum += values[idx]
		}

		result[i] = TeamInfo{
			team:   team,
			mmrSum: teamSum,
		}
	}

	return result
}

// Вспомогательная функция для генерации комбинаций с суммами MMR
func combinationHelperWithSum(values []float64, k int, start int, current []float64, result *[]TeamInfo) {
	if len(current) == k {
		// Создать копию текущей комбинации
		combination := make([]float64, k)
		copy(combination, current)

		// Рассчитать сумму
		sum := 0.0
		for _, mmr := range combination {
			sum += mmr
		}

		// Добавить в результат
		*result = append(*result, TeamInfo{
			team:   combination,
			mmrSum: sum,
		})
		return
	}

	for i := start; i < len(values); i++ {
		current = append(current, values[i])
		combinationHelperWithSum(values, k, i+1, current, result)
		current = current[:len(current)-1]
	}
}

// Найти команду с максимально возможным MMR (топ-5 игроков по MMR)
func findMaximumTeam(values []float64, k int) TeamInfo {
	// Создаем копию значений
	valuesCopy := make([]float64, len(values))
	copy(valuesCopy, values)

	// Сортируем в порядке убывания
	sort.Slice(valuesCopy, func(i, j int) bool {
		return valuesCopy[i] > valuesCopy[j]
	})

	// Берем первые k значений (топ-k игроков по MMR)
	topTeam := valuesCopy[:min(k, len(valuesCopy))]

	// Рассчитываем сумму MMR
	sum := 0.0
	for _, mmr := range topTeam {
		sum += mmr
	}

	return TeamInfo{
		team:   topTeam,
		mmrSum: sum,
	}
}

// Функция Min для целых чисел
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Функция Max для целых чисел
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
