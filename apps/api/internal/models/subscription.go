package models

import "time"

type SubscriptionStatus string

const (
	SubscriptionActive     SubscriptionStatus = "active"
	SubscriptionPastDue    SubscriptionStatus = "past_due"
	SubscriptionCancelled  SubscriptionStatus = "cancelled"
	SubscriptionIncomplete SubscriptionStatus = "incomplete"
	SubscriptionTrialing   SubscriptionStatus = "trialing"
)

type Subscription struct {
	ID                   int64              `db:"id" json:"id"`
	UserID               int64              `db:"user_id" json:"userId"`
	StripeCustomerID     string             `db:"stripe_customer_id" json:"-"`
	StripeSubscriptionID *string            `db:"stripe_subscription_id" json:"-"`
	StripePriceID        *string            `db:"stripe_price_id" json:"-"`
	Status               SubscriptionStatus `db:"status" json:"status"`
	CurrentPeriodStart   *time.Time         `db:"current_period_start" json:"currentPeriodStart,omitempty"`
	CurrentPeriodEnd     *time.Time         `db:"current_period_end" json:"currentPeriodEnd,omitempty"`
	CancelAtPeriodEnd    bool               `db:"cancel_at_period_end" json:"cancelAtPeriodEnd"`
	CancelledAt          *time.Time         `db:"cancelled_at" json:"cancelledAt,omitempty"`
	CreatedAt            time.Time          `db:"created_at" json:"createdAt"`
	UpdatedAt            time.Time          `db:"updated_at" json:"updatedAt"`
}
