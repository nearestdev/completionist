package bootstrap

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/GATEOPENERZ/completionist-api/internal/config"
	"github.com/GATEOPENERZ/completionist-api/internal/models"
	"github.com/GATEOPENERZ/completionist-api/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var usernameSanitizer = regexp.MustCompile(`[^a-z0-9_-]+`)

func EnsureDevAdmin(cfg *config.Config, userRepo *repository.UserRepository) error {
	if cfg.Environment != "development" {
		return nil
	}

	email := strings.TrimSpace(cfg.DevAdminEmail)
	password := cfg.DevAdminPassword
	if email == "" || password == "" {
		log.Printf("Dev admin bootstrap skipped: DEV_ADMIN_EMAIL and DEV_ADMIN_PASSWORD are required in development")
		return nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash dev admin password: %w", err)
	}

	existingUser, err := userRepo.FindByEmail(email)
	if err == nil {
		if err := userRepo.UpdateRoleAndPasswordByEmail(email, models.RoleAdmin, string(hashedPassword)); err != nil {
			return fmt.Errorf("failed to upsert dev admin user: %w", err)
		}
		log.Printf("Dev admin account upserted for existing user: %s", existingUser.Email)
		return nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("failed to look up dev admin user: %w", err)
	}

	username := strings.TrimSpace(cfg.DevAdminUsername)
	if username == "" {
		username = usernameFromEmail(email)
	}
	username, err = findAvailableUsername(userRepo, username)
	if err != nil {
		return fmt.Errorf("failed to pick dev admin username: %w", err)
	}

	user := &models.User{
		Username:     username,
		Email:        email,
		PasswordHash: string(hashedPassword),
		Role:         models.RoleAdmin,
	}
	if err := userRepo.Create(user); err != nil {
		return fmt.Errorf("failed to create dev admin user: %w", err)
	}

	log.Printf("Dev admin account created with email=%s username=%s", email, username)
	return nil
}

func usernameFromEmail(email string) string {
	parts := strings.Split(email, "@")
	base := strings.ToLower(strings.TrimSpace(parts[0]))
	base = usernameSanitizer.ReplaceAllString(base, "_")
	base = strings.Trim(base, "_")
	if base == "" {
		base = "dev_admin"
	}
	if len(base) > 240 {
		base = base[:240]
	}
	return base
}

func findAvailableUsername(userRepo *repository.UserRepository, base string) (string, error) {
	candidate := strings.TrimSpace(base)
	if candidate == "" {
		candidate = "dev_admin"
	}

	candidate = strings.ToLower(candidate)
	candidate = usernameSanitizer.ReplaceAllString(candidate, "_")
	candidate = strings.Trim(candidate, "_")
	if candidate == "" {
		candidate = "dev_admin"
	}

	root := candidate
	if len(root) > 240 {
		root = root[:240]
	}

	for i := 0; i < 1000; i++ {
		attempt := root
		if i > 0 {
			suffix := fmt.Sprintf("-%d", i)
			maxRootLen := 255 - len(suffix)
			trimmedRoot := root
			if len(trimmedRoot) > maxRootLen {
				trimmedRoot = trimmedRoot[:maxRootLen]
			}
			attempt = trimmedRoot + suffix
		}

		_, err := userRepo.FindByUsername(attempt)
		if errors.Is(err, sql.ErrNoRows) {
			return attempt, nil
		}
		if err != nil {
			return "", err
		}
	}

	return "", fmt.Errorf("could not find available username based on %q", root)
}
