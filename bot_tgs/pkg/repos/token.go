package repos

type Bucket string

const (
	AccessTokens Bucket = "access_tokens"
	RequestToken Bucket = "request_tokens"
)

type TokenRepository interface {
	Save(chatID int64, token string, bucket Bucket) error
	Get(chatID int64, bucket Bucket) (string, error)
}
