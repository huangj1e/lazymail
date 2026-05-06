package mail

import (
	"fmt"
	"io"
	"strings"

	imap "github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"

	"lazymail/internal/config"
	"lazymail/internal/domain"
)

// Client wraps an IMAP connection for a single account.
type Client struct {
	c       *imapclient.Client
	account config.Account
}

// Dial connects and logs in to the IMAP server described by acc.
func Dial(acc config.Account) (*Client, error) {
	addr := fmt.Sprintf("%s:%d", acc.IMAPHost, acc.IMAPPort)
	var (
		c   *imapclient.Client
		err error
	)
	if acc.TLS {
		c, err = imapclient.DialTLS(addr, nil)
	} else {
		c, err = imapclient.DialStartTLS(addr, nil)
	}
	if err != nil {
		return nil, fmt.Errorf("imap: dial %s: %w", addr, err)
	}

	if err := c.Login(acc.Username, acc.Password).Wait(); err != nil {
		c.Close()
		return nil, fmt.Errorf("imap: login %s: %w", acc.Username, err)
	}
	return &Client{c: c, account: acc}, nil
}

// Close logs out and terminates the connection.
func (cl *Client) Close() {
	_ = cl.c.Logout().Wait()
	cl.c.Close()
}

// ListFolders returns the mailbox names available on the server.
func (cl *Client) ListFolders() ([]string, error) {
	mailboxes, err := cl.c.List("", "%", nil).Collect()
	if err != nil {
		return nil, fmt.Errorf("imap: list: %w", err)
	}
	names := make([]string, 0, len(mailboxes))
	for _, mb := range mailboxes {
		names = append(names, mb.Mailbox)
	}
	return names, nil
}

// FetchMails fetches the latest `limit` messages from `folder`.
func (cl *Client) FetchMails(folder string, limit uint32) ([]domain.Mail, error) {
	mboxStatus, err := cl.c.Select(folder, nil).Wait()
	if err != nil {
		return nil, fmt.Errorf("imap: select %s: %w", folder, err)
	}

	total := mboxStatus.NumMessages
	if total == 0 {
		return nil, nil
	}

	var start uint32 = 1
	if total > limit {
		start = total - limit + 1
	}

	var seqSet imap.SeqSet
	seqSet.AddRange(start, total)
	fetchOptions := &imap.FetchOptions{
		Envelope: true,
		Flags:    true,
		BodySection: []*imap.FetchItemBodySection{
			{Specifier: imap.PartSpecifierText},
		},
	}

	messages, err := cl.c.Fetch(seqSet, fetchOptions).Collect()
	if err != nil {
		return nil, fmt.Errorf("imap: fetch: %w", err)
	}

	mails := make([]domain.Mail, 0, len(messages))
	for _, msg := range messages {
		if msg.Envelope == nil {
			continue
		}
		m := domain.Mail{
			ID:      fmt.Sprintf("%s-%d", folder, msg.UID),
			Subject: msg.Envelope.Subject,
			Date:    msg.Envelope.Date,
			Folder:  folder,
		}

		if len(msg.Envelope.From) > 0 {
			m.From = msg.Envelope.From[0].Addr()
		}
		for _, addr := range msg.Envelope.To {
			m.To = append(m.To, addr.Addr())
		}

		m.IsRead = hasFlag(msg.Flags, imap.FlagSeen)

		// extract text body from body section buffers
		for _, bs := range msg.BodySection {
			m.Body = strings.TrimSpace(string(bs.Bytes))
			break
		}

		mails = append(mails, m)
	}

	// reverse so newest is first
	for i, j := 0, len(mails)-1; i < j; i, j = i+1, j-1 {
		mails[i], mails[j] = mails[j], mails[i]
	}
	return mails, nil
}

// MarkRead sends a STORE FLAGS command to set/clear \Seen on a UID.
func (cl *Client) MarkRead(folder string, uid uint32, read bool) error {
	if _, err := cl.c.Select(folder, nil).Wait(); err != nil {
		return fmt.Errorf("imap: select: %w", err)
	}
	op := imap.StoreFlagsDel
	if read {
		op = imap.StoreFlagsAdd
	}
	storeFlags := &imap.StoreFlags{
		Op:    op,
		Flags: []imap.Flag{imap.FlagSeen},
	}
	var uidSet imap.UIDSet
	uidSet.AddNum(imap.UID(uid))
	return cl.c.Store(uidSet, storeFlags, nil).Close()
}

func hasFlag(flags []imap.Flag, target imap.Flag) bool {
	for _, f := range flags {
		if f == target {
			return true
		}
	}
	return false
}

// bodyReader is used internally; kept for future MIME parsing.
var _ io.Reader = (*strings.Reader)(nil)
