package model

import gl "gitlab.com/gitlab-org/api/client-go"

const (
	DefaultPerPage int64 = 20  // standard list page size
	MaxPerPage     int64 = 100 // used for sub-resource lists (notes, jobs)
)

// NewPaged wraps a slice with pagination metadata sourced from GitLab headers.
func NewPaged[T any](items []T, resp *gl.Response, page int64) Paged[T] {
	p := Paged[T]{
		Items:   items,
		Page:    page,
		PerPage: DefaultPerPage,
	}
	if resp != nil {
		p.Total = resp.TotalItems
		p.Truncated = resp.TotalItems > int64(len(items))+(page-1)*DefaultPerPage
	}
	return p
}

// ── Person helpers ────────────────────────────────────────────────────────────

func toPersonFromBasicUser(u *gl.BasicUser) PersonRef {
	if u == nil {
		return PersonRef{}
	}
	return PersonRef{ID: u.ID, Name: u.Name, Username: u.Username}
}

func toPersonsFromBasicUsers(users []*gl.BasicUser) []PersonRef {
	out := make([]PersonRef, len(users))
	for i, u := range users {
		out[i] = toPersonFromBasicUser(u)
	}
	return out
}

func toPersonFromIssueAuthor(a *gl.IssueAuthor) PersonRef {
	if a == nil {
		return PersonRef{}
	}
	return PersonRef{ID: a.ID, Name: a.Name, Username: a.Username}
}

func toPersonsFromIssueAssignees(assignees []*gl.IssueAssignee) []PersonRef {
	out := make([]PersonRef, len(assignees))
	for i, a := range assignees {
		out[i] = PersonRef{ID: a.ID, Name: a.Name, Username: a.Username}
	}
	return out
}

// ── Projects ──────────────────────────────────────────────────────────────────

func ToProjectRef(p *gl.Project) ProjectRef {
	return ProjectRef{
		ID:                p.ID,
		Name:              p.Name,
		PathWithNamespace: p.PathWithNamespace,
		WebURL:            p.WebURL,
	}
}

func ToProjectRefs(ps []*gl.Project) []ProjectRef {
	out := make([]ProjectRef, len(ps))
	for i, p := range ps {
		out[i] = ToProjectRef(p)
	}
	return out
}

func ToProject(p *gl.Project) Project {
	return Project{
		ID:                p.ID,
		Name:              p.Name,
		PathWithNamespace: p.PathWithNamespace,
		Description:       p.Description,
		Visibility:        string(p.Visibility),
		DefaultBranch:     p.DefaultBranch,
		WebURL:            p.WebURL,
	}
}

// ── Issues ────────────────────────────────────────────────────────────────────

func ToIssueRef(i *gl.Issue) IssueRef {
	return IssueRef{
		IID:    i.IID,
		Title:  i.Title,
		State:  i.State,
		Author: toPersonFromIssueAuthor(i.Author),
		Labels: i.Labels,
		WebURL: i.WebURL,
	}
}

func ToIssueRefs(issues []*gl.Issue) []IssueRef {
	out := make([]IssueRef, len(issues))
	for i, iss := range issues {
		out[i] = ToIssueRef(iss)
	}
	return out
}

func ToIssue(i *gl.Issue) Issue {
	return Issue{
		IID:         i.IID,
		Title:       i.Title,
		Description: i.Description,
		State:       i.State,
		Author:      toPersonFromIssueAuthor(i.Author),
		Assignees:   toPersonsFromIssueAssignees(i.Assignees),
		Labels:      i.Labels,
		CreatedAt:   i.CreatedAt,
		UpdatedAt:   i.UpdatedAt,
		WebURL:      i.WebURL,
	}
}

// ── Merge Requests ────────────────────────────────────────────────────────────

func ToMRRefFromBasic(m *gl.BasicMergeRequest) MRRef {
	return MRRef{
		IID:          m.IID,
		Title:        m.Title,
		State:        m.State,
		SourceBranch: m.SourceBranch,
		TargetBranch: m.TargetBranch,
		Author:       toPersonFromBasicUser(m.Author),
		WebURL:       m.WebURL,
	}
}

func ToMRRefsFromBasic(mrs []*gl.BasicMergeRequest) []MRRef {
	out := make([]MRRef, len(mrs))
	for i, m := range mrs {
		out[i] = ToMRRefFromBasic(m)
	}
	return out
}

func ToMRRef(m *gl.MergeRequest) MRRef {
	return MRRef{
		IID:          m.IID,
		Title:        m.Title,
		State:        m.State,
		SourceBranch: m.SourceBranch,
		TargetBranch: m.TargetBranch,
		Author:       toPersonFromBasicUser(m.Author),
		WebURL:       m.WebURL,
	}
}

