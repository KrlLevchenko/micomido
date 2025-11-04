package repository

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/KrlLevchenko/micomido/internal/models"
)

type MealRepository interface {
	GetMeals(ctx context.Context, from, to *time.Time) ([]models.Meal, error)
	CreateMeal(ctx context.Context, meal *models.Meal) error
	// DeleteMeal returns the number of affected rows and an error.
	DeleteMeal(ctx context.Context, id string) (int64, error)
	AddMealPhoto(ctx context.Context, mealID, photoID string) error
	DeleteMealPhoto(ctx context.Context, photoID string) (int64, error)
}

type mysqlMealRepository struct {
	db *sql.DB
}

func NewMysqlMealRepository(db *sql.DB) MealRepository {
	return &mysqlMealRepository{db: db}
}

func (r *mysqlMealRepository) GetMeals(ctx context.Context, from, to *time.Time) ([]models.Meal, error) {
	query := "SELECT id, at, comment FROM meal"
	args := []interface{}{}
	conditions := []string{}

	if from != nil {
		conditions = append(conditions, "at >= ?")
		args = append(args, *from)
	}
	if to != nil {
		conditions = append(conditions, "at <= ?")
		args = append(args, *to)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY at DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var meals []models.Meal
	for rows.Next() {
		var meal models.Meal
		if err := rows.Scan(&meal.ID, &meal.At, &meal.Comment); err != nil {
			return nil, err
		}
		meals = append(meals, meal)
	}

	return meals, rows.Err()
}

func (r *mysqlMealRepository) CreateMeal(ctx context.Context, meal *models.Meal) error {
	query := "INSERT INTO meal (id, comment, at) VALUES (?, ?, ?)"
	_, err := r.db.ExecContext(ctx, query, meal.ID, meal.Comment, meal.At)
	return err
}

func (r *mysqlMealRepository) DeleteMeal(ctx context.Context, id string) (int64, error) {
	query := "DELETE FROM meal WHERE id = ?"
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (r *mysqlMealRepository) AddMealPhoto(ctx context.Context, mealID, photoID string) error {
	query := "INSERT INTO meal_photo (id, meal_id) VALUES (?, ?)"
	_, err := r.db.ExecContext(ctx, query, photoID, mealID)
	return err
}

func (r *mysqlMealRepository) DeleteMealPhoto(ctx context.Context, photoID string) (int64, error) {
	query := "DELETE FROM meal_photo WHERE id = ?"
	result, err := r.db.ExecContext(ctx, query, photoID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
