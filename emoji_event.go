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
	WorkItem     *WorkItem          `json:"work_item,omitempty"`
	Note         *EmojiNote         `json:"note,omitempty"`
	MergeRequest *EmojiMergeRequest `json:"merge_request,omitempty"`
	Issue        *EmojiIssue        `json:"issue,omitempty"`
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

// EmojiNote represents a note/comment where an emoji was awarded
type EmojiNote struct {
	AuthorID         string               `json:"author_id"`
	ChangePosition   *gitlab.NotePosition `json:"change_position"`
	CommitID         string               `json:"commit_id"`
	CreatedAt        FlexibleTime         `json:"created_at"`
	DiscussionID     string               `json:"discussion_id"`
	ID               string               `json:"id"`
	Internal         bool                 `json:"internal"`
	LineCode         string               `json:"line_code"`
	Note             string               `json:"note"`
	NoteableID       string               `json:"noteable_id"`
	NoteableType     string               `json:"noteable_type"`
	OriginalPosition *gitlab.NotePosition `json:"original_position"`
	Position         *gitlab.NotePosition `json:"position"`
	ProjectID        string               `json:"project_id"`
	ResolvedAt       FlexibleTime         `json:"resolved_at"`
	ResolvedByID     int64                `json:"resolved_by_id"`
	ResolvedByPush   bool                 `json:"resolved_by_push"`
	StDiff           *gitlab.Diff         `json:"st_diff"`
	System           bool                 `json:"system"`
	Type             string               `json:"type"`
	UpdatedAt        FlexibleTime         `json:"updated_at"`
	UpdatedByID      int64                `json:"updated_by_id"`
	Description      string               `json:"description"`
	URL              string               `json:"url"`
}

// EmojiMergeRequest represents a merge request where an emoji was awarded
type EmojiMergeRequest struct {
	AssigneeID                  string              `json:"assignee_id"`
	AuthorID                    string              `json:"author_id"`
	CreatedAt                   FlexibleTime        `json:"created_at"`
	Description                 string              `json:"description"`
	Draft                       bool                `json:"draft"`
	HeadPipelineID              string              `json:"head_pipeline_id"`
	ID                          string              `json:"id"`
	IID                         string              `json:"iid"`
	LastEditedAt                FlexibleTime        `json:"last_edited_at"`
	LastEditedByID              int64               `json:"last_edited_by_id"`
	MergeCommitSha              string              `json:"merge_commit_sha"`
	MergeError                  string              `json:"merge_error"`
	MergeParams                 *gitlab.MergeParams `json:"merge_params"`
	MergeStatus                 string              `json:"merge_status"`
	MergeUserID                 int64               `json:"merge_user_id"`
	MergeWhenPipelineSucceeds   bool                `json:"merge_when_pipeline_succeeds"`
	MilestoneID                 int64               `json:"milestone_id"`
	SourceBranch                string              `json:"source_branch"`
	SourceProjectID             string              `json:"source_project_id"`
	StateID                     int                 `json:"state_id"`
	TargetBranch                string              `json:"target_branch"`
	TargetProjectID             string              `json:"target_project_id"`
	TimeEstimate                int                 `json:"time_estimate"`
	Title                       string              `json:"title"`
	UpdatedAt                   FlexibleTime        `json:"updated_at"`
	UpdatedByID                 string              `json:"updated_by_id"`
	PreparedAt                  FlexibleTime        `json:"prepared_at"`
	AssigneeIDs                 []int64             `json:"assignee_ids"`
	BlockingDiscussionsResolved bool                `json:"blocking_discussions_resolved"`
	DetailedMergeStatus         string              `json:"detailed_merge_status"`
	FirstContribution           bool                `json:"first_contribution"`
	HumanTimeChange             string              `json:"human_time_change"`
	HumanTimeEstimate           string              `json:"human_time_estimate"`
	HumanTotalTimeSpent         string              `json:"human_total_time_spent"`
	Labels                      []string            `json:"labels"`
	LastCommit                  *CommitInfo         `json:"last_commit"`
	ReviewerIDs                 []string            `json:"reviewer_ids"`
	Source                      *ProjectInfo        `json:"source"`
	State                       string              `json:"state"`
	System                      bool                `json:"system"`
	Target                      *ProjectInfo        `json:"target"`
	TimeChange                  int                 `json:"time_change"`
	TotalTimeSpent              int                 `json:"total_time_spent"`
	URL                         string              `json:"url"`
	WorkInProgress              bool                `json:"work_in_progress"`
	ApprovalRules               []ApprovalRule      `json:"approval_rules"`
}

