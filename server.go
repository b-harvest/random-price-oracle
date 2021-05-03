package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/mongo"
)

type Server struct {
	*echo.Echo
	cfg ServerConfig
	ss  *StoreService
}

func NewServer(cfg ServerConfig, ss *StoreService) *Server {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	s := &Server{Echo: e, cfg: cfg, ss: ss}
	s.RegisterRoutes()

	return s
}

func (s *Server) RegisterRoutes() {
	s.GET("/prices", s.GetPrices)
}

func (s *Server) GetPrices(c echo.Context) error {
	var req GetPricesRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	symbols := make(map[string]struct{})
	for _, symbol := range req.Symbols {
		symbol = NormalizeSymbol(symbol)
		if symbol == "" {
			continue
		}
		symbols[symbol] = struct{}{}
	}
	if len(symbols) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "no symbols specified")
	}
	resp := GetPricesResponse{
		Coins: make(map[string]GetPricesResponseCoin),
	}
	for symbol := range symbols {
		coin, err := s.ss.Coin(c.Request().Context(), symbol)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("no price for %q", symbol))
			}
			return fmt.Errorf("get price: %w", err)
		}
		resp.Coins[symbol] = GetPricesResponseCoin{
			Price:     coin.Price,
			UpdatedAt: coin.UpdatedAt,
		}
	}
	return c.JSON(http.StatusOK, resp)
}

func (s *Server) StartUpdater(ctx context.Context) error {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	ticker := time.NewTicker(s.cfg.PriceUpdateInterval)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			coins, err := s.ss.Coins(ctx)
			if err != nil {
				return fmt.Errorf("get coins: %w", err)
			}
			for _, coin := range coins {
				changeRatio := RandomPriceChangeRatio(r)
				nextPrice := coin.Price * changeRatio
				if err := s.ss.SetPrice(ctx, coin.Symbol, nextPrice); err != nil {
					return fmt.Errorf("set price: %w", err)
				}
			}
		}
	}
}

func (s *Server) ShutdownWithTimeout(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return s.Shutdown(ctx)
}
