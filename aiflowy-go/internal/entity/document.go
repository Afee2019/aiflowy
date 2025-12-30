package entity

import "time"

// Document 文档实体
type Document struct {
	ID           int64      `db:"id" json:"id,string"`
	CollectionID int64      `db:"collection_id" json:"collectionId,string"`
	DocumentType string     `db:"document_type" json:"documentType,omitempty"`
	DocumentPath string     `db:"document_path" json:"documentPath,omitempty"`
	Title        string     `db:"title" json:"title,omitempty"`
	Content      string     `db:"content" json:"content,omitempty"`
	ContentType  string     `db:"content_type" json:"contentType,omitempty"`
	Slug         string     `db:"slug" json:"slug,omitempty"`
	OrderNo      *int       `db:"order_no" json:"orderNo,omitempty"`
	Options      string     `db:"options" json:"options,omitempty"`
	Created      *time.Time `db:"created" json:"created,omitempty"`
	CreatedBy    *int64     `db:"created_by" json:"createdBy,string,omitempty"`
	Modified     *time.Time `db:"modified" json:"modified,omitempty"`
	ModifiedBy   *int64     `db:"modified_by" json:"modifiedBy,string,omitempty"`
}

// DocumentChunk 文档分块实体
type DocumentChunk struct {
	ID                   int64  `db:"id" json:"id,string"`
	DocumentID           int64  `db:"document_id" json:"documentId,string"`
	DocumentCollectionID int64  `db:"document_collection_id" json:"documentCollectionId,string"`
	Content              string `db:"content" json:"content,omitempty"`
	Sorting              int    `db:"sorting" json:"sorting,omitempty"`
}

// DocumentHistory 文档历史记录实体
type DocumentHistory struct {
	ID              int64      `db:"id" json:"id,string"`
	DocumentID      int64      `db:"document_id" json:"documentId,string"`
	OldTitle        string     `db:"old_title" json:"oldTitle,omitempty"`
	NewTitle        string     `db:"new_title" json:"newTitle,omitempty"`
	OldContent      string     `db:"old_content" json:"oldContent,omitempty"`
	NewContent      string     `db:"new_content" json:"newContent,omitempty"`
	OldDocumentType string     `db:"old_document_type" json:"oldDocumentType,omitempty"`
	NewDocumentType string     `db:"new_document_type" json:"newDocumentType,omitempty"`
	Created         *time.Time `db:"created" json:"created,omitempty"`
	CreatedBy       *int64     `db:"created_by" json:"createdBy,string,omitempty"`
}
