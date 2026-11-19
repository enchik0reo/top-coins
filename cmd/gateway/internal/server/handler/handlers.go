package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"
)

func (cr *CustomRouter) getCoinLastInfo() http.HandlerFunc {
	const fnName = "get last coin info"

	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := r.Body.Close(); err != nil {
				cr.log.Error("Can't close body", cr.log.Attr("error", err), cr.log.Attr("func", fnName))
			}
		}()

		id := r.URL.Query().Get("id")
		cr.log.Debug("Last coin info requested", cr.log.Attr("func", fnName), cr.log.Attr("id", id))

		ctx, cancel := context.WithTimeout(context.Background(), cr.timeout)
		defer cancel()

		coin, err := cr.coinClient.CoinLastInfo(ctx, id)
		if err != nil {
			cr.log.Error("Can't get last coin info", cr.log.Attr("error", err), cr.log.Attr("func", fnName))

			cr.mustSendResponseJSON(w,
				http.StatusInternalServerError,
				map[string]any{"error": http.StatusText(http.StatusInternalServerError)},
			)

			return
		}

		cr.log.Debug("Last coin info sent", cr.log.Attr("func", fnName), cr.log.Attr("id", id))

		cr.mustSendResponseJSON(w,
			http.StatusOK,
			map[string]any{"coin": coin},
		)
	}
}

func (cr *CustomRouter) getCoinHistory() http.HandlerFunc {
	const fnName = "get coin history"

	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := r.Body.Close(); err != nil {
				cr.log.Error("Can't close body", cr.log.Attr("error", err), cr.log.Attr("func", fnName))
			}
		}()

		id := r.URL.Query().Get("id")
		since := r.URL.Query().Get("since")
		cr.log.Debug("Coin history requested", cr.log.Attr("func", fnName), cr.log.Attr("id", id), cr.log.Attr("since", since))

		sinceTimestamp, err := strconv.Atoi(since)
		if err != nil {
			cr.log.Error("Can't convert time to int", cr.log.Attr("error", err), cr.log.Attr("func", fnName))

			cr.mustSendResponseJSON(w,
				http.StatusBadRequest,
				map[string]any{"error": http.StatusText(http.StatusBadRequest)},
			)

			return
		}

		sinceDate := time.Unix(int64(sinceTimestamp), 0)

		ctx, cancel := context.WithTimeout(context.Background(), cr.timeout)
		defer cancel()

		coins, err := cr.coinClient.CoinHistory(ctx, id, sinceDate)
		if err != nil {
			cr.log.Error("Can't get coin history", cr.log.Attr("error", err), cr.log.Attr("func", fnName))

			cr.mustSendResponseJSON(w,
				http.StatusInternalServerError,
				map[string]any{"error": http.StatusText(http.StatusInternalServerError)},
			)

			return
		}

		cr.log.Debug("Coin history sent", cr.log.Attr("func", fnName), cr.log.Attr("id", id), cr.log.Attr("count", len(coins)))

		cr.mustSendResponseJSON(w,
			http.StatusOK,
			map[string]any{"coins_history": coins},
		)
	}
}
