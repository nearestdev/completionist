package services

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/nearestdev/completionist/internal/models"
	"github.com/nearestdev/completionist/internal/repository"
	"github.com/stripe/stripe-go/v81"
	billingsession "github.com/stripe/stripe-go/v81/billingportal/session"
	checkoutsession "github.com/stripe/stripe-go/v81/checkout/session"
	"github.com/stripe/stripe-go/v81/webhook"
)

type StripeService struct {
	subRepo        *repository.SubscriptionRepository
	webhookSecret  string
	defaultPriceID string
	frontendURL    string
}

func NewStripeService(secretKey, webhookSecret, priceID, frontendURL string, subRepo *repository.SubscriptionRepository) *StripeService {
	stripe.Key = secretKey
	return &StripeService{
		subRepo:        subRepo,
		webhookSecret:  webhookSecret,
		defaultPriceID: priceID,
		frontendURL:    frontendURL,
	}
}

func (s *StripeService) IsConfigured() bool {
	return stripe.Key != "" && s.defaultPriceID != ""
}

func (s *StripeService) CreateCheckoutSession(userID int64, email string) (string, error) {
	sub, _ := s.subRepo.GetByUserID(userID)
	var customerID *string
	if sub != nil {
		customerID = &sub.StripeCustomerID
	}

	params := &stripe.CheckoutSessionParams{
		Mode: stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{Price: stripe.String(s.defaultPriceID), Quantity: stripe.Int64(1)},
		},
		SuccessURL: stripe.String(s.frontendURL + "/upgrade?success=true"),
		CancelURL:  stripe.String(s.frontendURL + "/upgrade?cancelled=true"),
		Metadata:   map[string]string{"user_id": fmt.Sprintf("%d", userID)},
	}
	if customerID != nil {
		params.Customer = customerID
	} else {
		params.CustomerEmail = stripe.String(email)
	}

	session, err := checkoutsession.New(params)
	if err != nil {
		return "", err
	}
	return session.URL, nil
}

func (s *StripeService) CreatePortalSession(userID int64) (string, error) {
	sub, err := s.subRepo.GetByUserID(userID)
	if err != nil || sub == nil {
		return "", fmt.Errorf("no subscription found")
	}

	params := &stripe.BillingPortalSessionParams{
		Customer:  stripe.String(sub.StripeCustomerID),
		ReturnURL: stripe.String(s.frontendURL + "/settings"),
	}
	session, err := billingsession.New(params)
	if err != nil {
		return "", err
	}
	return session.URL, nil
}

func (s *StripeService) HandleWebhook(r *http.Request) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	event, err := webhook.ConstructEvent(body, r.Header.Get("Stripe-Signature"), s.webhookSecret)
	if err != nil {
		return fmt.Errorf("webhook signature verification failed: %w", err)
	}

	switch event.Type {
	case "checkout.session.completed":
		return s.handleCheckoutCompleted(event.Data.Raw)
	case "customer.subscription.updated":
		return s.handleSubscriptionUpdated(event.Data.Raw)
	case "customer.subscription.deleted":
		return s.handleSubscriptionDeleted(event.Data.Raw)
	case "invoice.payment_failed":
		return s.handlePaymentFailed(event.Data.Raw)
	}
	return nil
}

func (s *StripeService) handleCheckoutCompleted(raw json.RawMessage) error {
	var session stripe.CheckoutSession
	if err := json.Unmarshal(raw, &session); err != nil {
		return err
	}

	userIDStr, ok := session.Metadata["user_id"]
	if !ok {
		return fmt.Errorf("missing user_id in metadata")
	}
	var userID int64
	fmt.Sscanf(userIDStr, "%d", &userID)

	sub := &models.Subscription{
		UserID:               userID,
		StripeCustomerID:     session.Customer.ID,
		StripeSubscriptionID: &session.Subscription.ID,
		Status:               models.SubscriptionActive,
	}

	existing, _ := s.subRepo.GetByUserID(userID)
	if existing != nil {
		existing.StripeSubscriptionID = sub.StripeSubscriptionID
		existing.Status = models.SubscriptionActive
		now := time.Now()
		existing.CurrentPeriodStart = &now
		if err := s.subRepo.Update(existing); err != nil {
			return err
		}
	} else {
		if err := s.subRepo.Create(sub); err != nil {
			return err
		}
	}

	return s.subRepo.UpdateUserRole(userID, models.RoleMember)
}

func (s *StripeService) handleSubscriptionUpdated(raw json.RawMessage) error {
	var sub stripe.Subscription
	if err := json.Unmarshal(raw, &sub); err != nil {
		return err
	}

	existing, err := s.subRepo.GetByStripeSubscriptionID(sub.ID)
	if err != nil || existing == nil {
		log.Printf("stripe: subscription %s not found locally", sub.ID)
		return nil
	}

	existing.Status = models.SubscriptionStatus(sub.Status)
	existing.CancelAtPeriodEnd = sub.CancelAtPeriodEnd
	start := time.Unix(sub.CurrentPeriodStart, 0)
	end := time.Unix(sub.CurrentPeriodEnd, 0)
	existing.CurrentPeriodStart = &start
	existing.CurrentPeriodEnd = &end

	return s.subRepo.Update(existing)
}

func (s *StripeService) handleSubscriptionDeleted(raw json.RawMessage) error {
	var sub stripe.Subscription
	if err := json.Unmarshal(raw, &sub); err != nil {
		return err
	}

	existing, err := s.subRepo.GetByStripeSubscriptionID(sub.ID)
	if err != nil || existing == nil {
		return nil
	}

	existing.Status = models.SubscriptionCancelled
	now := time.Now()
	existing.CancelledAt = &now
	if err := s.subRepo.Update(existing); err != nil {
		return err
	}

	return s.subRepo.UpdateUserRole(existing.UserID, models.RoleUser)
}

func (s *StripeService) handlePaymentFailed(raw json.RawMessage) error {
	var invoice struct {
		Subscription string `json:"subscription"`
	}
	if err := json.Unmarshal(raw, &invoice); err != nil {
		return err
	}

	existing, err := s.subRepo.GetByStripeSubscriptionID(invoice.Subscription)
	if err != nil || existing == nil {
		return nil
	}

	existing.Status = models.SubscriptionPastDue
	return s.subRepo.Update(existing)
}
