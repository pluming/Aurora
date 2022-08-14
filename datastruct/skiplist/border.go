package skiplist

import (
	"errors"
	"math"
	"strconv"
)

/*
 * ScoreBorder is a struct represents `min` `max` parameter of client command `ZRANGEBYSCORE`
 * can accept:
 *   int or float value, such as 2.718, 2, -2.718, -2 ...
 *   exclusive int or float value, such as (2.718, (2, (-2.718, (-2 ...
 *   infinity: +inf, -inf， inf(same as +inf)
 */

const (
	negativeInf int8 = -1
	positiveInf int8 = 1
)

// ScoreBorder represents range of a float value, including: <, <=, >, >=, +inf, -inf
type ScoreBorder struct {
	Inf     int8
	Value   float64
	Exclude bool
}

// if max.greater(score) then the score is within the upper border
// do not use min.greater()
func (border *ScoreBorder) greater(value float64) bool {
	if border.Inf == negativeInf {
		return false
	} else if border.Inf == positiveInf {
		return true
	}
	if border.Exclude {
		return border.Value > value
	}
	return border.Value >= value
}

//less Or maybe less and eq
func (border *ScoreBorder) less(value float64) bool {
	if border.Inf == negativeInf {
		return true
	} else if border.Inf == positiveInf {
		return false
	}
	if border.Exclude {
		return border.Value < value
	}
	return border.Value <= value
}

var PositiveInfBorder = &ScoreBorder{
	Inf:   positiveInf,
	Value: math.MaxInt64,
}

var NegativeInfBorder = &ScoreBorder{
	Inf:   negativeInf,
	Value: math.MinInt64,
}

// ParseScoreBorder creates ScoreBorder from client arguments
func ParseScoreBorder(s string) (*ScoreBorder, error) {
	if s == "inf" || s == "+inf" {
		return PositiveInfBorder, nil
	}
	if s == "-inf" {
		return NegativeInfBorder, nil
	}
	if s[0] == '(' {
		value, err := strconv.ParseFloat(s[1:], 64)
		if err != nil {
			return nil, errors.New("ERR min or max is not a float")
		}
		return &ScoreBorder{
			Inf:     0,
			Value:   value,
			Exclude: true,
		}, nil
	}
	value, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil, errors.New("ERR min or max is not a float")
	}
	return &ScoreBorder{
		Inf:     0,
		Value:   value,
		Exclude: false,
	}, nil
}
