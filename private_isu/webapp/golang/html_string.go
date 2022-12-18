package main

import "fmt"

func getPostsHTMLString(posts []Post) string {
	var html string
	html += `<div class="isu-posts">`
	for _, post := range posts {
		html += fmt.Sprintf(`
<div class="isu-post" id="pid_%d" data-created-at="%s"}}">
  <div class="isu-post-header">
    <a href="/@%s " class="isu-post-account-name">%s</a>
    <a href="/posts/%d" class="isu-post-permalink">
      <time class="timeago" datetime="%s"></time>
    </a>
  </div>
  <div class="isu-post-image">
    <img src="%s" class="isu-image">
  </div>
  <div class="isu-post-text">
    <a href="/@%s" class="isu-post-account-name">%s</a>
    %s
  </div>
  <div class="isu-post-comment">
    <div class="isu-post-comment-count">
      comments: <b>%d</b>
    </div>

	%s
    <div class="isu-comment-form">
      <form method="post" action="/comment">
        <input type="text" name="comment">
        <input type="hidden" name="post_id" value="%d">
        <input type="hidden" name="csrf_token" value="%s">
        <input type="submit" name="submit" value="submit">
      </form>
    </div>
  </div>
</div>
`,
			post.ID,
			post.CreatedAt.Format("2006-01-02T15:04:05-07:00"),
			post.User.AccountName,
			post.User.AccountName,
			post.ID,
			post.CreatedAt.Format("2006-01-02T15:04:05-07:00"),
			imageURL(post),
			post.User.AccountName,
			post.User.AccountName,
			post.Body,
			post.CommentCount,
			getPostCommentsHTMLString(post.Comments),
			post.ID,
			post.CSRFToken,
		)
	}
	html += `</div>`

	return html
}

func getPostCommentsHTMLString(comments []Comment) string {
	var html string
	for _, comment := range comments {
		html += fmt.Sprintf(`
	    <div class="isu-comment">
      <a href="/@%s" class="isu-comment-account-name">%s</a>
      <span class="isu-comment-text">%s</span>
    </div>
`,
			comment.User.AccountName,
			comment.User.AccountName,
			comment.Comment,
		)
	}

	return html
}