func ToMRRefs(mrs []*gl.MergeRequest) []MRRef {
	out := make([]MRRef, len(mrs))
	for i, m := range mrs {
		out[i] = ToMRRef(m)
	}
	return out
}

func ToMR(m *gl.MergeRequest) MR {
	mr := MR{
		IID:          m.IID,
		Title:        m.Title,
		Description:  m.Description,
		State:        m.State,
		SourceBranch: m.SourceBranch,
		TargetBranch: m.TargetBranch,
		Author:       toPersonFromBasicUser(m.Author),
		Labels:       m.Labels,
		MergeStatus:  m.DetailedMergeStatus,
		SHA:          m.SHA,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
		WebURL:       m.WebURL,
	}
	if m.Assignees != nil {
		mr.Assignees = toPersonsFromBasicUsers(m.Assignees)
	}
	return mr
}

// ── Pipelines ─────────────────────────────────────────────────────────────────

func ToPipelineRef(p *gl.PipelineInfo) PipelineRef {
	return PipelineRef{
		ID:     p.ID,
		Status: p.Status,
		Ref:    p.Ref,
		SHA:    p.SHA,
		WebURL: p.WebURL,
	}
}

func ToPipelineRefs(ps []*gl.PipelineInfo) []PipelineRef {
	out := make([]PipelineRef, len(ps))
	for i, p := range ps {
		out[i] = ToPipelineRef(p)
	}
	return out
}

func ToPipeline(p *gl.Pipeline) Pipeline {
	return Pipeline{
		ID:        p.ID,
		Status:    p.Status,
		Ref:       p.Ref,
		SHA:       p.SHA,
		CreatedAt: p.CreatedAt,
		Duration:  int(p.Duration),
		WebURL:    p.WebURL,
	}
}

// ── Jobs ──────────────────────────────────────────────────────────────────────

func ToJob(j *gl.Job) Job {
	job := Job{
		ID:       j.ID,
		Name:     j.Name,
		Stage:    j.Stage,
		Status:   j.Status,
		WebURL:   j.WebURL,
		Duration: j.Duration,
	}
	if j.Pipeline.ID != 0 {
		job.Pipeline = PipelineRef{
			ID:     j.Pipeline.ID,
			Status: j.Pipeline.Status,
			Ref:    j.Pipeline.Ref,
			SHA:    j.Pipeline.Sha,
		}
	}
	return job
}

func ToJobs(jobs []*gl.Job) []Job {
	out := make([]Job, len(jobs))
	for i, j := range jobs {
		out[i] = ToJob(j)
	}
	return out
}

// ── Notes ─────────────────────────────────────────────────────────────────────

func ToNote(n *gl.Note) Note {
	return Note{
		ID:        n.ID,
		Body:      n.Body,
		Author:    PersonRef{ID: n.Author.ID, Name: n.Author.Name, Username: n.Author.Username},
		CreatedAt: n.CreatedAt,
		System:    n.System,
	}
}

func ToNotes(notes []*gl.Note) []Note {
	out := make([]Note, len(notes))
	for i, n := range notes {
		out[i] = ToNote(n)
	}
	return out
}

// ── Branches ──────────────────────────────────────────────────────────────────

func ToBranch(b *gl.Branch) Branch {
	br := Branch{
		Name:      b.Name,
		Protected: b.Protected,
		Default:   b.Default,
	}
	if b.Commit != nil {
		br.Commit = CommitRef{
			ID:         b.Commit.ID,
			ShortID:    b.Commit.ShortID,
			Title:      b.Commit.Title,
			AuthorName: b.Commit.AuthorName,
			CreatedAt:  b.Commit.CreatedAt,
		}
	}
	return br
}

func ToBranches(bs []*gl.Branch) []Branch {
	out := make([]Branch, len(bs))
	for i, b := range bs {
		out[i] = ToBranch(b)
	}
	return out
}

// ── Commits ───────────────────────────────────────────────────────────────────

func ToCommitRef(c *gl.Commit) CommitRef {
	return CommitRef{
		ID:         c.ID,
		ShortID:    c.ShortID,
		Title:      c.Title,
		AuthorName: c.AuthorName,
		CreatedAt:  c.CreatedAt,
	}
}

func ToCommitRefs(cs []*gl.Commit) []CommitRef {
	out := make([]CommitRef, len(cs))
	for i, c := range cs {
		out[i] = ToCommitRef(c)
	}
	return out
}

// ── Users ─────────────────────────────────────────────────────────────────────

func ToUser(u *gl.User) User {
	return User{
		ID:       u.ID,
		Username: u.Username,
		Name:     u.Name,
		Email:    u.Email,
		State:    u.State,
		WebURL:   u.WebURL,
	}
}
