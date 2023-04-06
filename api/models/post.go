package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Post struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Title     string    `gorm:"size:255;not null;unique" json:"title"`
	Content   string    `gorm:"size:255;not null;" json:"content"`
	Author    User      `json:"author"`
	AuthorID  uint32    `gorm:"not null" json:"author_id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (post *Post) Prepare() {
	post.ID = 0
	post.Title = html.EscapeString(strings.TrimSpace(post.Title))
	post.Content = html.UnescapeString(strings.TrimSpace(post.Content))
	post.Author = User{}
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()
}

func (post *Post) ValidatePost() error {
	if post.Title == "" {
		return errors.New("Title required")
	}
	if post.Content == "" {
		return errors.New("Title is required")
	}
	if post.AuthorID < 1 {
		return errors.New("Author is required")
	}

	return nil
}

func (post *Post) SavePost(db *gorm.DB) (*Post, error) {
	var err error
	err = db.Debug().Model(&Post{}).Create(&post).Error
	if err != nil {
		return &Post{}, err
	}

	if post.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", post.AuthorID).Take(&post.AuthorID).Error
		if err != nil {
			return &Post{}, err
		}
	}

	return post, nil
}

func (post *Post) FIndAllPosts(db *gorm.DB) (*[]Post, error) {
	var err error
	posts := []Post{}
	err = db.Debug().Model(&Post{}).Limit(100).Find(&posts).Error
	if err != nil {
		return &[]Post{}, err
	}

	if len(posts) > 0 {
		for i := range posts {
			err := db.Debug().Model(&User{}).Where("id = ?", posts[i].AuthorID).Take(&posts[i].Author).Error
			if err != nil {
				return &[]Post{}, err
			}
		}
	}
	return &posts, nil
}

func (post *Post) FIndPostByID(db gorm.DB, postid uint64) (*Post, error) {
	var err error
	err = db.Debug().Model(&Post{}).Where("id = ?", postid).Take(&post).Error
	if err != nil {
		return &Post{}, err
	}

	if post.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", post.AuthorID).Take(&post.Author).Error
		if err != nil {
			return &Post{}, err
		}
	}

	return post, nil
}

func (post *Post) UpdatePost(db gorm.DB) (*Post, error) {
	var err error
	err = db.Debug().Model(&Post{}).Where("id = ?", post.ID).Updates(Post{Title: post.Title, Content: post.Content, UpdatedAt: time.Now()}).Error
	if err != nil {
		return &Post{}, err
	}

	if post.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", post.AuthorID).Take(&post.Author).Error
		if err != nil {
			return &Post{}, err
		}
	}

	return post, nil
}

func (post *Post) DeletePost(db gorm.DB, postid uint64, userid uint32) (int64, error) {
	db = *db.Debug().Model(&Post{}).Where("id = ?", postid, userid).Take(&Post{}).Delete(&Post{})
	if db.Error != nil {
		if gorm.ErrRecordNotFound(db.Error) {
			return 0, errors.New("Post Not Found")
		}
		return 0, db.Error
	}

	return db.RowsAffected, nil
}
