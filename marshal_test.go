package phpcrudapi

import (
	"testing"
)

type author struct {
	ID        int    `json:"id" collection:"authors"`
	Name      string `json:"name"`
	CreatedAt Time   `json:"created_at"`
	UpdatedAt Time   `json:"updated_at"`
}

func TestUnmarshalTable(t *testing.T) {
	simple := `{"authors":
		{"columns":[
			"id","name","created_at","updated_at"
		],"records":[
			[1,"test","2018-09-08 11:19:17.361572","2018-09-08 11:19:17.361572"]
		]
		}
	}`
	var results []*author
	err := Unmarshal([]byte(simple[:]), &results)
	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 1 {
		t.Fatalf("expected one row, got: %d", len(results))
	}

	if results[0].Name != "test" {
		t.Errorf("expected test, got: %q", results[0].Name)
	}

	t.Logf("%#v", results[0])
}

type post struct {
	ID       int        `json:"id" collection:"posts"`
	UserID   int        `json:"user_id"`
	Category *category  `json:"category_id" relation:"belongs-to"`
	Content  string     `json:"content"`
	Tags     []*tag     `relation:"many-to-many" via:"post_tags"`
	Comments []*comment `relation:"one-to-many"`
}

type category struct {
	ID   int    `json:"id" collection:"categories"`
	Name string `json:"name"`
}

type tag struct {
	ID   int    `json:"id" collection:"tags"`
	Name string `json:"name"`
}

type comment struct {
	ID      int    `json:"id" collection:"comments"`
	Message string `json:"message"`
}

func TestUnmarshalRelations(t *testing.T) {
	var posts []*post
	err := Unmarshal([]byte(relationExample[:]), &posts)
	if err != nil {
		t.Fatal(err)
	}

	// test post
	if len(posts) != 1 {
		t.Fatalf("expected one post, got: %d", len(posts))
	}
	if posts[0].Content != "blog started" {
		t.Errorf("expected blog started, got: %q", posts[0].Content)
	}

	// test post category
	if posts[0].Category == nil {
		t.Fatalf("expected post to have a category")
	}
	if posts[0].Category.Name != "announcement" {
		t.Errorf("expected blog category announcement, got: %q", posts[0].Category.Name)
	}

	// test comments
	if len(posts[0].Comments) != 2 {
		t.Fatalf("expected two comments, got: %d", len(posts[0].Comments))
	}
	if posts[0].Comments[0].Message != "great" ||
		posts[0].Comments[1].Message != "fantastic" {
		t.Errorf("expected blog great and fantastic, got: %q and %q",
			posts[0].Comments[0].Message, posts[0].Comments[1].Message)
	}

	t.Logf("%#v", posts[0])
}

var relationExample = `{
    "posts": {
        "columns": [
            "id",
            "user_id",
            "category_id",
            "content"
        ],
        "records": [
            [
                1,
                1,
                1,
                "blog started"
            ]
        ]
    },
    "post_tags": {
        "relations": {
            "post_id": "posts.id"
        },
        "columns": [
            "id",
            "post_id",
            "tag_id"
        ],
        "records": [
            [
                1,
                1,
                1
            ],
            [
                2,
                1,
                2
            ]
        ]
    },
    "categories": {
        "relations": {
            "id": "posts.category_id"
        },
        "columns": [
            "id",
            "name"
        ],
        "records": [
            [
                1,
                "announcement"
            ]
        ]
    },
    "tags": {
        "relations": {
            "id": "post_tags.tag_id"
        },
        "columns": [
            "id",
            "name"
        ],
        "records": [
            [
                1,
                "funny"
            ],
            [
                2,
                "important"
            ]
        ]
    },
    "comments": {
        "relations": {
            "post_id": "posts.id"
        },
        "columns": [
            "id",
            "post_id",
            "message"
        ],
        "records": [
            [
                1,
                1,
                "great"
            ],
            [
                2,
                1,
                "fantastic"
            ]
        ]
    }
}`
