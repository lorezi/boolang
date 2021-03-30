package controllers

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/lorezi/boolang/pkg/metric"
)

var (
	codes = []int{
		http.StatusOK,
		http.StatusCreated,
		http.StatusBadRequest,
		http.StatusNotFound,
		http.StatusInternalServerError,
		http.StatusServiceUnavailable,
	}

	clients   = []string{"client_1", "client_2", "client_3"}
	amounts   = []string{"5.51", "-2.02", "14.00", "-17.00", "71.17", "59.01", "9.99", "-19.55", "8.00"}
	durations = []int{32, 43, 55, 67, 76, 88, 99, 101, 106, 109, 111, 120, 131, 135, 148, 160, 180, 203, 210, 211}
)

type BalanceUpdate struct {
	metric metric.Metric
}

func NewBalanceUpdate(metric metric.Metric) BalanceUpdate {
	return BalanceUpdate{
		metric: metric,
	}
}

func (b BalanceUpdate) Handle(w http.ResponseWriter, r *http.Request) {
	rand.Seed(time.Now().UnixNano())

	// Used for histogram
	start := time.Now()

	// randomly pick a client
	client := clients[rand.Intn(len(clients))]

	// randomly decide how long it would take to respond.
	duration := durations[rand.Intn(len(durations))]

	time.Sleep(time.Duration(duration) * time.Millisecond)

	// randomly decide what response code we would return
	code := codes[rand.Intn(len(codes))]

	if code != http.StatusOK {
		b.metric.HttpResponseCounter.WithLabelValues("balance_update", fmt.Sprintf("%d", code)).Inc()
		dur := float64(time.Since(start).Microseconds())
		b.metric.ResponseDurationHistogram.WithLabelValues("balance_update").Observe(dur)
		dump(code, dur, "n/a", client)
		return
	}

	// randomly decide what amount would go in/out of balance
	amount := amounts[rand.Intn(len(amounts))]

	if amount[0:1] == "-" {
		a, _ := strconv.ParseFloat(amount[1:], 64)
		b.metric.BalanceGauge.Sub(a)
		b.metric.BalanceActivityCounter.WithLabelValues("down", client).Inc()
	} else {
		a, _ := strconv.ParseFloat(amount, 64)
		b.metric.BalanceGauge.Add(a)
		b.metric.BalanceActivityCounter.WithLabelValues("up", client).Inc()
	}

	b.metric.HttpResponseCounter.WithLabelValues("balance_update", fmt.Sprintf("%d", code)).Inc()
	dur := float64(time.Since(start).Milliseconds())
	b.metric.ResponseDurationHistogram.WithLabelValues("balance_update").Observe(dur)
	dump(code, dur, amount, client)

}

func dump(code int, dur float64, amount string, client string) {
	fmt.Println("CODE:", code, "AMOUNT:", amount, "DURATION:", dur, "CLIENT:", client)
}
