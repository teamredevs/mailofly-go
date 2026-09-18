package mailofly

import (
	"net/http"
	"strconv"
)

// Client is the Mailofly REST API client.
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client

	Identities *IdentitiesService
	Accounts   *IdentitiesService
	Contacts   *ContactsService
	Templates *TemplatesService
	Segments  *SegmentsService
	Campaigns *CampaignsService
	Compose   *ComposeService
	Emails    *EmailsService
	Batch     *BatchService
	MailLogs  *MailLogsService
}

// Options configures a Client.
type Options struct {
	APIKey     string
	BaseURL    string
	HTTPClient *http.Client
}

// New creates a Mailofly client.
func New(opts Options) (*Client, error) {
	key, err := mustAPIKey(opts.APIKey)
	if err != nil {
		return nil, err
	}
	c := &Client{
		apiKey:     key,
		baseURL:    normalizeBaseURL(opts.BaseURL),
		httpClient: opts.HTTPClient,
	}
	c.Identities = &IdentitiesService{c}
	c.Accounts = c.Identities
	c.Contacts = &ContactsService{c}
	c.Templates = &TemplatesService{c}
	c.Segments = &SegmentsService{c: c, Contacts: &SegmentContactsService{c}}
	c.Campaigns = &CampaignsService{c}
	c.Compose = &ComposeService{c}
	c.Emails = &EmailsService{c}
	c.Batch = &BatchService{c}
	c.MailLogs = &MailLogsService{c}
	return c, nil
}

// Discovery performs unauthenticated GET /.
func Discovery(baseURL string) (any, error) {
	discPath := "/"
	if APIPrefix != "" {
		discPath = APIPrefix
	}
	return doRequest(nil, normalizeBaseURL(baseURL), discPath, http.MethodGet, "", nil, nil)
}

func (c *Client) req(path, method string, body any, query map[string]string) (any, error) {
	return doRequest(c.httpClient, c.baseURL, APIPrefix+path, method, c.apiKey, body, query)
}

// IdentitiesService manages sending identities (creation is dashboard-only).
type IdentitiesService struct{ c *Client }

func (s *IdentitiesService) List() (any, error) {
	return s.c.req("/identities", http.MethodGet, nil, nil)
}
func (s *IdentitiesService) Get(id string) (any, error) {
	return s.c.req("/identities/"+enc(id), http.MethodGet, nil, nil)
}
func (s *IdentitiesService) Update(id string, body map[string]any) (any, error) {
	return s.c.req("/identities/"+enc(id), http.MethodPatch, body, nil)
}
func (s *IdentitiesService) Delete(id string) (any, error) {
	return s.c.req("/identities/"+enc(id), http.MethodDelete, nil, nil)
}

// ContactsService manages contacts.
type ContactsService struct{ c *Client }

func (s *ContactsService) List(segmentID string) (any, error) {
	var q map[string]string
	if segmentID != "" {
		q = map[string]string{"segment_id": segmentID}
	}
	return s.c.req("/contacts", http.MethodGet, nil, q)
}
func (s *ContactsService) Create(body map[string]any) (any, error) {
	return s.c.req("/contacts", http.MethodPost, body, nil)
}
func (s *ContactsService) Get(id string) (any, error) {
	return s.c.req("/contacts/"+enc(id), http.MethodGet, nil, nil)
}
func (s *ContactsService) Update(id string, body map[string]any) (any, error) {
	return s.c.req("/contacts/"+enc(id), http.MethodPatch, body, nil)
}
func (s *ContactsService) Delete(id string) (any, error) {
	return s.c.req("/contacts/"+enc(id), http.MethodDelete, nil, nil)
}

// TemplatesService manages templates.
type TemplatesService struct{ c *Client }

func (s *TemplatesService) List() (any, error) {
	return s.c.req("/templates", http.MethodGet, nil, nil)
}
func (s *TemplatesService) Create(body map[string]any) (any, error) {
	return s.c.req("/templates", http.MethodPost, body, nil)
}
func (s *TemplatesService) Get(id string) (any, error) {
	return s.c.req("/templates/"+enc(id), http.MethodGet, nil, nil)
}
func (s *TemplatesService) Update(id string, body map[string]any) (any, error) {
	return s.c.req("/templates/"+enc(id), http.MethodPatch, body, nil)
}
func (s *TemplatesService) Delete(id string) (any, error) {
	return s.c.req("/templates/"+enc(id), http.MethodDelete, nil, nil)
}

// SegmentContactsService manages segment membership.
type SegmentContactsService struct{ c *Client }

func (s *SegmentContactsService) List(segmentID string) (any, error) {
	return s.c.req("/segments/"+enc(segmentID)+"/contacts", http.MethodGet, nil, nil)
}
func (s *SegmentContactsService) Add(segmentID string, body map[string]any) (any, error) {
	return s.c.req("/segments/"+enc(segmentID)+"/contacts", http.MethodPost, body, nil)
}
func (s *SegmentContactsService) Remove(segmentID, contactID string) (any, error) {
	return s.c.req("/segments/"+enc(segmentID)+"/contacts/"+enc(contactID), http.MethodDelete, nil, nil)
}

