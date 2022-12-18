package main

import "fmt"

func getPostsHTMLString(posts []Post) string {
	var html string
	html += `<div class="isu-posts">`
	for i := range posts {
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
			posts[i].ID,
			posts[i].CreatedAt.Format("2006-01-02T15:04:05-07:00"),
			posts[i].User.AccountName,
			posts[i].User.AccountName,
			posts[i].ID,
			posts[i].CreatedAt.Format("2006-01-02T15:04:05-07:00"),
			imageURL(posts[i]),
			posts[i].User.AccountName,
			posts[i].User.AccountName,
			posts[i].Body,
			posts[i].CommentCount,
			getPostCommentsHTMLString(posts[i].Comments),
			posts[i].ID,
			posts[i].CSRFToken,
		)
	}
	html += `</div>`

	return html
}

func getPostCommentsHTMLString(comments []Comment) string {
	var html string
	for i := range comments {
		html += fmt.Sprintf(`
	    <div class="isu-comment">
      <a href="/@%s" class="isu-comment-account-name">%s</a>
      <span class="isu-comment-text">%s</span>
    </div>
`,
			comments[i].User.AccountName,
			comments[i].User.AccountName,
			comments[i].Comment,
		)
	}

	return html
}
