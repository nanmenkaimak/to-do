package authdb

import "golang.org/x/crypto/bcrypt"

var usersPasswords = map[string][]byte{
	"nanmenkaimak": []byte("$2a$14$4f1ke6Zbp0LbkzPtMzgyeOoauDxHl.kZ0iHC20szNFCNXdBfkhkei"),
}

// VerifyUserPass verifies that username/password is a valid pair matching
// our userPasswords "database".
func VerifyUserPass(username string, password string) bool {
	wantPass, ok := usersPasswords[username]
	if !ok {
		return false
	}

	if cmperr := bcrypt.CompareHashAndPassword(wantPass, []byte(password)); cmperr == nil {
		return true
	}

	return false
}