// SegmentsService manages segments.
type SegmentsService struct {
	c        *Client
	Contacts *SegmentContactsService
}

func (s *SegmentsService) List() (any, error) {
	return s.c.req("/segments", http.MethodGet, nil, nil)
}
func (s *SegmentsService) Create(body map[string]any) (any, error) {
	return s.c.req("/segments", http.MethodPost, body, nil)
}
func (s *SegmentsService) Get(id string) (any, error) {
	return s.c.req("/segments/"+enc(id), http.MethodGet, nil, nil)
}
func (s *SegmentsService) Update(id string, body map[string]any) (any, error) {
	return s.c.req("/segments/"+enc(id), http.MethodPatch, body, nil)
}
func (s *SegmentsService) Delete(id string) (any, error) {
	return s.c.req("/segments/"+enc(id), http.MethodDelete, nil, nil)
}

// CampaignsService manages campaigns.
type CampaignsService struct{ c *Client }

func (s *CampaignsService) List() (any, error) {
	return s.c.req("/campaigns", http.MethodGet, nil, nil)
}
func (s *CampaignsService) Create(body map[string]any) (any, error) {
	return s.c.req("/campaigns", http.MethodPost, body, nil)
}
func (s *CampaignsService) Get(id string) (any, error) {
	return s.c.req("/campaigns/"+enc(id), http.MethodGet, nil, nil)
}
func (s *CampaignsService) Update(id string, body map[string]any) (any, error) {
	return s.c.req("/campaigns/"+enc(id), http.MethodPatch, body, nil)
}
func (s *CampaignsService) Delete(id string) (any, error) {
	return s.c.req("/campaigns/"+enc(id), http.MethodDelete, nil, nil)
}
func (s *CampaignsService) Runs(id string) (any, error) {
	return s.c.req("/campaigns/"+enc(id)+"/runs", http.MethodGet, nil, nil)
}
func (s *CampaignsService) Send(id string, body map[string]any) (any, error) {
	if body == nil {
		body = map[string]any{"send_now": true}
	}
	return s.c.req("/campaigns/"+enc(id)+"/send", http.MethodPost, body, nil)
}

// ComposeService sends one-off email (deprecated — use EmailsService).
type ComposeService struct{ c *Client }

func (s *ComposeService) Send(params map[string]any) (any, error) {
	return s.c.req("/emails", http.MethodPost, params, nil)
}

// EmailsService sends transactional email.
type EmailsService struct{ c *Client }

func (s *EmailsService) Send(params map[string]any) (any, error) {
	return s.c.req("/emails", http.MethodPost, params, nil)
}

func (s *EmailsService) Get(id string) (any, error) {
	return s.c.req("/emails/"+enc(id), http.MethodGet, nil, nil)
}

func (s *EmailsService) List(opts *EmailsListOptions) (any, error) {
	q := map[string]string{}
	if opts != nil {
		if opts.Limit > 0 {
			q["limit"] = strconv.Itoa(opts.Limit)
		}
		if opts.After != "" {
			q["after"] = opts.After
		}
		if opts.Before != "" {
			q["before"] = opts.Before
		}
	}
	return s.c.req("/emails", http.MethodGet, nil, q)
}

// EmailsListOptions filters list emails (Resend cursor pagination).
type EmailsListOptions struct {
	Limit  int
	After  string
	Before string
}

// BatchService sends multiple emails (POST /emails/batch).
type BatchService struct{ c *Client }

func (s *BatchService) Send(emails []map[string]any) (any, error) {
	return s.c.req("/emails/batch", http.MethodPost, emails, nil)
}

// MailLogsListOptions filters mail logs.
type MailLogsListOptions struct {
	Page          int
	PageSize      int
	CampaignID    string
	AccountID     string
	CampaignRunID string
	Status        string
}

// MailLogsService queries mail logs.
type MailLogsService struct{ c *Client }

func (s *MailLogsService) List(opts *MailLogsListOptions) (any, error) {
	q := map[string]string{}
	if opts != nil {
		if opts.Page > 0 {
			q["page"] = strconv.Itoa(opts.Page)
		}
		if opts.PageSize > 0 {
			q["page_size"] = strconv.Itoa(opts.PageSize)
		}
		if opts.CampaignID != "" {
			q["campaign_id"] = opts.CampaignID
		}
		if opts.AccountID != "" {
			q["account_id"] = opts.AccountID
		}
		if opts.CampaignRunID != "" {
			q["campaign_run_id"] = opts.CampaignRunID
		}
		if opts.Status != "" {
			q["status"] = opts.Status
		}
	}
	if len(q) == 0 {
		q = nil
	}
	return s.c.req("/mail-logs", http.MethodGet, nil, q)
}
