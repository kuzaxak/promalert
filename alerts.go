package main

import (
	"github.com/prometheus/prometheus/promql"
	"log"
	"strconv"
	"time"
)

func GetPlotTimeRange(alert Alert) (time.Time, time.Duration) {
	var queryTime time.Time
	var duration time.Duration
	if alert.StartsAt.Second() > alert.EndsAt.Second() {
		queryTime = alert.StartsAt
		duration = time.Minute * 20
	} else {
		queryTime = alert.EndsAt
		duration = queryTime.Sub(alert.StartsAt)

		if duration < time.Minute*20 {
			duration = time.Minute * 20
		}
	}
	log.Printf("Querying Time %v Duration: %v", queryTime, duration)
	return queryTime, duration
}

func GetPlotExpr(alertFormula string) []PlotExpr {
	expr, _ := promql.ParseExpr(alertFormula)
	if parenExpr, ok := expr.(*promql.ParenExpr); ok {
		expr = parenExpr.Expr
		log.Printf("Removing redundant brackets: %v", expr.String())
	}

	if binaryExpr, ok := expr.(*promql.BinaryExpr); ok {
		var alertOperator string

		switch binaryExpr.Op {
		case promql.ItemLAND:
			log.Printf("Logical condition, drawing sides separately")
			return append(GetPlotExpr(binaryExpr.LHS.String()), GetPlotExpr(binaryExpr.RHS.String())...)
		case promql.ItemLTE, promql.ItemLSS:
			alertOperator = "<"
		case promql.ItemGTE, promql.ItemGTR:
			alertOperator = ">"
		default:
			log.Printf("Unexpected operator: %v", binaryExpr.Op.String())
			alertOperator = ">"
		}

		alertLevel, _ := strconv.ParseFloat(binaryExpr.RHS.String(), 64)
		return []PlotExpr{PlotExpr{
			Formula:  binaryExpr.LHS.String(),
			Operator: alertOperator,
			Level:    alertLevel,
		}}
	} else {
		log.Printf("Non binary excpression: %v", alertFormula)
		return nil
	}
}
