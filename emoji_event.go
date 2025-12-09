package gitlabwebhook

import (
	"encoding/json"
	"fmt"
	"time"

	gitlab "gitlab.com/gitlab-org/api/client-go"
)

// FlexibleTime is a custom time type that can unmarshal from multiple time formats
type FlexibleTime struct {
	time.Time
}

// UnmarshalJSON implements json.Unmarshaler to handle multiple time formats
func (ft *FlexibleTime) UnmarshalJSON(b []byte) error {
	// Remove quotes from the JSON string
	s := string(b)
	if len(s) < 2 {
		return fmt.Errorf("invalid time value: %s", s)
	}
	s = s[1 : len(s)-1] // Remove surrounding quotes

	// Try multiple time formats
	formats := []string{
		time.RFC3339,                  // "2006-01-02T15:04:05Z07:00"
		time.RFC3339Nano,              // "2006-01-02T15:04:05.999999999Z07:00"
		"2006-01-02 15:04:05 MST",     // "2025-12-09 20:44:22 UTC"
		"2006-01-02 15:04:05.999 MST", // With milliseconds
	}

	var lastErr error
	for _, format := range formats {
		t, err := time.Parse(format, s)
		if err == nil {
			ft.Time = t
			return nil
		}
		lastErr = err
	}

	return fmt.Errorf("unable to parse time %q: %w", s, lastErr)
}

// MarshalJSON implements json.Marshaler
func (ft FlexibleTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(ft.Time.Format(time.RFC3339))
}

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
	// These used to include an "Issue" but now include a "WorkItem"
	WorkItem     *WorkItem            `json:"work_item,omitempty"`
	Note         *gitlab.Note         `json:"note,omitempty"`
	MergeRequest *gitlab.MergeRequest `json:"merge_request,omitempty"`
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
	UserID        int          `json:"user_id"`
	CreatedAt     FlexibleTime `json:"created_at"`
	ID            int          `json:"id"`
	Name          string       `json:"name"`
	AwardableType string       `json:"awardable_type"`
	AwardableID   int          `json:"awardable_id"`
	UpdatedAt     FlexibleTime `json:"updated_at"`
	AwardedOnURL  string       `json:"awarded_on_url"`
}

// WorkItem represents a GitLab work item (issue)
// These are not yet defined in the gitlab client library, so we define a
// simplified version here as they seem a bit different than issues
type WorkItem struct {
	AuthorID         int           `json:"author_id"`
	ClosedAt         *FlexibleTime `json:"closed_at"`
	Confidential     bool          `json:"confidential"`
	CreatedAt        FlexibleTime  `json:"created_at"`
	Description      string        `json:"description"`
	DiscussionLocked *bool         `json:"discussion_locked"`
	DueDate          *string       `json:"due_date"`
	ID               int           `json:"id"`
	IID              int           `json:"iid"`
	LastEditedAt     *FlexibleTime `json:"last_edited_at"`
	LastEditedByID   *int          `json:"last_edited_by_id"`
	MilestoneID      *int          `json:"milestone_id"`
	MovedToID        *int          `json:"moved_to_id"`
	DuplicatedToID   *int          `json:"duplicated_to_id"`
	ProjectID        int           `json:"project_id"`
	RelativePosition int           `json:"relative_position"`
	StateID          int           `json:"state_id"`
	TimeEstimate     int           `json:"time_estimate"`
	Title            string        `json:"title"`
	UpdatedAt        FlexibleTime  `json:"updated_at"`
	UpdatedByID      int           `json:"updated_by_id"`
	Weight           *int          `json:"weight"`
	HealthStatus     *string       `json:"health_status"`
	Type             string        `json:"type"`
	URL              string        `json:"url"`
	TotalTimeSpent   int           `json:"total_time_spent"`
	TimeChange       int           `json:"time_change"`
	AssigneeIDs      []int         `json:"assignee_ids"`
	AssigneeID       int           `json:"assignee_id"`
	Labels           []string      `json:"labels"`
	State            string        `json:"state"`
	Severity         string        `json:"severity"`
	// ... add other fields as needed
}
