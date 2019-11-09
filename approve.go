package main

// Scope is a localized scope for the approve screen
type Scope struct {
	Icon        string
	Name        string
	Description string
}

var scopes = map[string]Scope{
	"user.basic": Scope{
		Icon:        "user",
		Name:        "User",
		Description: "Read access to full name, email and avatar",
	},
	"user.read": Scope{
		Icon:        "user",
		Name:        "User",
		Description: "Read access to full name, email and avatar and 2FA status",
	},
	"user": Scope{
		Icon:        "user",
		Name:        "User",
		Description: "Read and write access to full name, email, avatar and 2FA status",
	},
	"auth": Scope{
		Icon:        "auth",
		Name:        "Auth",
		Description: "Change your password and enable/disable 2FA",
	},
	"challenge.attempt.read": Scope{
		Icon:        "challenge",
		Name:        "Challenge Attempt",
		Description: "Read access to your past challenge attempts",
	},
	"challenge.attempt.write": Scope{
		Icon:        "challenge",
		Name:        "Challenge Attempt",
		Description: "Submit new challenge attempts",
	},
	"challenge.attempt": Scope{
		Icon:        "challenge",
		Name:        "Challenge Attempt",
		Description: "Submit new challenge attempts access your past challenge attempts",
	},
	"challenge": Scope{
		Icon:        "challenge",
		Name:        "Challenge",
		Description: "Submit new challenge attempts access your past challenge attempts",
	},
	"admin.user": Scope{
		Icon:        "admin",
		Name:        "User Admin",
		Description: "Access all user data on the platform",
	},
	"admin.attempts": Scope{
		Icon:        "admin",
		Name:        "Challenge Attempts Admin",
		Description: "Access all challenge attempts on the platform",
	},
	"admin.challenges": Scope{
		Icon:        "admin",
		Name:        "Challenge Admin",
		Description: "Access all challenges on the platform",
	},
	"admin": Scope{
		Icon:        "admin",
		Name:        "Admin",
		Description: "Do everything. Be everyone.",
	},
}
