package repository

import (
	"database/sql"

	"github.com/GATEOPENERZ/completionist-api/internal/models"
	"github.com/jmoiron/sqlx"
)

type SubscriptionRepository struct {
	DB *sqlx.DB
}

func NewSubscriptionRepository(db *sqlx.DB) *SubscriptionRepository {
	return &SubscriptionRepository{DB: db}
}

func (r *SubscriptionRepository) Create(sub *models.Subscription) error {
	return r.DB.QueryRowx(`
		INSERT INTO subscriptions (user_id, stripe_customer_id, stripe_subscription_id, stripe_price_id, status, current_period_start, current_period_end)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, user_id, stripe_customer_id, stripe_subscription_id, stripe_price_id, status, current_period_start, current_period_end, cancel_at_period_end, cancelled_at, created_at, updated_at
	`, sub.UserID, sub.StripeCustomerID, sub.StripeSubscriptionID, sub.StripePriceID, sub.Status, sub.CurrentPeriodStart, sub.CurrentPeriodEnd).StructScan(sub)
}

func (r *SubscriptionRepository) GetByUserID(userID int64) (*models.Subscription, error) {
	var sub models.Subscription
	err := r.DB.QueryRowx(`
		SELECT id, user_id, stripe_customer_id, stripe_subscription_id, stripe_price_id, status, current_period_start, current_period_end, cancel_at_period_end, cancelled_at, created_at, updated_at
		FROM subscriptions WHERE user_id = $1
	`, userID).StructScan(&sub)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

func (r *SubscriptionRepository) GetByStripeCustomerID(customerID string) (*models.Subscription, error) {
	var sub models.Subscription
	err := r.DB.QueryRowx(`
		SELECT id, user_id, stripe_customer_id, stripe_subscription_id, stripe_price_id, status, current_period_start, current_period_end, cancel_at_period_end, cancelled_at, created_at, updated_at
		FROM subscriptions WHERE stripe_customer_id = $1
	`, customerID).StructScan(&sub)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

func (r *SubscriptionRepository) GetByStripeSubscriptionID(subID string) (*models.Subscription, error) {
	var sub models.Subscription
	err := r.DB.QueryRowx(`
		SELECT id, user_id, stripe_customer_id, stripe_subscription_id, stripe_price_id, status, current_period_start, current_period_end, cancel_at_period_end, cancelled_at, created_at, updated_at
		FROM subscriptions WHERE stripe_subscription_id = $1
	`, subID).StructScan(&sub)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

func (r *SubscriptionRepository) Update(sub *models.Subscription) error {
	_, err := r.DB.Exec(`
		UPDATE subscriptions SET
			stripe_subscription_id = $2, stripe_price_id = $3, status = $4,
			current_period_start = $5, current_period_end = $6,
			cancel_at_period_end = $7, cancelled_at = $8, updated_at = NOW()
		WHERE id = $1
	`, sub.ID, sub.StripeSubscriptionID, sub.StripePriceID, sub.Status, sub.CurrentPeriodStart, sub.CurrentPeriodEnd, sub.CancelAtPeriodEnd, sub.CancelledAt)
	return err
}

func (r *SubscriptionRepository) UpdateUserRole(userID int64, role models.UserRole) error {
	_, err := r.DB.Exec(`UPDATE users SET role = $2, updated_at = NOW() WHERE id = $1`, userID, role)
	return err
}
