package website

import "time"

// Website A web property that is controlled by CURT.
type Website struct {
	ID               int       `json:"id,omitempty" xml:"id,omitempty"`
	URL              string    `json:"url,omitempty" xml:"url,omitempty"`
	Description      string    `json:"description,omitempty" xml:"description,omitempty"`
	Menus            []Menu    `json:"menus,omitempty" xml:"menus,omitempty"`
	Content          []Content `json:"contents,omitempty" xml:"contents,omitempty"`
	BrandIdentifiers []int     `json:"brand_identifiers,omitempty" xml:"brand_identifiers,omitempty"`
}

// Menu A list of Content for a given Website.
type Menu struct {
	ID                    int       `json:"id,omitempty" xml:"id,omitempty"`
	Name                  string    `json:"name,omitempty" xml:"name,omitempty"`
	IsPrimary             bool      `json:"isPrimary,omitempty" xml:"isPrimary,omitempty"`
	Active                bool      `json:"active,omitempty" xml:"active,omitempty"`
	DisplayName           string    `json:"displayName,omitempty" xml:"displayName,omitempty"`
	RequireAuthentication bool      `json:"requireAuthentication,omitempty" xml:"requireAuthentication,omitempty"`
	ShowOnSitemap         bool      `json:"showOnSitemap,omitempty" xml:"showOnSitemap,omitempty"`
	Sort                  int       `json:"sort,omitempty" xml:"sort,omitempty"`
	WebsiteID             int       `json:"websiteID,omitempty" xml:"websiteID,omitempty"`
	Contents              []Content `json:"contents,omitempty" xml:"contents,omitempty"`
}

// Content Properties that are used to generate a content page.
type Content struct {
	ID                    int               `json:"id,omitempty" xml:"id,omitempty"`
	Type                  string            `json:"type,omitempty" xml:"type,omitempty"`
	Title                 string            `json:"title,omitempty" xml:"title,omitempty"`
	CreatedDate           time.Time         `json:"createdDate,omitempty" xml:"createdDate,omitempty"`
	LastModified          time.Time         `json:"lastModified,omitempty" xml:"lastModified,omitempty"`
	MetaTitle             string            `json:"metaTitle,omitempty" xml:"metaTitle,omitempty"`
	MetaDescription       string            `json:"metaDescription,omitempty" xml:"metaDescription,omitempty"`
	Keywords              string            `json:"keywords,omitempty" xml:"keywords,omitempty"`
	IsPrimary             bool              `json:"isPrimary,omitempty" xml:"isPrimary,omitempty"`
	Published             bool              `json:"published,omitempty" xml:"published,omitempty"`
	Active                bool              `json:"active,omitempty" xml:"active,omitempty"`
	Slug                  string            `json:"slug,omitempty" xml:"slug,omitempty"`
	RequireAuthentication bool              `json:"requireAuthentication,omitempty" xml:"requireAuthentication,omitempty"`
	Canonical             string            `json:"canonical,omitempty" xml:"canonical,omitempty"`
	WebsiteID             int               `json:"websiteID,omitempty" xml:"websiteID,omitempty"`
	ContentRevisions      []ContentRevision `json:"contentRevisions,omitempty" xml:"contentRevisions,omitempty"`
	MenuSort              int               `json:"menuSort,omitempty" xml:"menuSort,omitempty"`
	MenuTitle             string            `json:"menuTitle,omitempty" xml:"menuTitle,omitempty"`
	MenuLink              string            `json:"menuLink,omitempty" xml:"menuLink,omitempty"`
	ParentID              int               `json:"parentID,omitempty" xml:"parentID,omitempty"`
	LinkTarget            bool              `json:"linkTarget,omitempty" xml:"v,omitempty"`
}

// ContentRevision Tracks history of Content.
type ContentRevision struct {
	ID          int       `json:"id,omitempty" xml:"id,omitempty"`
	ContentID   int       `json:"contentID,omitempty" xml:"contentID,omitempty"`
	Text        string    `json:"text,omitempty" xml:"text,omitempty"`
	CreatedDate time.Time `json:"createdDate,omitempty" xml:"createdDate,omitempty"`
	Active      bool      `json:"active,omitempty" xml:"active,omitempty"`
}
