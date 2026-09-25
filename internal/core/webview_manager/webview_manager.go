package webviewmanager

import (
	"fmt"
	"sync"
	"time"

	"github.com/jchv/go-webview-selector"
)

type WebViewManager struct {
	WebView webview.WebView
	mu      sync.Mutex
	queue   []string
	cond    *sync.Cond
	closed  bool
}

func NewWebViewManager(debug bool) *WebViewManager {

	wm := &WebViewManager{
		queue: make([]string, 0),
	}
	wm.WebView = webview.New(debug)
	wm.cond = sync.NewCond(&wm.mu)
	return wm
}

// Безопасный вызов Eval из любого потока
func (wm *WebViewManager) SafeEval(js string) {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	if wm.closed {
		return
	}

	wm.queue = append(wm.queue, js)
	wm.cond.Signal()
}

// Запускает обработчик в главном потоке
func (wm *WebViewManager) Run() {
	go func() {
		ticker := time.NewTicker(1 * time.Microsecond)
		defer ticker.Stop()

		for {
			wm.mu.Lock()
			for len(wm.queue) == 0 && !wm.closed {
				wm.cond.Wait()
			}

			if wm.closed {
				wm.mu.Unlock()
				return
			}

			// Забираем все задачи из очереди
			queue := wm.queue
			wm.queue = make([]string, 0)
			wm.mu.Unlock()

			// Выполняем все задачи
			for _, js := range queue {
				// logger.Debug("L", "Exec ", js)
				wm.WebView.Dispatch(func() {
					jscode := fmt.Sprintf(`try{%s}catch(e){console.log(e);}`, js)
					wm.WebView.Eval(jscode)
				})
			}
		}
	}()
	wm.WebView.Run()
}

func (wm *WebViewManager) Close() {
	wm.mu.Lock()
	wm.closed = true
	wm.mu.Unlock()
	wm.cond.Signal()
	wm.WebView.Destroy()
}