// CommitInfo represents information about a commit
type CommitInfo struct {
	ID        string        `json:"id"`
	Message   string        `json:"message"`
	Title     string        `json:"title"`
	Timestamp FlexibleTime  `json:"timestamp"`
	URL       string        `json:"url"`
	Author    *CommitAuthor `json:"author"`
}

// CommitAuthor represents the author of a commit
type CommitAuthor struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// ProjectInfo represents basic project information
type ProjectInfo struct {
	ID                any    `json:"id"` // TODO: Can be string or int depending on the webhook
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

// ApprovalRule represents an approval rule for a merge request
type ApprovalRule struct {
	ID                                  string       `json:"id"`
	ApprovalsRequired                   int          `json:"approvals_required"`
	Name                                string       `json:"name"`
	RuleType                            string       `json:"rule_type"`
	ReportType                          string       `json:"report_type"`
	MergeRequestID                      string       `json:"merge_request_id"`
	Section                             string       `json:"section"`
	ModifiedFromProjectRule             bool         `json:"modified_from_project_rule"`
	OrchestrationPolicyIdx              any          `json:"orchestration_policy_idx"` // TODO: not sure what the type should be here
	VulnerabilitiesAllowed              int          `json:"vulnerabilities_allowed"`
	Scanners                            []string     `json:"scanners"`
	SeverityLevels                      []string     `json:"severity_levels"`
	VulnerabilityStates                 []string     `json:"vulnerability_states"`
	SecurityOrchestrationPolicyConfigID any          `json:"security_orchestration_policy_configuration_id"`
	ScanResultPolicyID                  int64        `json:"scan_result_policy_id"` // TODO: it is a guess that this ID types are int64
	ApplicablePostMerge                 any          `json:"applicable_post_merge"` // TODO: not sure what the type should be here
	ProjectID                           string       `json:"project_id"`
	ApprovalPolicyRuleID                int64        `json:"approval_policy_rule_id"` // TODO: it is a guess that this ID types are int64
	UpdatedAt                           FlexibleTime `json:"updated_at"`
	CreatedAt                           FlexibleTime `json:"created_at"`
}

// EmojiIssue represents an issue where an emoji was awarded
type EmojiIssue struct {
	AuthorID         string       `json:"author_id"`
	ClosedAt         FlexibleTime `json:"closed_at"`
	Confidential     bool         `json:"confidential"`
	CreatedAt        FlexibleTime `json:"created_at"`
	Description      string       `json:"description"`
	DiscussionLocked bool         `json:"discussion_locked"`
	DueDate          FlexibleTime `json:"due_date"`
	ID               string       `json:"id"`
	IID              string       `json:"iid"`
	LastEditedAt     FlexibleTime `json:"last_edited_at"`
	LastEditedByID   string       `json:"last_edited_by_id"`

	MilestoneID               int64        `json:"milestone_id"` // TODO: not sure if these should be *int64
	MovedToID                 int64        `json:"moved_to_id"`
	DuplicatedToID            int64        `json:"duplicated_to_id"`
	ProjectID                 string       `json:"project_id"`
	RelativePosition          string       `json:"relative_position"`
	StateID                   int          `json:"state_id"`
	TimeEstimate              int          `json:"time_estimate"`
	Title                     string       `json:"title"`
	UpdatedAt                 FlexibleTime `json:"updated_at"`
	UpdatedByID               string       `json:"updated_by_id"`
	Weight                    int64        `json:"weight"`
	HealthStatus              string       `json:"health_status"`
	Type                      string       `json:"type"`
	URL                       string       `json:"url"`
	TotalTimeSpent            int          `json:"total_time_spent"`
	TimeChange                int          `json:"time_change"`
	HumanTotalTimeSpent       string       `json:"human_total_time_spent"`
	HumanTimeChange           string       `json:"human_time_change"`
	HumanTimeEstimate         string       `json:"human_time_estimate"`
	AssigneeIDs               []string     `json:"assignee_ids"`
	AssigneeID                string       `json:"assignee_id"`
	Labels                    []string     `json:"labels"`
	State                     string       `json:"state"`
	Severity                  string       `json:"severity"`
	CustomerRelationsContacts []any        `json:"customer_relations_contacts"` // TODO: not sure what type this should be
	Status                    *IssueStatus `json:"status"`
}

// IssueStatus represents the status of an issue
type IssueStatus struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Category    string `json:"category"`
	Description string `json:"description"`
	Color       string `json:"color"`
}
