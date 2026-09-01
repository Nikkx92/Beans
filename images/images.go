package images

import (
	"embed"
	"image"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"voicer/internal/storage"

	"gioui.org/op/paint"
	"gioui.org/widget"
)

//go:embed icons
var IconsFS embed.FS

//go:embed flags
var FlagsFS embed.FS

func flag(i image.Image) *widget.Image {
	return &widget.Image{
		Src: paint.NewImageOp(i),
		Fit: widget.Fill,
	}
}
func icon(i image.Image) *widget.Image {
	return &widget.Image{
		Src: paint.NewImageOp(i),
		Fit: widget.Fill,
	}
}

func LoadFlags(fsys embed.FS) map[string]*widget.Image {
	entries, err := fsys.ReadDir("flags")
	if err != nil {
		slog.Info(err.Error())
	}

	flags := make(map[string]*widget.Image, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		f, err := fsys.Open("flags/" + entry.Name())
		if err != nil {
			slog.Info(err.Error())
			continue
		}

		img, _, err := image.Decode(f)
		f.Close()
		if err != nil {
			slog.Info(err.Error())
			continue
		}
		title := strings.ToUpper(strings.TrimSuffix(entry.Name(), ".png"))

		flags[title] = flag(img)

	}

	return flags
}

func LoadAiCreatedFlags(m map[string]*widget.Image) {
	// 1. Получаем базовый путь к хранилищу
	p, err := storage.Path()
	if err != nil {
		slog.Error("Ошибка получения пути: " + err.Error())
		return
	}

	flagsDirPath := filepath.Join(p, "flags")

	// 2. Читаем содержимое директории
	files, err := os.ReadDir(flagsDirPath)
	if err != nil {
		slog.Error("Ошибка чтения папки с флагами: " + err.Error())
		return
	}

	// Инициализируем карту для хранения результата
	//flagsMap := make(map[string]*widget.Image)

	// 3. Перебираем все файлы в папке
	for _, file := range files {
		// Пропускаем папки, если они там есть, и обрабатываем только файлы
		if file.IsDir() {
			continue
		}

		fileName := file.Name()

		// Проверяем, что файл имеет расширение .png (в нижнем регистре)
		if !strings.HasSuffix(strings.ToLower(fileName), ".png") {
			continue
		}

		// Открываем конкретный файл флага
		filePath := filepath.Join(flagsDirPath, fileName)
		f, err := os.Open(filePath)
		if err != nil {
			slog.Info("Не удалось открыть файл " + fileName + ": " + err.Error())
			continue // Переходим к следующему флагу, если этот не открылся
		}

		// Декодируем изображение
		img, _, err := image.Decode(f)
		f.Close() // Закрываем сразу, не откладывая через defer, чтобы не копить открытые файлы в цикле
		if err != nil {
			slog.Info("Ошибка декодирования " + fileName + ": " + err.Error())
			continue
		}

		// Превращаем в формат виджета Gio (используя вашу внутреннюю функцию flag)
		widgetImg := flag(img)

		// Вырезаем расширение .png из имени файла, чтобы использовать как чистый ключ (например, "ru", "us")
		flagKey := strings.TrimSuffix(fileName, filepath.Ext(fileName))

		// Сохраняем в карту
		m[flagKey] = widgetImg
	}
}

func LoadIcons(fsys embed.FS) map[string]*widget.Image {
	entries, err := fsys.ReadDir("icons")
	if err != nil {
		slog.Info(err.Error())
	}

	icons := make(map[string]*widget.Image, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		f, err := fsys.Open("icons/" + entry.Name())
		if err != nil {
			slog.Info(err.Error())
			continue
		}

		img, _, err := image.Decode(f)
		f.Close()
		if err != nil {
			slog.Info(err.Error())
			continue
		}

		icons[strings.TrimSuffix(entry.Name(), ".png")] = icon(img)

	}

	return icons
}

func LoadImage(s string) *widget.Image {
	s = strings.TrimSpace(s)
	p, err := storage.Path()
	if err != nil {
		slog.Info(err.Error())
	}

	//f, err := os.Open("./images/flags/" + s + ".png")
	f, err := os.Open(p + "/flags/" + s + ".png")
	if err != nil {
		slog.Info(err.Error())
		return nil
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		slog.Info(err.Error())
		return nil
	}
	i := flag(img)
	return i
}
