package handler

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/top-coins/pkg/models"
)

var ErrInvalidCoinName = errors.New("invalid coin name")

func (cr *CustomRouter) getCoinLastInfo() http.HandlerFunc {
	const fnName = "get last coin info"

	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := r.Body.Close(); err != nil {
				cr.log.Error("Can't close body", cr.log.Attr("error", err), cr.log.Attr("func", fnName))
			}
		}()

		rawCoinName := r.URL.Query().Get("name")
		cr.log.Debug("Last coin info requested", cr.log.Attr("func", fnName), cr.log.Attr("name", rawCoinName))

		ctx, cancel := context.WithTimeout(context.Background(), cr.timeout)
		defer cancel()

		parsedCoinName, err := parseCoinName(rawCoinName)
		if err != nil {
			cr.log.Error("Can't parse coin name", cr.log.Attr("error", err), cr.log.Attr("func", fnName), cr.log.Attr("coin_name", rawCoinName))

			cr.mustSendResponseJSON(w,
				http.StatusBadRequest,
				map[string]any{"error": http.StatusText(http.StatusBadRequest)},
			)

			return
		}

		coin, err := cr.coinClient.CoinLastInfo(ctx, string(parsedCoinName))
		if err != nil {
			cr.log.Error("Can't get last coin info", cr.log.Attr("error", err), cr.log.Attr("func", fnName))

			cr.mustSendResponseJSON(w,
				http.StatusInternalServerError,
				map[string]any{"error": http.StatusText(http.StatusInternalServerError)},
			)

			return
		}

		cr.log.Debug("Last coin info sent", cr.log.Attr("func", fnName), cr.log.Attr("name", parsedCoinName))

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

		rawCoinName := r.URL.Query().Get("name")
		rawSinceTime := r.URL.Query().Get("since")
		cr.log.Debug("Coin history requested", cr.log.Attr("func", fnName), cr.log.Attr("name", rawCoinName), cr.log.Attr("since", rawSinceTime))

		ctx, cancel := context.WithTimeout(context.Background(), cr.timeout)
		defer cancel()

		parsedCoinName, err := parseCoinName(rawCoinName)
		if err != nil {
			cr.log.Error("Can't parse coin name", cr.log.Attr("error", err), cr.log.Attr("func", fnName), cr.log.Attr("coin_name", rawCoinName))

			cr.mustSendResponseJSON(w,
				http.StatusBadRequest,
				map[string]any{"error": http.StatusText(http.StatusBadRequest)},
			)

			return
		}

		parcedSinceTime, err := parseSinceTime(rawSinceTime)
		if err != nil {
			cr.log.Error("Can't parse since time", cr.log.Attr("error", err), cr.log.Attr("func", fnName), cr.log.Attr("since_time", rawSinceTime))

			cr.mustSendResponseJSON(w,
				http.StatusBadRequest,
				map[string]any{"error": http.StatusText(http.StatusBadRequest)},
			)

			return
		}

		coins, err := cr.coinClient.CoinHistory(ctx, string(parsedCoinName), parcedSinceTime)
		if err != nil {
			cr.log.Error("Can't get coin history", cr.log.Attr("error", err), cr.log.Attr("func", fnName))

			cr.mustSendResponseJSON(w,
				http.StatusInternalServerError,
				map[string]any{"error": http.StatusText(http.StatusInternalServerError)},
			)

			return
		}

		cr.log.Debug("Coin history sent", cr.log.Attr("func", fnName), cr.log.Attr("name", parsedCoinName), cr.log.Attr("count", len(coins)))

		cr.mustSendResponseJSON(w,
			http.StatusOK,
			map[string]any{"coins_history": coins},
		)
	}
}

func parseCoinName(rawName string) (models.CoinName, error) {
	const fnName = "parse coin name"

	normalized := models.CoinName(strings.ToUpper(strings.TrimSpace(rawName)))
	if !normalized.IsValid() {
		return "", errors.Wrap(ErrInvalidCoinName, fnName)
	}
	return normalized, nil
}

func parseSinceTime(strTime string) (time.Time, error) {
	const fnName = "parse since time"

	sinceTimestamp, err := strconv.Atoi(strTime)
	if err != nil {
		return time.Time{}, errors.Wrap(err, fnName)
	}

	sinceDate := time.Unix(int64(sinceTimestamp), 0)

	return sinceDate, nil
}
