package expression

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/MonoBear123/MyApi/back/model"
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
func AgentIDKey(id uint64) string {
	return fmt.Sprintf("agent:%d", id)

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
func (r *RedisRepo) AgentALLFind(ctx context.Context) []model.Requert {
	keys, err := r.Client.Keys(ctx, "agent:*").Result()
	if err != nil {

		return nil
	}

	var allStatuses []model.Requert

	for _, key := range keys {
		// Получаем значение по ключу
		value, err := r.Client.Get(ctx, key).Result()
		if err != nil {

			continue
		}

		// Распаковываем JSON в структуру AgentStatus
		var agentStatus model.Requert
		err = json.Unmarshal([]byte(value), &agentStatus)
		if err != nil {

			continue
		}

		allStatuses = append(allStatuses, agentStatus)
	}
	return allStatuses
}

func (r *RedisRepo) AgentInsert(ctx context.Context, agent model.Requert) error {

	key := AgentIDKey(agent.Id)
	res, err := json.Marshal(agent)
	if err != nil {
		return err
	}
	err = r.Client.Set(ctx, key, res, 0).Err()
	if err != nil {
		return err
	}
	fmt.Println("выражение добавлено в базу данных")
	return nil
}
func (r *RedisRepo) DeleteExpression(name string) error {
	err := r.Client.Del(context.Background(), name).Err()
	if err != nil {

		return err
	}
	return nil
}
func (r *RedisRepo) DeleteAgent(name string) error {
	err := r.Client.Del(context.Background(), name).Err()
	if err != nil {

		return err
	}
	return nil
}
