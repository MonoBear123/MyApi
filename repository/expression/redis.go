package expression

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/MonoBear123/MyApi/model"
	"github.com/redis/go-redis/v9"
)

type RedisRepo struct {
	Client *redis.Client
}

func (r *RedisRepo) Insert(ctx context.Context, expression model.EXpression) error {
	Data, err := json.Marshal(expression)
	if err != nil {
		return err
	}
	key := ExpressionIDKey(expression.ExpressinID)
	res := r.Client.SetNX(ctx, key, string(Data), 0)
	if err = res.Err(); err != nil {
		return err
	}
	return nil
}
func ExpressionIDKey(id uint64) string {
	return fmt.Sprint(id)

}
func (r *RedisRepo) ExpressionFind(ctx context.Context, id uint64) (model.EXpression, error) {
	key := ExpressionIDKey(id)
	value, err := r.Client.Get(ctx, key).Result()
	if err != nil {
		return model.EXpression{}, err
	}
	var out model.EXpression
	err = json.Unmarshal([]byte(value), &out)
	if err != nil {
		return model.EXpression{}, err
	}
	return out, nil
}
