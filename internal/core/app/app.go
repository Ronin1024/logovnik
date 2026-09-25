package app

import (
	"fmt"
	"go_logovnik/internal/core/logovnik"
	webviewmanager "go_logovnik/internal/core/webview_manager"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"go_logovnik/pkg/logger"

	"github.com/jchv/go-webview-selector"
)

var logov *logovnik.Logovnik

func Run() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	// Определяем порт (можно использовать фиксированный или динамический)
	port := "8084"

	// Получаем путь к папке frontend
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	exeDir := filepath.Dir(exePath)
	frontendDir := filepath.Join(exeDir, "frontend")

	// Создаем HTTP-сервер с файловым сервером
	fs := http.FileServer(http.Dir(frontendDir))
	// Создаем обработчик, который добавляет заголовки для отключения кэша
	noCacheHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Устанавливаем заголовки, запрещающие кэширование
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")

		// Передаем управление файловому серверу
		fs.ServeHTTP(w, r)
	})

	http.Handle("/", noCacheHandler)

	// Запускаем сервер в горутине
	go func() {
		logger.Info("MAIN", "Сервер запущен на http://localhost:"+port)
		logger.Info("MAIN", "Папка frontend: "+frontendDir)

		if err := http.ListenAndServe(":"+port, nil); err != nil {
			logger.Error("MAIN", "Ошибка запуска сервера: ", err)
		}
	}()

	// Даем серверу время на запуск
	time.Sleep(100 * time.Millisecond)

	// Создаем WebView и направляем его на сервер
	debug := true
	wm := webviewmanager.NewWebViewManager(debug)

	wm.WebView.Bind("check_master_key", func(master_key string) string {
		logov, err = logovnik.NewLogovnik("data.dat", master_key)
		if err != nil {
			return "fail"
		}

		return "ok"

	})

	wm.WebView.Bind("delete_item", func(id string) string {
		err := logov.DeleteItem(id)
		if err != nil {
			return "fail"
		}
		return "ok"
	})

	wm.WebView.Bind("add_item", func(item string) string {
		// logger.Info("add_item ", item)
		resp, err := logov.AddItem(item)
		if err != nil {
			return "fail"
		}
		return resp
	})

	wm.WebView.Bind("update_item", func(item string) string {
		// logger.Info("update_item ", item)
		resp, err := logov.UpdateItem(item)
		if err != nil {
			return "fail"
		}
		return resp
	})

	wm.WebView.Bind("get_items", func() string {
		// logger.Info("APP", "get_items")
		items, err := logov.GetItems()
		if err != nil {
			return `{"status":"error"}`
		}
		// logger.Info("APP", items)
		return items
	})

	wm.WebView.SetTitle("Логовник Webview app")
	// w, h := common.GetScreenResolution()
	wm.WebView.SetSize(800, 750, webview.HintNone)
	wm.WebView.Navigate(fmt.Sprintf("http://localhost:%s", port))

	// Запускаем WebView (этот вызов блокирует главный поток)
	wm.Run()
	if logov != nil {
		logov.Close()
	}
	defer wm.Close()
}
