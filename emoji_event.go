package gitlabwebhook

import (
	"time"
)

// The gitlab client library (gitlab.com/gitlab-org/api/client-go) does not currently have a struct for EmojiEvent,
// so we define it here.

// EmojiEvent represents an emoji/award event from GitLab
type EmojiEvent struct {
	ObjectKind string           `json:"object_kind"`
	EventType  string           `json:"event_type"`
	User       *EmojiUser       `json:"user"`
	ProjectID  int              `json:"project_id"`
	Project    *EmojiProject    `json:"project"`
	ObjectAttr *EmojiAttributes `json:"object_attributes"`
}

// EmojiUser represents the user who awarded the emoji
type EmojiUser struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatar_url"`
	Email     string `json:"email"`
}

// EmojiProject represents the project where the emoji was awarded
type EmojiProject struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	WebURL            string `json:"web_url"`
	AvatarURL         string `json:"avatar_url"`
	GitSSHURL         string `json:"git_ssh_url"`
	GitHTTPURL        string `json:"git_http_url"`
	Namespace         string `json:"namespace"`
	VisibilityLevel   int    `json:"visibility_level"`
	PathWithNamespace string `json:"path_with_namespace"`
	DefaultBranch     string `json:"default_branch"`
	CIConfigPath      string `json:"ci_config_path"`
	Homepage          string `json:"homepage"`
	URL               string `json:"url"`
	SSHURL            string `json:"ssh_url"`
	HTTPURL           string `json:"http_url"`
}

// EmojiAttributes represents the attributes of the emoji/award
type EmojiAttributes struct {
	UserID        int       `json:"user_id"`
	CreatedAt     time.Time `json:"created_at"`
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	AwardableType string    `json:"awardable_type"`
	AwardableID   int       `json:"awardable_id"`
	UpdatedAt     time.Time `json:"updated_at"`
	AwardedOnURL  string    `json:"awarded_on_url"`
}
